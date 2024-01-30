package canvas

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvasauto"
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
	return fmt.Sprintf("URL: %s\n", c.api.Server)
}

func (c *Client) reqMiddlwareFunc() canvasauto.RequestEditorFn {
	return func(req *http.Request, ctx context.Context) error {
		req.Header.Add("Authorization", "Bearer "+c.token)
		req.Header.Set("Accept", "application/json")
		req.WithContext(c.ctx)
		slog.Info("Sending HTTP Request", slog.Any("req", req.URL))
		return nil
	}
}

func MustNewClient(opts *ClientOpts) *Client {
	client := &Client{
		token: opts.Token,
		ctx:   opts.Ctx,
	}

	canvas, err := canvasauto.NewClient(
		opts.Url.String(),
		canvasauto.WithRequestEditorFn(client.reqMiddlwareFunc()),
	)
	if err != nil {
		panic(err)
	}

	client.api = canvas
	return client
}
