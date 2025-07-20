package markdown

import (
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

	// Create map of all original cards
	allOriginalCards := make(map[string]*parsedCard)
	for _, list := range originalBoard.Lists {
		for _, card := range list.cards {
			allOriginalCards[card.id] = card
		}
	}

	// Create map of all the edited cards to check for moves
	allEditedCards := make(map[string]*parsedCard)
	for _, list := range editedBoard.Lists {
		for _, card := range list.cards {
			allEditedCards[card.id] = card
		}
	}

	for labelID, originalLabel := range originalLabelsMap {
		if editedLabel, exists := editedLabelsMap[labelID]; exists {
			if originalLabel.Name != editedLabel.Name {
				actions = append(actions, UpdateLabelName{
					ID:      originalLabel.ID,
					OldName: originalLabel.Name,
					NewName: editedLabel.Name,
				})
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

					actions = append(actions, checkCardProperties(originalCard, editedCard)...)
				} else {
					// Checks to see if the cardID exists in another list, if it does exist it's moved to another list
					if movedCard, movedToAnotherList := allEditedCards[cardID]; movedToAnotherList {
						actions = append(actions, MoveCardAction{
							CardID:   cardID,
							Name:     movedCard.name,
							FromList: originalList.id,
							ToList:   movedCard.listID,
							Position: movedCard.position,
						})

						actions = append(actions, checkCardProperties(originalCard, movedCard)...)
					} else {
						actions = append(actions, DeleteCardAction{
							CardID: originalCard.id,
							Name:   originalCard.name,
						})
					}
				}
			}

			for cardID, editedCard := range editedCardsMap {
				if _, exists := originalCardsMap[cardID]; !exists {
					if _, existedAnywhere := allOriginalCards[cardID]; !existedAnywhere {
						isComplete := strings.ToLower(editedCard.isComplete) == "x"

						actions = append(actions, CreateCardAction{
							ListID:      editedList.id,
							Name:        editedCard.name,
							Position:    editedCard.position,
							IsCompleted: isComplete,
						})
					}
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

func checkCardProperties(originalCard, editedCard *parsedCard) []TrelloAction {
	var actions []TrelloAction

	if originalCard.name != editedCard.name {
		actions = append(actions, UpdateCardNameAction{
			CardID:  originalCard.id,
			OldName: originalCard.name,
			NewName: editedCard.name,
		})
	}

	if originalCard.isComplete != editedCard.isComplete {
		isComplete := strings.ToLower(editedCard.isComplete) == "x"
		actions = append(actions, UpdateCardIsCompletedAction{
			CardID:     originalCard.id,
			IsComplete: isComplete,
			Name:       editedCard.name,
		})
	}

	// TODO: Add label

	// TODO: add duedate check

	return actions
}
