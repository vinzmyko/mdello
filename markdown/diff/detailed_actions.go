package diff

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
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

		changedFields := make(map[string]any)
		oldValues := make(map[string]string)
		newValues := make(map[string]string)

		for fieldName, newValue := range editedSection.Fields {
			oldValue, existed := originalSection.Fields[fieldName]

			if !existed || oldValue != newValue {
				normalisedField := normaliseFieldNames(fieldName)

				apiFieldName, err := getAPIFieldName(editedSection.ObjectType, normalisedField)
				if err != nil {
					fmt.Printf("\nWarning: %v, skipping field %s\n", err, normalisedField)
					continue
				}

				apiValue, err := convertValueForAPI(apiFieldName, newValue)
				if err != nil {
					fmt.Printf("\nWarning: failed to convert value for %s: %v, skipping\n", apiFieldName, err)
					continue
				}

				changedFields[apiFieldName] = apiValue
				oldValues[normalisedField] = oldValue
				newValues[normalisedField] = newValue

				if normalisedField == "Description" {
					fmt.Printf("\n[%s] %s.%s updated to\n%s",
						editedSection.ObjectType, editedSection.ObjectName, normalisedField, newValue)
				} else {
					fmt.Printf("\n[%s] %s.%s changed: '%s' -> '%s'",
						editedSection.ObjectType, editedSection.ObjectName, normalisedField, oldValue, newValue)
				}

			}
		}

		if len(changedFields) > 0 {
			action, err := markdown.CreateBulkUpdateAction(editedSection, changedFields, oldValues, newValues)
			if err != nil {
				return nil, err
			}
			actions = append(actions, action)
		}
	}

	return actions, nil
}

func normaliseFieldNames(fieldName string) string {
	fieldName = strings.TrimSuffix(fieldName, ":")
	fieldName = strings.TrimSpace(fieldName)

	switch fieldName {
	case "Start Date":
		return "Start"
	case "Due Date":
		return "Due"
	default:
		return fieldName
	}
}

func getAPIFieldName(objectType, markdownFieldName string) (string, error) {
	var mapping map[string]string

	switch objectType {
	case "board":
		mapping = markdown.BoardFieldMappings
	case "card":
		mapping = markdown.CardFieldMappings
	case "list":
		mapping = markdown.ListFieldMappings
	default:
		return "", fmt.Errorf("unknown object type: %s", objectType)
	}

	apiField, exists := mapping[markdownFieldName]
	if !exists {
		return "", fmt.Errorf("unknown field %s for %s", markdownFieldName, objectType)
	}

	return apiField, nil
}

func convertValueForAPI(apiField, stringValue string) (any, error) {
	boolFields := map[string]bool{
		// Board
		"closed":                    true,
		"prefs/selfJoin":            true,
		"prefs/cardCovers":          true,
		"prefs/hideVotes":           true,
		"prefs/calendarFeedEnabled": true,

		// List
		"subscribed": true,

		// Card
		"dueComplete": true,
	}

	if boolFields[apiField] {
		if stringValue == "" {
			return nil, nil
		}
		return strconv.ParseBool(stringValue)
	}

	if apiField == "pos" {
		if stringValue == "" {
			return nil, nil
		}
		return stringValue, nil
	}

	if stringValue == "" {
		return nil, nil
	}

	return stringValue, nil
}
