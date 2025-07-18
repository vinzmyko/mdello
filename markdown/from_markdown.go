package markdown

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	boardRegex = regexp.MustCompile(`^# (.+?)(?:\{([^}]+)\})?$`)
	listRegex  = regexp.MustCompile(`^## (.+?)(?:\{([^}]+)\})?$`)

	// TODO: Still need to do the card field editting
	cardRegex      = regexp.MustCompile(`^- \[([ xX]?)\] (.+)$`)
	cardLabelRegex = regexp.MustCompile(`@(\w+)`)
	cardDueRegex   = regexp.MustCompile(`due:(\S+(?:\s+\S+)?)`)
	cardIDRegex    = regexp.MustCompile(`\{([^}]+)\}`)
)

func FromMarkdown(r io.Reader, boardSession *BoardSession) (*ParsedBoard, error) {
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

		if name, id := extractNameAndID(boardRegex, line); name != "" {
			resolvedBoardID, err := boardSession.ResolveShortID(id)
			if err != nil {
				return nil, fmt.Errorf("Failed to convert board shortID back to trelloID: %w", err)
			}
			parsedData.Name = name
			parsedData.ID = resolvedBoardID
			continue // This line is a board and has been processed go next line
		}

		if name, id := extractNameAndID(listRegex, line); name != "" {
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

		card, err := parseCardLine(line, boardSession)
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

func parseCardLine(line string, boardSession *BoardSession) (*parsedCard, error) {
	matches := cardRegex.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil, fmt.Errorf("Missing checkbox or card text")
	}
	cardIsCompleted := matches[1]
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
		isComplete: cardIsCompleted,
		labels:     cardLabels,
		dueDate:    dueDate,
	}

	return card, nil
}

func extractNameAndID(re *regexp.Regexp, line string) (name, id string) {
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
