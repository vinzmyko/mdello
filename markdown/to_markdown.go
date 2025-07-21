package markdown

import (
	"fmt"
	"strings"
	"time"

	"github.com/vinzmyko/mdello/config"
	"github.com/vinzmyko/mdello/trello"
)

func ToMarkdown(trelloClient *trello.TrelloClient, configuration *config.Config, board *trello.Board) (string, *BoardSession, error) {
	session, err := NewBoardSession(board, trelloClient)
	if err != nil {
		return "", nil, err
	}

	var markdown strings.Builder

	markdown.WriteString(fmt.Sprintf("# %s {%s}\n", board.Name, session.GetShortID(board.ID)))
	for _, label := range board.Labels {
		if label.Name == "" {
			continue
		}
		markdownLabelName := strings.ReplaceAll(label.Name, " ", "~")
		markdown.WriteString(fmt.Sprintf("@%s:%s {%s}\n", markdownLabelName, label.Colour, session.GetShortID(label.ID)))
	}

	lists, _ := trelloClient.GetLists(board.ID)
	for _, list := range lists {
		markdown.WriteString(fmt.Sprintf("\n## %s {%s}", list.Name, session.GetShortID(list.ID)))
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
				markdownLabel := strings.ReplaceAll(label.Name, " ", "~")
				labels.WriteString(fmt.Sprintf(" @%s", markdownLabel))
			}
			var dueDateStr string
			if card.Due != nil && *card.Due != "" {
				dueDateStr = fmt.Sprintf(" due:%s", formatDate(*card.Due, configuration))
			}

			markdown.WriteString(fmt.Sprintf("\n- %s %s%s%s {%s}",
				checkbox, card.Name, labels.String(), dueDateStr, session.GetShortID(card.ID)))
		}
		markdown.WriteString("\n")
	}

	return markdown.String(), session, nil
}

func formatDate(due string, configuration *config.Config) string {
	parsedTime, err := time.Parse(time.RFC3339, due)
	if err != nil {
		fmt.Println("\nError parsing due date: skipping")
		return ""
	}

	return parsedTime.Local().Format(configuration.DateFormat)
}

func parseMarkdownDate(dateStr string, configuration *config.Config) (string, error) {
	parsedTime, err := time.ParseInLocation(configuration.DateFormat, dateStr, time.Local)
	if err != nil {
		return "", fmt.Errorf("failed to parse date %s with format %s: %w",
			dateStr, configuration.DateFormat, err)
	}

	return parsedTime.UTC().Format(time.RFC3339), nil
}
