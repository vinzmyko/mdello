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

		case OTList:

		case OTCard:

		}
	}
	return content.String()
}

func GenerateDetailedBoardContent(action detailedTrelloAction, trelloClient *trello.TrelloClient) (string, error) {
	var content strings.Builder

	board, err := trelloClient.GetBoard(string(action.ObjectID))
	if err != nil {
		return "", fmt.Errorf("Failed to get board when generating detailed board content: %w", err)
	}

	content.WriteString(generateSectionHeader(action.ObjectName, string(action.ObjectType)))

	content.WriteString("\n## Description\n")
	content.WriteString("=== DESCRIPTION START ===\n")
	if board.Desc != "" {
		content.WriteString(board.Desc)
	}
	content.WriteString("\n=== DESCRIPTION END ===\n")

	content.WriteString("\n## Board Settings\n")
	content.WriteString(fmt.Sprintf("Name: %s\n", board.Name))
	content.WriteString(fmt.Sprintf("Closed: %t  # [true|false] – whether the board is archived/deactivated\n", board.Closed))

	content.WriteString("\n## Permissions\n")
	content.WriteString(fmt.Sprintf("Permission Level: %s  # [org|private|public] – board visibility\n", board.Prefs.PermissionLevel))
	content.WriteString(fmt.Sprintf("Voting: %s  # [disabled|members|observers|org|public] – who can vote on cards\n", board.Prefs.Voting))
	content.WriteString(fmt.Sprintf("Comments: %s  # [disabled|members|observers|org|public] – who can comment\n", board.Prefs.Comments))
	content.WriteString(fmt.Sprintf("Invitations: %s  # [admins|members] – who can invite new members\n", board.Prefs.Invitations))

	content.WriteString("\n## Display Preferences\n")
	content.WriteString(fmt.Sprintf("Self Join: %t  # allow users to join without invite\n", board.Prefs.SelfJoin))
	content.WriteString(fmt.Sprintf("Card Covers: %t  # display images on card fronts\n", board.Prefs.CardCovers))
	content.WriteString(fmt.Sprintf("Hide Votes: %t  # hide member vote counts on cards\n", board.Prefs.HideVotes))
	content.WriteString(fmt.Sprintf("Card Aging: %s  # [regular|pirate] – visual aging of neglected cards\n", board.Prefs.CardAging))
	content.WriteString(fmt.Sprintf("Calendar Feed: %t  # enable .ics feed for calendar integrations\n", board.Prefs.CalendarFeedEnabled))

	content.WriteString("\n\n")
	return content.String(), nil
}

func GenerateDetailedListContent(action detailedTrelloAction, trelloClient *trello.TrelloClient) (string, error) {
	var content strings.Builder

	list, err := trelloClient.GetList(string(action.ObjectID))
	if err != nil {
		return "", fmt.Errorf("Failed to get list when generating detailed list content: %w", err)
	}

	content.WriteString(generateSectionHeader(action.ObjectName, string(action.ObjectType)))

	content.WriteString("\n## List Settings\n")
	content.WriteString(fmt.Sprintf("Name: %s\n", list.Name))
	content.WriteString(fmt.Sprintf("Closed: %t  # [true|false] - Archive this list\n", list.Closed))
	content.WriteString(fmt.Sprintf("Position: %d\n", int(list.Pos)))
	content.WriteString(fmt.Sprintf("Subscribed: %t  # [true|false] - Get notifications for this list\n", list.Subscribed))

	content.WriteString("\n\n")

	return content.String(), nil
}

func GenerateDetailedCardContent(action detailedTrelloAction, trelloClient *trello.TrelloClient, cfg *config.Config) (string, error) {
	var content strings.Builder

	card, err := trelloClient.GetCard(string(action.ObjectID))
	if err != nil {
		return "", fmt.Errorf("Failed to get card when generating detailed card content: %w", err)
	}

	content.WriteString(generateSectionHeader(action.ObjectName, string(action.ObjectType)))

	content.WriteString("\n## Description\n")
	content.WriteString("=== DESCRIPTION START ===\n")
	if card.Desc != "" {
		content.WriteString(card.Desc)
	}
	content.WriteString("\n=== DESCRIPTION END ===\n")

	content.WriteString("\n## Basic Settings\n")
	content.WriteString(fmt.Sprintf("Name: %s\n", card.Name))
	content.WriteString(fmt.Sprintf("Closed: %t  # – archived state\n", card.Closed))
	content.WriteString(fmt.Sprintf("Position: %d \n", int(card.Pos)))
	content.WriteString(fmt.Sprintf("Subscribed: %t  # [true|false] – receive notifications on card activity\n", card.Subscribed))

	content.WriteString("\n## Scheduling\n")

	empty := ""
	if card.Badges.Start != &empty {
		formattedStart := formatDate(*card.Badges.Start, cfg)
		content.WriteString(fmt.Sprintf("Start: %s \n", formattedStart))
	} else {
		content.WriteString("Start Date:   # (leave empty for no start date)\n")
	}

	if card.Due != &empty {
		formattedDue := formatDate(*card.Due, cfg)
		content.WriteString(fmt.Sprintf("Due: %s  # due date\n", formattedDue))
	} else {
		content.WriteString("Due Date:   # (leave empty for no due date)\n")
	}


	content.WriteString(fmt.Sprintf("Due Complete: %t  # [true|false] – if card is marked done\n", card.Badges.DueComplete))

	content.WriteString("\n\n")

	return content.String(), nil
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
