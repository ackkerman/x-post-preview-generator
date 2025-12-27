package render

import "image/color"

// TweetData holds the values to render.
type TweetData struct {
	Text     string
	Icon     string
	Name     string
	Handle   string
	Date     string
	Location string
	CTA      string
	Verified bool
}

// RenderOptions controls output sizes, fonts, and theme.
type RenderOptions struct {
	Width        int
	Padding      int
	AvatarSize   int
	Gap          int
	FontPath     string
	BoldFontPath string
	FontFamily   string
	Theme        Theme
}

// Theme defines color values for the card.
type Theme struct {
	Background string
	Border     string
	Divider    string
	Text       string
	Muted      string
	Accent     string
	AvatarBg   string
	AvatarText string
}

// DefaultOptions returns base rendering options.
func DefaultOptions() RenderOptions {
	return RenderOptions{
		Width:      960,
		Padding:    32,
		AvatarSize: 64,
		Gap:        16,
		FontFamily: "\"Helvetica Neue\", \"SF Pro Text\", \"SF Pro Display\", \"Segoe UI\", Roboto, \"Noto Sans JP\", Arial, sans-serif",
		Theme:      LightTheme(),
	}
}

// LightTheme is the default X-like light theme.
func LightTheme() Theme {
	return Theme{
		Background: "#FFFFFF",
		Border:     "#E1E8ED",
		Divider:    "#E6ECF0",
		Text:       "#0F1419",
		Muted:      "#536471",
		Accent:     "#1DA1F2",
		AvatarBg:   "#D1D9DE",
		AvatarText: "#FFFFFF",
	}
}

// DarkTheme is an alternative dark theme.
func DarkTheme() Theme {
	return Theme{
		Background: "#15202B",
		Border:     "#253341",
		Divider:    "#2F3B44",
		Text:       "#E7E9EA",
		Muted:      "#8B98A5",
		Accent:     "#1DA1F2",
		AvatarBg:   "#3E5463",
		AvatarText: "#FFFFFF",
	}
}

func normalizeOptions(opts RenderOptions) RenderOptions {
	def := DefaultOptions()
	if opts.Width <= 0 {
		opts.Width = def.Width
	}
	if opts.Padding <= 0 {
		opts.Padding = def.Padding
	}
	if opts.AvatarSize <= 0 {
		opts.AvatarSize = def.AvatarSize
	}
	if opts.Gap <= 0 {
		opts.Gap = def.Gap
	}
	if opts.FontFamily == "" {
		opts.FontFamily = def.FontFamily
	}
	if opts.Theme.Background == "" {
		opts.Theme = def.Theme
	}
	return opts
}

func colorFromHex(hex string) (color.NRGBA, error) {
	return parseHexColor(hex)
}
