package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vinzmyko/mdello/config"
)

var rootCmd = &cobra.Command{
	Use:   "mdello",
	Short: "Markdown-driven Trello CLI for efficient board management",
	Long:  "mdello is a command-line tool that lets you manage Trello boards using markdown as the primary interface.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to mdello! Use --help to see available commands")
	},
}

var apiKey string
var cfg *config.Config

func Execute(trelloAPIKey string, configuration *config.Config) {
	apiKey = trelloAPIKey
	cfg = configuration

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(boardsCmd)
	rootCmd.AddCommand(boardCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
