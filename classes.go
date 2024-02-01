package main

import (
	"fmt"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
)

func classes() {
	fmt.Println("Courses:")
	for _, course := range client.ListCourses() {
		fmt.Println(canvas.CourseString(course))
	}
}
