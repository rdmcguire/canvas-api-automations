package main

import (
	"flag"
	"fmt"

	"gitea.libretechconsulting.com/50W/canvas-api-automations/pkg/canvas"
)

func modules() {
	requireArg()
	courseID := flag.Args()[1]
	fmt.Printf("Modules in course %s\n", courseID)
	for _, module := range client.ListModules(courseID) {
		fmt.Println(canvas.ModuleString(module))
	}
}
