package markdown

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/vinzmyko/mdello/trello"
)

type ParsedBoard struct {
	ID           string
	Name         string
	Lists        []*ParsedList
	Labels       []*trello.Label
	DetailedEdit bool
}

type ParsedList struct {
	ID           string
	Name         string
	MarkdownIdx  int
	Cards        []*ParsedCard
	DetailedEdit bool
}

type ParsedCard struct {
	ID           string
	ListID       string
	Name         string
	Position     int
	IsComplete   string
	Labels       []string
	DueDate      string
	DetailedEdit bool
}

type DiffResult struct {
	QuickActions    []TrelloAction
	DetailedActions []DetailedTrelloAction
}

type DetailedMarkdownParser struct {
	Reader   *bufio.Scanner
	Current  *DetailedTrelloAction
	Sections []DetailedTrelloAction
}

type DetailedTrelloAction struct {
	ObjectType string
	ObjectID   string
	ObjectName string
	Fields     map[string]string
}

type ParserState int

const (
	StateHeader ParserState = iota
	StateField
	StateDescription
	StateIgnore
)

func (p *DetailedMarkdownParser) Parse() ([]DetailedTrelloAction, error) {
	state := StateIgnore
	var descriptionBuilder strings.Builder

	for p.Reader.Scan() {
		line := p.Reader.Text()
		trimmed := strings.TrimSpace(line)

		switch {
		case strings.Contains(line, "# EDITING"):
			state = StateHeader
			p.FinaliseCurrent()
			p.StartNewSection(line)

		case trimmed == "=== DESCRIPTION START ===":
			state = StateDescription
			descriptionBuilder.Reset()

		case trimmed == "=== DESCRIPTION END ===":
			if p.Current != nil {
				p.Current.Fields["Description"] = descriptionBuilder.String()
			}
			state = StateField

		case state == StateDescription:
			if descriptionBuilder.Len() > 0 {
				descriptionBuilder.WriteString("\n")
			}
			descriptionBuilder.WriteString(line)

		case state == StateField && strings.Contains(trimmed, ":"):
			p.ParseFieldLine(trimmed)

		case strings.HasPrefix(trimmed, "##"):
			state = StateField
		}
	}

	p.FinaliseCurrent()
	return p.Sections, p.Reader.Err()
}

func (p *DetailedMarkdownParser) StartNewSection(headerLine string) {
    re := regexp.MustCompile(`# EDITING (\w+): (.+?) \{([^}]+)\}`)
    matches := re.FindStringSubmatch(headerLine)

    if len(matches) >= 4 {
        p.Current = &DetailedTrelloAction{
            ObjectType: strings.ToLower(matches[1]),
            ObjectName: matches[2],
            ObjectID:   matches[3],
            Fields:     make(map[string]string),
        }
    }
}

func (p *DetailedMarkdownParser) ParseFieldLine(line string) {
	if p.Current == nil {
		return
	}

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}

	fieldName := strings.TrimSpace(parts[0])
	fieldValue := strings.TrimSpace(parts[1])

	if commentIndex := strings.Index(fieldValue, "#"); commentIndex != -1 {
		fieldValue = strings.TrimSpace(fieldValue[:commentIndex])
	}

	if fieldValue != "" {
		p.Current.Fields[fieldName] = fieldValue
	}
}

func (p *DetailedMarkdownParser) FinaliseCurrent() {
	if p.Current != nil {
		p.Sections = append(p.Sections, *p.Current)
		p.Current = nil
	}
}
