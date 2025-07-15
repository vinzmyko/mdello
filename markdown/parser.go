// Markdown → BoardParam objects
// On error just close and tell user which lines was wrong
package markdown

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/vinzmyko/mdello/trello"
)

type ParsedBoard struct {
	Board *trello.Board
	Lists []*parsedList
}

type parsedList struct {
	id       string
	name     string
	position int
	cards    []*parsedCard
}

type parsedCard struct {
	id       string
	name     string
	position int
	status   string
	labels   []string
	dueDate  string
	listId   string
}

var (
	boardRegex = regexp.MustCompile(`^# (.+?)(?:\{([^}]+)\})?$`)
	listRegex  = regexp.MustCompile(`^## (.+?)(?:\{([^}]+)\})?$`)

	cardRegex      = regexp.MustCompile(`^- \[([ xX])\] (.+)$`)
	cardLabelRegex = regexp.MustCompile(`@(\w+)`)
	cardDueRegex   = regexp.MustCompile(`due:(\S+(?:\s+\S+)?)`)
	cardIDRegex    = regexp.MustCompile(`\{([^}]+)\}`)
)

func ParseMarkdown(r io.Reader) (*ParsedBoard, error) {
	scanner := bufio.NewScanner(r)
	lineNum := 0

	parsedData := &ParsedBoard{
		Board: &trello.Board{},
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
			parsedData.Board.Name = name
			parsedData.Board.ID = id
			continue // This line is a board and has been processed go next line
		}

		if name, id := parseWithRegex(listRegex, line); name != "" {
			newList := &parsedList{
				id:       id,
				name:     name,
				position: listPosition,
				cards:    make([]*parsedCard, 0),
			}
			parsedData.Lists = append(parsedData.Lists, newList)
			currentList = newList
			listPosition++
			continue
		}

		card, err := parseCard(line)
		if err != nil {
			return nil, err
		}
		if card != nil {
			if currentList == nil {
				return nil, fmt.Errorf("line %d: found card before any list", lineNum)
			}
			card.position = len(currentList.cards)
			card.listId = currentList.id

			currentList.cards = append(currentList.cards, card)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading markdown source: %w", err)
	}

	return parsedData, nil
}

func parseCard(line string) (*parsedCard, error) {
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

	// Removes all matches of regexp pattern
	cleanText := cardText
	cleanText = cardLabelRegex.ReplaceAllString(cleanText, "")
	cleanText = cardDueRegex.ReplaceAllString(cleanText, "")
	cleanText = cardIDRegex.ReplaceAllString(cleanText, "")
	cleanText = strings.TrimSpace(cleanText)

	var card = &parsedCard{
		id:      id,
		name:    cleanText,
		status:  cardStatus,
		labels:  cardLabels,
		dueDate: dueDate,
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

func DetectChanges(originalBoard, edittedBoard *ParsedBoard) {
	if originalBoard == edittedBoard {
		fmt.Println("THEY ARE THE SAME")
	} else {
		fmt.Println("THEY ARE DIFFERENT")
	}
}
