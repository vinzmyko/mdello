package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/vinzmyko/mdello/config"
	"github.com/vinzmyko/mdello/trello"
	"golang.org/x/term"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise mdello with your Trello token",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg != nil && cfg.Token != "" {
			fmt.Print("Configuration already exists. Overwrite? (y/N): ")

			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println() // For consistent UX line breaks on Ctrl+c & wrong input
				printCancelled()
				return
			}
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "y" && response != "yes" {
				printCancelled()
				return
			}
		}

		fmt.Println("Enter your Trello Token: ")

		tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Println()
			printCancelled()
			return
		}
		fmt.Println()

		token := strings.TrimSpace(string(tokenBytes))

		trelloClient, err := trello.NewTrelloClient(apiKey, token)
		if err != nil {
			fmt.Println("\nThe provided token was invalid. Please try again.")
			fmt.Println()
			return
		}
		fmt.Println("Token registered successfully!")
		boards, err := trelloClient.GetBoards()
		if err != nil {
			fmt.Printf("\nCould not access boards: %v", err)
			fmt.Println()
			return
		}
		if len(boards) < 1 {
			fmt.Println("\nUser has no boards")
			fmt.Println()
			return
		}

		boardOptions := make([]string, 0, len(boards))
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
			printStepCancelled("Current board")
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
			fmt.Println("\nError: selected board not found")
			fmt.Println()
			return
		}
		fmt.Printf("\n%s selected.", selectedBoardName)

		var selectedDateFormatDisplay string
		boardPrompt = &survey.Select{
			Message: "Select a date format:",
			Options: config.GetDisplayOptions(),
			VimMode: true,
		}
		err = survey.AskOne(boardPrompt, &selectedDateFormatDisplay)
		if err != nil {
			printStepCancelled("Date format")
			return
		}
		actualDateFormat, found := config.GetFormatFromDisplay(selectedDateFormatDisplay)
		if !found {
			fmt.Println("\nError: Invalid date format selected")
			fmt.Println()
			return
		}
		fmt.Printf("\n%s selected.", selectedDateFormatDisplay)

		cfg := &config.Config{
			Token:          token,
			CurrentBoardID: selectedBoard.ID,
			DateFormat:     actualDateFormat,
		}
		fmt.Println("\nAll options have successfully been saved!")
		fmt.Println()

		err = config.SaveConfig(*cfg)
		if err != nil {
			fmt.Printf("\nError: Could not save configuration: %v\n", err)
			fmt.Println()
			return
		}
	},
}

func printCancelled() {
	fmt.Println()
	fmt.Println("'mdello init' cancelled.")
	fmt.Println()
}

func printStepCancelled(selection string) {
	fmt.Println()
	fmt.Printf("%s selection cancelled.\n", selection)
	fmt.Println("Previous options will not be saved.")
	fmt.Println()
}
