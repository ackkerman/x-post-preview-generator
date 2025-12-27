package render

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"strings"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

//go:embed assets/icons/*.svg
var iconFS embed.FS

func iconSVG(name string) (string, error) {
	data, err := iconFS.ReadFile(fmt.Sprintf("assets/icons/%s.svg", name))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func iconSVGWithColor(name string, color string) (string, error) {
	svg, err := iconSVG(name)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(svg, "currentColor", color), nil
}

func rasterizeIcon(name string, color string, size int) (image.Image, error) {
	svg, err := iconSVGWithColor(name, color)
	if err != nil {
		return nil, err
	}
	icon, err := oksvg.ReadIconStream(bytes.NewReader([]byte(svg)))
	if err != nil {
		return nil, err
	}
	w := float64(size)
	icon.SetTarget(0, 0, w, w)
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	scanner := rasterx.NewScannerGV(size, size, img, img.Bounds())
	raster := rasterx.NewDasher(size, size, scanner)
	icon.Draw(raster, 1.0)
	return img, nil
}

func iconElement(name string, x float64, y float64, size float64, color string) (string, error) {
	svg, err := iconSVGWithColor(name, color)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(svg, "<svg") {
		return "", fmt.Errorf("invalid svg icon: %s", name)
	}
	replacement := fmt.Sprintf("<svg x=\"%.2f\" y=\"%.2f\" width=\"%.2f\" height=\"%.2f\"", x, y, size, size)
	return strings.Replace(svg, "<svg", replacement, 1), nil
}
