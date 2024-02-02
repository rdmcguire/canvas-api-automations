package main

import (
	"flag"
	"fmt"
	"os"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/cmd"
	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
)

var (
	oldCmd string
	client *canvas.Client
)

func main() {
	cmd.Execute()
}

func oldmain() {
	switch oldCmd {

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

// func init() {
// 	logLevel := flag.String("logLevel", "info", "Log Level (error|warn|info|debug|trace)")
// 	flag.Parse()

// 	if len(flag.Args()) < 1 {
// 		help()
// 	} else {
// 		oldCmd = flag.Args()[0]
// 	}

// 	switch strings.ToLower(*logLevel) {
// 	case "error":
// 		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
// 	case "warn":
// 		zerolog.SetGlobalLevel(zerolog.WarnLevel)
// 	case "info":
// 		zerolog.SetGlobalLevel(zerolog.InfoLevel)
// 	case "debug":
// 		zerolog.SetGlobalLevel(zerolog.DebugLevel)
// 	case "trace":
// 		zerolog.SetGlobalLevel(zerolog.TraceLevel)
// 	}
// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

// 	token := os.Getenv("CANVAS_TOKEN")
// 	canvasUrl, err := url.Parse(os.Getenv("CANVAS_URL"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	client = canvas.MustNewClient(&canvas.ClientOpts{
// 		Url:   canvasUrl,
// 		Token: token,
// 		Ctx:   context.Background(),
// 	})
// }
