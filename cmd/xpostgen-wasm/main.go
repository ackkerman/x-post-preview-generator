//go:build js && wasm

package main

import (
	"encoding/json"
	"strings"
	"syscall/js"

	"github.com/ackkerman/x-post-preview-generator/internal/render"
)

type renderPayload struct {
	Text      string `json:"text"`
	Name      string `json:"name"`
	Handle    string `json:"handle"`
	Verified  bool   `json:"verified"`
	AvatarURL string `json:"avatarUrl"`
	Date      string `json:"date"`
	LikeCount string `json:"likeCount"`
	CTA       string `json:"cta"`
	Width     string `json:"width"`
	Mode      string `json:"mode"`
}

func renderSVG(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return newRejectedPromise("missing payload")
	}

	payloadJSON := strings.TrimSpace(args[0].String())
	if payloadJSON == "" {
		return newRejectedPromise("empty payload")
	}

	promiseConstructor := js.Global().Get("Promise")
	handler := js.FuncOf(func(_ js.Value, promiseArgs []js.Value) any {
		resolve := promiseArgs[0]
		reject := promiseArgs[1]

		go func() {
			svg, err := renderFromJSON(payloadJSON)
			if err != nil {
				reject.Invoke(err.Error())
				return
			}
			resolve.Invoke(svg)
		}()

		return nil
	})

	promise := promiseConstructor.New(handler)
	handler.Release()
	return promise
}

func renderFromJSON(payloadJSON string) (string, error) {
	var payload renderPayload
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return "", err
	}

	text := strings.TrimSpace(payload.Text)
	name := strings.TrimSpace(payload.Name)
	handle := strings.TrimSpace(payload.Handle)

	if text == "" || name == "" || handle == "" {
		return "", errMissingFields()
	}

	data := render.TweetData{
		Text:      text,
		Icon:      strings.TrimSpace(payload.AvatarURL),
		Name:      name,
		Handle:    handle,
		Date:      strings.TrimSpace(payload.Date),
		CTA:       strings.TrimSpace(payload.CTA),
		Verified:  payload.Verified,
		Simple:    strings.EqualFold(payload.Mode, "simple"),
		LikeCount: strings.TrimSpace(payload.LikeCount),
	}

	opts := render.DefaultOptions()
	if strings.EqualFold(payload.Width, "tight") {
		opts.WidthMode = "tight"
	} else {
		opts.WidthMode = "fixed"
	}

	return render.RenderSVG(data, opts)
}

func newRejectedPromise(message string) js.Value {
	promiseConstructor := js.Global().Get("Promise")
	handler := js.FuncOf(func(_ js.Value, args []js.Value) any {
		reject := args[1]
		reject.Invoke(message)
		return nil
	})
	promise := promiseConstructor.New(handler)
	handler.Release()
	return promise
}

type missingFieldsError struct{}

func (missingFieldsError) Error() string {
	return "text, name, handle are required"
}

func errMissingFields() error {
	return missingFieldsError{}
}

func main() {
	js.Global().Set("xpostgenRender", js.FuncOf(renderSVG))
	select {}
}
