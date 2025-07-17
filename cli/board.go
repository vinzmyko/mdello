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
			fmt.Println("No changes made. Aborting.")
			return
		}

		// check edited content
		reader := bytes.NewReader(editedContent)

		editedBoard, err := markdown.FromMarkdown(reader, boardSession)
		if err != nil {
			fmt.Printf("Error parsing markdown: %v\n", err)
			return
		}

		actions := markdown.Diff(originalBoard, editedBoard)

		if len(actions) == 0 {
			fmt.Println("No logical changes detected.")
			return
		}

		fmt.Printf("\nDetected %d change(s):\n", len(actions))

		// TODO: Should I add a confirmation step to confirm the changes? Will decide later

		fmt.Println("\n Applying changes...")
		for _, act := range actions {
			err := act.Apply(trelloClient)
			fmt.Printf("Change: %s\n", act.Description())
			if err != nil {
				// TODO: Decide how to handle partial failures. Stop or continue?
				fmt.Printf("\nError applying change '%s': %v\n", act.Description(), err)
				return
			}
		}

		fmt.Println("\nBoard updated successfully!")
	},
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
