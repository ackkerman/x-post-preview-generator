"use client";

import { useEffect, useState } from "react";
import { Copy, Download, RefreshCcw, Sparkles } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import { languageOptions, translations, type Language } from "@/lib/i18n";
import { type PostConfig } from "@/lib/post-config";
import { renderSvg } from "@/lib/wasm";

type ExportFormat = "svg" | "png";

function parseSvgSize(svgMarkup: string) {
  const document = new DOMParser().parseFromString(svgMarkup, "image/svg+xml");
  const svg = document.querySelector("svg");
  if (!svg) {
    throw new Error("SVG parse failed");
  }

  const widthAttr = svg.getAttribute("width");
  const heightAttr = svg.getAttribute("height");
  const width = widthAttr ? Number.parseFloat(widthAttr) : Number.NaN;
  const height = heightAttr ? Number.parseFloat(heightAttr) : Number.NaN;

  if (Number.isFinite(width) && Number.isFinite(height)) {
    return { width, height };
  }

  const viewBox = svg.getAttribute("viewBox");
  if (viewBox) {
    const parts = viewBox.split(/[\s,]+/).map((value) => Number.parseFloat(value));
    if (parts.length === 4 && parts.every((value) => Number.isFinite(value))) {
      return { width: parts[2], height: parts[3] };
    }
  }

  return { width: 960, height: 540 };
}

async function svgToPngBlob(svgMarkup: string) {
  const { width, height } = parseSvgSize(svgMarkup);
  const svgBlob = new Blob([svgMarkup], { type: "image/svg+xml" });
  const url = URL.createObjectURL(svgBlob);

  return new Promise<Blob>((resolve, reject) => {
    const img = new Image();
    img.onload = () => {
      const canvas = document.createElement("canvas");
      canvas.width = Math.max(1, Math.floor(width));
      canvas.height = Math.max(1, Math.floor(height));
      const context = canvas.getContext("2d");
      if (!context) {
        URL.revokeObjectURL(url);
        reject(new Error("Canvas is not available"));
        return;
      }
      context.drawImage(img, 0, 0, canvas.width, canvas.height);
      canvas.toBlob((blob) => {
        URL.revokeObjectURL(url);
        if (!blob) {
          reject(new Error("Failed to export PNG"));
          return;
        }
        resolve(blob);
      }, "image/png");
    };
    img.onerror = () => {
      URL.revokeObjectURL(url);
      reject(new Error("Failed to load SVG for PNG export"));
    };
    img.src = url;
  });
}

function stripXmlDeclaration(svgMarkup: string) {
  return svgMarkup.replace(/^<\?xml[^>]*>\s*/i, "");
}

const defaultConfig: PostConfig = {
  text: "Just shipped a preview generator that turns raw text into a polished X post mock. What should we build next?",
  name: "Rina Sato",
  handle: "rinasato",
  verified: true,
  avatarUrl: "",
  date: "9:32 AM · Aug 15, 2024",
  likeCount: "12.8K",
  cta: "Read 120 replies",
  width: "tight",
  mode: "classic"
};

const presets = [
  {
    key: "launchDay",
    value: {
      ...defaultConfig,
      text: "We just shipped v2.0. Cleaner layout, faster renders, and export-ready previews. Big thanks to everyone who tested it!",
      cta: "See the roadmap",
      likeCount: "4.2K"
    }
  },
  {
    key: "minimal",
    value: {
      ...defaultConfig,
      mode: "simple",
      verified: false,
      likeCount: "",
      cta: ""
    }
  }
] as const;

export default function Home() {
  const [language, setLanguage] = useState<Language>("en");
  const [config, setConfig] = useState<PostConfig>(defaultConfig);
  const [exportFormat, setExportFormat] = useState<ExportFormat>("svg");
  const [copyStatus, setCopyStatus] = useState<"idle" | "success">("idle");
  const [exportError, setExportError] = useState("");
  const [svgMarkup, setSvgMarkup] = useState<string>("");
  const [previewStatus, setPreviewStatus] = useState<"idle" | "loading" | "ready" | "error">("idle");
  const [previewError, setPreviewError] = useState<string>("");
  const t = translations[language];

  useEffect(() => {
    let cancelled = false;
    setPreviewStatus("loading");
    setPreviewError("");

    renderSvg(config)
      .then((svg) => {
        if (cancelled) return;
        setSvgMarkup(svg);
        setPreviewStatus("ready");
      })
      .catch((error) => {
        if (cancelled) return;
        setPreviewStatus("error");
        setPreviewError(error instanceof Error ? error.message : "");
      });

    return () => {
      cancelled = true;
    };
  }, [config]);

  const handleCopy = async () => {
    setExportError("");
    if (!svgMarkup) {
      setExportError(t.errors.previewNotReady);
      return;
    }

    let blob: Blob;
    let mimeType: string;
    if (exportFormat === "svg") {
      blob = new Blob([svgMarkup], { type: "image/svg+xml" });
      mimeType = "image/svg+xml";
    } else {
      try {
        blob = await svgToPngBlob(svgMarkup);
      } catch {
        setExportError(t.errors.pngExportFailed);
        return;
      }
      mimeType = "image/png";
    }

    if (!navigator.clipboard?.write) {
      if (exportFormat === "svg" && navigator.clipboard?.writeText) {
        try {
          await navigator.clipboard.writeText(svgMarkup);
          setCopyStatus("success");
          toast.success(t.toast.copySuccess.replace("{format}", t.options.svg));
          setTimeout(() => setCopyStatus("idle"), 1600);
        } catch {
          setCopyStatus("idle");
          setExportError(t.errors.copyFailed);
        }
        return;
      }
      setExportError(t.errors.clipboardUnavailable);
      return;
    }

    try {
      await navigator.clipboard.write([new ClipboardItem({ [mimeType]: blob })]);
      setCopyStatus("success");
      const formatLabel = exportFormat === "svg" ? t.options.svg : t.options.png;
      toast.success(t.toast.copySuccess.replace("{format}", formatLabel));
      setTimeout(() => setCopyStatus("idle"), 1600);
    } catch {
      if (exportFormat === "svg" && navigator.clipboard?.writeText) {
        try {
          await navigator.clipboard.writeText(svgMarkup);
          setCopyStatus("success");
          setTimeout(() => setCopyStatus("idle"), 1600);
          return;
        } catch {
          setCopyStatus("idle");
          setExportError(t.errors.copyFailed);
          return;
        }
      }
      setCopyStatus("idle");
      setExportError(t.errors.copyFailed);
    }
  };

  const handleDownload = async () => {
    setExportError("");
    if (!svgMarkup) {
      setExportError(t.errors.previewNotReady);
      return;
    }

    let blob: Blob;
    if (exportFormat === "svg") {
      blob = new Blob([svgMarkup], { type: "image/svg+xml" });
    } else {
      try {
        blob = await svgToPngBlob(svgMarkup);
      } catch {
        setExportError(t.errors.pngExportFailed);
        return;
      }
    }

    try {
      const url = URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `x-post-preview.${exportFormat}`;
      link.click();
      URL.revokeObjectURL(url);
    } catch {
      setExportError(t.errors.downloadFailed);
    }
  };

  return (
    <div className="relative">
      <div className="pointer-events-none absolute inset-0 grid-overlay opacity-30" />
      <div className="mx-auto flex min-h-screen max-w-6xl flex-col gap-10 px-6 py-12">
        <header className="flex flex-col gap-4">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div className="flex items-center gap-3 text-sm font-semibold uppercase tracking-[0.3em] text-muted">
              <Sparkles className="h-4 w-4 text-accent" />
              {t.app.name}
            </div>
            <div className="flex items-center gap-3">
              <Label htmlFor="language" className="text-xs font-semibold uppercase tracking-[0.3em] text-muted">
                {t.language.label}
              </Label>
              <Select
                value={language}
                onValueChange={(value) => setLanguage(value as Language)}
              >
                <SelectTrigger id="language" className="h-9 w-40">
                  <SelectValue placeholder={t.language.placeholder} />
                </SelectTrigger>
                <SelectContent>
                  {languageOptions.map((option) => (
                    <SelectItem key={option.value} value={option.value}>
                      {option.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          <h1 className="max-w-2xl text-4xl font-semibold leading-tight text-ink md:text-5xl">
            {t.app.title}
          </h1>
          <p className="max-w-2xl text-base text-muted">{t.app.description}</p>
        </header>

        <main className="grid gap-8 lg:grid-cols-[minmax(0,1.05fr)_minmax(0,0.95fr)]">
        
          <section className="space-y-4">
            <div className="sticky top-6 space-y-4">
              <div className="panel flex min-h-[280px] items-center justify-center p-4">
                {previewStatus === "loading" ? (
                  <div className="text-sm text-muted">{t.preview.loading}</div>
                ) : null}
                {previewStatus === "error" ? (
                  <div className="space-y-2 text-center text-sm text-muted">
                    <p>{t.preview.failed}</p>
                    <p className="text-xs">{previewError}</p>
                    <p className="text-xs">{t.preview.rebuildHint}</p>
                  </div>
                ) : null}
                {previewStatus === "ready" ? (
                  <div
                    className="
                      w-full
                      [&>svg]:h-auto
                      [&>svg]:w-full
                      [&>svg]:drop-shadow-[0_4px_12px_rgba(0,0,0,0.25)]
                    "
                    dangerouslySetInnerHTML={{ __html: stripXmlDeclaration(svgMarkup) }}
                  />
                ) : null}
              </div>
            </div>
          </section>
        
          <Card className="animate-fade-up">
            <CardHeader>
              <CardTitle>{t.controls.title}</CardTitle>
              <CardDescription>{t.controls.description}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="text">{t.labels.text}</Label>
                <Textarea
                  id="text"
                  value={config.text}
                  onChange={(event) => setConfig({ ...config, text: event.target.value })}
                />
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="name">{t.labels.name}</Label>
                  <Input
                    id="name"
                    value={config.name}
                    onChange={(event) => setConfig({ ...config, name: event.target.value })}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="handle">{t.labels.handle}</Label>
                  <Input
                    id="handle"
                    value={config.handle}
                    onChange={(event) => setConfig({ ...config, handle: event.target.value })}
                  />
                </div>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="avatar">{t.labels.avatar}</Label>
                  <Input
                    id="avatar"
                    placeholder="https://"
                    value={config.avatarUrl}
                    onChange={(event) => setConfig({ ...config, avatarUrl: event.target.value })}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="date">{t.labels.date}</Label>
                  <Input
                    id="date"
                    value={config.date}
                    onChange={(event) => setConfig({ ...config, date: event.target.value })}
                  />
                </div>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="likes">{t.labels.likes}</Label>
                  <Input
                    id="likes"
                    value={config.likeCount}
                    onChange={(event) => setConfig({ ...config, likeCount: event.target.value })}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="cta">{t.labels.cta}</Label>
                  <Input
                    id="cta"
                    value={config.cta}
                    onChange={(event) => setConfig({ ...config, cta: event.target.value })}
                  />
                </div>
              </div>

              <div className="grid gap-4 md:grid-cols-3">
                <div className="space-y-2">
                  <Label>{t.labels.layout}</Label>
                  <Select
                    value={config.mode}
                    onValueChange={(value) =>
                      setConfig({ ...config, mode: value as PostConfig["mode"] })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={t.placeholders.layout} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="classic">{t.options.classic}</SelectItem>
                      <SelectItem value="simple">{t.options.simple}</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>{t.labels.width}</Label>
                  <Select
                    value={config.width}
                    onValueChange={(value) =>
                      setConfig({ ...config, width: value as PostConfig["width"] })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue placeholder={t.placeholders.width} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="tight">{t.options.tight}</SelectItem>
                      <SelectItem value="wide">{t.options.wide}</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="flex items-center justify-between gap-3 rounded-2xl border border-border bg-white/70 px-4 py-3">
                  <div>
                    <Label htmlFor="verified">{t.labels.verified}</Label>
                    <p className="text-xs text-muted">{t.labels.verifiedHint}</p>
                  </div>
                  <Switch
                    id="verified"
                    checked={config.verified}
                    onCheckedChange={(checked) => setConfig({ ...config, verified: checked })}
                  />
                </div>
              </div>

              <Separator />

              <div className="space-y-3">
                <Label>{t.labels.presets}</Label>
                <div className="flex flex-wrap gap-3">
                  {presets.map((preset) => (
                    <Button
                      key={preset.key}
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => setConfig(preset.value)}
                    >
                      {t.presets[preset.key]}
                    </Button>
                  ))}
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => setConfig(defaultConfig)}
                  >
                    <RefreshCcw className="h-4 w-4" />
                    {t.buttons.reset}
                  </Button>
                </div>
              </div>

              <Separator />

              <div className="space-y-3">
                <Label>{t.labels.exportFormat}</Label>
                <Select
                  value={exportFormat}
                  onValueChange={(value) => setExportFormat(value as ExportFormat)}
                >
                  <SelectTrigger>
                    <SelectValue placeholder={t.placeholders.format} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="svg">{t.options.svg}</SelectItem>
                    <SelectItem value="png">{t.options.png}</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="flex flex-wrap gap-3">
                <Button type="button" onClick={handleCopy} variant="secondary">
                  <Copy className="h-4 w-4" />
                  {copyStatus === "success" ? t.buttons.copied : t.buttons.copy}
                </Button>
                <Button type="button" onClick={handleDownload} variant="outline">
                  <Download className="h-4 w-4" />
                  {t.buttons.download}
                </Button>
              </div>
              {exportError ? <p className="text-xs text-rose-600">{exportError}</p> : null}
            </CardContent>
          </Card>
        </main>
      </div>
    </div>
  );
}
