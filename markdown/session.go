package markdown

import (
	"crypto/sha256"
	"fmt"

	"github.com/vinzmyko/mdello/trello"
)

const SHORT_ID_LENGTH = 5

type BoardSession struct {
	board    *trello.Board
	idMapper *idMapper
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

type idMapper struct {
	shortToFull map[string]string
	fullToShort map[string]string
}

func (m *idMapper) addMapping(fullID string) {
	shortID := generateShortID(fullID)
	m.shortToFull[shortID] = fullID
	m.fullToShort[fullID] = shortID
}

func generateShortID(trelloID string) string {
	hash := sha256.Sum256([]byte(trelloID))
	return fmt.Sprintf("%x", hash)[:SHORT_ID_LENGTH]
}
