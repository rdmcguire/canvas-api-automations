package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
)

var (
	cmd    string
	client *canvas.Client
)

func main() {
	switch cmd {

	case "students":
		if len(flag.Args()) < 2 {
			help()
		}
		users := client.ListStudentsInCourse(flag.Args()[1])
		fmt.Println("First Name,Last Name,Email Address,Student ID")
		for _, user := range users {
			nameParts := strings.Split(*user.SortableName, ", ")
			if len(nameParts) < 2 {
				continue
			}
			fmt.Printf("%s,%s,%s,%d\n", nameParts[1], nameParts[0], *user.Email, user.Id)
		}

	case "classes":
		fmt.Println("Courses:")
		for _, course := range client.ListCourses() {
			fmt.Println(canvas.CourseString(course))
		}

	default:
		help()

	}
}

func help() {
	fmt.Printf("Usage: %s (classes | students [classId])\n", os.Args[0])
	os.Exit(1)
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
