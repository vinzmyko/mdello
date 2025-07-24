package diff

import (
	"fmt"
	"strings"

	"github.com/vinzmyko/mdello/config"
	"github.com/vinzmyko/mdello/markdown"
	"github.com/vinzmyko/mdello/trello"
)

func QuickActionsDiff(originalBoard, editedBoard *markdown.ParsedBoard, cfg *config.Config) (*markdown.DiffResult, error) {
	quickActions := make([]markdown.TrelloAction, 0)
	detailedTrelloAction := make([]markdown.DetailedTrelloAction, 0)

	if originalBoard.Name != editedBoard.Name {
		quickActions = append(quickActions, markdown.UpdateBoardNameAction{
			BoardID: originalBoard.ID,
			OldName: originalBoard.Name,
			NewName: editedBoard.Name,
		})
	}

	if editedBoard.DetailedEdit {
		detailedTrelloAction = append(detailedTrelloAction, createDetailedAction("board", originalBoard.ID, editedBoard.Name))
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
	allOriginalCards := make(map[string]*markdown.ParsedCard)
	for _, list := range originalBoard.Lists {
		for _, card := range list.Cards {
			allOriginalCards[card.ID] = card
		}
	}

	// Create map of all the edited cards to check for moves
	allEditedCards := make(map[string]*markdown.ParsedCard)
	for _, list := range editedBoard.Lists {
		for _, card := range list.Cards {
			allEditedCards[card.ID] = card
		}
	}

	for labelID, originalLabel := range originalLabelsMap {
		if editedLabel, exists := editedLabelsMap[labelID]; exists {
			if originalLabel.Name != editedLabel.Name {
				quickActions = append(quickActions, markdown.UpdateLabelName{
					ID:      originalLabel.ID,
					OldName: originalLabel.Name,
					NewName: editedLabel.Name,
				})
			}

			if originalLabel.Colour != editedLabel.Colour {
				quickActions = append(quickActions, markdown.UpdateLabelColour{
					ID:        originalLabel.ID,
					Name:      originalLabel.Name,
					OldColour: originalLabel.Colour,
					NewColour: editedLabel.Colour,
				})
			}

		} else {
			quickActions = append(quickActions, markdown.DeleteLabelAction{
				ID:   originalLabel.ID,
				Name: originalLabel.Name,
			})
		}
	}

	for labelID, editedLabel := range editedLabelsMap {
		if _, exists := originalLabelsMap[labelID]; !exists {
			quickActions = append(quickActions, markdown.CreateLabelAction{
				BoardID: originalBoard.ID,
				Name:    editedLabel.Name,
				Colour:  editedLabel.Colour,
			})
		}
	}

	originalListsMap := make(map[string]*markdown.ParsedList)
	for _, list := range originalBoard.Lists {
		originalListsMap[list.ID] = list
	}

	editedListMap := make(map[string]*markdown.ParsedList)
	for _, list := range editedBoard.Lists {
		editedListMap[list.ID] = list
	}

	for listID, originalList := range originalListsMap {
		if editedList, exists := editedListMap[listID]; exists {
			if originalList.Name != editedList.Name {
				quickActions = append(quickActions, markdown.UpdateListNameAction{
					ListID:  originalList.ID,
					OldName: originalList.Name,
					NewName: editedList.Name,
				})
			}

			if originalList.MarkdownIdx != editedList.MarkdownIdx {
				quickActions = append(quickActions, markdown.UpdateListPositionAction{
					ListID:      originalList.ID,
					Name:        originalList.Name,
					OldPosition: originalList.MarkdownIdx,
					NewPosition: editedList.MarkdownIdx,
				})
			}

			if editedList.DetailedEdit {
				detailedTrelloAction = append(detailedTrelloAction, createDetailedAction("list", originalList.ID, editedList.Name))
			}

			originalCardsMap := make(map[string]*markdown.ParsedCard)
			for _, card := range originalList.Cards {
				originalCardsMap[card.ID] = card
			}

			editedCardsMap := make(map[string]*markdown.ParsedCard)
			for _, card := range editedList.Cards {
				editedCardsMap[card.ID] = card
			}

			for cardID, originalCard := range originalCardsMap {
				if editedCard, exists := editedCardsMap[cardID]; exists {
					if originalCard.Position != editedCard.Position {
						quickActions = append(quickActions, markdown.UpdateCardPositionAction{
							CardID:      originalCard.ID,
							Name:        originalCard.Name,
							OldPosition: originalCard.Position,
							NewPosition: editedCard.Position,
						})
					}

					cardActions, err := checkCardProperties(originalCard, editedCard, cfg)
					if err != nil {
						return nil, fmt.Errorf("error checking card properties for card '%s': %w", originalCard.Name, err)
					}
					quickActions = append(quickActions, cardActions...)

					if editedCard.DetailedEdit {
						detailedTrelloAction = append(detailedTrelloAction, createDetailedAction("card", editedCard.ID, editedCard.Name))
					}

				} else {
					// Checks to see if the cardID exists in another list, if it does exist it's moved to another list
					if movedCard, movedToAnotherList := allEditedCards[cardID]; movedToAnotherList {
						quickActions = append(quickActions, markdown.MoveCardAction{
							CardID:   cardID,
							Name:     movedCard.Name,
							FromList: originalList.ID,
							ToList:   movedCard.ListID,
							Position: movedCard.Position,
						})

						cardActions, err := checkCardProperties(originalCard, movedCard, cfg)
						if err != nil {
							return nil, fmt.Errorf("error checking card properties for card '%s': %w", originalCard.Name, err)
						}
						quickActions = append(quickActions, cardActions...)

					} else {
						quickActions = append(quickActions, markdown.DeleteCardAction{
							CardID: originalCard.ID,
							Name:   originalCard.Name,
						})
					}
				}
			}

			for cardID, editedCard := range editedCardsMap {
				if _, exists := originalCardsMap[cardID]; !exists {
					if _, existedAnywhere := allOriginalCards[cardID]; !existedAnywhere {
						isComplete := strings.ToLower(editedCard.IsComplete) == "x"

						// Composition when creating a new card
						quickActions = append(quickActions, markdown.CreateCardAction{
							ListName:    editedList.Name,
							Name:        editedCard.Name,
							Position:    editedCard.Position,
							IsCompleted: isComplete,
						})

						for _, labelName := range editedCard.Labels {
							quickActions = append(quickActions, markdown.AddCardLabelAction{
								CardID:    cardID,
								CardName:  editedCard.Name,
								LabelName: labelName,
							})
						}

						if editedCard.DueDate != "" {
							rfcDateFormat, err := markdown.ParseMarkdownDate(editedCard.DueDate, cfg)
							if err != nil {
								return nil, fmt.Errorf("invalid due date format: %w", err)
							}
							quickActions = append(quickActions, markdown.UpdateCardDueDate{
								CardID: cardID,
								Name:   editedCard.Name,
								Due:    rfcDateFormat,
								Cfg:    cfg,
							})
						}
					}
				}
			}

		} else {
			quickActions = append(quickActions, markdown.ArchiveListAction{
				ListID: originalList.ID,
				Name:   originalList.Name,
				Value:  true,
			})
		}
	}

	for listID, editedList := range editedListMap {
		if _, exists := originalListsMap[listID]; !exists {
			quickActions = append(quickActions, markdown.CreateListAction{
				BoardID:  originalBoard.ID,
				Name:     editedList.Name,
				Position: editedList.MarkdownIdx,
			})

			for _, editedCard := range editedList.Cards {
				isComplete := strings.ToLower(editedCard.IsComplete) == "x"

				quickActions = append(quickActions, markdown.CreateCardAction{
					ListName:    editedList.Name,
					Name:        editedCard.Name,
					Position:    editedCard.Position,
					IsCompleted: isComplete,
				})

				for _, labelName := range editedCard.Labels {
					quickActions = append(quickActions, markdown.AddCardLabelAction{
						CardID:    editedCard.ID,
						CardName:  editedCard.Name,
						LabelName: labelName,
					})
				}

				if editedCard.DueDate != "" {
					rfcDateFormat, err := markdown.ParseMarkdownDate(editedCard.DueDate, cfg)
					if err != nil {
						return nil, fmt.Errorf("invalid due date format: %w", err)
					}
					quickActions = append(quickActions, markdown.UpdateCardDueDate{
						CardID: editedCard.ID,
						Name:   editedCard.Name,
						Due:    rfcDateFormat,
						Cfg:    cfg,
					})
				}
			}
		}
	}

	diffResult := markdown.DiffResult{
		QuickActions:    quickActions,
		DetailedActions: detailedTrelloAction,
	}

	return &diffResult, nil
}

func checkCardProperties(originalCard, editedCard *markdown.ParsedCard, cfg *config.Config) ([]markdown.TrelloAction, error) {
	var actions []markdown.TrelloAction

	if originalCard.Name != editedCard.Name {
		actions = append(actions, markdown.UpdateCardNameAction{
			CardID:  originalCard.ID,
			OldName: originalCard.Name,
			NewName: editedCard.Name,
		})
	}

	if originalCard.IsComplete != editedCard.IsComplete {
		isComplete := strings.ToLower(editedCard.IsComplete) == "x"
		actions = append(actions, markdown.UpdateCardIsCompletedAction{
			CardID:     originalCard.ID,
			IsComplete: isComplete,
			Name:       editedCard.Name,
		})
	}

	originalCardLabelsMap := make(map[string]bool)
	for _, label := range originalCard.Labels {
		originalCardLabelsMap[label] = true
	}

	editedCardLabelsMap := make(map[string]bool)
	for _, label := range editedCard.Labels {
		editedCardLabelsMap[label] = true
	}

	// In original but not in edited
	for _, label := range originalCard.Labels {
		if !editedCardLabelsMap[label] {
			actions = append(actions, markdown.DeleteCardLabelAction{
				CardID:    originalCard.ID,
				CardName:  editedCard.Name,
				LabelName: label,
			})
		}
	}

	// In edited but not in original
	for _, label := range editedCard.Labels {
		if !originalCardLabelsMap[label] {
			actions = append(actions, markdown.AddCardLabelAction{
				CardID:    originalCard.ID,
				CardName:  editedCard.Name,
				LabelName: label,
			})
		}
	}

	if originalCard.DueDate != editedCard.DueDate {
		if editedCard.DueDate == "" {
			actions = append(actions, markdown.DeleteCardDueDate{
				CardID: originalCard.ID,
				Name:   editedCard.Name,
			})
		} else {
			rfcDateFormat, err := markdown.ParseMarkdownDate(editedCard.DueDate, cfg)
			if err != nil {
				return nil, fmt.Errorf("invalid due date format: %w", err)
			}
			actions = append(actions, markdown.UpdateCardDueDate{
				CardID: originalCard.ID,
				Cfg:    cfg,
				Name:   editedCard.Name,
				Due:    rfcDateFormat,
			})
		}
	}

	return actions, nil
}

func createDetailedAction(objectType string, objectID string, objectName string) markdown.DetailedTrelloAction {
	return markdown.DetailedTrelloAction{
		ObjectType: objectType,
		ObjectID:   objectID,
		ObjectName: objectName,
	}
}
