package render

import (
	"fmt"
	"io"
	"strings"
)

// RenderToWriter dispatches rendering based on the format.
func RenderToWriter(w io.Writer, data TweetData, opts RenderOptions, format string) error {
	format = normalizeFormat(format)
	switch format {
	case "png", "jpg", "jpeg", "gif":
		img, err := RenderImage(data, opts)
		if err != nil {
			return err
		}
		return EncodeImage(w, img, format)
	case "svg":
		svg, err := RenderSVG(data, opts)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, svg)
		return err
	case "html":
		html, err := RenderHTML(data, opts)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, html)
		return err
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func normalizeFormat(format string) string {
	lower := strings.ToLower(strings.TrimSpace(format))
	if lower == "" {
		return ""
	}
	if lower == "jpg" {
		return "jpeg"
	}
	return lower
}
