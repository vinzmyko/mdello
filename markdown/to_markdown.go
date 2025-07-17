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
				labels.WriteString(fmt.Sprintf(" @%s", label.Name))
			}
			var dueDateStr string
			if card.Due != nil && *card.Due != "" {
				dueDateStr = fmt.Sprintf(" due:%s", formatDate(card, configuration))
			}

			markdown.WriteString(fmt.Sprintf("\n- %s %s%s%s {%s}",
				checkbox, card.Name, labels.String(), dueDateStr, session.GetShortID(card.ID)))
		}
		markdown.WriteString("\n")
	}

	return markdown.String(), session, nil
}

func formatDate(card trello.Card, configuration *config.Config) string {
	parsedTime, err := time.Parse(time.RFC3339, *card.Due)
	if err != nil {
		fmt.Printf("\nError parsing %s's due date: skipping", card.Name)
		return ""
	}

	return parsedTime.Format(configuration.DateFormat)
}
