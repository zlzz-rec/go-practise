package main

import (
	"fmt"
	"week2/errservice"
)

func main() {
	controller := errservice.TeacherController{}

	controller.Query(5)
	fmt.Printf("----------------------------------------------------\n")
	controller.Query(-5)
}
