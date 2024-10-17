package main

import (
	"addressdb/cmd/cli/commands"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "bloom-cli"}

	// Add commands
	rootCmd.AddCommand(commands.EncodeCmd)
	rootCmd.AddCommand(commands.CheckCmd)
	rootCmd.AddCommand(commands.BatchCheckCmd)
	rootCmd.AddCommand(commands.AddressGenCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
