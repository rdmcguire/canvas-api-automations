package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
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
	debug := flag.Bool("debug", false, "Debug Output")
	flag.Parse()

	if len(flag.Args()) < 1 {
		help()
	} else {
		cmd = flag.Args()[0]
	}

	lvl := new(slog.LevelVar)
	if *debug {
		lvl.Set(slog.LevelDebug)
	} else {
		lvl.Set(slog.LevelInfo)
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)

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
