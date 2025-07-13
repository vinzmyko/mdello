package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/vinzmyko/mdello/trello"
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

		config := Config{
			Token:        token,
			CurrentBoard: selectedBoard,
		}

		saveConfig(config)
	},
}

var boardsCmd = &cobra.Command{
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
		if err != nil {
			fmt.Println(err)
		}
		boards, err := trelloClient.GetBoards()
		if err != nil {
			fmt.Println(err)
		}

		if config.CurrentBoard != nil {
			fmt.Printf("Current board: %s\n\n", config.CurrentBoard.Name)
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

		newConfig := Config{
			Token:        config.Token,
			CurrentBoard: selectedBoard,
		}

		saveConfig(newConfig)
	},
}

var boardCmd = &cobra.Command{
	Use:   "board",
	Short: "Edit current board via markdown file",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := loadConfig()
		if err != nil {
			fmt.Println("No config found. Please run 'mdello init'.")
			os.Exit(1)
			return
		}
		if config.CurrentBoard == nil {
			fmt.Println("No current board set. Please run 'mdello boards' to select a board.")
			return
		}

		currentBoard := config.CurrentBoard

		tempFile, err := os.CreateTemp("", "mdello-board-*.md")
		if err != nil {
			fmt.Printf("Error creating temp file: %v\n", err)
			return
		}
		defer os.Remove(tempFile.Name())

		// Write some placeholder content
		originalContent := fmt.Sprintf("# %s\n\n## To Do\n- [ ] Example task\n", currentBoard.Name)
		tempFile.WriteString(originalContent)
		tempFile.Close()

		// Get file info before editing
		beforeStat, err := os.Stat(tempFile.Name())
		if err != nil {
			fmt.Printf("Error getting file stats: %v\n", err)
			return
		}

		editor, err := getEditor()
		if err != nil {
			fmt.Printf("Error getting editor: %v\n", err)
			return
		}

		editorCmd := exec.Command(editor, tempFile.Name())
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		err = editorCmd.Run()
		if err != nil {
			fmt.Printf("Error opening editor: %v\n", err)
			return
		}

		afterStat, err := os.Stat(tempFile.Name())
		if err != nil {
			fmt.Printf("Error getting file stats after edit: %v\n", err)
			return
		}

		editedContent, err := os.ReadFile(tempFile.Name())
		if err != nil {
			fmt.Printf("Error reading edited file: %v\n", err)
			return
		}

		if string(editedContent) == originalContent {
			fmt.Println("No changes made. Aborting.")
			return
		}

		if !afterStat.ModTime().After(beforeStat.ModTime()) {
			fmt.Println("File was not saved. Aborting.")
			return
		}

		fmt.Println("Board edited successfully!")
		fmt.Printf("Changes detected: %d bytes\n", len(editedContent))
	},
}

func Execute() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(boardsCmd)
	rootCmd.AddCommand(boardCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getEditor() (string, error) {
	cmd := exec.Command("git", "var", "GIT_EDITOR")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Git editor: %w", err)
	}

	editor := strings.TrimSpace(string(output))
	if editor == "" {
		return "", fmt.Errorf("no editor configured")
	}

	return editor, nil
}
