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

func GenerateDetailedMarkdown(detailedActions []detailedTrelloAction) string {
	var content strings.Builder
	for _, detailedAction := range detailedActions {
		switch detailedAction.ObjectType {
		case OTBoard:
			content.WriteString(GenerateDetailedBoardContent(detailedAction))
		case OTList:
			content.WriteString(GenerateDetailedListContent(detailedAction))
		case OTCard:
			content.WriteString(GenerateDetailedCardContent(detailedAction))
		}
	}
	return content.String()
}

func GenerateDetailedBoardContent(action detailedTrelloAction) string {
	var content strings.Builder

	content.WriteString(generateSectionHeader(action.ObjectName, string(action.ObjectType)))

	return content.String()
}

func GenerateDetailedListContent(action detailedTrelloAction) string {
	var content strings.Builder

	content.WriteString(generateSectionHeader(action.ObjectName, string(action.ObjectType)))

	return content.String()
}

func GenerateDetailedCardContent(action detailedTrelloAction) string {
	var content strings.Builder

	content.WriteString(generateSectionHeader(action.ObjectName, string(action.ObjectType)))

	return content.String()
}

const sectionHeaderSeparatorLength = 77

func generateSectionHeader(title, objectType string) string {
	separator := strings.Repeat("=", sectionHeaderSeparatorLength)
	return fmt.Sprintf("# %s\n# EDITING %s: %s\n# %s\n",
		separator, strings.ToUpper(objectType), title, separator)
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
