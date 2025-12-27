package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ackkerman/x-post-preview-generator/internal/render"
)

func main() {
	opts := render.DefaultOptions()

	text := flag.String("text", "", "ツイート本文")
	icon := flag.String("icon", "", "アイコン画像パスまたはURL")
	name := flag.String("name", "", "表示名")
	handle := flag.String("id", "", "ユーザーID (@なし可)")
	date := flag.String("date", "", "日付(任意)")
	location := flag.String("location", "", "現在地(任意)")
	cta := flag.String("cta", "Explore what's happening on Twitter", "CTAボタン文言(空で非表示)")
	noCTA := flag.Bool("no-cta", false, "CTAを非表示にする")
	verified := flag.Bool("verified", false, "認証バッジを表示する")
	simple := flag.Bool("simple", false, "Simpleモード(フッター非表示)")
	likeCount := flag.String("like-count", "0", "Like件数表示")
	output := flag.String("output", "tweet.png", "出力ファイルパス")
	format := flag.String("format", "", "出力形式: png|jpg|jpeg|gif|svg|html (省略時は拡張子から推定)")
	width := flag.Int("width", opts.Width, "出力幅(px)")
	widthMode := flag.String("width-mode", opts.WidthMode, "横幅モード: fixed|tight")
	padding := flag.Int("padding", opts.Padding, "余白(px)")
	theme := flag.String("theme", "light", "テーマ: light|dark")
	fontPath := flag.String("font", "", "本文フォントのパス(.ttf/.otf)")
	fontBoldPath := flag.String("font-bold", "", "太字フォントのパス(.ttf/.otf)")
	fontFamily := flag.String("font-family", opts.FontFamily, "HTML/SVG用のfont-family")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "X投稿プレビュー生成CLI\n\n")
		fmt.Fprintf(os.Stderr, "必須: -text, -name, -id\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if strings.TrimSpace(*text) == "" || strings.TrimSpace(*name) == "" || strings.TrimSpace(*handle) == "" {
		flag.Usage()
		os.Exit(2)
	}

	selectedTheme, err := parseTheme(*theme)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmtValue := *format
	if strings.TrimSpace(fmtValue) == "" {
		fmtValue = inferFormat(*output)
	}
	if fmtValue == "" {
		fmtValue = "png"
	}

	data := render.TweetData{
		Text:      *text,
		Icon:      *icon,
		Name:      *name,
		Handle:    *handle,
		Date:      *date,
		Location:  *location,
		CTA:       *cta,
		Verified:  *verified,
		Simple:    *simple,
		LikeCount: *likeCount,
	}
	if *noCTA {
		data.CTA = ""
	}

	opts.Width = *width
	opts.WidthMode = *widthMode
	opts.Padding = *padding
	opts.FontPath = *fontPath
	opts.BoldFontPath = *fontBoldPath
	opts.FontFamily = *fontFamily
	opts.Theme = selectedTheme

	if err := writeOutput(*output, data, opts, fmtValue); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func writeOutput(path string, data render.TweetData, opts render.RenderOptions, format string) error {
	if path == "-" {
		return render.RenderToWriter(os.Stdout, data, opts, format)
	}
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return render.RenderToWriter(file, data, opts, format)
}

func inferFormat(output string) string {
	ext := strings.ToLower(filepath.Ext(output))
	switch ext {
	case ".png":
		return "png"
	case ".jpg", ".jpeg":
		return "jpeg"
	case ".gif":
		return "gif"
	case ".svg":
		return "svg"
	case ".html", ".htm":
		return "html"
	default:
		return ""
	}
}

func parseTheme(value string) (render.Theme, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "light":
		return render.LightTheme(), nil
	case "dark":
		return render.DarkTheme(), nil
	default:
		return render.Theme{}, fmt.Errorf("unknown theme: %s", value)
	}
}
