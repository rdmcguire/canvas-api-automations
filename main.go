package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
)

func main() {
	token := os.Getenv("CANVAS_TOKEN")
	canvasUrl, err := url.Parse(os.Getenv("CANVAS_URL"))
	if err != nil {
		panic(err)
	}

	client := canvas.MustNewClient(&canvas.ClientOpts{
		Url:   canvasUrl,
		Token: token,
		Ctx:   context.Background(),
	})

	fmt.Println("Courses:")
	for _, course := range client.ListCourses() {
		fmt.Printf("Course %d: %s, [%v]\n", *course.Id, *course.Name, *course.WorkflowState)
	}

	users := client.ListStudentsInCourse(os.Getenv("CANVAS_COURSE_ID"))
	fmt.Println("First Name,Last Name,Email Address,Student ID")
	for _, user := range users {
		nameParts := strings.Split(*user.SortableName, ", ")
		if len(nameParts) < 2 {
			continue
		}
		fmt.Printf("%s,%s,%s,%d\n", nameParts[1], nameParts[0], *user.Email, user.Id)
	}
}
