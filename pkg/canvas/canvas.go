package canvas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
	"github.com/tomnomnom/linkheader"
)

type Client struct {
	token string
	ctx   context.Context
	api   *canvasauto.Client
}

type ClientOpts struct {
	Ctx   context.Context
	Url   *url.URL
	Token string
}

func (c *Client) String() string {
	return fmt.Sprintf("URL: %s", c.api.Server)
}

func MustNewClient(opts *ClientOpts) *Client {
	client := &Client{
		token: opts.Token,
		ctx:   opts.Ctx,
	}

	// Create client with request and response middleware
	canvas, err := canvasauto.NewClient(
		opts.Url.String(),
		canvasauto.WithRequestEditorFn(client.reqMiddlwareFunc()),
		canvasauto.WithHTTPClient(&http.Client{Transport: ClientRoundTripper{Ctx: opts.Ctx}}),
	)
	if err != nil {
		panic(err)
	}

	client.api = canvas
	return client
}

func isLastPage(r *http.Response) bool {
	links := linkheader.Parse(r.Header.Get("link"))
	for _, link := range links {
		if link.Rel == "next" {
			return false
		}
	}
	return true
}
