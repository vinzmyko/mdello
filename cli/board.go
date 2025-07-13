package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vinzmyko/mdello/config"
)

var boardCmd = &cobra.Command{
	Use:   "board",
	Short: "Edit current board via markdown file",
	Run: func(cmd *cobra.Command, args []string) {
		configurtaion, err := config.LoadConfig()
		if err != nil {
			fmt.Println("No config found. Please run 'mdello init'.")
			os.Exit(1)
			return
		}
		if configurtaion.CurrentBoard == nil {
			fmt.Println("No current board set. Please run 'mdello boards' to select a board.")
			return
		}

		currentBoard := configurtaion.CurrentBoard

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
