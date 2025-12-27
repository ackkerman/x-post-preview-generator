package render

import (
	"bytes"
	"html/template"
	"strings"
)

type htmlView struct {
	Width         int
	Padding       int
	AvatarSize    int
	Gap           int
	Name          string
	Handle        string
	DateLine      string
	Text          template.HTML
	CTA           string
	Verified      bool
	ShowFooter    bool
	AvatarDataURI string
	Initials      string
	FontFamily    string
	Background    string
	Border        string
	Divider       string
	TextColor     string
	MutedColor    string
	AccentColor   string
	AvatarBg      string
	AvatarText    string
	TwitterIcon   template.HTML
	VerifiedIcon  template.HTML
	InfoIcon      template.HTML
	Actions       []htmlAction
}

const htmlTemplate = `<!doctype html>
<html lang="ja">
<head>
  <meta charset="utf-8" />
  <title>X Post Preview</title>
  <style>
    :root {
      --bg: {{.Background}};
      --border: {{.Border}};
      --divider: {{.Divider}};
      --text: {{.TextColor}};
      --muted: {{.MutedColor}};
      --accent: {{.AccentColor}};
      --avatar-bg: {{.AvatarBg}};
      --avatar-text: {{.AvatarText}};
    }
    body {
      margin: 0;
      padding: 0;
      background: var(--bg);
      font-family: {{.FontFamily}};
      color: var(--text);
    }
    .card {
      width: {{.Width}}px;
      box-sizing: border-box;
      padding: {{.Padding}}px;
      border: 1.5px solid var(--border);
      border-radius: 20px;
      background: var(--bg);
    }
    .header {
      display: flex;
      align-items: flex-start;
      justify-content: space-between;
    }
    .header-left {
      display: flex;
      gap: {{.Gap}}px;
      align-items: flex-start;
    }
    .avatar {
      width: {{.AvatarSize}}px;
      height: {{.AvatarSize}}px;
      border-radius: 999px;
      background: var(--avatar-bg);
      color: var(--avatar-text);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 28px;
      font-weight: 700;
      overflow: hidden;
      flex: none;
    }
    .avatar img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      display: block;
    }
    .name-row {
      display: flex;
      align-items: center;
      gap: 6px;
    }
    .name {
      font-size: 28px;
      font-weight: 700;
      line-height: 1.2;
    }
    .verified {
      color: var(--accent);
    }
    .handle {
      font-size: 22px;
      color: var(--muted);
      margin-top: 4px;
    }
    .twitter {
      color: var(--accent);
    }
    .text {
      margin-top: 16px;
      font-size: 28px;
      line-height: 1.45;
      white-space: pre-wrap;
      word-break: keep-all;
      overflow-wrap: break-word;
    }
    .date-row {
      margin-top: 16px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      color: var(--muted);
      font-size: 22px;
    }
    .divider {
      margin-top: 12px;
      border-top: 1px solid var(--divider);
    }
    .actions {
      margin-top: 14px;
      display: flex;
      gap: 32px;
      align-items: center;
      color: var(--muted);
      font-size: 20px;
    }
    .action {
      display: inline-flex;
      align-items: center;
      gap: 8px;
    }
    .icon svg {
      width: 22px;
      height: 22px;
      display: block;
    }
    .twitter svg {
      width: 30px;
      height: 30px;
      display: block;
    }
    .verified svg {
      width: 20px;
      height: 20px;
      display: block;
    }
    .info svg {
      width: 20px;
      height: 20px;
      display: block;
    }
    .cta {
      margin-top: 16px;
      height: 44px;
      border-radius: 999px;
      border: 1px solid var(--divider);
      display: flex;
      align-items: center;
      justify-content: center;
      color: var(--accent);
      font-weight: 600;
      font-size: 20px;
    }
  </style>
</head>
<body>
  <div class="card">
    <div class="header">
      <div class="header-left">
        <div class="avatar">
          {{if .AvatarDataURI}}
            <img src="{{.AvatarDataURI}}" alt="avatar" />
          {{else}}
            {{.Initials}}
          {{end}}
        </div>
        <div class="header-text">
          <div class="name-row">
            <div class="name">{{.Name}}</div>
            {{if .Verified}}<div class="verified icon">{{.VerifiedIcon}}</div>{{end}}
          </div>
          <div class="handle">{{.Handle}}</div>
        </div>
      </div>
      <div class="twitter icon">{{.TwitterIcon}}</div>
    </div>
    <div class="text">{{.Text}}</div>
    {{if .ShowFooter}}
      {{if .DateLine}}
      <div class="date-row">
        <div>{{.DateLine}}</div>
        <div class="info">{{.InfoIcon}}</div>
      </div>
      <div class="divider"></div>
      {{else}}
      <div class="divider" style="margin-top: 16px;"></div>
      {{end}}
      <div class="actions">
        {{range .Actions}}
        <div class="action icon">{{.Icon}}<span>{{.Label}}</span></div>
        {{end}}
      </div>
      {{if .CTA}}
      <div class="cta">{{.CTA}}</div>
      {{end}}
    {{end}}
  </div>
</body>
</html>
`

// RenderHTML returns the tweet preview as HTML.
func RenderHTML(data TweetData, opts RenderOptions) (string, error) {
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
	icons, err := loadHTMLIcons()
	if err != nil {
		return "", err
	}

	view := htmlView{
		Width:         layout.Width,
		Padding:       opts.Padding,
		AvatarSize:    opts.AvatarSize,
		Gap:           opts.Gap,
		Name:          data.Name,
		Handle:        buildHandleLine(data),
		DateLine:      buildDateLine(data),
		Text:          formatHTMLText(data.Text),
		CTA:           strings.TrimSpace(data.CTA),
		Verified:      data.Verified,
		ShowFooter:    !data.Simple,
		AvatarDataURI: avatar,
		Initials:      initials(data.Name),
		FontFamily:    opts.FontFamily,
		Background:    opts.Theme.Background,
		Border:        opts.Theme.Border,
		Divider:       opts.Theme.Divider,
		TextColor:     opts.Theme.Text,
		MutedColor:    opts.Theme.Muted,
		AccentColor:   opts.Theme.Accent,
		AvatarBg:      opts.Theme.AvatarBg,
		AvatarText:    opts.Theme.AvatarText,
		TwitterIcon:   icons.Twitter,
		VerifiedIcon:  icons.Verified,
		InfoIcon:      icons.Info,
	}
	view.Actions = buildHTMLActions(data, icons)

	tmpl, err := template.New("tweet").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, view); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func formatHTMLText(text string) template.HTML {
	if strings.TrimSpace(text) == "" {
		return template.HTML(template.HTMLEscapeString(text))
	}
	segments := strings.Split(text, "\n")
	var builder strings.Builder
	for i, segment := range segments {
		if i > 0 {
			builder.WriteByte('\n')
		}
		if segment == "" {
			continue
		}
		tokens := budouxTokens(segment)
		for j, token := range tokens {
			builder.WriteString(template.HTMLEscapeString(token))
			if j < len(tokens)-1 {
				builder.WriteString("<wbr>")
			}
		}
	}
	return template.HTML(builder.String())
}

type htmlIcons struct {
	Twitter  template.HTML
	Verified template.HTML
	Reply    template.HTML
	Like     template.HTML
	Link     template.HTML
	Info     template.HTML
}

type htmlAction struct {
	Icon  template.HTML
	Label string
}

func loadHTMLIcons() (htmlIcons, error) {
	icons := htmlIcons{}
	var err error
	if icons.Twitter, err = iconHTML("twitter"); err != nil {
		return icons, err
	}
	if icons.Verified, err = iconHTML("verified"); err != nil {
		return icons, err
	}
	if icons.Reply, err = iconHTML("reply"); err != nil {
		return icons, err
	}
	if icons.Like, err = iconHTML("like"); err != nil {
		return icons, err
	}
	if icons.Link, err = iconHTML("link"); err != nil {
		return icons, err
	}
	if icons.Info, err = iconHTML("info"); err != nil {
		return icons, err
	}
	return icons, nil
}

func iconHTML(name string) (template.HTML, error) {
	icon, err := iconSVG(name)
	if err != nil {
		return "", err
	}
	return template.HTML(icon), nil
}

func buildHTMLActions(data TweetData, icons htmlIcons) []htmlAction {
	actions := buildActions(data)
	out := make([]htmlAction, 0, len(actions))
	for _, action := range actions {
		var icon template.HTML
		switch action.IconName {
		case "reply":
			icon = icons.Reply
		case "like":
			icon = icons.Like
		case "link":
			icon = icons.Link
		default:
			continue
		}
		out = append(out, htmlAction{
			Icon:  icon,
			Label: action.Label,
		})
	}
	return out
}
