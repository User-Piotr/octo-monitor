package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "octo-cli",
	Short: "A CLI for interacting with the Octo API",
	Long:  `A Command Line Interface for interacting with the Octo API`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the Octo CLI!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
