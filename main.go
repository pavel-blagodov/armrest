package main

import (
	"fmt"
	"os"

	"github.com/pavel-blagodov/armrest/cmd/armrest/app"
)

func main() {
	err := app.NewRootCommand().Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
