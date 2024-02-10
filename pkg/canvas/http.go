package canvas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"github.com/henvic/httpretty"
	"github.com/pterm/pterm"
)

var pretty = httpretty.Logger{
	Colors:         true,
	RequestHeader:  true,
	RequestBody:    true,
	ResponseHeader: true,
	ResponseBody:   true,
	Time:           true,
	Formatters:     []httpretty.Formatter{&httpretty.JSONFormatter{}},
}

type ClientRoundTripper struct {
	Ctx context.Context
}

func (c ClientRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	readOnly := r.Context().Value("readOnly")
	// Need to fix this broken awful auto-generated code
	// API seemingly does not support JSON body, and requires
	// form data. Irritatingly, json body does seem to work with curl
	if r.Method == http.MethodPut {
		body := make(map[string]string, 0)
		json.NewDecoder(r.Body).Decode(&body)

		formData := url.Values{}
		for k, val := range body {
			formData.Add(k, val)
		}

		newReq, err := http.NewRequestWithContext(c.Ctx,
			"PUT",
			r.URL.String(),
			bytes.NewBufferString(formData.Encode()),
		)
		if err != nil {
			return nil, err
		}

		newReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		newReq.Header.Set("Authorization", r.Header.Clone().Get("Authorization"))
		r = newReq
	}

	// Debug request
	if log.Trace().Enabled() {
		fmt.Println("\n" + pterm.Info.Sprintf("> %s %s REQUEST",
			strings.ToUpper(r.URL.Scheme),
			r.Method))
		pretty.PrintRequest(r)
		fmt.Println(pterm.Info.Sprint("> END HTTP REQUEST"))
	}

	// Perform the request
	if r.Method != http.MethodGet && readOnly.(bool) {
		log.Warn().Msg("Global readOnly enabled, printing request only!")
		pretty.PrintRequest(r)
		return nil, nil
	}

	resp, err := http.DefaultTransport.RoundTrip(r)

	// Debug response
	if log.Trace().Enabled() {
		fmt.Println("\n" + pterm.Info.Sprint("< HTTP RESPONSE"))
		pretty.PrintResponse(resp)
		fmt.Println(pterm.Info.Sprint("< END HTTP RESPONSE") + "\n")
	}

	if err == nil && resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, errors.New(fmt.Sprintf("Received non-200 response: %s",
			string(body)))
	}

	return resp, err
}

func (c *Client) reqMiddlwareFunc() canvasauto.RequestEditorFn {
	return func(req *http.Request, ctx context.Context) error {
		req.Header.Add("Authorization", "Bearer "+c.token)
		req.WithContext(c.ctx)
		return nil
	}
}
