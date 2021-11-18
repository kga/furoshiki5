package main

import (
	"fmt"
	"os"
)

const version = "0.0.0"

func main() {
	args := os.Args

	if len(args) > 1 {
		switch args[1] {
		case "version":
			commandVersion()
		}
	}
}

func commandVersion() {
	fmt.Printf("furoshiki5 version %s", version)
}
