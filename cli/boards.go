package cli

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/vinzmyko/mdello/trello"
)

var boardsCmd = &cobra.Command{
	Use:   "boards",
	Short: "Get all current users boards",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg == nil || cfg.Token == "" {
			fmt.Println("No valid configuration found. Please run 'mdello init'.")
			return
		}

		trelloClient, err := trello.NewTrelloClient(apiKey, cfg.Token)
		if err != nil {
			fmt.Println(err)
			return
		}
		boards, err := trelloClient.GetBoards()
		if err != nil {
			fmt.Println(err)
			return
		}

		if cfg.CurrentBoard != nil {
			fmt.Printf("Current board: %s\n\n", cfg.CurrentBoard.Name)
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

		cfg.UpdateCurrentBoard(selectedBoard)

		err = cfg.Save()
		if err != nil {
			fmt.Printf("Error saving configuration: %v\n", err)
			return
		}
		fmt.Printf("Current board updated to: %s\n", selectedBoard.Name)
	},
}
