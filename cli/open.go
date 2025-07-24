package cli

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open current board in Trello via default browser",
	Run: func(cmd *cobra.Command, args []string) {
		if cfg == nil || cfg.CurrentBoardID == "" {
			fmt.Println("No current board set. Run 'mdello init' first.")
			return
		}
		url := fmt.Sprintf("https://trello.com/b/%s", cfg.CurrentBoardID)
		err := browser.OpenURL(url)
		if err != nil {
			fmt.Printf("Error opening browser: %v\n", err)
			return
		}

		fmt.Printf("Opening board in browser: %s\n", url)
	},
}
