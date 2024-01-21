package main

import (
	"fmt"
	"os"

	app "github.com/pavel-blagodov/armrest/cmd/armrest"
)

func main() {
	err := app.NewRootCommand().Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
