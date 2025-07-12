package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "mdello",
	Short: "Manage Trello boards through the terminal",
	Long:  "A CLI tool to manage Trello boards using markdown files as the source of truth",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to mdello! Use --help to see available commands")
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise mdello with your Trello token",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter your Trello Token: ")
		reader := bufio.NewReader(os.Stdin)
		token, _ := reader.ReadString('\n')
		token = strings.TrimSpace(token)

		saveConfig(token)
		fmt.Println("Token saved successfully!")
	},
}

func Execute() {
	rootCmd.AddCommand(initCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validateConfig() (*Config, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("no config found, please run 'mdello init' first")
	}
	return config, nil
}
