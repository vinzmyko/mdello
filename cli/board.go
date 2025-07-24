package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vinzmyko/mdello/markdown"
	"github.com/vinzmyko/mdello/markdown/diff"
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
		if err != nil {
			fmt.Printf("Error creating trello client: %v\n", err)
			return
		}

		currentBoard, err := cfg.GetCurrentBoard(trelloClient)
		if err != nil {
			fmt.Printf("Error could not access current board: %v", err)
			return
		}

		safeName := strings.ReplaceAll(currentBoard.Name, " ", "~")
		safeName = strings.ReplaceAll(safeName, "/", "~")

		originalContent, boardSession, err := markdown.ToMarkdown(trelloClient, cfg, currentBoard)
		if err != nil {
			fmt.Printf("Converting to markdown failed: %v", err)
			return
		}

		// Parse original content for comparison
		originalReader := strings.NewReader(originalContent)
		originalBoard, err := markdown.FromMarkdown(originalReader, boardSession)
		if err != nil {
			fmt.Printf("Error parsing original markdown: %v\n", err)
			return
		}

		// First editor session
		editedContent, err := openEditorForContent(originalContent, fmt.Sprintf("mdello-%s", safeName))
		if err != nil {
			fmt.Printf("Error with editor: %v\n", err)
			return
		}

		if editedContent == originalContent {
			fmt.Println("No changes made.")
			return
		}

		// Parse edited content
		reader := bytes.NewReader([]byte(editedContent))
		editedBoard, err := markdown.FromMarkdown(reader, boardSession)
		if err != nil {
			fmt.Printf("Error parsing edited markdown: %v\n", err)
			return
		}

		diffResult, err := diff.QuickActionsDiff(originalBoard, editedBoard, cfg)
		if err != nil {
			fmt.Printf("Failed to analyse differences between original and edited content: %v", err)
			return
		}

		if len(diffResult.QuickActions) == 0 && len(diffResult.DetailedActions) == 0 {
			fmt.Println("No logical changes detected.")
			return
		}

		if len(diffResult.QuickActions) > 0 {
			fmt.Printf("Applying %d quick change(s)...\n", len(diffResult.QuickActions))
			applyActionsInOrder(diffResult.QuickActions, trelloClient, &markdown.ActionContext{BoardID: cfg.CurrentBoardID})
			fmt.Println("Quick changes applied!")
		}

		if len(diffResult.DetailedActions) > 0 {
			fmt.Printf("\nFound %d item(s) marked for detailed editing.\n", len(diffResult.DetailedActions))
			fmt.Println("Generating detailed editor...")

			detailedContent, err := markdown.GenerateDetailedMarkdown(diffResult.DetailedActions, trelloClient, cfg)
			if err != nil {
				fmt.Printf("\nFailed to generate detailed markdown file: %v", err)
			}

			//detailedEditedContent, err := openEditorForContent(detailedContent, fmt.Sprintf("mdello-%s-detailed", safeName))
			_, err = openEditorForContent(detailedContent, fmt.Sprintf("mdello-%s-detailed", safeName))
			if err != nil {
				fmt.Printf("\n Failed to open editor for detailed edit editor: %v", err)
			}

			fmt.Println("Detailed changes applied!")
		}
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

func openEditorForContent(content string, filePrefix string) (string, error) {
	tempFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.md", filePrefix))
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(content); err != nil {
		return "", fmt.Errorf("error writing to temp file: %w", err)
	}
	tempFile.Close()

	editor, err := getEditor()
	if err != nil {
		return "", fmt.Errorf("error getting editor: %w", err)
	}

	// Open editor
	editorCmd := exec.Command(editor, tempFile.Name())
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return "", fmt.Errorf("error opening editor: %w", err)
	}

	editedContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("error reading edited file: %w", err)
	}

	return string(editedContent), nil
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
