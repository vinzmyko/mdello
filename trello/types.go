package trello

// From "Update a Board" endpoint
type Board struct {
	Name                     *string `json:"name,omitempty"`                      // The new name for the board. 1 to 16384 characters long.
	Desc                     *string `json:"desc,omitempty"`                      // A new description for the board, 0 to 16384 characters long
	Closed                   *bool   `json:"closed,omitempty"`                    // Whether the board is closed
	Subscribed               *string `json:"subscribed,omitempty"`                // Whether the acting user is subscribed to the board
	IdOrganisation           *string `json:"idOrganization,omitempty"`            // The id of the Workspace the board should be moved to
	PrefsPermissionLevel     *string `json:"prefs/permissionLevel,omitempty"`     // One of: org, private, public
	PrefsSelfJoin            *bool   `json:"prefs/selfJoin,omitempty"`            // Whether Workspace members can join the board themselves
	PrefsCardCovers          *bool   `json:"prefs/cardCovers,omitempty"`          // Whether card covers should be displayed on this board
	PrefsHideVotes           *bool   `json:"prefs/hideVotes,omitempty"`           // Determines whether the Voting Power-Up should hide who voted on cards or not.
	PrefsInvitations         *string `json:"prefs/invitations,omitempty"`         // Who can invite people to this board. One of: admins, members
	PrefsVoting              *string `json:"prefs/voting,omitempty"`              // Who can vote on this board. One of disabled, members, observers, org, public
	PrefsComments            *string `json:"prefs/comments,omitempty"`            // Who can comment on cards on this board. One of: disabled, members, observers, org, public
	PrefsBackground          *string `json:"prefs/background,omitempty"`          // The id of a custom background or one of: blue, orange, green, red, purple, pink, lime, sky, grey
	PrefsCardAging           *string `json:"prefs/cardAging,omitempty"`           // One of: pirate, regular
	PrefsCalendarFeedEnabled *bool   `json:"prefs/calendarFeedEnabled,omitempty"` // Determines whether the calendar feed is enabled or not.
	LabelNamesGreen          *string `json:"labelNames/green,omitempty"`          // Name for the green label. 1 to 16384 characters long
	LabelNamesYellow         *string `json:"labelNames/yellow,omitempty"`
	LabelNamesOrange         *string `json:"labelNames/orange,omitempty"`
	LabelNamesRed            *string `json:"labelNames/red,omitempty"`
	LabelNamesPurple         *string `json:"labelNames/purple,omitempty"`
	LabelNamesBlue           *string `json:"labelNames/blue,omitempty"`
}
