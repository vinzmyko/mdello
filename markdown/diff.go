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
					if originalCard.name != editedCard.name {
						actions = append(actions, UpdateCardNameAction{
							CardID:  originalCard.id,
							OldName: originalCard.name,
							NewName: editedCard.name,
						})
					}

					// TODO: Check for status change

					if originalCard.position != editedCard.position {
						actions = append(actions, UpdateCardPositionAction{
							CardID:      originalCard.id,
							Name:        originalCard.name,
							OldPosition: originalCard.position,
							NewPosition: editedCard.position,
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
					actions = append(actions, CreateCardAction{
						ListID:   originalList.id,
						Name:     editedCard.name,
						Position: editedCard.position,
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
