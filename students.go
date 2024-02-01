package main

import (
	"flag"
	"fmt"
	"strings"
)

func students() {
	requireArg()
	users := client.ListStudentsInCourse(flag.Args()[1])
	fmt.Println("First Name,Last Name,Email Address,Student ID")
	for _, user := range users {
		nameParts := strings.Split(*user.SortableName, ", ")
		if len(nameParts) < 2 {
			continue
		}
		fmt.Printf("%s,%s,%s,%d\n", nameParts[1], nameParts[0], *user.Email, user.Id)
	}
}
