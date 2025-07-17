// Markdown → BoardParam objects
// On error just close and tell user which lines was wrong
package markdown

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type ParsedBoard struct {
	ID    string
	Name  string
	Lists []*parsedList
}

type parsedList struct {
	id          string
	name        string
	markdownIdx int
	cards       []*parsedCard
}

type parsedCard struct {
	id         string
	name       string
	position   int
	isComplete string
	labels     []string
	dueDate    string
}

var (
	boardRegex = regexp.MustCompile(`^# (.+?)(?:\{([^}]+)\})?$`)
	listRegex  = regexp.MustCompile(`^## (.+?)(?:\{([^}]+)\})?$`)

	// TODO: Try and implement it for the cards
	cardRegex      = regexp.MustCompile(`^- \[([ xX])\] (.+)$`)
	cardLabelRegex = regexp.MustCompile(`@(\w+)`)
	cardDueRegex   = regexp.MustCompile(`due:(\S+(?:\s+\S+)?)`)
	cardIDRegex    = regexp.MustCompile(`\{([^}]+)\}`)
)

func ParseMarkdown(r io.Reader, boardSession *BoardSession) (*ParsedBoard, error) {
	scanner := bufio.NewScanner(r)
	lineNum := 0

	parsedData := &ParsedBoard{
		Lists: make([]*parsedList, 0),
	}

	var currentList *parsedList
	listPosition := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if name, id := parseWithRegex(boardRegex, line); name != "" {
			resolvedBoardID, err := boardSession.ResolveShortID(id)
			if err != nil {
				return nil, fmt.Errorf("Failed to convert board shortID back to trelloID: %w", err)
			}
			parsedData.Name = name
			parsedData.ID = resolvedBoardID
			continue // This line is a board and has been processed go next line
		}

		if name, id := parseWithRegex(listRegex, line); name != "" {
			resolvedListID, err := boardSession.ResolveShortID(id)
			if err != nil {
				return nil, fmt.Errorf("Failed to convert board shortID back to trelloID: %w", err)
			}
			newList := &parsedList{
				id:          resolvedListID,
				name:        name,
				markdownIdx: listPosition,
				cards:       make([]*parsedCard, 0),
			}
			parsedData.Lists = append(parsedData.Lists, newList)
			currentList = newList
			listPosition++
			continue
		}

		card, err := parseCard(line, boardSession)
		if err != nil {
			return nil, err
		}
		if card != nil {
			if currentList == nil {
				return nil, fmt.Errorf("line %d: found card before any list", lineNum)
			}
			card.position = len(currentList.cards)

			currentList.cards = append(currentList.cards, card)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading markdown source: %w", err)
	}

	return parsedData, nil
}

func parseCard(line string, boardSession *BoardSession) (*parsedCard, error) {
	matches := cardRegex.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil, fmt.Errorf("Missing checkbox or card text")
	}
	cardStatus := matches[1]
	cardText := matches[2]

	var cardLabels []string
	var dueDate string
	var id string

	labelMatches := cardLabelRegex.FindAllStringSubmatch(cardText, -1)
	for _, match := range labelMatches {
		cardLabels = append(cardLabels, match[1])
	}

	if dueMatch := cardDueRegex.FindStringSubmatch(cardText); len(dueMatch) > 1 {
		dueDate = dueMatch[1]
	}

	if idMatch := cardIDRegex.FindStringSubmatch(cardText); len(idMatch) > 1 {
		id = idMatch[1]
	}

	resolvedCardID, err := boardSession.ResolveShortID(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert card shortID back to trelloID: %w", err)
	}

	// Removes all matches of regexp pattern
	cleanText := cardText
	cleanText = cardLabelRegex.ReplaceAllString(cleanText, "")
	cleanText = cardDueRegex.ReplaceAllString(cleanText, "")
	cleanText = cardIDRegex.ReplaceAllString(cleanText, "")
	cleanText = strings.TrimSpace(cleanText)

	var card = &parsedCard{
		id:         resolvedCardID,
		name:       cleanText,
		isComplete: cardStatus,
		labels:     cardLabels,
		dueDate:    dueDate,
	}

	return card, nil
}

func parseWithRegex(re *regexp.Regexp, line string) (name, id string) {
	matches := re.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", ""
	}
	name = strings.TrimSpace(matches[1])
	if len(matches) > 2 && matches[2] != "" {
		id = strings.TrimSpace(matches[2])
	}
	return name, id
}

func DetectChanges(originalBoard, editedBoard *ParsedBoard) []TrelloAction {
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
					NewName: editedList.name,
				})
			}

			if originalList.markdownIdx != editedList.markdownIdx {
				// TODO: Add UpdateListPositionAction
				actions = append(actions, UpdateListPositionAction{
					ListID:      originalList.id,
					Name:        originalList.name,
					OldPosition: originalList.markdownIdx,
					NewPosition: editedList.markdownIdx,
				})
			}

			// --- CARD CHANGE DETECTION ---
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

		} else {
			// List was deleted
			// TODO: Add DeleteListAction
		}
	}

	// Find new lists
	for id, edList := range originalListsMap {
		if _, exists := originalListsMap[id]; !exists {
			actions = append(actions, CreateListAction{
				BoardID:  originalBoard.ID,
				Name:     edList.name,
				Position: edList.markdownIdx,
			})
		}
	}

	return actions
}
