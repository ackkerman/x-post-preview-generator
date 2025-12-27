package render

import (
	"math"
	"strings"
)

type ActionLayout struct {
	IconName string
	Label    string
	IconX    float64
	IconY    float64
	IconSize float64
	LabelX   float64
	LabelY   float64
}

type Layout struct {
	Width          int
	Height         int
	Padding        float64
	AvatarSize     float64
	AvatarX        float64
	AvatarY        float64
	HeaderGap      float64
	NameX          float64
	NameY          float64
	Verified       bool
	VerifiedX      float64
	VerifiedY      float64
	VerifiedSize   float64
	HandleX        float64
	HandleY        float64
	TwitterX       float64
	TwitterY       float64
	TwitterSize    float64
	TextX          float64
	TextY          float64
	TextLines      []string
	TextLineHeight float64
	DateX          float64
	DateY          float64
	InfoX          float64
	InfoY          float64
	InfoSize       float64
	DividerY       float64
	Actions        []ActionLayout
	ShowFooter     bool
	CTA            string
	CtaX           float64
	CtaY           float64
	CtaWidth       float64
	CtaHeight      float64
	CtaTextX       float64
	CtaTextY       float64
	NameLine       string
	HandleLine     string
	DateLine       string
}

func buildHandleLine(data TweetData) string {
	return normalizeHandle(data.Handle)
}

func buildDateLine(data TweetData) string {
	parts := []string{}
	if data.Date != "" {
		parts = append(parts, data.Date)
	}
	if data.Location != "" {
		parts = append(parts, data.Location)
	}
	return strings.Join(parts, " · ")
}

func buildActions(data TweetData) []ActionLayout {
	likeLabel := strings.TrimSpace(data.LikeCount)
	if likeLabel == "" {
		likeLabel = "0"
	}
	return []ActionLayout{
		{IconName: "like", Label: likeLabel},
		{IconName: "reply", Label: "Reply"},
		{IconName: "link", Label: "Copy link"},
	}
}

func computeLayout(data TweetData, opts RenderOptions, fonts FontSet) Layout {
	opts = normalizeOptions(opts)
	padding := float64(opts.Padding)
	avatarSize := float64(opts.AvatarSize)
	gap := float64(opts.Gap)
	width := float64(opts.Width)

	twitterSize := 30.0
	infoSize := 20.0
	actionIconSize := 22.0
	ctaHeight := 44.0
	verifiedSize := 30.0
	verifiedGap := 6.0
	showFooter := !data.Simple
	minWidth := 600.0

	if strings.EqualFold(opts.WidthMode, "tight") {
		width = computeTightWidth(data, opts, fonts)
		if width < minWidth {
			width = minWidth
		}
		if opts.Width > 0 && float64(opts.Width) > minWidth && width > float64(opts.Width) {
			width = float64(opts.Width)
		}
	}

	textStartX := padding + avatarSize + gap
	headerAvailableWidth := width - padding - textStartX - twitterSize - 8
	if headerAvailableWidth < 1 {
		headerAvailableWidth = 1
	}

	nameAvailableWidth := headerAvailableWidth
	if data.Verified {
		nameAvailableWidth = math.Max(1, nameAvailableWidth-verifiedSize-verifiedGap)
	}
	nameLine := ellipsize(data.Name, nameAvailableWidth, fonts.Name)
	handleLine := ellipsize(buildHandleLine(data), headerAvailableWidth, fonts.Handle)

	textAvailableWidth := width - padding - textStartX
	if textAvailableWidth < 1 {
		textAvailableWidth = 1
	}

	textLines := wrapText(data.Text, textAvailableWidth, fonts.Text)

	nameHeight := lineHeight(fonts.Name)
	handleHeight := lineHeight(fonts.Handle)
	textHeight := lineHeight(fonts.Text)
	metaHeight := lineHeight(fonts.Meta)
	actionHeight := lineHeight(fonts.Action)
	ctaFontHeight := lineHeight(fonts.CTA)

	textLineHeight := textHeight * 1.45

	headerTextHeight := nameHeight + 4 + handleHeight
	headerHeight := math.Max(avatarSize, headerTextHeight)

	nameY := padding + nameHeight
	handleY := nameY + 4 + handleHeight
	textY := padding + headerHeight + 16 + textHeight
	verifiedX := 0.0
	verifiedY := 0.0
	if data.Verified {
		nameWidth := measureString(fonts.Name, nameLine)
		verifiedX = textStartX + nameWidth + verifiedGap
		verifiedY = nameY - nameHeight + (nameHeight-verifiedSize)/2
	}

	textBlockHeight := textHeight
	if len(textLines) > 1 {
		textBlockHeight = float64(len(textLines)-1)*textLineHeight + textHeight
	}

	if !showFooter {
		height := int(math.Ceil(padding + headerHeight + 16 + textBlockHeight + padding))
		return Layout{
			Width:          int(width),
			Height:         height,
			Padding:        padding,
			AvatarSize:     avatarSize,
			AvatarX:        padding,
			AvatarY:        padding,
			HeaderGap:      gap,
			NameX:          textStartX,
			NameY:          nameY,
			Verified:       data.Verified,
			VerifiedX:      verifiedX,
			VerifiedY:      verifiedY,
			VerifiedSize:   verifiedSize,
			HandleX:        textStartX,
			HandleY:        handleY,
			TwitterX:       width - padding - twitterSize,
			TwitterY:       padding,
			TwitterSize:    twitterSize,
			TextX:          textStartX,
			TextY:          textY,
			TextLines:      textLines,
			TextLineHeight: textLineHeight,
			ShowFooter:     false,
			NameLine:       nameLine,
			HandleLine:     handleLine,
		}
	}

	dateAvailableWidth := textAvailableWidth - infoSize - 8
	if dateAvailableWidth < 1 {
		dateAvailableWidth = textAvailableWidth
	}
	dateLine := ellipsize(buildDateLine(data), dateAvailableWidth, fonts.Meta)
	dateY := 0.0
	infoX := 0.0
	infoY := 0.0
	dividerY := 0.0
	cursorY := padding + headerHeight + 16 + textBlockHeight

	if dateLine != "" {
		dateRowHeight := math.Max(metaHeight, infoSize)
		dateY = cursorY + 16 + (dateRowHeight-metaHeight)/2 + metaHeight
		infoX = width - padding - infoSize
		infoY = dateY - metaHeight/2 - infoSize/2
		dividerY = dateY + 12
		cursorY = dividerY
	} else {
		dividerY = cursorY + 16
		cursorY = dividerY
	}

	actionsTop := cursorY + 16
	actionRowHeight := math.Max(actionHeight, actionIconSize)
	actionsBaseline := actionsTop + (actionRowHeight-actionHeight)/2 + actionHeight
	cursorY = actionsTop + actionRowHeight

	actions := buildActions(data)
	actionX := padding
	for i := range actions {
		labelWidth := measureString(fonts.Action, actions[i].Label)
		actions[i].IconSize = actionIconSize
		actions[i].IconX = actionX
		actions[i].IconY = actionsBaseline - actionHeight/2 - actionIconSize/2
		actions[i].LabelX = actionX + actionIconSize + 8
		actions[i].LabelY = actionsBaseline
		actionX += actionIconSize + 8 + labelWidth + 32
	}

	cta := strings.TrimSpace(data.CTA)
	ctaWidth := width - padding*2
	ctaX := padding
	ctaY := 0.0
	ctaTextX := 0.0
	ctaTextY := 0.0
	if cta != "" {
		ctaY = cursorY + 16
		ctaTextWidth := measureString(fonts.CTA, cta)
		ctaTextX = ctaX + (ctaWidth-ctaTextWidth)/2
		ctaTextY = ctaY + (ctaHeight+ctaFontHeight)/2 - 2
		cursorY = ctaY + ctaHeight
	}

	height := int(math.Ceil(cursorY + padding))

	return Layout{
		Width:          int(width),
		Height:         height,
		Padding:        padding,
		AvatarSize:     avatarSize,
		AvatarX:        padding,
		AvatarY:        padding,
		HeaderGap:      gap,
		NameX:          textStartX,
		NameY:          nameY,
		Verified:       data.Verified,
		VerifiedX:      verifiedX,
		VerifiedY:      verifiedY,
		VerifiedSize:   verifiedSize,
		HandleX:        textStartX,
		HandleY:        handleY,
		TwitterX:       width - padding - twitterSize,
		TwitterY:       padding,
		TwitterSize:    twitterSize,
		TextX:          textStartX,
		TextY:          textY,
		TextLines:      textLines,
		TextLineHeight: textLineHeight,
		DateX:          textStartX,
		DateY:          dateY,
		InfoX:          infoX,
		InfoY:          infoY,
		InfoSize:       infoSize,
		DividerY:       dividerY,
		Actions:        actions,
		ShowFooter:     true,
		CTA:            cta,
		CtaX:           ctaX,
		CtaY:           ctaY,
		CtaWidth:       ctaWidth,
		CtaHeight:      ctaHeight,
		CtaTextX:       ctaTextX,
		CtaTextY:       ctaTextY,
		NameLine:       nameLine,
		HandleLine:     handleLine,
		DateLine:       dateLine,
	}
}

func computeTightWidth(data TweetData, opts RenderOptions, fonts FontSet) float64 {
	padding := float64(opts.Padding)
	avatarSize := float64(opts.AvatarSize)
	gap := float64(opts.Gap)
	twitterSize := 30.0
	verifiedSize := 20.0
	verifiedGap := 6.0
	infoSize := 20.0
	actionIconSize := 22.0
	actionGap := 8.0
	actionSpacing := 32.0

	nameWidth := measureString(fonts.Name, strings.TrimSpace(data.Name))
	if data.Verified {
		nameWidth += verifiedSize + verifiedGap
	}
	handleWidth := measureString(fonts.Handle, normalizeHandle(data.Handle))
	headerTextWidth := math.Max(nameWidth, handleWidth)
	headerWidth := padding + avatarSize + gap + headerTextWidth + twitterSize + 8 + padding

	textWidth := 0.0
	for _, line := range strings.Split(data.Text, "\n") {
		textWidth = math.Max(textWidth, measureString(fonts.Text, line))
	}
	textBlockWidth := padding + avatarSize + gap + textWidth + padding

	maxWidth := math.Max(headerWidth, textBlockWidth)

	if !data.Simple {
		dateLine := buildDateLine(data)
		if dateLine != "" {
			dateWidth := measureString(fonts.Meta, dateLine)
			dateRowWidth := padding + avatarSize + gap + dateWidth + 8 + infoSize + padding
			maxWidth = math.Max(maxWidth, dateRowWidth)
		}

		actions := buildActions(data)
		actionsWidth := padding
		for i, action := range actions {
			labelWidth := measureString(fonts.Action, action.Label)
			actionsWidth += actionIconSize + actionGap + labelWidth
			if i < len(actions)-1 {
				actionsWidth += actionSpacing
			}
		}
		actionsWidth += padding
		maxWidth = math.Max(maxWidth, actionsWidth)

		cta := strings.TrimSpace(data.CTA)
		if cta != "" {
			ctaTextWidth := measureString(fonts.CTA, cta)
			ctaWidth := padding*2 + ctaTextWidth + 40
			maxWidth = math.Max(maxWidth, ctaWidth)
		}
	}

	return math.Ceil(maxWidth)
}
