package markdown

// User-friendly name -> Trello API parameter
var BoardFieldMappings = map[string]string{
	"Name":             "name",
	"Closed":           "closed",
	"Description":      "desc",
	"Permission Level": "prefs/permissionLevel",
	"Self Join":        "prefs/selfJoin",
	"Card Covers":      "prefs/cardCovers",
	"Hide Votes":       "prefs/hideVotes",
	"Invitations":      "prefs/invitations",
	"Voting":           "prefs/voting",
	"Comments":         "prefs/comments",
	"Card Aging":       "prefs/cardAging",
	"Calendar Feed":    "prefs/calendarFeedEnabled",
}

var ListFieldMappings = map[string]string{
	"Name":       "name",
	"Closed":     "closed",
	"Position":   "pos",
	"Subscribed": "subscribed",
}

var CardFieldMappings = map[string]string{
	"Name":       "name",
	"Closed":     "closed",
	"Position":   "pos",
	"Subscribed": "subscribed",

	"Description": "desc",

	"Start":        "start",
	"Due":          "due",
	"Due Complete": "dueComplete",
}
