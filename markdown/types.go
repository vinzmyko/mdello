package markdown

import "github.com/vinzmyko/mdello/trello"

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

type DetailedTrelloAction struct {
	ObjectType ObjectType
	ObjectID   string
	ObjectName string
}

type ObjectType string

const (
	OTBoard ObjectType = "board"
	OTList  ObjectType = "list"
	OTCard  ObjectType = "card"
)
