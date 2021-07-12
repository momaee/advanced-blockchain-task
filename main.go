package main

import (
	"backend_task/cmd"
	"fmt"

	"github.com/logrusorgru/aurora"
)

func main() {

	if err := cmd.Execute(); err != nil {
		fmt.Printf("cannot run the app, why? %v\n", aurora.Red(err))
	}
}
