package diff

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/vinzmyko/mdello/config"
	"github.com/vinzmyko/mdello/markdown"
)

func ParseDetailedMarkdown(r io.Reader) ([]markdown.DetailedTrelloAction, error) {
	parser := &markdown.DetailedMarkdownParser{
		Reader:   bufio.NewScanner(r),
		Sections: make([]markdown.DetailedTrelloAction, 0),
	}

	return parser.Parse()
}

func DetailedActionsDiff(originalContent, editedContent string, cfg *config.Config) ([]markdown.TrelloAction, error) {
	originalSections, err := ParseDetailedMarkdown(strings.NewReader(originalContent))
	if err != nil {
		return nil, fmt.Errorf("parsing original content: %w", err)
	}

	editedSections, err := ParseDetailedMarkdown(strings.NewReader(editedContent))
	if err != nil {
		return nil, fmt.Errorf("parsing edited content: %w", err)
	}

	return generateDetailedActions(originalSections, editedSections)
}

func generateDetailedActions(original, edited []markdown.DetailedTrelloAction) ([]markdown.TrelloAction, error) {
	var actions []markdown.TrelloAction

	originalMap := make(map[string]markdown.DetailedTrelloAction)
	for _, section := range original {
		originalMap[section.ObjectName] = section
	}

	for _, editedSection := range edited {
		originalSection := originalMap[editedSection.ObjectName]

		for fieldName, newValue := range editedSection.Fields {
			oldValue, existed := originalSection.Fields[fieldName]

			if !existed {
				fmt.Printf("\n[%s] %s.%s added: '%s'",
					editedSection.ObjectType, editedSection.ObjectName, fieldName, newValue)
			} else if oldValue != newValue {
				fmt.Printf("\n[%s] %s.%s changed: '%s' -> '%s'",
					editedSection.ObjectType, editedSection.ObjectName, fieldName, oldValue, newValue)
			}
		}
	}

	return actions, nil
}
