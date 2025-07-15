package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/vinzmyko/mdello/config"
	"github.com/vinzmyko/mdello/trello"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise mdello with your Trello token",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg != nil && cfg.Token != "" {
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

		trelloClient, err := trello.NewTrelloClient(apiKey, token)
		if err != nil {
			fmt.Println("Token was invalid please try again")
			return
		}
		fmt.Println("Token saved successfully!")
		boards, err := trelloClient.GetBoards()
		if err != nil {
			fmt.Printf("Could not access boards: %v", err)
			return
		}
		if len(boards) < 1 {
			fmt.Println("User has no boards")
			return
		}

		var boardOptions []string
		for _, board := range boards {
			boardOptions = append(boardOptions, board.Name)
		}

		var selectedBoardName string
		boardPrompt := &survey.Select{
			Message: "Select a board:",
			Options: boardOptions,
			VimMode: true,
		}
		err = survey.AskOne(boardPrompt, &selectedBoardName)
		if err != nil {
			fmt.Println("\nBoard selection cancelled.")
			return
		}

		var selectedBoard *trello.Board
		for _, board := range boards {
			if board.Name == selectedBoardName {
				selectedBoard = &board
				break
			}
		}
		if selectedBoard == nil {
			fmt.Println("Error: selected board not found")
			return
		}

		var selectedDateFormatDisplay string
		boardPrompt = &survey.Select{
			Message: "Select a date format:",
			Options: config.GetDisplayOptions(),
			VimMode: true,
		}
		err = survey.AskOne(boardPrompt, &selectedDateFormatDisplay)
		if err != nil {
			fmt.Println("\nDate format selection cancelled.")
		}
		actualDateFormat, found := config.GetFormatFromDisplay(selectedDateFormatDisplay)
		if !found {
			return
		}

		cfg := &config.Config{
			Token:        token,
			CurrentBoard: selectedBoard,
			DateFormat:   actualDateFormat,
		}

		config.SaveConfig(*cfg)
	},
}
