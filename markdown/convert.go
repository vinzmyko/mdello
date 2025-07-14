// Trello API JSON objects → Markdown
package markdown

import (
	"fmt"
	"strings"
	"time"

	"github.com/vinzmyko/mdello/config"
	"github.com/vinzmyko/mdello/trello"
)

func ConvertToMarkdown(trelloClient *trello.TrelloClient, configuration *config.Config, board *trello.Board) (string, error) {
	var markdown strings.Builder
	markdown.WriteString(fmt.Sprintf("# %s\n", board.Name))

	lists, _ := trelloClient.GetLists(board.ID)

	for _, list := range lists {
		markdown.WriteString(fmt.Sprintf("\n## %s", list.Name))
		cards, _ := trelloClient.GetCards(list.ID)

		for _, card := range cards {
			var checkbox string
			if card.Badges.DueComplete {
				checkbox = "[x]"
			} else {
				checkbox = "[ ]"
			}
			var labels strings.Builder
			for _, label := range card.Labels {
				labels.WriteString(fmt.Sprintf(" @%s", label.Name))
			}
			var dueDateStr string
			if card.Due != nil && *card.Due != "" {
				dueDateStr = fmt.Sprintf(" due:%s", formatDate(card, configuration))
			}

			markdown.WriteString(fmt.Sprintf("\n- %s %s%s%s", checkbox, card.Name, labels.String(), (dueDateStr)))
		}
		markdown.WriteString("\n")
	}

	return markdown.String(), nil
}

func formatDate(card trello.Card, configuration *config.Config) string {
	parsedTime, err := time.Parse(time.RFC3339, *card.Due)
	if err != nil {
		fmt.Printf("\nError parsing %s's due date: skipping", card.Name)
		return ""
	}

	dateFormat := configuration.DateFormat
	switch {
	case strings.Contains(dateFormat, "MM-DD-YYYY"):
		return parsedTime.Format("01-02-2006 15:04")
	case strings.Contains(dateFormat, "DD-MM-YYYY"):
		return parsedTime.Format("02-01-2006 15:04")
	default:
		return parsedTime.Format("2006-01-02 15:04")
	}
}
