package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	cmd    string
	client *canvas.Client
)

func main() {
	switch cmd {

	case "students":
		students()
	case "classes":
		classes()
	case "modules":
		modules()
	case "assignments":
		assignments()

	default:
		help()

	}
}

func help() {
	fmt.Printf("Usage: %s (classes | students [classId] | modules [classId] | assignments [file])\n", os.Args[0])
	os.Exit(1)
}

func requireArg() {
	if len(flag.Args()) < 2 {
		help()
	}
}

func init() {
	logLevel := flag.String("logLevel", "info", "Log Level (error|warn|info|debug|trace)")
	flag.Parse()

	if len(flag.Args()) < 1 {
		help()
	} else {
		cmd = flag.Args()[0]
	}

	switch strings.ToLower(*logLevel) {
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	token := os.Getenv("CANVAS_TOKEN")
	canvasUrl, err := url.Parse(os.Getenv("CANVAS_URL"))
	if err != nil {
		panic(err)
	}

	client = canvas.MustNewClient(&canvas.ClientOpts{
		Url:   canvasUrl,
		Token: token,
		Ctx:   context.Background(),
	})
}
