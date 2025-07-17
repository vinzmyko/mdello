package markdown

import (
	"fmt"

	"github.com/vinzmyko/mdello/trello"
)

type TrelloAction interface {
	Apply(client *trello.TrelloClient) error
	Description() string
}

// === BOARD ACTIONS ===

type UpdateBoardNameAction struct {
	BoardID string
	OldName string
	NewName string
}

func (act UpdateBoardNameAction) Apply(t *trello.TrelloClient) error {
	params := &trello.UpdateBoardParams{
		ID:   act.BoardID,
		Name: &act.NewName,
	}

	_, err := t.UpdateBoard(params)
	return err
}

func (act UpdateBoardNameAction) Description() string {
	return fmt.Sprintf("Board name changed from '%s' to '%s'", act.OldName, act.NewName)
}

// === LIST ACTIONS ===

type CreateListAction struct {
	BoardID  string
	Name     string
	Position int
}

func (act CreateListAction) Apply(t *trello.TrelloClient) error {
	posStr := fmt.Sprintf("%d", act.Position)

	params := &trello.CreateListParams{
		IdBoard: act.BoardID,
		Name:    act.Name,
		Pos:     &posStr,
	}
	_, err := t.CreateList(params)
	return err
}

func (act CreateListAction) Description() string {
	return fmt.Sprintf("Create new list: '%s'", act.Name)
}

// Need to handle due date, labels after we see it working
type UpdateListNameAction struct {
	ListID  string
	OldName string
	NewName string
}

func (act UpdateListNameAction) Apply(t *trello.TrelloClient) error {
	params := &trello.UpdateListParams{
		ID:   act.ListID,
		Name: &act.NewName,
	}
	_, err := t.UpdateList(params)
	return err
}

func (act UpdateListNameAction) Description() string {
	return fmt.Sprintf("List name changed from '%s' to '%s'", act.OldName, act.NewName)
}

type UpdateListPositionAction struct {
	ListID      string
	Name        string
	OldPosition int
	NewPosition int
}

func (act UpdateListPositionAction) Apply(t *trello.TrelloClient) error {
	posStr := fmt.Sprintf("%d.0", act.NewPosition)

	params := &trello.UpdateListParams{
		ID:  act.ListID,
		Pos: &posStr,
	}
	_, err := t.UpdateList(params)
	return err
}

func (act UpdateListPositionAction) Description() string {
	return fmt.Sprintf("List '%s' position moved from '%d' to '%d'", act.Name, act.OldPosition, act.NewPosition)
}

// === CARD ACTIONS ===

type UpdateCardPosition struct {
	CardID      string
	Name        string
	OldPosition int
	NewPosition int
}

func (act UpdateCardPosition) Apply(t *trello.TrelloClient) error {
	posStr := fmt.Sprintf("%d.0", act.NewPosition)

	params := &trello.UpdateCardParams{
		ID:  act.CardID,
		Pos: &posStr,
	}
	_, err := t.UpdateCard(params)
	return err
}

func (act UpdateCardPosition) Description() string {
	return fmt.Sprintf("Card '%s' position moved from '%d' to '%d'", act.Name, act.OldPosition, act.NewPosition)
}
