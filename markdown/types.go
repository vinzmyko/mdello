package markdown

import "github.com/vinzmyko/mdello/trello"

type ParsedBoard struct {
	ID           string
	Name         string
	Lists        []*parsedList
	Labels       []*trello.Label
	DetailedEdit bool
}

type parsedList struct {
	id           string
	name         string
	markdownIdx  int
	cards        []*parsedCard
	detailedEdit bool
}

type parsedCard struct {
	id           string
	listID       string
	name         string
	position     int
	isComplete   string
	labels       []string
	dueDate      string
	detailedEdit bool
}

type DiffResult struct {
	QuickActions    []TrelloAction
	DetailedActions []detailedTrelloAction
}

type detailedTrelloAction struct {
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
