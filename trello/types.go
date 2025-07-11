package trello

type CreateBoard struct {
	Name                 string  `json:"name"`                            // Required. The new name for the board. 1 to 16384 characters long
	DefaultLabels        *bool   `json:"defaultLabels,omitempty"`         // Determines whether to use the default set of labels. Default: true
	DefaultLists         *bool   `json:"defaultLists,omitempty"`          // Determines whether to add the default set of lists (To Do, Doing, Done). Ignored if idBoardSource is provided. Default: true
	Desc                 *string `json:"desc,omitempty"`                  // A new description for the board, 0 to 16384 characters long
	IdOrganisation       *string `json:"idOrganization,omitempty"`        // The id or name of the Workspace the board should belong to. Pattern: ^[0-9a-fA-F]{24}$
	IdBoardSource        *string `json:"idBoardSource,omitempty"`         // The id of a board to copy into the new board. Pattern: ^[0-9a-fA-F]{24}$
	KeepFromSource       *string `json:"keepFromSource,omitempty"`        // To keep cards from the original board pass in 'cards'. Valid values: cards, none. Default: none
	PowerUps             *string `json:"powerUps,omitempty"`              // The Power-Ups that should be enabled. Valid values: all, calendar, cardAging, recap, voting
	PrefsPermissionLevel *string `json:"prefs_permissionLevel,omitempty"` // The permissions level of the board. Valid values: org, private, public. Default: private
	PrefsVoting          *string `json:"prefs_voting,omitempty"`          // Who can vote on this board. Valid values: disabled, members, observers, org, public. Default: disabled
	PrefsComments        *string `json:"prefs_comments,omitempty"`        // Who can comment on cards. Valid values: disabled, members, observers, org, public. Default: members
	PrefsInvitations     *string `json:"prefs_invitations,omitempty"`     // What types of members can invite users to join. Valid values: members, admins. Default: members
	PrefsSelfJoin        *bool   `json:"prefs_selfJoin,omitempty"`        // Determines whether users can join the boards themselves or must be invited. Default: true
	PrefsCardCovers      *bool   `json:"prefs_cardCovers,omitempty"`      // Determines whether card covers are enabled. Default: true
	PrefsBackground      *string `json:"prefs_background,omitempty"`      // The id of a custom background or colour. Valid values: blue, orange, green, red, purple, pink, lime, sky, grey. Default: blue
	PrefsCardAging       *string `json:"prefs_cardAging,omitempty"`       // The type of card aging. Valid values: pirate, regular. Default: regular
}

type GetBoard struct {
	ID                     string  `json:"id"`                                // Required. The ID of the board to retrieve. Must match pattern: ^[0-9a-fA-F]{24}$
	Actions                *string `json:"actions,omitempty"`                 // Nested resource for actions. Default: all
	BoardStars             *string `json:"boardStars,omitempty"`              // Valid values: mine or none. Default: none
	Cards                  *string `json:"cards,omitempty"`                   // Nested resource for cards. Default: none
	CardPluginData         *bool   `json:"card_pluginData,omitempty"`         // Include card pluginData with response when used with cards param. Default: false
	Checklists             *string `json:"checklists,omitempty"`              // Nested resource for checklists. Default: none
	CustomFields           *bool   `json:"customFields,omitempty"`            // Nested resource for custom fields. Default: false
	Fields                 *string `json:"fields,omitempty"`                  // Board fields to include. Valid: all or comma-separated list of: closed, dateLastActivity, dateLastView, desc, descData, idMemberCreator, idOrganization, invitations, invited, labelNames, memberships, name, pinned, powerUps, prefs, shortLink, shortUrl, starred, subscribed, url. Default: name,desc,descData,closed,idOrganization,pinned,url,shortUrl,prefs,labelNames
	Labels                 *string `json:"labels,omitempty"`                  // Nested resource for labels
	Lists                  *string `json:"lists,omitempty"`                   // Nested resource for lists. Default: open
	Members                *string `json:"members,omitempty"`                 // Nested resource for members. Default: none
	Memberships            *string `json:"memberships,omitempty"`             // Nested resource for memberships. Default: none
	PluginData             *bool   `json:"pluginData,omitempty"`              // Whether pluginData for this board should be returned. Default: false
	Organisation           *bool   `json:"organization,omitempty"`            // Nested resource for organisations. Default: false
	OrganisationPluginData *bool   `json:"organization_pluginData,omitempty"` // Include organisation pluginData with response when used with organisation param. Default: false
	MyPrefs                *bool   `json:"myPrefs,omitempty"`                 // Include user preferences for this board. Default: false
	Tags                   *bool   `json:"tags,omitempty"`                    // Include collections/tags that the board belongs to. Default: false
}

type UpdateBoard struct {
	ID                       string  `json:"id"`                                  // Required. The ID of the board to retrieve. Must match pattern: ^[0-9a-fA-F]{24}$
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

type CreateList struct {
	Name         string  `json:"name"`                   // Required. Name for the list
	IdBoard      string  `json:"idBoard"`                // Required. The long ID of the board the list should be created on. Pattern: ^[0-9a-fA-F]{24}$
	IdListSource *string `json:"idListSource,omitempty"` // ID of the List to copy into the new List. Pattern: ^[0-9a-fA-F]{24}$
	Pos          *string `json:"pos,omitempty"`          // Position of the list. Valid values: top, bottom, or a positive floating point number
}

type GetList struct {
	ID     string  `json:"id"`               // Required. The ID of the list
	Fields *string `json:"fields,omitempty"` // All or a comma separated list of List field names. Default: name,closed,idBoard,pos
}

type UpdateList struct {
	ID         string  `json:"id"`                   // Required. The ID of the list
	Name       *string `json:"name,omitempty"`       // New name for the list
	Closed     *bool   `json:"closed,omitempty"`     // Whether the list should be closed (archived)
	IdBoard    *string `json:"idBoard,omitempty"`    // ID of a board the list should be moved to. Pattern: ^[0-9a-fA-F]{24}$
	Pos        *string `json:"pos,omitempty"`        // New position for the list: top, bottom, or a positive floating point number
	Subscribed *bool   `json:"subscribed,omitempty"` // Whether the active member is subscribed to this list
}

type ArchiveList struct {
	ID    string `json:"id"`    // Required. The ID of the list. Pattern: ^[0-9a-fA-F]{24}$
	Value bool   `json:"value"` // Required. Set to true to close (archive) the list, false to unarchive
}
