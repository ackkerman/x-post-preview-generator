package render

import (
	"math"
	"strings"
	"unicode"

	budoux "github.com/sg0hsmt/budoux-go"
	"github.com/sg0hsmt/budoux-go/models"
	"golang.org/x/image/font"
	"golang.org/x/text/width"
)

var japaneseModel = models.DefaultJapaneseModel()

const budouxThreshold = -800

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
	if text == "" {
		return 0
	}
	fallback := fallbackAdvance(face)
	var total float64
	var prev rune
	hasPrev := false
	for _, r := range text {
		if hasPrev {
			total += float64(face.Kern(prev, r)) / 64.0
		}
		if isWideRune(r) {
			if adv, ok := face.GlyphAdvance(r); ok {
				total += math.Max(float64(adv)/64.0, fallback)
			} else {
				total += fallback
			}
		} else if adv, ok := face.GlyphAdvance(r); ok {
			total += float64(adv) / 64.0
		} else {
			total += fallback
		}
		prev = r
		hasPrev = true
	}
	return total
}

func lineHeight(face font.Face) float64 {
	m := face.Metrics()
	return float64(m.Height) / 64.0
}

func fontAscentDescent(face font.Face) (float64, float64) {
	m := face.Metrics()
	return float64(m.Ascent) / 64.0, float64(m.Descent) / 64.0
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
		segment = strings.TrimRightFunc(segment, unicode.IsSpace)
		if segment == "" {
			lines = append(lines, "")
			continue
		}
		if containsJapanese(segment) {
			lines = append(lines, wrapTokens(budoux.ParseWithThreshold(japaneseModel, segment, budouxThreshold), maxWidth, face)...)
			continue
		}
		lines = append(lines, wrapRunes(segment, maxWidth, face)...)
	}
	return lines
}

func containsJapanese(text string) bool {
	for _, r := range text {
		if unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) || unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func isWideRune(r rune) bool {
	switch width.LookupRune(r).Kind() {
	case width.EastAsianWide, width.EastAsianFullwidth, width.EastAsianAmbiguous:
		return true
	default:
		return false
	}
}

func fallbackAdvance(face font.Face) float64 {
	return float64(face.Metrics().Height) / 64.0
}

func budouxTokens(segment string) []string {
	if containsJapanese(segment) {
		return budoux.ParseWithThreshold(japaneseModel, segment, budouxThreshold)
	}
	return []string{segment}
}

func wrapTokens(tokens []string, maxWidth float64, face font.Face) []string {
	var lines []string
	var current strings.Builder

	appendLine := func(line string) {
		line = strings.TrimRightFunc(line, unicode.IsSpace)
		if line == "" {
			return
		}
		lines = append(lines, line)
	}

	for _, token := range tokens {
		if token == "" {
			continue
		}
		if current.Len() == 0 {
			token = strings.TrimLeftFunc(token, unicode.IsSpace)
			if token == "" {
				continue
			}
			if measureString(face, token) <= maxWidth {
				current.WriteString(token)
				continue
			}
			tokenLines := wrapRunes(token, maxWidth, face)
			if len(tokenLines) == 0 {
				continue
			}
			for i := 0; i < len(tokenLines)-1; i++ {
				appendLine(tokenLines[i])
			}
			current.Reset()
			current.WriteString(tokenLines[len(tokenLines)-1])
			continue
		}

		candidate := current.String() + token
		if measureString(face, candidate) <= maxWidth {
			current.WriteString(token)
			continue
		}

		appendLine(current.String())
		current.Reset()
		token = strings.TrimLeftFunc(token, unicode.IsSpace)
		if token == "" {
			continue
		}
		if measureString(face, token) <= maxWidth {
			current.WriteString(token)
			continue
		}
		tokenLines := wrapRunes(token, maxWidth, face)
		if len(tokenLines) == 0 {
			continue
		}
		for i := 0; i < len(tokenLines)-1; i++ {
			appendLine(tokenLines[i])
		}
		current.Reset()
		current.WriteString(tokenLines[len(tokenLines)-1])
	}

	if current.Len() > 0 {
		line := strings.TrimSpace(current.String())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		return []string{""}
	}
	return lines
}

func wrapRunes(segment string, maxWidth float64, face font.Face) []string {
	segmentRunes := []rune(strings.TrimRightFunc(segment, unicode.IsSpace))
	if len(segmentRunes) == 0 {
		return []string{""}
	}

	var lines []string
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
		} else if i > start {
			breakAt = i - 1
		}
		line := strings.TrimRightFunc(string(segmentRunes[start:breakAt+1]), unicode.IsSpace)
		if line == "" && breakAt < i {
			line = string(segmentRunes[start:i])
		}
		if line == "" {
			line = string(segmentRunes[start : i+1])
		}
		lines = append(lines, line)
		start = breakAt + 1
		i = start - 1
		lastBreak = -1
	}

	if start < len(segmentRunes) {
		lines = append(lines, strings.TrimSpace(string(segmentRunes[start:])))
	}
	return lines
}
