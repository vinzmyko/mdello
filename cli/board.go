package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vinzmyko/mdello/markdown"
	"github.com/vinzmyko/mdello/trello"
)

var boardCmd = &cobra.Command{
	Use:   "board",
	Short: "Edit current board via markdown file",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg == nil || cfg.Token == "" || cfg.CurrentBoardID == "" {
			fmt.Println("No valid cfg found. Please run 'mdello init'.")
			return
		}

		trelloClient, err := trello.NewTrelloClient(apiKey, cfg.Token)

		currentBoard, err := cfg.GetCurrentBoard(trelloClient)
		if err != nil {
			fmt.Printf("Error could not access current board: %v", err)
			return
		}

		safeName := strings.ReplaceAll(currentBoard.Name, " ", "~")
		safeName = strings.ReplaceAll(safeName, "/", "~")
		tempFileName := fmt.Sprintf("mdello-%s-*.md", safeName)
		tempFile, err := os.CreateTemp("", tempFileName)
		if err != nil {
			fmt.Printf("Error creating temp file: %v\n", err)
			return
		}
		defer os.Remove(tempFile.Name())

		originalContent, boardSession, err := markdown.ToMarkdown(trelloClient, cfg, currentBoard)
		if err != nil {
			fmt.Printf("Coverting to markdown failed: %v", err)
		}

		originalReader := strings.NewReader(originalContent)
		originalBoard, err := markdown.FromMarkdown(originalReader, boardSession)
		if err != nil {
			fmt.Printf("Error parsing original markdown: %v\n", err)
			return
		}

		tempFile.WriteString(originalContent)
		tempFile.Close()

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

		editedContent, err := os.ReadFile(tempFile.Name())
		if err != nil {
			fmt.Printf("Error reading edited file: %v\n", err)
			return
		}

		if string(editedContent) == originalContent {
			fmt.Println("No changes made.")
			return
		}

		// check edited content
		reader := bytes.NewReader(editedContent)

		editedBoard, err := markdown.FromMarkdown(reader, boardSession)
		if err != nil {
			fmt.Printf("Error parsing markdown: %v\n", err)
			return
		}

		actions, err := markdown.Diff(originalBoard, editedBoard, cfg)
		if err != nil {
			fmt.Printf("Failed to analyse differences between original and edited content: %v", err)
			return
		}

		if len(actions) == 0 {
			fmt.Println("No logical changes detected.")
			return
		}

		fmt.Printf("\nDetected %d change(s):\n", len(actions))

		fmt.Println("\nApplying changes...")
		applyActionsInOrder(actions, trelloClient, &markdown.ActionContext{BoardID: cfg.CurrentBoardID})
		fmt.Println("\nBoard updated successfully!")
	},
}

func applyActionsInOrder(actions []markdown.TrelloAction, client *trello.TrelloClient, ctx *markdown.ActionContext) error {
	// ===== BOARD AND LABEL CHANGES =====
	var remainingActions []markdown.TrelloAction

	for _, action := range actions {
		switch action.(type) {
		case markdown.CreateLabelAction, markdown.UpdateBoardNameAction, markdown.UpdateLabelName, markdown.UpdateLabelColour, markdown.DeleteLabelAction:
			fmt.Printf("Change: %s\n", action.Description())
			if err := action.Apply(client, ctx); err != nil {
				return err
			}
		default:
			remainingActions = append(remainingActions, action)
		}
	}

	// ===== LIST CHANGES - CREATE FIRST, THEN MOVE =====
	var listMoveActions []markdown.TrelloAction
	var cardActions []markdown.TrelloAction

	for _, action := range remainingActions {
		switch action.(type) {
		case markdown.CreateListAction, markdown.UpdateListNameAction, markdown.ArchiveListAction:
			// Create and modify lists first
			fmt.Printf("Change: %s\n", action.Description())
			if err := action.Apply(client, ctx); err != nil {
				return err
			}
		case markdown.UpdateListPositionAction:
			// Save position updates for later
			listMoveActions = append(listMoveActions, action)
		default:
			cardActions = append(cardActions, action)
		}
	}

	for _, action := range listMoveActions {
		fmt.Printf("Change: %s\n", action.Description())
		if err := action.Apply(client, ctx); err != nil {
			return err
		}
	}

	// ===== CARD CHANGES =====
	for _, action := range cardActions {
		fmt.Printf("Change: %s\n", action.Description())
		if err := action.Apply(client, ctx); err != nil {
			return err
		}
	}

	return nil
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
