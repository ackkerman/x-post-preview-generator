export const translations = {
  en: {
    language: {
      label: "Language",
      placeholder: "Select language"
    },
    app: {
      name: "X Post Preview Studio",
      title: "Shape your next X post before it hits the feed.",
      description:
        "Tune voice, layout, and engagement cues in one place. Export the preview and move fast with a consistent look."
    },
    preview: {
      loading: "Rendering preview...",
      failed: "Preview failed.",
      rebuildHint: "Run `make ui-wasm` to rebuild the wasm bundle."
    },
    controls: {
      title: "Post controls",
      description: "Update the inputs to see the preview refresh instantly."
    },
    labels: {
      text: "Post text",
      name: "Display name",
      handle: "Handle",
      avatar: "Avatar image URL",
      date: "Date line",
      likes: "Like count",
      cta: "CTA text",
      layout: "Layout",
      width: "Width",
      verified: "Verified",
      verifiedHint: "Show the badge",
      presets: "Quick presets",
      exportFormat: "Export format"
    },
    placeholders: {
      layout: "Select layout",
      width: "Select width",
      format: "Select format"
    },
    options: {
      classic: "Classic",
      simple: "Simple",
      tight: "Tight",
      wide: "Wide",
      svg: "SVG (vector)",
      png: "PNG (raster)"
    },
    presets: {
      launchDay: "Launch Day",
      minimal: "Minimal"
    },
    buttons: {
      reset: "Reset",
      copy: "Copy",
      copied: "Copied",
      download: "Download"
    },
    toast: {
      copySuccess: "Copied {format} to clipboard."
    },
    errors: {
      previewNotReady: "Preview is not ready",
      clipboardUnavailable: "Clipboard API is not available",
      pngExportFailed: "Failed to export PNG",
      copyFailed: "Failed to copy preview",
      downloadFailed: "Failed to download preview"
    }
  },
  ja: {
    language: {
      label: "言語",
      placeholder: "言語を選択"
    },
    app: {
      name: "X投稿プレビュースタジオ",
      title: "投稿前にXのプレビューを整えよう。",
      description: "文体やレイアウト、反応のニュアンスをまとめて調整。プレビューを書き出して素早く共有できます。"
    },
    preview: {
      loading: "プレビューを生成中...",
      failed: "プレビュー生成に失敗しました。",
      rebuildHint: "`make ui-wasm` を実行してWasmを再生成してください。"
    },
    controls: {
      title: "投稿設定",
      description: "入力を更新するとプレビューが即時に反映されます。"
    },
    labels: {
      text: "投稿本文",
      name: "表示名",
      handle: "ハンドル",
      avatar: "アバター画像URL",
      date: "日時行",
      likes: "いいね数",
      cta: "CTA文言",
      layout: "レイアウト",
      width: "幅",
      verified: "認証",
      verifiedHint: "バッジを表示",
      presets: "クイックプリセット",
      exportFormat: "書き出し形式"
    },
    placeholders: {
      layout: "レイアウトを選択",
      width: "幅を選択",
      format: "形式を選択"
    },
    options: {
      classic: "クラシック",
      simple: "シンプル",
      tight: "タイト",
      wide: "ワイド",
      svg: "SVG（ベクター）",
      png: "PNG（ラスター）"
    },
    presets: {
      launchDay: "ローンチデイ",
      minimal: "ミニマル"
    },
    buttons: {
      reset: "リセット",
      copy: "コピー",
      copied: "コピー済み",
      download: "ダウンロード"
    },
    toast: {
      copySuccess: "{format} をクリップボードにコピーしました。"
    },
    errors: {
      previewNotReady: "プレビューの準備ができていません。",
      clipboardUnavailable: "クリップボードAPIが利用できません。",
      pngExportFailed: "PNGの書き出しに失敗しました。",
      copyFailed: "プレビューのコピーに失敗しました。",
      downloadFailed: "プレビューのダウンロードに失敗しました。"
    }
  }
} as const;

export type Language = keyof typeof translations;

export const languageOptions: Array<{ value: Language; label: string }> = [
  { value: "en", label: "English" },
  { value: "ja", label: "日本語" }
];
