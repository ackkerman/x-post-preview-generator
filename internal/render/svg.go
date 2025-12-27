package render

import (
	"bytes"
	"encoding/xml"
	"math"
	"strings"
	"text/template"
)

type svgLine struct {
	X    float64
	Y    float64
	Text string
}

type svgAction struct {
	Icon   string
	Label  string
	LabelX float64
	LabelY float64
}

type svgView struct {
	Width         int
	Height        int
	DividerX1     float64
	DividerX2     float64
	Background    string
	Border        string
	Divider       string
	TextColor     string
	MutedColor    string
	AccentColor   string
	AvatarBg      string
	AvatarText    string
	FontFamily    string
	CornerRadius  float64
	StrokeWidth   float64
	AvatarX       float64
	AvatarY       float64
	AvatarSize    float64
	NameX         float64
	NameY         float64
	HandleX       float64
	HandleY       float64
	NameLine      string
	HandleLine    string
	VerifiedIcon  string
	TextLines     []svgLine
	DateX         float64
	DateY         float64
	DateLine      string
	DividerY      float64
	TwitterIcon   string
	InfoIcon      string
	Actions       []svgAction
	CTA           string
	CtaX          float64
	CtaY          float64
	CtaWidth      float64
	CtaHeight     float64
	CtaTextX      float64
	CtaTextY      float64
	AvatarDataURI string
	Initials      string
}

const svgTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="{{.Width}}" height="{{.Height}}" viewBox="0 0 {{.Width}} {{.Height}}" role="img" aria-label="X post preview">
  <rect x="1" y="1" width="{{addInt .Width -2}}" height="{{addInt .Height -2}}" rx="{{.CornerRadius}}" ry="{{.CornerRadius}}" fill="{{.Background}}" stroke="{{.Border}}" stroke-width="{{.StrokeWidth}}" />
  {{if .AvatarDataURI}}
  <defs>
    <clipPath id="avatar-clip">
      <circle cx="{{add .AvatarX (div .AvatarSize 2)}}" cy="{{add .AvatarY (div .AvatarSize 2)}}" r="{{div .AvatarSize 2}}" />
    </clipPath>
  </defs>
  <image href="{{.AvatarDataURI}}" x="{{.AvatarX}}" y="{{.AvatarY}}" width="{{.AvatarSize}}" height="{{.AvatarSize}}" clip-path="url(#avatar-clip)" preserveAspectRatio="xMidYMid slice" />
  {{else}}
  <circle cx="{{add .AvatarX (div .AvatarSize 2)}}" cy="{{add .AvatarY (div .AvatarSize 2)}}" r="{{div .AvatarSize 2}}" fill="{{.AvatarBg}}" />
  <text x="{{add .AvatarX (div .AvatarSize 2)}}" y="{{add .AvatarY (div .AvatarSize 2)}}" fill="{{.AvatarText}}" font-family="{{.FontFamily}}" font-size="28" font-weight="700" text-anchor="middle" dominant-baseline="central">{{escape .Initials}}</text>
  {{end}}

  <text x="{{.NameX}}" y="{{.NameY}}" fill="{{.TextColor}}" font-family="{{.FontFamily}}" font-size="28" font-weight="700">{{escape .NameLine}}</text>
  {{if .VerifiedIcon}}{{.VerifiedIcon}}{{end}}
  <text x="{{.HandleX}}" y="{{.HandleY}}" fill="{{.MutedColor}}" font-family="{{.FontFamily}}" font-size="22">{{escape .HandleLine}}</text>

  {{.TwitterIcon}}

  {{range .TextLines}}
  <text x="{{.X}}" y="{{.Y}}" fill="{{$.TextColor}}" font-family="{{$.FontFamily}}" font-size="28">{{escape .Text}}</text>
  {{end}}

  {{if .DateLine}}
  <text x="{{.DateX}}" y="{{.DateY}}" fill="{{.MutedColor}}" font-family="{{.FontFamily}}" font-size="22">{{escape .DateLine}}</text>
  {{.InfoIcon}}
  {{end}}

  <line x1="{{.DividerX1}}" y1="{{.DividerY}}" x2="{{.DividerX2}}" y2="{{.DividerY}}" stroke="{{.Divider}}" stroke-width="1" />

  {{range .Actions}}
  {{.Icon}}
  <text x="{{.LabelX}}" y="{{.LabelY}}" fill="{{$.MutedColor}}" font-family="{{$.FontFamily}}" font-size="20">{{escape .Label}}</text>
  {{end}}

  {{if .CTA}}
  <rect x="{{.CtaX}}" y="{{.CtaY}}" width="{{.CtaWidth}}" height="{{.CtaHeight}}" rx="{{div .CtaHeight 2}}" ry="{{div .CtaHeight 2}}" fill="{{.Background}}" stroke="{{.Divider}}" stroke-width="1" />
  <text x="{{.CtaTextX}}" y="{{.CtaTextY}}" fill="{{.AccentColor}}" font-family="{{.FontFamily}}" font-size="20" font-weight="600">{{escape .CTA}}</text>
  {{end}}
</svg>
`

// RenderSVG returns the tweet preview as SVG markup.
func RenderSVG(data TweetData, opts RenderOptions) (string, error) {
	opts = normalizeOptions(opts)
	fonts, err := loadFontSet(opts)
	if err != nil {
		return "", err
	}
	defer fonts.Close()

	layout := computeLayout(data, opts, fonts)
	avatar, err := avatarDataURI(data.Icon)
	if err != nil {
		return "", err
	}

	lines := make([]svgLine, len(layout.TextLines))
	for i, line := range layout.TextLines {
		lines[i] = svgLine{
			X:    layout.TextX,
			Y:    layout.TextY + float64(i)*layout.TextLineHeight,
			Text: line,
		}
	}

	actions := make([]svgAction, 0, len(layout.Actions))
	for _, action := range layout.Actions {
		icon, err := iconElement(action.IconName, action.IconX, action.IconY, action.IconSize, opts.Theme.Muted)
		if err != nil {
			return "", err
		}
		actions = append(actions, svgAction{
			Icon:   icon,
			Label:  action.Label,
			LabelX: action.LabelX,
			LabelY: action.LabelY,
		})
	}

	twitterIcon, err := iconElement("twitter", layout.TwitterX, layout.TwitterY, layout.TwitterSize, opts.Theme.Accent)
	if err != nil {
		return "", err
	}

	infoIcon := ""
	if layout.DateLine != "" {
		infoIcon, err = iconElement("info", layout.InfoX, layout.InfoY, layout.InfoSize, opts.Theme.Muted)
		if err != nil {
			return "", err
		}
	}
	verifiedIcon := ""
	if layout.Verified {
		verifiedIcon, err = iconElement("verified", layout.VerifiedX, layout.VerifiedY, layout.VerifiedSize, opts.Theme.Accent)
		if err != nil {
			return "", err
		}
	}

	corner := math.Min(20, float64(layout.Height)/12)
	view := svgView{
		Width:         layout.Width,
		Height:        layout.Height,
		Background:    opts.Theme.Background,
		Border:        opts.Theme.Border,
		Divider:       opts.Theme.Divider,
		TextColor:     opts.Theme.Text,
		MutedColor:    opts.Theme.Muted,
		AccentColor:   opts.Theme.Accent,
		AvatarBg:      opts.Theme.AvatarBg,
		AvatarText:    opts.Theme.AvatarText,
		FontFamily:    sanitizeFontFamily(opts.FontFamily),
		CornerRadius:  corner,
		StrokeWidth:   1.5,
		AvatarX:       layout.AvatarX,
		AvatarY:       layout.AvatarY,
		AvatarSize:    layout.AvatarSize,
		NameX:         layout.NameX,
		NameY:         layout.NameY,
		HandleX:       layout.HandleX,
		HandleY:       layout.HandleY,
		NameLine:      layout.NameLine,
		HandleLine:    layout.HandleLine,
		VerifiedIcon:  verifiedIcon,
		TextLines:     lines,
		DateX:         layout.DateX,
		DateY:         layout.DateY,
		DateLine:      layout.DateLine,
		DividerX1:     layout.Padding,
		DividerX2:     float64(layout.Width) - layout.Padding,
		DividerY:      layout.DividerY,
		TwitterIcon:   twitterIcon,
		InfoIcon:      infoIcon,
		Actions:       actions,
		CTA:           layout.CTA,
		CtaX:          layout.CtaX,
		CtaY:          layout.CtaY,
		CtaWidth:      layout.CtaWidth,
		CtaHeight:     layout.CtaHeight,
		CtaTextX:      layout.CtaTextX,
		CtaTextY:      layout.CtaTextY,
		AvatarDataURI: avatar,
		Initials:      initials(data.Name),
	}

	funcs := template.FuncMap{
		"escape": func(s string) string {
			var buf bytes.Buffer
			_ = xml.EscapeText(&buf, []byte(s))
			return buf.String()
		},
		"div":    func(a float64, b float64) float64 { return a / b },
		"add":    func(a float64, b float64) float64 { return a + b },
		"addInt": func(a int, b int) int { return a + b },
	}

	tmpl, err := template.New("svg").Funcs(funcs).Parse(svgTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, view); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func sanitizeFontFamily(value string) string {
	if value == "" {
		return "sans-serif"
	}
	return strings.ReplaceAll(value, "\"", "'")
}
