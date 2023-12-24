package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "armrest",
	Short: "My CLI tool",
	Long:  "A simple CLI tool written in Go",
	Run: func(cmd *cobra.Command, args []string) {
		// Add your logic here
		fmt.Println("Hello from mycli!")
	},
}

// Execute is the entry point for the CLI tool
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add any command-line flags or subcommands here
	// Example: rootCmd.PersistentFlags().StringVar(&flagVar, "flagName", "defaultValue", "Description")

	// Add any additional initialization code here
}