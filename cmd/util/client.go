package util

import (
	"context"
	"net/url"
	"os"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"github.com/spf13/cobra"
)

// Retrieves the Canvas client from the command context
func Client(cmd *cobra.Command) *canvas.Client {
	client := cmd.Context().Value("client").(*canvas.Client)
	return client
}

// Prepares the canvas client and sets it on the
// global command context at key "client".
// Easily retrieve with util.Client()
func SetClient(cmd *cobra.Command, args []string) {
	log := Logger(cmd)

	if _, set := os.LookupEnv("CANVAS_TOKEN"); !set {
		log.Fatal().Msg("Must set CANVAS_TOKEN in environment")
	}

	host, _ := cmd.Flags().GetString("canvasUrl")
	if host == "" {
		host = os.Getenv("CANVAS_URL")
	}

	canvasUrl, err := url.Parse(host)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to locate canvas url by flag or env")
	}

	client := canvas.MustNewClient(&canvas.ClientOpts{
		Ctx:   cmd.Context(),
		Url:   canvasUrl,
		Token: os.Getenv("CANVAS_TOKEN"),
	})

	cmd.SetContext(context.WithValue(cmd.Context(), "client", client))
}
