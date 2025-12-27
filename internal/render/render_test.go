package render

import (
	"strings"
	"testing"
)

func TestWrapTextBreaksLines(t *testing.T) {
	fonts, err := loadFontSet(DefaultOptions())
	if err != nil {
		t.Fatalf("loadFontSet: %v", err)
	}
	defer fonts.Close()

	maxWidth := measureString(fonts.Text, "hello")
	lines := wrapText("hello world", maxWidth, fonts.Text)
	if len(lines) < 2 {
		t.Fatalf("expected wrapped lines, got %v", lines)
	}
}

func TestWrapTextNoSpaces(t *testing.T) {
	fonts, err := loadFontSet(DefaultOptions())
	if err != nil {
		t.Fatalf("loadFontSet: %v", err)
	}
	defer fonts.Close()

	maxWidth := measureString(fonts.Text, "hello")
	lines := wrapText("helloworldhelloworld", maxWidth, fonts.Text)
	if len(lines) < 2 {
		t.Fatalf("expected no-space text to wrap, got %v", lines)
	}
}

func TestEllipsize(t *testing.T) {
	fonts, err := loadFontSet(DefaultOptions())
	if err != nil {
		t.Fatalf("loadFontSet: %v", err)
	}
	defer fonts.Close()

	maxWidth := measureString(fonts.Meta, "short")
	result := ellipsize("this is a long name", maxWidth, fonts.Meta)
	if result == "this is a long name" {
		t.Fatalf("expected ellipsized string, got %s", result)
	}
	if !strings.HasSuffix(result, "...") {
		t.Fatalf("expected ellipsis, got %s", result)
	}
}

func TestRenderOutputs(t *testing.T) {
	data := TweetData{
		Text:     "Hello X preview",
		Name:     "Example User",
		Handle:   "example",
		Date:     "2025-01-01",
		CTA:      "Read 16K replies",
		Verified: true,
	}

	img, err := RenderImage(data, DefaultOptions())
	if err != nil {
		t.Fatalf("RenderImage: %v", err)
	}
	if img.Bounds().Dx() == 0 || img.Bounds().Dy() == 0 {
		t.Fatalf("image has invalid bounds")
	}

	svg, err := RenderSVG(data, DefaultOptions())
	if err != nil {
		t.Fatalf("RenderSVG: %v", err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Fatalf("svg output missing <svg>")
	}

	html, err := RenderHTML(data, DefaultOptions())
	if err != nil {
		t.Fatalf("RenderHTML: %v", err)
	}
	if !strings.Contains(html, "<html") {
		t.Fatalf("html output missing <html>")
	}
}

func TestIconsLoad(t *testing.T) {
	icon, err := iconSVG("reply")
	if err != nil {
		t.Fatalf("iconSVG: %v", err)
	}
	if !strings.Contains(icon, "<svg") {
		t.Fatalf("icon svg missing <svg>")
	}
	img, err := rasterizeIcon("like", "#000000", 16)
	if err != nil {
		t.Fatalf("rasterizeIcon: %v", err)
	}
	if img.Bounds().Dx() != 16 || img.Bounds().Dy() != 16 {
		t.Fatalf("unexpected icon size")
	}
}

func TestSimpleModeHTML(t *testing.T) {
	data := TweetData{
		Text:   "Simple preview",
		Name:   "Example User",
		Handle: "example",
		Simple: true,
	}
	html, err := RenderHTML(data, DefaultOptions())
	if err != nil {
		t.Fatalf("RenderHTML: %v", err)
	}
	if strings.Contains(html, "Reply") || strings.Contains(html, "Copy link") {
		t.Fatalf("simple mode should hide footer actions")
	}
}

func TestTightWidthMode(t *testing.T) {
	opts := DefaultOptions()
	opts.WidthMode = "tight"
	opts.Width = 800
	data := TweetData{
		Text:   "short",
		Name:   "Example",
		Handle: "example",
	}
	img, err := RenderImage(data, opts)
	if err != nil {
		t.Fatalf("RenderImage: %v", err)
	}
	if img.Bounds().Dx() >= opts.Width {
		t.Fatalf("expected tight width to be smaller than %d", opts.Width)
	}
}

func TestTightWidthMin(t *testing.T) {
	opts := DefaultOptions()
	opts.WidthMode = "tight"
	opts.Width = 300
	data := TweetData{
		Text:   "short",
		Name:   "Example",
		Handle: "example",
	}
	img, err := RenderImage(data, opts)
	if err != nil {
		t.Fatalf("RenderImage: %v", err)
	}
	if img.Bounds().Dx() < 400 {
		t.Fatalf("expected minimum width of 400, got %d", img.Bounds().Dx())
	}
}
