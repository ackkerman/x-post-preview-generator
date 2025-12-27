package render

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

func parseHexColor(hex string) (color.NRGBA, error) {
	trimmed := strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(trimmed) != 6 && len(trimmed) != 8 {
		return color.NRGBA{}, fmt.Errorf("invalid hex color: %s", hex)
	}
	value, err := strconv.ParseUint(trimmed, 16, 32)
	if err != nil {
		return color.NRGBA{}, fmt.Errorf("invalid hex color: %s", hex)
	}
	if len(trimmed) == 6 {
		return color.NRGBA{
			R: uint8(value >> 16),
			G: uint8((value >> 8) & 0xFF),
			B: uint8(value & 0xFF),
			A: 0xFF,
		}, nil
	}
	return color.NRGBA{
		R: uint8(value >> 24),
		G: uint8((value >> 16) & 0xFF),
		B: uint8((value >> 8) & 0xFF),
		A: uint8(value & 0xFF),
	}, nil
}
