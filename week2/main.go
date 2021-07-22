package main

import (
	"feedcoin.one/darwin/errservice"
	"fmt"
)

func main() {
	controller := errservice.TeacherController{}

	controller.Query(5)
	fmt.Printf("----------------------------------------------------\n")
	controller.Query(-5)
}
