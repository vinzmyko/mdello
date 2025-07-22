package markdown

import (
	"fmt"
	"strings"

	"github.com/vinzmyko/mdello/config"
	"github.com/vinzmyko/mdello/trello"
)

func Diff(originalBoard, editedBoard *ParsedBoard, cfg *config.Config) (*DiffResult, error) {
	quickActions := make([]TrelloAction, 0)
	detailedTrelloAction := make([]detailedTrelloAction, 0)

	if originalBoard.Name != editedBoard.Name {
		quickActions = append(quickActions, UpdateBoardNameAction{
			BoardID: originalBoard.ID,
			OldName: originalBoard.Name,
			NewName: editedBoard.Name,
		})
	}

	if editedBoard.DetailedEdit {
		detailedTrelloAction = append(detailedTrelloAction, createDetailedAction(OTBoard, originalBoard.ID, editedBoard.Name))
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
				quickActions = append(quickActions, UpdateLabelName{
					ID:      originalLabel.ID,
					OldName: originalLabel.Name,
					NewName: editedLabel.Name,
				})
			}

			if originalLabel.Colour != editedLabel.Colour {
				quickActions = append(quickActions, UpdateLabelColour{
					ID:        originalLabel.ID,
					Name:      originalLabel.Name,
					OldColour: originalLabel.Colour,
					NewColour: editedLabel.Colour,
				})
			}

		} else {
			quickActions = append(quickActions, DeleteLabelAction{
				ID:   originalLabel.ID,
				Name: originalLabel.Name,
			})
		}
	}

	for labelID, editedLabel := range editedLabelsMap {
		if _, exists := originalLabelsMap[labelID]; !exists {
			quickActions = append(quickActions, CreateLabelAction{
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
				quickActions = append(quickActions, UpdateListNameAction{
					ListID:  originalList.id,
					OldName: originalList.name,
					NewName: editedList.name,
				})
			}

			if originalList.markdownIdx != editedList.markdownIdx {
				quickActions = append(quickActions, UpdateListPositionAction{
					ListID:      originalList.id,
					Name:        originalList.name,
					OldPosition: originalList.markdownIdx,
					NewPosition: editedList.markdownIdx,
				})
			}

			if editedList.detailedEdit {
				detailedTrelloAction = append(detailedTrelloAction, createDetailedAction(OTList, originalList.id, editedList.name))
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
						quickActions = append(quickActions, UpdateCardPositionAction{
							CardID:      originalCard.id,
							Name:        originalCard.name,
							OldPosition: originalCard.position,
							NewPosition: editedCard.position,
						})
					}

					cardActions, err := checkCardProperties(originalCard, editedCard, cfg)
					if err != nil {
						return nil, fmt.Errorf("error checking card properties for card '%s': %w", originalCard.name, err)
					}
					quickActions = append(quickActions, cardActions...)

					if editedCard.detailedEdit {
						detailedTrelloAction = append(detailedTrelloAction, createDetailedAction(OTCard, editedCard.id, editedCard.name))
					}

				} else {
					// Checks to see if the cardID exists in another list, if it does exist it's moved to another list
					if movedCard, movedToAnotherList := allEditedCards[cardID]; movedToAnotherList {
						quickActions = append(quickActions, MoveCardAction{
							CardID:   cardID,
							Name:     movedCard.name,
							FromList: originalList.id,
							ToList:   movedCard.listID,
							Position: movedCard.position,
						})

						cardActions, err := checkCardProperties(originalCard, movedCard, cfg)
						if err != nil {
							return nil, fmt.Errorf("error checking card properties for card '%s': %w", originalCard.name, err)
						}
						quickActions = append(quickActions, cardActions...)

					} else {
						quickActions = append(quickActions, DeleteCardAction{
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

						// Composition when creating a new card
						quickActions = append(quickActions, CreateCardAction{
							ListName:    editedList.name,
							Name:        editedCard.name,
							Position:    editedCard.position,
							IsCompleted: isComplete,
						})

						for _, labelName := range editedCard.labels {
							quickActions = append(quickActions, AddCardLabelAction{
								CardID:    cardID,
								CardName:  editedCard.name,
								LabelName: labelName,
							})
						}

						if editedCard.dueDate != "" {
							rfcDateFormat, err := parseMarkdownDate(editedCard.dueDate, cfg)
							if err != nil {
								return nil, fmt.Errorf("invalid due date format: %w", err)
							}
							quickActions = append(quickActions, UpdateCardDueDate{
								CardID: cardID,
								Name:   editedCard.name,
								Due:    rfcDateFormat,
								Cfg:    cfg,
							})
						}
					}
				}
			}

		} else {
			quickActions = append(quickActions, ArchiveListAction{
				ListID: originalList.id,
				Name:   originalList.name,
				Value:  true,
			})
		}
	}

	for listID, editedList := range editedListMap {
		if _, exists := originalListsMap[listID]; !exists {
			quickActions = append(quickActions, CreateListAction{
				BoardID:  originalBoard.ID,
				Name:     editedList.name,
				Position: editedList.markdownIdx,
			})

			for _, editedCard := range editedList.cards {
				isComplete := strings.ToLower(editedCard.isComplete) == "x"

				quickActions = append(quickActions, CreateCardAction{
					ListName:    editedList.name,
					Name:        editedCard.name,
					Position:    editedCard.position,
					IsCompleted: isComplete,
				})

				for _, labelName := range editedCard.labels {
					quickActions = append(quickActions, AddCardLabelAction{
						CardID:    editedCard.id,
						CardName:  editedCard.name,
						LabelName: labelName,
					})
				}

				if editedCard.dueDate != "" {
					rfcDateFormat, err := parseMarkdownDate(editedCard.dueDate, cfg)
					if err != nil {
						return nil, fmt.Errorf("invalid due date format: %w", err)
					}
					quickActions = append(quickActions, UpdateCardDueDate{
						CardID: editedCard.id,
						Name:   editedCard.name,
						Due:    rfcDateFormat,
						Cfg:    cfg,
					})
				}
			}
		}
	}

	diffResult := DiffResult{
		QuickActions:    quickActions,
		DetailedActions: detailedTrelloAction,
	}

	return &diffResult, nil
}

func checkCardProperties(originalCard, editedCard *parsedCard, cfg *config.Config) ([]TrelloAction, error) {
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

	originalCardLabelsMap := make(map[string]bool)
	for _, label := range originalCard.labels {
		originalCardLabelsMap[label] = true
	}

	editedCardLabelsMap := make(map[string]bool)
	for _, label := range editedCard.labels {
		editedCardLabelsMap[label] = true
	}

	// In original but not in edited
	for _, label := range originalCard.labels {
		if !editedCardLabelsMap[label] {
			actions = append(actions, DeleteCardLabelAction{
				CardID:    originalCard.id,
				CardName:  editedCard.name,
				LabelName: label,
			})
		}
	}

	// In edited but not in original
	for _, label := range editedCard.labels {
		if !originalCardLabelsMap[label] {
			actions = append(actions, AddCardLabelAction{
				CardID:    originalCard.id,
				CardName:  editedCard.name,
				LabelName: label,
			})
		}
	}

	if originalCard.dueDate != editedCard.dueDate {
		if editedCard.dueDate == "" {
			actions = append(actions, DeleteCardDueDate{
				CardID: originalCard.id,
				Name:   editedCard.name,
			})
		} else {
			rfcDateFormat, err := parseMarkdownDate(editedCard.dueDate, cfg)
			if err != nil {
				return nil, fmt.Errorf("invalid due date format: %w", err)
			}
			actions = append(actions, UpdateCardDueDate{
				CardID: originalCard.id,
				Cfg:    cfg,
				Name:   editedCard.name,
				Due:    rfcDateFormat,
			})
		}
	}

	return actions, nil
}

func createDetailedAction(objectType ObjectType, objectID string, objectName string) detailedTrelloAction {
	return detailedTrelloAction{
		ObjectType: objectType,
		ObjectID:   objectID,
		ObjectName: objectName,
	}
}
