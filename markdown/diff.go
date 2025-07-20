package markdown

import (
	"fmt"
	"strings"

	"github.com/vinzmyko/mdello/trello"
)

func Diff(originalBoard, editedBoard *ParsedBoard) []TrelloAction {
	actions := make([]TrelloAction, 0)

	if originalBoard.Name != editedBoard.Name {
		actions = append(actions, UpdateBoardNameAction{
			BoardID: originalBoard.ID,
			OldName: originalBoard.Name,
			NewName: editedBoard.Name,
		})
	}

	originalLabelsMap := make(map[string]*trello.Label)
	for _, label := range originalBoard.Labels {
		originalLabelsMap[label.ID] = label
	}

	editedLabelsMap := make(map[string]*trello.Label)
	for _, label := range editedBoard.Labels {
		editedLabelsMap[label.ID] = label
	}

	for labelID, originalLabel := range originalLabelsMap {
		if editedLabel, exists := editedLabelsMap[labelID]; exists {
			if originalLabel.Name != editedLabel.Name {
				actions = append(actions, UpdateLabelName{
					ID:      originalLabel.ID,
					OldName: originalLabel.Name,
					NewName: editedLabel.Name,
				})
				fmt.Printf("\n OriginalLabelID: %s", originalLabel.ID)
				fmt.Printf("\n Editted LabelID: %s", editedLabel.ID)
			}

			if originalLabel.Colour != editedLabel.Colour {
				actions = append(actions, UpdateLabelColour{
					ID:        originalLabel.ID,
					Name:      originalLabel.Name,
					OldColour: originalLabel.Colour,
					NewColour: editedLabel.Colour,
				})
			}

		} else {
			actions = append(actions, DeleteLabelAction{
				ID:   originalLabel.ID,
				Name: originalLabel.Name,
			})
		}
	}

	for labelID, editedLabel := range editedLabelsMap {
		if _, exists := originalLabelsMap[labelID]; !exists {
			actions = append(actions, CreateLabelAction{
				BoardID: originalBoard.ID,
				Name:    editedLabel.Name,
				Colour:  editedLabel.Colour,
			})
		}
	}

	originalListsMap := make(map[string]*parsedList)
	for _, list := range originalBoard.Lists {
		originalListsMap[list.id] = list
	}

	editedListMap := make(map[string]*parsedList)
	for _, list := range editedBoard.Lists {
		editedListMap[list.id] = list
	}

	for listID, originalList := range originalListsMap {
		if editedList, exists := editedListMap[listID]; exists {
			if originalList.name != editedList.name {
				actions = append(actions, UpdateListNameAction{
					ListID:  originalList.id,
					OldName: originalList.name,
					NewName: editedList.name,
				})
			}

			if originalList.markdownIdx != editedList.markdownIdx {
				actions = append(actions, UpdateListPositionAction{
					ListID:      originalList.id,
					Name:        originalList.name,
					OldPosition: originalList.markdownIdx,
					NewPosition: editedList.markdownIdx,
				})
			}

			originalCardsMap := make(map[string]*parsedCard)
			for _, card := range originalList.cards {
				originalCardsMap[card.id] = card
			}

			editedCardsMap := make(map[string]*parsedCard)
			for _, card := range editedList.cards {
				editedCardsMap[card.id] = card
			}

			for cardID, originalCard := range originalCardsMap {
				if editedCard, exists := editedCardsMap[cardID]; exists {
					if originalCard.position != editedCard.position {
						actions = append(actions, UpdateCardPositionAction{
							CardID:      originalCard.id,
							Name:        originalCard.name,
							OldPosition: originalCard.position,
							NewPosition: editedCard.position,
						})
					}

					if originalCard.name != editedCard.name {
						actions = append(actions, UpdateCardNameAction{
							CardID:  originalCard.id,
							OldName: originalCard.name,
							NewName: editedCard.name,
						})
					}

					if originalCard.isComplete != editedCard.isComplete {
						isComplete := strings.ToLower(editedCard.isComplete) == ""

						actions = append(actions, UpdateCardIsCompletedAction{
							CardID:     originalCard.id,
							IsComplete: isComplete,
							Name:       editedCard.name,
						})
					}

					// TODO: Check for label change

					// TODO: Check for due date change
				} else {
					actions = append(actions, DeleteCardAction{
						CardID: originalCard.id,
						Name:   originalCard.name,
					})
				}
			}

			for cardID, editedCard := range editedCardsMap {
				if _, exists := originalCardsMap[cardID]; !exists {
					isComplete := strings.ToLower(editedCard.isComplete) == "x"

					actions = append(actions, CreateCardAction{
						ListID:      originalList.id,
						Name:        editedCard.name,
						Position:    editedCard.position,
						IsCompleted: isComplete,
					})
				}
			}

		} else {
			actions = append(actions, ArchiveListAction{
				ListID: originalList.id,
				Name:   originalList.name,
				Value:  true,
			})
		}
	}

	for listID, editedList := range editedListMap {
		if _, exists := originalListsMap[listID]; !exists {
			actions = append(actions, CreateListAction{
				BoardID:  originalBoard.ID,
				Name:     editedList.name,
				Position: editedList.markdownIdx,
			})
		}
	}

	return actions
}
