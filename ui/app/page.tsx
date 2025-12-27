"use client";

import { useEffect, useMemo, useState } from "react";
import { Copy, Download, RefreshCcw, Sparkles } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import { type PostConfig } from "@/lib/post-config";
import { renderSvg } from "@/lib/wasm";

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

const presets: Array<{ label: string; value: PostConfig }> = [
  {
    label: "Launch Day",
    value: {
      ...defaultConfig,
      text: "We just shipped v2.0. Cleaner layout, faster renders, and export-ready previews. Big thanks to everyone who tested it!",
      cta: "See the roadmap",
      likeCount: "4.2K"
    }
  },
  {
    label: "Minimal",
    value: {
      ...defaultConfig,
      mode: "simple",
      verified: false,
      likeCount: "",
      cta: ""
    }
  }
];

export default function Home() {
  const [config, setConfig] = useState<PostConfig>(defaultConfig);
  const [copied, setCopied] = useState(false);
  const [svgMarkup, setSvgMarkup] = useState<string>("");
  const [previewStatus, setPreviewStatus] = useState<"idle" | "loading" | "ready" | "error">("idle");
  const [previewError, setPreviewError] = useState<string>("");

  const payload = useMemo(() => JSON.stringify(config, null, 2), [config]);

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
        setPreviewError(error instanceof Error ? error.message : "Failed to render preview");
      });

    return () => {
      cancelled = true;
    };
  }, [config]);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(payload);
      setCopied(true);
      setTimeout(() => setCopied(false), 1600);
    } catch {
      setCopied(false);
    }
  };

  const handleDownload = () => {
    const blob = new Blob([payload], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = "x-post-config.json";
    link.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="relative">
      <div className="pointer-events-none absolute inset-0 grid-overlay opacity-30" />
      <div className="mx-auto flex min-h-screen max-w-6xl flex-col gap-10 px-6 py-12">
        <header className="flex flex-col gap-4">
          <div className="flex items-center gap-3 text-sm font-semibold uppercase tracking-[0.3em] text-muted">
            <Sparkles className="h-4 w-4 text-accent" />
            X Post Preview Studio
          </div>
          <h1 className="max-w-2xl text-4xl font-semibold leading-tight text-ink md:text-5xl">
            Shape your next X post before it hits the feed.
          </h1>
          <p className="max-w-2xl text-base text-muted">
            Tune voice, layout, and engagement cues in one place. Export the config and move fast with a consistent preview.
          </p>
        </header>

        <main className="grid gap-8 lg:grid-cols-[minmax(0,1.05fr)_minmax(0,0.95fr)]">
          <section className="space-y-4">
            <div className="sticky top-6 space-y-4">
              <div className="panel flex min-h-[280px] items-center justify-center p-4">
                {previewStatus === "loading" ? (
                  <div className="text-sm text-muted">Rendering preview...</div>
                ) : null}
                {previewStatus === "error" ? (
                  <div className="space-y-2 text-center text-sm text-muted">
                    <p>Preview failed.</p>
                    <p className="text-xs">{previewError}</p>
                    <p className="text-xs">Run `make ui-wasm` to rebuild the wasm bundle.</p>
                  </div>
                ) : null}
                {previewStatus === "ready" ? (
                  <div
                    className="w-full [&>svg]:h-auto [&>svg]:w-full"
                    dangerouslySetInnerHTML={{ __html: svgMarkup }}
                  />
                ) : null}
              </div>
            </div>
          </section>
        
          <Card className="animate-fade-up">
            <CardHeader>
              <CardTitle>Post controls</CardTitle>
              <CardDescription>Update the inputs to see the preview refresh instantly.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="text">Post text</Label>
                <Textarea
                  id="text"
                  value={config.text}
                  onChange={(event) => setConfig({ ...config, text: event.target.value })}
                />
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="name">Display name</Label>
                  <Input
                    id="name"
                    value={config.name}
                    onChange={(event) => setConfig({ ...config, name: event.target.value })}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="handle">Handle</Label>
                  <Input
                    id="handle"
                    value={config.handle}
                    onChange={(event) => setConfig({ ...config, handle: event.target.value })}
                  />
                </div>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="avatar">Avatar image URL</Label>
                  <Input
                    id="avatar"
                    placeholder="https://"
                    value={config.avatarUrl}
                    onChange={(event) => setConfig({ ...config, avatarUrl: event.target.value })}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="date">Date line</Label>
                  <Input
                    id="date"
                    value={config.date}
                    onChange={(event) => setConfig({ ...config, date: event.target.value })}
                  />
                </div>
              </div>

              <div className="grid gap-4 md:grid-cols-2">
                <div className="space-y-2">
                  <Label htmlFor="likes">Like count</Label>
                  <Input
                    id="likes"
                    value={config.likeCount}
                    onChange={(event) => setConfig({ ...config, likeCount: event.target.value })}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="cta">CTA text</Label>
                  <Input
                    id="cta"
                    value={config.cta}
                    onChange={(event) => setConfig({ ...config, cta: event.target.value })}
                  />
                </div>
              </div>

              <div className="grid gap-4 md:grid-cols-3">
                <div className="space-y-2">
                  <Label>Layout</Label>
                  <Select
                    value={config.mode}
                    onValueChange={(value) =>
                      setConfig({ ...config, mode: value as PostConfig["mode"] })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select layout" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="classic">Classic</SelectItem>
                      <SelectItem value="simple">Simple</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>Width</Label>
                  <Select
                    value={config.width}
                    onValueChange={(value) =>
                      setConfig({ ...config, width: value as PostConfig["width"] })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Select width" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="tight">Tight</SelectItem>
                      <SelectItem value="wide">Wide</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="flex items-center justify-between gap-3 rounded-2xl border border-border bg-white/70 px-4 py-3">
                  <div>
                    <Label htmlFor="verified">Verified</Label>
                    <p className="text-xs text-muted">Show the badge</p>
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
                <Label>Quick presets</Label>
                <div className="flex flex-wrap gap-3">
                  {presets.map((preset) => (
                    <Button
                      key={preset.label}
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => setConfig(preset.value)}
                    >
                      {preset.label}
                    </Button>
                  ))}
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => setConfig(defaultConfig)}
                  >
                    <RefreshCcw className="h-4 w-4" />
                    Reset
                  </Button>
                </div>
              </div>

              <Separator />

              <div className="flex flex-wrap gap-3">
                <Button type="button" onClick={handleCopy} variant="secondary">
                  <Copy className="h-4 w-4" />
                  {copied ? "Copied" : "Copy JSON"}
                </Button>
                <Button type="button" onClick={handleDownload} variant="outline">
                  <Download className="h-4 w-4" />
                  Download JSON
                </Button>
              </div>
            </CardContent>
          </Card>
        </main>
      </div>
    </div>
  );
}
