package markdown

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/vinzmyko/mdello/trello"
)

var (
	boardRegex        = regexp.MustCompile(`^# (.+?)(?:\{([^}]+)\})?(!)?$`)
	labelPatternRegex = regexp.MustCompile(`@([^:]+):([^:\s]+)(?:\s*\{([^}]+)\})?`)

	listRegex = regexp.MustCompile(`^## (.+?)(?:\{([^}]+)\})?(!)?$`)

	cardRegex      = regexp.MustCompile(`^- \[([ xX]?)\] (.+?)(!)?$`)
	cardLabelRegex = regexp.MustCompile(`@([^:\s{]+)`)
	cardDueRegex   = regexp.MustCompile(`due:([\d-]+(?:\s+[\d:]+)?)`)
	cardIDRegex    = regexp.MustCompile(`\{([^}]+)\}`)
)

func FromMarkdown(r io.Reader, boardSession *BoardSession) (*ParsedBoard, error) {
	scanner := bufio.NewScanner(r)
	lineNum := 0

	parsedBoard := &ParsedBoard{
		Lists:  make([]*ParsedList, 0),
		Labels: make([]*trello.Label, 0),
	}

	var currentList *ParsedList
	listPosition := 0
	inLabelSection := false

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if name, id, detailedEdit := parseHeadingFields(boardRegex, line); name != "" {
			resolvedBoardID, err := boardSession.ResolveShortID(id)
			if err != nil {
				return nil, fmt.Errorf("Failed to convert board shortID back to trelloID: %w", err)
			}
			parsedBoard.Name = name
			parsedBoard.ID = resolvedBoardID
			parsedBoard.DetailedEdit = detailedEdit
			inLabelSection = true
			continue // This line is a board and has been processed go next line
		}

		// Before any lists are defined
		if inLabelSection {
			if label, err := extractBoardLabels(line, boardSession); err != nil {
				return nil, fmt.Errorf("error parsing board label: %w", err)
			} else if label != nil {
				parsedBoard.Labels = append(parsedBoard.Labels, label)
				continue
			}
		}

		if name, id, detailedEdit := parseHeadingFields(listRegex, line); name != "" {
			inLabelSection = false
			resolvedListID, err := boardSession.ResolveShortID(id)
			if err != nil {
				return nil, fmt.Errorf("Failed to convert board shortID back to trelloID: %w", err)
			}
			newList := &ParsedList{
				ID:           resolvedListID,
				Name:         name,
				MarkdownIdx:  listPosition,
				Cards:        make([]*ParsedCard, 0),
				DetailedEdit: detailedEdit,
			}
			parsedBoard.Lists = append(parsedBoard.Lists, newList)
			currentList = newList
			listPosition++

			continue
		}

		// Parse only after we have a list
		if currentList != nil {
			card, err := parseCardLine(line, currentList.ID, boardSession)
			if err != nil {
				return nil, err
			}
			if card != nil {
				card.Position = len(currentList.Cards)
				currentList.Cards = append(currentList.Cards, card)
				continue
			}
		}

		if inLabelSection {
			return nil, fmt.Errorf("line %d: unexpected content in label section: %s", lineNum, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading markdown source: %w", err)
	}

	return parsedBoard, nil
}

func parseCardLine(line string, listID string, boardSession *BoardSession) (*ParsedCard, error) {
	matches := cardRegex.FindStringSubmatch(line)
	if len(matches) < 3 {
		return nil, fmt.Errorf("%s: Missing checkbox or card text", line)
	}
	cardIsCompleted := matches[1]
	cardText := matches[2]

	var cardLabels []string
	var dueDate string
	var id string

	if dueMatch := cardDueRegex.FindStringSubmatch(cardText); len(dueMatch) > 1 {
		dueDate = dueMatch[1]
	}

	if idMatch := cardIDRegex.FindStringSubmatch(cardText); len(idMatch) > 1 {
		id = idMatch[1]
	}
	var detailedEdit bool
	if len(matches) > 3 && matches[3] == "!" {
		detailedEdit = true
	}

	tempText := cardText
	tempText = cardDueRegex.ReplaceAllString(tempText, "")
	tempText = cardIDRegex.ReplaceAllString(tempText, "")

	if invalidLabelPattern := regexp.MustCompile(`@\w+\s+\w`).FindString(tempText); invalidLabelPattern != "" {
		return nil, fmt.Errorf("invalid label format: labels cannot contain spaces. Use ~ for spaces (e.g., @front~end for 'front end')")
	}

	labelMatches := cardLabelRegex.FindAllStringSubmatch(cardText, -1)
	for _, match := range labelMatches {
		cardLabels = append(cardLabels, match[1])
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

	var card = &ParsedCard{
		ID:           resolvedCardID,
		Name:         cleanText,
		ListID:       listID,
		IsComplete:   cardIsCompleted,
		Labels:       cardLabels,
		DueDate:      dueDate,
		DetailedEdit: detailedEdit,
	}

	return card, nil
}

func parseHeadingFields(re *regexp.Regexp, line string) (name, id string, detailedEdit bool) {
	matches := re.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", "", false
	}
	name = strings.TrimSpace(matches[1])
	if len(matches) > 2 && matches[2] != "" {
		id = strings.TrimSpace(matches[2])
	}
	if len(matches) > 3 && matches[3] == "!" {
		detailedEdit = true
	}
	return name, id, detailedEdit
}

func extractBoardLabels(line string, boardSession *BoardSession) (*trello.Label, error) {
	if !strings.HasPrefix(strings.TrimSpace(line), "@") {
		return nil, nil
	}

	labelMatch := labelPatternRegex.FindStringSubmatch(line)
	if len(labelMatch) < 3 {
		return nil, fmt.Errorf("invalid label format on line: %s", line)
	}

	labelName := strings.TrimSpace(labelMatch[1])
	labelColour := strings.TrimSpace(labelMatch[2])

	var labelShortID string
	if len(labelMatch) > 3 && labelMatch[3] != "" {
		labelShortID = strings.TrimSpace(labelMatch[3])
	}

	resolvedLabelID, err := boardSession.ResolveShortID(labelShortID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve label ID for %s: %w", labelName, err)
	}

	label := &trello.Label{
		ID:     resolvedLabelID,
		Name:   labelName,
		Colour: labelColour,
	}

	return label, nil
}
