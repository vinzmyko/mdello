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
	return fmt.Sprintf(`Board "%s" name changed from "%s" to "%s"`, act.BoardID, act.OldName, act.NewName)
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
	return fmt.Sprintf(`Created list "%s" at position %d`, act.Name, act.Position)
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
	return fmt.Sprintf(`List "%s" renamed to "%s"`, act.OldName, act.NewName)
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
	return fmt.Sprintf(`List "%s" moved from position %d to %d`, act.Name, act.OldPosition, act.NewPosition)
}

type ArchiveListAction struct {
	ListID string
	Name   string
	Value  bool
}

func (act ArchiveListAction) Apply(t *trello.TrelloClient) error {
	val := false
	params := &trello.ArchiveListParams{
		ID:    act.ListID,
		Value: &val,
	}
	_, err := t.ArchiveList(params)
	return err
}

func (act ArchiveListAction) Description() string {
	if act.Value {
		return fmt.Sprintf(`List "%s" archived`, act.Name)
	} else {
		return fmt.Sprintf(`List "%s" unarchived`, act.Name)
	}
}

// === CARD ACTIONS ===

type CreateCardAction struct {
	ListID   string
	Name     string
	Position int
}

// TODO need to make it so that user can add in duedate labels, and is completed
func (act CreateCardAction) Apply(t *trello.TrelloClient) error {
	pos := fmt.Sprintf("%d", act.Position)
	params := &trello.CreateCardParams{
		IdList: act.ListID,
		Name:   &act.Name,
		Pos:    &pos,
	}
	_, err := t.CreateCard(params)
	return err
}

func (act CreateCardAction) Description() string {
	return fmt.Sprintf(`Create card "%s" at position %d`, act.Name, act.Position)
}

type UpdateCardNameAction struct {
	CardID  string
	OldName string
	NewName string
}

func (act UpdateCardNameAction) Apply(t *trello.TrelloClient) error {
	params := &trello.UpdateCardParams{
		ID:   act.CardID,
		Name: &act.NewName,
	}
	_, err := t.UpdateCard(params)
	return err
}

func (act UpdateCardNameAction) Description() string {
	return fmt.Sprintf(`Card "%s" renamed to "%s"`, act.OldName, act.NewName)
}

type UpdateCardPositionAction struct {
	CardID      string
	Name        string
	OldPosition int
	NewPosition int
}

func (act UpdateCardPositionAction) Apply(t *trello.TrelloClient) error {
	posStr := fmt.Sprintf("%d.0", act.NewPosition)

	params := &trello.UpdateCardParams{
		ID:  act.CardID,
		Pos: &posStr,
	}
	_, err := t.UpdateCard(params)
	return err
}

func (act UpdateCardPositionAction) Description() string {
	return fmt.Sprintf(`Card "%s" moved from position %d to %d`, act.Name, act.OldPosition, act.NewPosition)
}

type DeleteCardAction struct {
	Name   string
	CardID string
}

func (act DeleteCardAction) Apply(t *trello.TrelloClient) error {
	err := t.DeleteCard(act.CardID)
	return err
}

func (act DeleteCardAction) Description() string {
	return fmt.Sprintf(`List "%s" deleted`, act.Name)
}
