package render

import (
	"strings"
	"unicode"

	"golang.org/x/image/font"
)

func normalizeHandle(handle string) string {
	trimmed := strings.TrimSpace(handle)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "@") {
		return trimmed
	}
	return "@" + trimmed
}

func initials(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "?"
	}
	parts := strings.Fields(trimmed)
	if len(parts) == 1 {
		runes := []rune(parts[0])
		if len(runes) <= 2 {
			return string(runes)
		}
		return string(runes[:2])
	}
	first := []rune(parts[0])
	second := []rune(parts[1])
	if len(first) == 0 {
		return "?"
	}
	if len(second) == 0 {
		return string(first[0])
	}
	return string([]rune{first[0], second[0]})
}

func measureString(face font.Face, text string) float64 {
	d := &font.Drawer{Face: face}
	return float64(d.MeasureString(text)) / 64.0
}

func lineHeight(face font.Face) float64 {
	m := face.Metrics()
	return float64(m.Height) / 64.0
}

func ellipsize(text string, maxWidth float64, face font.Face) string {
	if text == "" || measureString(face, text) <= maxWidth {
		return text
	}
	ellipsis := "..."
	ellipsisWidth := measureString(face, ellipsis)
	if ellipsisWidth > maxWidth {
		return ""
	}
	var builder strings.Builder
	for _, r := range []rune(text) {
		candidate := builder.String() + string(r)
		if measureString(face, candidate)+ellipsisWidth > maxWidth {
			break
		}
		builder.WriteRune(r)
	}
	result := strings.TrimRightFunc(builder.String(), unicode.IsSpace)
	if result == "" {
		return ""
	}
	return result + ellipsis
}

func wrapText(text string, maxWidth float64, face font.Face) []string {
	if strings.TrimSpace(text) == "" {
		return []string{""}
	}

	var lines []string
	for _, segment := range strings.Split(text, "\n") {
		segmentRunes := []rune(strings.TrimRightFunc(segment, unicode.IsSpace))
		if len(segmentRunes) == 0 {
			lines = append(lines, "")
			continue
		}

		start := 0
		lastBreak := -1
		for i := 0; i < len(segmentRunes); i++ {
			r := segmentRunes[i]
			if unicode.IsSpace(r) {
				lastBreak = i
			}
			current := string(segmentRunes[start : i+1])
			if measureString(face, current) <= maxWidth {
				continue
			}

			breakAt := i
			if lastBreak >= start {
				breakAt = lastBreak
			}
			line := strings.TrimRightFunc(string(segmentRunes[start:breakAt+1]), unicode.IsSpace)
			if line == "" && breakAt < i {
				line = string(segmentRunes[start:i])
			}
			lines = append(lines, line)
			start = breakAt + 1
			i = start - 1
			lastBreak = -1
		}

		if start < len(segmentRunes) {
			lines = append(lines, strings.TrimSpace(string(segmentRunes[start:])))
		}
	}
	return lines
}
