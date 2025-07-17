package markdown

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
