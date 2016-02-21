package main

import (
	"fmt"
	"os"
)

// Version of the app
var Version = "dev-build"

func main() {
	if err := vasukiCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
