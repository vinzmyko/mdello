package cli

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/vinzmyko/mdello/config"
	"github.com/vinzmyko/mdello/trello"
)

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Get all current users boards",
	Run: func(cmd *cobra.Command, args []string) {
		configuration, err := config.LoadConfig()
		if err != nil {
			fmt.Println("No config found. Please run 'mdello init'.")
			os.Exit(1)
			return
		}

		trelloClient, err := trello.NewTrelloClient(apiKey, configuration.Token)
		if err != nil {
			fmt.Println(err)
		}
		boards, err := trelloClient.GetBoards()
		if err != nil {
			fmt.Println(err)
		}

		if configuration.CurrentBoard != nil {
			fmt.Printf("Current board: %s\n\n", configuration.CurrentBoard.Name)
		} else {
			fmt.Println("No current board set")
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

		newConfig := config.Config{
			Token:        configuration.Token,
			CurrentBoard: selectedBoard,
		}

		config.SaveConfig(newConfig)
	},
}
