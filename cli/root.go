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

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(boardsCmd)
	rootCmd.AddCommand(boardCmd)
	rootCmd.AddCommand(openCmd)

	rootCmd.SetUsageTemplate(
		`Usage:
  {{.CommandPath}} [command]{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
