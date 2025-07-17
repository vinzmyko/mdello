// Trello API JSON objects → Markdown
package markdown

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/vinzmyko/mdello/config"
	"github.com/vinzmyko/mdello/trello"
)

type BoardSession struct {
	board    *trello.Board
	idMapper *idMapper
}

type idMapper struct {
	shortToFull map[string]string
	fullToShort map[string]string
}

func NewBoardSession(board *trello.Board, trelloClient *trello.TrelloClient) (*BoardSession, error) {
	session := &BoardSession{
		board: board,
		idMapper: &idMapper{
			shortToFull: make(map[string]string),
			fullToShort: make(map[string]string),
		},
	}

	err := session.buildIDMapping(trelloClient)
	return session, err
}

func (s *BoardSession) GetShortID(fullID string) string {
	return s.idMapper.fullToShort[fullID]
}

func (s *BoardSession) ResolveShortID(shortID string) (string, error) {
	if fullID, exists := s.idMapper.shortToFull[shortID]; exists {
		return fullID, nil
	}
	return "", fmt.Errorf("short ID %s not found", shortID)
}

func ConvertToMarkdown(trelloClient *trello.TrelloClient, configuration *config.Config, board *trello.Board) (string, *BoardSession, error) {
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

const SHORT_ID_LENGTH = 5

func shortID(trelloID string) string {
	hash := sha256.Sum256([]byte(trelloID))
	return fmt.Sprintf("%x", hash)[:SHORT_ID_LENGTH]
}

func formatDate(card trello.Card, configuration *config.Config) string {
	parsedTime, err := time.Parse(time.RFC3339, *card.Due)
	if err != nil {
		fmt.Printf("\nError parsing %s's due date: skipping", card.Name)
		return ""
	}

	return parsedTime.Format(configuration.DateFormat)
}

func (s *BoardSession) buildIDMapping(trelloClient *trello.TrelloClient) error {
	s.idMapper.addMapping(s.board.ID)

	lists, err := trelloClient.GetLists(s.board.ID)
	if err != nil {
		return err
	}

	for _, list := range lists {
		s.idMapper.addMapping(list.ID)

		cards, err := trelloClient.GetCards(list.ID)
		if err != nil {
			return err
		}

		for _, card := range cards {
			s.idMapper.addMapping(card.ID)
		}
	}

	return nil
}

func (m *idMapper) addMapping(fullID string) {
	shortID := shortID(fullID)
	m.shortToFull[shortID] = fullID
	m.fullToShort[fullID] = shortID
}
