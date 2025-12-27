package render

import (
	"fmt"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

const (
	nameFontSize     = 28
	handleFontSize   = 22
	textFontSize     = 28
	metaFontSize     = 22
	actionFontSize   = 20
	ctaFontSize      = 20
	initialsFontSize = 28
)

type FontSet struct {
	Name     font.Face
	Handle   font.Face
	Meta     font.Face
	Text     font.Face
	Action   font.Face
	CTA      font.Face
	Initials font.Face
}

func (f FontSet) Close() {
	closeFace(f.Name)
	closeFace(f.Handle)
	closeFace(f.Meta)
	closeFace(f.Text)
	closeFace(f.Action)
	closeFace(f.CTA)
	closeFace(f.Initials)
}

func closeFace(face font.Face) {
	if closer, ok := face.(interface{ Close() error }); ok {
		_ = closer.Close()
	}
}

func loadFontSet(opts RenderOptions) (FontSet, error) {
	regularFont, err := loadFont(opts.FontPath, goregular.TTF)
	if err != nil {
		return FontSet{}, err
	}
	boldFontPath := opts.BoldFontPath
	if boldFontPath == "" && opts.FontPath != "" {
		boldFontPath = opts.FontPath
	}
	boldFont, err := loadFont(boldFontPath, gobold.TTF)
	if err != nil {
		return FontSet{}, err
	}

	nameFace, err := newFace(boldFont, nameFontSize)
	if err != nil {
		return FontSet{}, err
	}
	handleFace, err := newFace(regularFont, handleFontSize)
	if err != nil {
		return FontSet{}, err
	}
	metaFace, err := newFace(regularFont, metaFontSize)
	if err != nil {
		return FontSet{}, err
	}
	textFace, err := newFace(regularFont, textFontSize)
	if err != nil {
		return FontSet{}, err
	}
	actionFace, err := newFace(regularFont, actionFontSize)
	if err != nil {
		return FontSet{}, err
	}
	ctaFace, err := newFace(boldFont, ctaFontSize)
	if err != nil {
		return FontSet{}, err
	}
	initialsFace, err := newFace(boldFont, initialsFontSize)
	if err != nil {
		return FontSet{}, err
	}

	return FontSet{
		Name:     nameFace,
		Handle:   handleFace,
		Meta:     metaFace,
		Text:     textFace,
		Action:   actionFace,
		CTA:      ctaFace,
		Initials: initialsFace,
	}, nil
}

func loadFont(path string, fallback []byte) (*opentype.Font, error) {
	if path == "" {
		return opentype.Parse(fallback)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read font: %w", err)
	}
	parsed, err := opentype.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}
	return parsed, nil
}

func newFace(otf *opentype.Font, size float64) (font.Face, error) {
	return opentype.NewFace(otf, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
