package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vinzmyko/mdello/trello"
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
		if _, err := loadConfig(); err == nil {
			fmt.Print("Configuration already exists. Overwrite? (y/N): ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "y" && response != "yes" {
				fmt.Println("'mdello init' cancelled.")
				return
			}
		}

		fmt.Println("Enter your Trello Token: ")
		reader := bufio.NewReader(os.Stdin)
		token, _ := reader.ReadString('\n')
		token = strings.TrimSpace(token)

		trelloClient, err := trello.NewTrelloClient(trelloAPIKey, token)
		if err != nil {
			fmt.Println("Token was invalid please try again")
			return
		}
		fmt.Println("Token saved successfully!")
		var currentBoard *string = nil
		boards, err := trelloClient.GetBoards()
		if err != nil {
			fmt.Printf("Could not access boards: %v", err)
			return
		}
		if len(boards) < 1 {
			fmt.Println("User has no boards")
			return
		} else {
			// ask user to select a default board
			fmt.Println("Select a board")
			for _, board := range boards {
				fmt.Println(board.Name)
			}
		}

		config := Config{
			Token:        token,
			CurrentBoard: currentBoard,
		}

		saveConfig(config)
	},
}

var getBoardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Get all current users boards",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("No config found. Please run 'mdello init'.")
			os.Exit(1)
			return
		}

		trelloClient, err := trello.NewTrelloClient(trelloAPIKey, config.Token)
		boards, err := trelloClient.GetBoards()
		if err != nil {
			fmt.Println(err)
		}

		for _, board := range boards {
			fmt.Printf("%s\n", board.Name)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(getBoardsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
