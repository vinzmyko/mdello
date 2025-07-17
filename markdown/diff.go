package markdown

func Diff(originalBoard, editedBoard *ParsedBoard) []TrelloAction {
	actions := make([]TrelloAction, 0)

	if originalBoard.Name != editedBoard.Name {
		actions = append(actions, UpdateBoardNameAction{
			BoardID: originalBoard.ID,
			OldName: originalBoard.Name,
			NewName: editedBoard.Name,
		})
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
			// If lists exists for both, check for modifications (name, position)
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
					// TODO Check for name changes

					// TODO Check for status change

					if originalCard.position != editedCard.position {
						actions = append(actions, UpdateCardPosition{
							CardID:      originalCard.id,
							Name:        originalCard.name,
							OldPosition: originalCard.position,
							NewPosition: editedCard.position,
						})
					}

					// TODO Check for label change

					// TODO Check for due date change
				} else {
					// Card does not exist in the editted list thus deleted
					// TODO Add DeleteCardAction
				}
			}

			// for cardID, editedCard := range editedCardsMap {
			// 	if _, exists := originalCardsMap[cardID]; !exists {
			// 		// TODO CreateCardAction
			// 	}
			// }

			// for cardID, originalCard := range originalCardsMap {
			// 	if _, exists := editedCardsMap[cardID]; !exists {
			// 		// TODO DeleteCardAction
			// 	}
			// }

		} else {
			// List was deleted
			// TODO: Add DeleteListAction
		}
	}

	// Find new lists
	for listID, editedList := range editedListMap {
		if _, exists := originalListsMap[listID]; !exists {
			actions = append(actions, CreateListAction{
				BoardID:  originalBoard.ID,
				Name:     editedList.name,
				Position: editedList.markdownIdx,
			})
		}
	}

	// Delete new lists
	// for listID, originalList := range originalListMap {
	// 	if _, exists := editedListMap[listID]; !exists {
	// 		actions = append(actions, DeleteListAction{
	//
	// 		})
	// 	}
	// }

	return actions
}
