package markdown

import (
	"errors"
	"fmt"
	"strings"

	"github.com/vinzmyko/mdello/config"
	"github.com/vinzmyko/mdello/trello"
)

type TrelloAction interface {
	Apply(client *trello.TrelloClient, ctx *ActionContext) error
	Description() string
}

type ActionContext struct {
	BoardID string
}

// === BOARD ACTIONS ===

type UpdateBoardNameAction struct {
	BoardID string
	OldName string
	NewName string
}

func (act UpdateBoardNameAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
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

// === LABEL ACTIONS ===
type CreateLabelAction struct {
	BoardID string
	Name    string
	Colour  string
}

func (act CreateLabelAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	params := &trello.CreateLabelParams{
		BoardID: act.BoardID,
		Name:    act.Name,
		Colour:  act.Colour,
	}
	_, err := t.CreateLabel(params)
	return err
}

func (act CreateLabelAction) Description() string {
	return fmt.Sprintf(`Created label "%s" with colour "%s"`, act.Name, act.Colour)
}

type UpdateLabelName struct {
	ID      string
	OldName string
	NewName string
}

func (act UpdateLabelName) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	params := &trello.UpdateLabelParams{
		ID:   act.ID,
		Name: &act.NewName,
	}
	_, err := t.UpdateLabel(params)
	return err
}

func (act UpdateLabelName) Description() string {
	return fmt.Sprintf(`Label "%s" renamed to "%s"`, act.OldName, act.NewName)
}

type UpdateLabelColour struct {
	ID        string
	Name      string
	OldColour string
	NewColour string
}

func (act UpdateLabelColour) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	params := &trello.UpdateLabelParams{
		ID:     act.ID,
		Colour: &act.NewColour,
	}
	_, err := t.UpdateLabel(params)
	return err
}

func (act UpdateLabelColour) Description() string {
	return fmt.Sprintf(`Label "%s" colour changed from "%s" to "%s"`, act.Name, act.OldColour, act.NewColour)
}

type DeleteLabelAction struct {
	ID   string
	Name string
}

func (act DeleteLabelAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	err := t.DeleteLabel(act.ID)
	return err
}

func (act DeleteLabelAction) Description() string {
	return fmt.Sprintf(`Label "%s" deleted`, act.Name)
}

// === LIST ACTIONS ===

type CreateListAction struct {
	BoardID  string
	Name     string
	Position int
}

func (act CreateListAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	posStr := fmt.Sprintf("%d.0", act.Position)

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

func (act UpdateListNameAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
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

func (act UpdateListPositionAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
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

func (act ArchiveListAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
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
	ListName    string
	Name        string
	Position    int
	IsCompleted bool
}

func (act CreateCardAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	lists, err := t.GetLists(ctx.BoardID)
	if err != nil {
		return errors.New("Failed to get board lists")
	}
	pos := fmt.Sprintf("%d", act.Position)

	var targetedListID string
	for _, list := range lists {
		if list.Name == act.ListName {
			targetedListID = list.ID
			break
		}
	}
	if targetedListID == "" {
		return fmt.Errorf(`List "%s" not found on board`, act.ListName)
	}

	params := &trello.CreateCardParams{
		IdList:      targetedListID,
		Name:        &act.Name,
		Pos:         &pos,
		DueComplete: &act.IsCompleted,
	}
	_, err = t.CreateCard(params)
	return err
}

func (act CreateCardAction) Description() string {
	return fmt.Sprintf(`Create card "%s" at position %d`, act.Name, act.Position)
}

type MoveCardAction struct {
	CardID   string
	Name     string
	FromList string
	ToList   string
	Position int
}

func (act MoveCardAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	posStr := fmt.Sprintf("%d.0", act.Position)

	params := &trello.UpdateCardParams{
		ID:     act.CardID,
		IdList: &act.ToList,
		Pos:    &posStr,
	}
	_, err := t.UpdateCard(params)
	return err
}

func (act MoveCardAction) Description() string {
	return fmt.Sprintf(`Card "%s" moved from list "%s" to "%s"`, act.Name, act.FromList, act.ToList)
}

type UpdateCardNameAction struct {
	CardID  string
	OldName string
	NewName string
}

func (act UpdateCardNameAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
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

func (act UpdateCardPositionAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
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

type UpdateCardIsCompletedAction struct {
	CardID     string
	Name       string
	IsComplete bool
}

func (act UpdateCardIsCompletedAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	params := &trello.UpdateCardParams{
		ID:          act.CardID,
		DueComplete: &act.IsComplete,
	}
	_, err := t.UpdateCard(params)
	return err
}

func (act UpdateCardIsCompletedAction) Description() string {
	if act.IsComplete {
		return fmt.Sprintf(`Card "%s" status is complete`, act.Name)
	} else {
		return fmt.Sprintf(`Card "%s" status is not complete`, act.Name)
	}
}

type AddCardLabelAction struct {
	CardID    string
	CardName  string
	LabelName string
}

func (act AddCardLabelAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	board, err := t.GetBoard(ctx.BoardID)
	if err != nil {
		return fmt.Errorf("failed to get board: %w", err)
	}

	var labelID string
	for _, label := range board.Labels {
		labelNameMarkdown := strings.ReplaceAll(label.Name, " ", "~")
		if labelNameMarkdown == act.LabelName {
			labelID = label.ID
			break
		}
	}
	if labelID == "" {
		return fmt.Errorf("label '%s' not found on board", act.LabelName)
	}

	// Find the card by name (since ID might be sentinel)
	lists, err := t.GetLists(ctx.BoardID)
	if err != nil {
		return fmt.Errorf("failed to get lists: %w", err)
	}

	for _, list := range lists {
		cards, err := t.GetCards(list.ID)
		if err != nil {
			continue // Skip this list if we can't get cards
		}

		for _, card := range cards {
			if card.Name == act.CardName {
				params := &trello.AddCardLabelParams{
					ID:      card.ID,
					Name:    &act.CardName,
					LabelID: labelID,
				}
				return t.AddCardLabel(params)
			}
		}
	}

	return fmt.Errorf("card '%s' not found on board", act.CardName)
}

func (act AddCardLabelAction) Description() string {
	return fmt.Sprintf(`Added label "%s" to card "%s"`, act.LabelName, act.CardName)
}

type UpdateCardDueDate struct {
	CardID string
	Name   string
	Due    string
	Cfg    *config.Config
}

func (act UpdateCardDueDate) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	// Find the card by name, since ID might be sentinel ID
	lists, err := t.GetLists(ctx.BoardID)
	if err != nil {
		return fmt.Errorf("failed to get lists: %w", err)
	}

	for _, list := range lists {
		cards, err := t.GetCards(list.ID)
		if err != nil {
			continue // Skip this list if we can't get cards
		}

		for _, card := range cards {
			if card.Name == act.Name {
				params := &trello.UpdateCardParams{
					ID:  card.ID,
					Due: &act.Due,
				}
				_, err := t.UpdateCard(params)
				return err
			}
		}
	}

	return fmt.Errorf("card '%s' not found on board", act.Name)
}

func (act UpdateCardDueDate) Description() string {
	formatDate := formatDate(act.Due, act.Cfg)
	return fmt.Sprintf(`Card "%s" due date updated to "%s"`, act.Name, formatDate)
}

type DeleteCardDueDate struct {
	CardID string
	Name   string
}

func (act DeleteCardDueDate) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	emptyString := ""
	params := &trello.UpdateCardParams{
		ID:  act.CardID,
		Due: &emptyString,
	}
	_, err := t.UpdateCard(params)
	return err
}

func (act DeleteCardDueDate) Description() string {
	return fmt.Sprintf(`Card "%s" due date removed`, act.Name)
}

type DeleteCardLabelAction struct {
	CardID    string
	CardName  string
	LabelID   string
	LabelName string
}

func (act DeleteCardLabelAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	params := &trello.DeleteCardLabelParams{
		ID:      act.CardID,
		Name:    &act.CardName,
		LabelID: act.LabelID,
	}
	err := t.DeleteCardLabel(params)
	return err
}

func (act DeleteCardLabelAction) Description() string {
	return fmt.Sprintf(`Removed label "%s" from card "%s"`, act.LabelName, act.CardName)
}

type DeleteCardAction struct {
	Name   string
	CardID string
}

func (act DeleteCardAction) Apply(t *trello.TrelloClient, ctx *ActionContext) error {
	err := t.DeleteCard(act.CardID)
	return err
}

func (act DeleteCardAction) Description() string {
	return fmt.Sprintf(`List "%s" deleted`, act.Name)
}
