package render

import (
	"fmt"
	"image"
	"image/color"
	imagedraw "image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/fogleman/gg"
	xdraw "golang.org/x/image/draw"
)

// RenderImage renders the tweet preview into an RGBA image.
func RenderImage(data TweetData, opts RenderOptions) (*image.RGBA, error) {
	opts = normalizeOptions(opts)
	fonts, err := loadFontSet(opts)
	if err != nil {
		return nil, err
	}
	defer fonts.Close()

	layout := computeLayout(data, opts, fonts)
	ctx := gg.NewContext(layout.Width, layout.Height)

	bg, err := colorFromHex(opts.Theme.Background)
	if err != nil {
		return nil, err
	}
	border, err := colorFromHex(opts.Theme.Border)
	if err != nil {
		return nil, err
	}
	divider, err := colorFromHex(opts.Theme.Divider)
	if err != nil {
		return nil, err
	}
	text, err := colorFromHex(opts.Theme.Text)
	if err != nil {
		return nil, err
	}
	muted, err := colorFromHex(opts.Theme.Muted)
	if err != nil {
		return nil, err
	}
	accent, err := colorFromHex(opts.Theme.Accent)
	if err != nil {
		return nil, err
	}
	avatarBg, err := colorFromHex(opts.Theme.AvatarBg)
	if err != nil {
		return nil, err
	}
	avatarText, err := colorFromHex(opts.Theme.AvatarText)
	if err != nil {
		return nil, err
	}

	ctx.SetColor(bg)
	ctx.Clear()

	ctx.SetColor(border)
	ctx.SetLineWidth(2)
	corner := math.Min(20, float64(layout.Height)/12)
	ctx.DrawRoundedRectangle(1, 1, float64(layout.Width-2), float64(layout.Height-2), corner)
	ctx.Stroke()

	drawAvatar(ctx, data, layout, fonts, avatarBg, avatarText)

	ctx.SetFontFace(fonts.Name)
	ctx.SetColor(text)
	ctx.DrawString(layout.NameLine, layout.NameX, layout.NameY)

	if layout.Verified {
		icon, err := rasterizeIcon("verified", opts.Theme.Accent, int(layout.VerifiedSize))
		if err == nil {
			ctx.DrawImage(icon, int(layout.VerifiedX), int(layout.VerifiedY))
		}
	}

	ctx.SetFontFace(fonts.Handle)
	ctx.SetColor(muted)
	ctx.DrawString(layout.HandleLine, layout.HandleX, layout.HandleY)

	ctx.SetFontFace(fonts.Text)
	ctx.SetColor(text)
	y := layout.TextY
	for _, line := range layout.TextLines {
		ctx.DrawString(line, layout.TextX, y)
		y += layout.TextLineHeight
	}

	if layout.DateLine != "" {
		ctx.SetFontFace(fonts.Meta)
		ctx.SetColor(muted)
		ctx.DrawString(layout.DateLine, layout.DateX, layout.DateY)

		infoIcon, err := rasterizeIcon("info", opts.Theme.Muted, int(layout.InfoSize))
		if err == nil {
			ctx.DrawImage(infoIcon, int(layout.InfoX), int(layout.InfoY))
		}
	}

	ctx.SetColor(divider)
	ctx.SetLineWidth(1)
	ctx.DrawLine(layout.Padding, layout.DividerY, float64(layout.Width)-layout.Padding, layout.DividerY)
	ctx.Stroke()

	for _, action := range layout.Actions {
		icon, err := rasterizeIcon(action.IconName, opts.Theme.Muted, int(action.IconSize))
		if err == nil {
			ctx.DrawImage(icon, int(action.IconX), int(action.IconY))
		}
		ctx.SetFontFace(fonts.Action)
		ctx.SetColor(muted)
		ctx.DrawString(action.Label, action.LabelX, action.LabelY)
	}

	if layout.CTA != "" {
		ctx.SetColor(bg)
		ctx.DrawRoundedRectangle(layout.CtaX, layout.CtaY, layout.CtaWidth, layout.CtaHeight, layout.CtaHeight/2)
		ctx.FillPreserve()
		ctx.SetColor(divider)
		ctx.SetLineWidth(1)
		ctx.Stroke()

		ctx.SetFontFace(fonts.CTA)
		ctx.SetColor(accent)
		ctx.DrawString(layout.CTA, layout.CtaTextX, layout.CtaTextY)
	}

	twitterIcon, err := rasterizeIcon("twitter", opts.Theme.Accent, int(layout.TwitterSize))
	if err == nil {
		ctx.DrawImage(twitterIcon, int(layout.TwitterX), int(layout.TwitterY))
	}

	img := ctx.Image()
	rgba := image.NewRGBA(img.Bounds())
	imagedraw.Draw(rgba, img.Bounds(), img, image.Point{}, imagedraw.Src)
	return rgba, nil
}

// EncodeImage writes the image to the writer with the given format.
func EncodeImage(w io.Writer, img image.Image, format string) error {
	switch strings.ToLower(format) {
	case "png":
		return png.Encode(w, img)
	case "jpg", "jpeg":
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
	case "gif":
		return gif.Encode(w, img, &gif.Options{NumColors: 256})
	default:
		return fmt.Errorf("unsupported image format: %s", format)
	}
}

func drawAvatar(ctx *gg.Context, data TweetData, layout Layout, fonts FontSet, bg color.Color, fg color.Color) {
	if data.Icon != "" {
		img, err := loadImage(data.Icon)
		if err == nil {
			square := cropSquare(img)
			size := int(layout.AvatarSize)
			resized := image.NewRGBA(image.Rect(0, 0, size, size))
			xdraw.CatmullRom.Scale(resized, resized.Bounds(), square, square.Bounds(), xdraw.Over, nil)

			ctx.Push()
			ctx.DrawCircle(layout.AvatarX+layout.AvatarSize/2, layout.AvatarY+layout.AvatarSize/2, layout.AvatarSize/2)
			ctx.Clip()
			ctx.DrawImage(resized, int(layout.AvatarX), int(layout.AvatarY))
			ctx.Pop()
			return
		}
	}

	ctx.SetColor(bg)
	ctx.DrawCircle(layout.AvatarX+layout.AvatarSize/2, layout.AvatarY+layout.AvatarSize/2, layout.AvatarSize/2)
	ctx.Fill()

	ctx.SetFontFace(fonts.Initials)
	ctx.SetColor(fg)
	ctx.DrawStringAnchored(initials(data.Name), layout.AvatarX+layout.AvatarSize/2, layout.AvatarY+layout.AvatarSize/2, 0.5, 0.5)
}

func loadImage(pathOrURL string) (image.Image, error) {
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		resp, err := http.Get(pathOrURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, fmt.Errorf("failed to fetch icon: %s", resp.Status)
		}
		img, _, err := image.Decode(resp.Body)
		return img, err
	}

	file, err := os.Open(pathOrURL)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

func cropSquare(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	size := width
	if height < size {
		size = height
	}
	x0 := bounds.Min.X + (width-size)/2
	y0 := bounds.Min.Y + (height-size)/2

	out := image.NewRGBA(image.Rect(0, 0, size, size))
	imagedraw.Draw(out, out.Bounds(), img, image.Point{X: x0, Y: y0}, imagedraw.Src)
	return out
}
