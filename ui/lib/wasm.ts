import type { PostConfig } from "@/lib/post-config";

type GoRuntime = {
  importObject: WebAssembly.Imports;
  run: (instance: WebAssembly.Instance) => Promise<void>;
};

declare global {
  interface Window {
    Go?: new () => GoRuntime;
    xpostgenRender?: (payload: string) => Promise<string>;
  }
}

let goRuntime: GoRuntime | null = null;
let wasmReady: Promise<void> | null = null;

function loadScript(src: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const existing = document.querySelector(`script[src="${src}"]`);
    if (existing) {
      resolve();
      return;
    }
    const script = document.createElement("script");
    script.src = src;
    script.async = true;
    script.onload = () => resolve();
    script.onerror = () => reject(new Error(`Failed to load ${src}`));
    document.head.appendChild(script);
  });
}

async function instantiateWasm(url: string, go: GoRuntime) {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error("Failed to load wasm module");
  }

  if ("instantiateStreaming" in WebAssembly) {
    const cloned = response.clone();
    try {
      return await WebAssembly.instantiateStreaming(cloned, go.importObject);
    } catch {
      const buffer = await response.arrayBuffer();
      return await WebAssembly.instantiate(buffer, go.importObject);
    }
  }

  const buffer = await response.arrayBuffer();
  return await WebAssembly.instantiate(buffer, go.importObject);
}

function waitForRenderFn(timeoutMs = 5000) {
  return new Promise<void>((resolve, reject) => {
    const startedAt = Date.now();
    const timer = window.setInterval(() => {
      if (window.xpostgenRender) {
        window.clearInterval(timer);
        resolve();
        return;
      }
      if (Date.now() - startedAt > timeoutMs) {
        window.clearInterval(timer);
        reject(new Error("xpostgenRender not available"));
      }
    }, 50);
  });
}

export async function ensureWasmReady() {
  if (typeof window === "undefined") {
    throw new Error("Wasm runtime is not available on the server");
  }

  if (window.xpostgenRender) {
    return;
  }

  if (!wasmReady) {
    wasmReady = (async () => {
      if (!window.Go) {
        await loadScript("/wasm/wasm_exec.js");
      }
      if (!window.Go) {
        throw new Error("Go runtime not found");
      }

      goRuntime = new window.Go();
      const { instance } = await instantiateWasm("/wasm/xpostgen.wasm", goRuntime);
      void goRuntime.run(instance);
      await waitForRenderFn();
    })().catch((error) => {
      wasmReady = null;
      throw error;
    });
  }

  return wasmReady;
}

export async function renderSvg(config: PostConfig) {
  await ensureWasmReady();
  if (!window.xpostgenRender) {
    throw new Error("xpostgenRender not ready");
  }
  return window.xpostgenRender(JSON.stringify(config));
}
