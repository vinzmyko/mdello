package trello

// ========== BOARD ==========

type Board struct {
	ID                string       `json:"id"`
	Name              string       `json:"name"`
	Desc              string       `json:"desc"`
	DescData          DescData     `json:"descData"`
	Closed            bool         `json:"closed"`
	IdMemberCreator   string       `json:"idMemberCreator"`
	IdOrganisation    string       `json:"idOrganization"`
	Pinned            bool         `json:"pinned"`
	Url               string       `json:"url"`
	ShortUrl          string       `json:"shortUrl"`
	Prefs             BoardPrefs   `json:"prefs"`
	Labels            []Label      `json:"labels"`
	LabelNames        LabelNames   `json:"labelNames"`
	Limits            BoardLimits  `json:"limits"`
	Starred           bool         `json:"starred"`
	Memberships       []Membership `json:"memberships"`
	ShortLink         string       `json:"shortLink"`
	Subscribed        bool         `json:"subscribed"`
	PowerUps          []string     `json:"powerUps"`
	PremiumFeatures   []string     `json:"premiumFeatures"`
	DateLastActivity  string       `json:"dateLastActivity"`
	DateLastView      string       `json:"dateLastView"`
	IdTags            []string     `json:"idTags"`
	DatePluginDisable string       `json:"datePluginDisable"`
	CreationMethod    string       `json:"creationMethod"`
	IxUpdate          string       `json:"ixUpdate"`
	TemplateGallery   string       `json:"templateGallery"`
	EnterpriseOwned   bool         `json:"enterpriseOwned"`
}

type DescData struct {
	Emoji map[string]any `json:"emoji,omitempty"`
}

type BoardPrefs struct {
	PermissionLevel        string                 `json:"permissionLevel"`
	Invitations            string                 `json:"invitations"`
	HideVotes              bool                   `json:"hideVotes"`
	Voting                 string                 `json:"voting"`
	Comments               string                 `json:"comments"`
	SelfJoin               bool                   `json:"selfJoin"`
	CardCovers             bool                   `json:"cardCovers"`
	IsTemplate             bool                   `json:"isTemplate"`
	CardAging              string                 `json:"cardAging"`
	CalendarFeedEnabled    bool                   `json:"calendarFeedEnabled"`
	Background             string                 `json:"background"`
	BackgroundImage        string                 `json:"backgroundImage"`
	BackgroundImageScaled  []BackgroundImageScale `json:"backgroundImageScaled"`
	BackgroundTile         bool                   `json:"backgroundTile"`
	BackgroundBrightness   string                 `json:"backgroundBrightness"`
	BackgroundBottomColour string                 `json:"backgroundBottomColor"`
	BackgroundTopColour    string                 `json:"backgroundTopColor"`
	CanBePublic            bool                   `json:"canBePublic"`
	CanBeEnterprise        bool                   `json:"canBeEnterprise"`
	CanBeOrg               bool                   `json:"canBeOrg"`
	CanBePrivate           bool                   `json:"canBePrivate"`
	CanInvite              bool                   `json:"canInvite"`
}

type BackgroundImageScale struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Url    string `json:"url"`
}

type LabelNames struct {
	Green  string `json:"green"`
	Yellow string `json:"yellow"`
	Orange string `json:"orange"`
	Red    string `json:"red"`
	Purple string `json:"purple"`
	Blue   string `json:"blue"`
	Sky    string `json:"sky"`
	Lime   string `json:"lime"`
	Pink   string `json:"pink"`
	Black  string `json:"black"`
}

type BoardLimits struct {
	Attachments AttachmentLimits `json:"attachments"`
}

type AttachmentLimits struct {
	PerBoard LimitDetail `json:"perBoard"`
}

type LimitDetail struct {
	Status    string `json:"status"`
	DisableAt int    `json:"disableAt"`
	WarnAt    int    `json:"warnAt"`
}

type Membership struct {
	ID          string `json:"id"`
	IdMember    string `json:"idMember"`
	MemberType  string `json:"memberType"`
	Unconfirmed bool   `json:"unconfirmed"`
	Deactivated bool   `json:"deactivated"`
}

type CreateBoardParams struct {
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

type GetBoardParams struct {
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

type UpdateBoardParams struct {
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
}

// ========== LABEL ==========

type Label struct {
	ID      string `json:"id"`
	IDBoard string `json:"idBoard"`
	Name    string `json:"name"`
	Colour  string `json:"color"`
	Uses    int    `json:"uses"`
}

type CreateLabelParams struct {
	BoardID string `json:"idBoard"`
	Name    string `json:"name"`
	Colour  string `json:"color"`
}

type GetLabelParams struct {
	ID     string  `json:"id"`
	Fields *string `json:"fields,omitempty"`
}

type UpdateLabelParams struct {
	ID     string  `json:"id"`
	Name   *string `json:"name,omitempty"`
	Colour *string `json:"color,omitempty"`
}

// ========== LIST ==========

type List struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Closed     bool       `json:"closed"`
	Colour     *string    `json:"color"`
	IdBoard    string     `json:"idBoard"`
	Pos        float64    `json:"pos"`
	Subscribed bool       `json:"subscribed"`
	SoftLimit  *string    `json:"softLimit"`
	Type       *string    `json:"type"`
	Datasource Datasource `json:"datasource"`
}

type Datasource struct {
	Filter bool `json:"filter"`
}

type CreateListParams struct {
	Name         string  `json:"name"`                   // Required. Name for the list
	IdBoard      string  `json:"idBoard"`                // Required. The long ID of the board the list should be created on. Pattern: ^[0-9a-fA-F]{24}$
	IdListSource *string `json:"idListSource,omitempty"` // ID of the List to copy into the new List. Pattern: ^[0-9a-fA-F]{24}$
	Pos          *string `json:"pos,omitempty"`          // Position of the list. Valid values: top, bottom, or a positive floating point number
}

type GetListParams struct {
	ID     string  `json:"id"`               // Required. The ID of the list
	Fields *string `json:"fields,omitempty"` // All or a comma separated list of List field names. Default: name,closed,idBoard,pos
}

type UpdateListParams struct {
	ID         string  `json:"id"`                   // Required. The ID of the list
	Name       *string `json:"name,omitempty"`       // New name for the list
	Closed     *bool   `json:"closed,omitempty"`     // Whether the list should be closed (archived)
	IdBoard    *string `json:"idBoard,omitempty"`    // ID of a board the list should be moved to. Pattern: ^[0-9a-fA-F]{24}$
	Pos        *string `json:"pos,omitempty"`        // New position for the list: top, bottom, or a positive floating point number
	Subscribed *bool   `json:"subscribed,omitempty"` // Whether the active member is subscribed to this list
}

type ArchiveListParams struct {
	ID    string `json:"id"`    // Required. The ID of the list. Pattern: ^[0-9a-fA-F]{24}$
	Value *bool  `json:"value"` // Required. Set to true to close (archive) the list, false to unarchive
}

// ========== CARD ==========

type Card struct {
	ID                    string       `json:"id"`
	Address               *string      `json:"address"`
	Badges                CardBadges   `json:"badges"`
	CheckItemStates       []string     `json:"checkItemStates"`
	Closed                bool         `json:"closed"`
	Coordinates           *string      `json:"coordinates"`
	CreationMethod        *string      `json:"creationMethod"`
	DateLastActivity      string       `json:"dateLastActivity"`
	Desc                  string       `json:"desc"`
	DescData              CardDescData `json:"descData"`
	Due                   *string      `json:"due"`
	DueReminder           *int64       `json:"dueReminder"`
	IdBoard               string       `json:"idBoard"`
	IdChecklists          []string     `json:"idChecklists"`
	IdLabels              []string     `json:"idLabels"`
	IdList                string       `json:"idList"`
	IdMembers             []string     `json:"idMembers"`
	IdMembersVoted        []string     `json:"idMembersVoted"`
	IdShort               int          `json:"idShort"`
	Labels                []CardLabel  `json:"labels"`
	Limits                CardLimits   `json:"limits"`
	LocationName          *string      `json:"locationName"`
	ManualCoverAttachment bool         `json:"manualCoverAttachment"`
	Name                  string       `json:"name"`
	Pos                   float64      `json:"pos"`
	ShortLink             string       `json:"shortLink"`
	ShortUrl              string       `json:"shortUrl"`
	Subscribed            bool         `json:"subscribed"`
	Url                   string       `json:"url"`
	Cover                 *CardCover   `json:"cover"`
}

type CardBadges struct {
	AttachmentsByType  AttachmentsByType `json:"attachmentsByType"`
	Location           bool              `json:"location"`
	Votes              int               `json:"votes"`
	ViewingMemberVoted bool              `json:"viewingMemberVoted"`
	Subscribed         bool              `json:"subscribed"`
	Fogbugz            *string           `json:"fogbugz"`
	CheckItems         int               `json:"checkItems"`
	CheckItemsChecked  int               `json:"checkItemsChecked"`
	Comments           int               `json:"comments"`
	Attachments        int               `json:"attachments"`
	Description        bool              `json:"description"`
	Due                *string           `json:"due"`
	Start              *string           `json:"start"`
	DueComplete        bool              `json:"dueComplete"`
}

type AttachmentsByType struct {
	Trello TrelloAttachments `json:"trello"`
}

type TrelloAttachments struct {
	Board int `json:"board"`
	Card  int `json:"card"`
}

type CardDescData struct {
	Emoji map[string]any `json:"emoji"`
}

type CardLabel struct {
	ID             string `json:"id"`
	IdBoard        string `json:"idBoard"`
	IdOrganization string `json:"idOrganization"`
	Name           string `json:"name"`
	NodeId         string `json:"nodeId"`
	Color          string `json:"color"`
	Uses           int    `json:"uses"`
}

type CardLimits struct {
	Attachments AttachmentLimits `json:"attachments"`
}

type CardCover struct {
	Color                *string `json:"color"`                // Valid values: pink, yellow, lime, blue, black, orange, red, purple, sky, green
	IdUploadedBackground *bool   `json:"idUploadedBackground"` // Whether this uses an uploaded background
	Size                 *string `json:"size"`                 // Valid values: normal, full
	Brightness           *string `json:"brightness"`           // Valid values: dark, light
	IsTemplate           *bool   `json:"isTemplate"`           // Whether this is a template
	Url                  *string `json:"url,omitempty"`
	IdAttachment         *string `json:"idAttachment,omitempty"`
}

type CreateCardParams struct {
	IdList         string    `json:"idList"`                   // Required. The ID of the list the card should be created in. Pattern: ^[0-9a-fA-F]{24}$
	Name           *string   `json:"name,omitempty"`           // The name for the card
	Desc           *string   `json:"desc,omitempty"`           // The description for the card
	Pos            *string   `json:"pos,omitempty"`            // The position of the new card. Valid values: top, bottom, or a positive float
	Due            *string   `json:"due,omitempty"`            // A due date for the card. Format: date
	Start          *string   `json:"start,omitempty"`          // The start date of a card, or null. Format: date
	DueComplete    *bool     `json:"dueComplete,omitempty"`    // Whether the status of the card is complete
	IdMembers      *[]string `json:"idMembers,omitempty"`      // Comma-separated list of member IDs to add to the card
	IdLabels       *[]string `json:"idLabels,omitempty"`       // Comma-separated list of label IDs to add to the card
	UrlSource      *string   `json:"urlSource,omitempty"`      // A URL starting with http:// or https://. The URL will be attached to the card upon creation
	FileSource     *string   `json:"fileSource,omitempty"`     // Binary file source
	MimeType       *string   `json:"mimeType,omitempty"`       // The mimeType of the attachment. Max length 256
	IdCardSource   *string   `json:"idCardSource,omitempty"`   // The ID of a card to copy into the new card. Pattern: ^[0-9a-fA-F]{24}$
	KeepFromSource *string   `json:"keepFromSource,omitempty"` // Properties to copy when using idCardSource. Valid values: all, attachments, checklists, comments, customFields, due, start, labels, members, stickers. Default: all
	Address        *string   `json:"address,omitempty"`        // For use with/by the Map View
	LocationName   *string   `json:"locationName,omitempty"`   // For use with/by the Map View
	Coordinates    *string   `json:"coordinates,omitempty"`    // For use with/by the Map View. Should take the form latitude,longitude
}

type GetCardParams struct {
	ID                string  `json:"id"`                           // Required. The ID of the Card. Pattern: ^[0-9a-fA-F]{24}$
	Fields            *string `json:"fields,omitempty"`             // All or a comma-separated list of fields. Defaults: badges, checkItemStates, closed, dateLastActivity, desc, descData, due, start, idBoard, idChecklists, idLabels, idList, idMembers, idShort, idAttachmentCover, manualCoverAttachment, labels, name, pos, shortUrl, url
	Actions           *string `json:"actions,omitempty"`            // See the Actions Nested Resource
	Attachments       *string `json:"attachments,omitempty"`        // Valid values: true, false, or cover. Default: false
	AttachmentFields  *string `json:"attachment_fields,omitempty"`  // All or a comma-separated list of attachment fields. Default: all
	Members           *bool   `json:"members,omitempty"`            // Whether to return member objects for members on the card. Default: false
	MemberFields      *string `json:"member_fields,omitempty"`      // All or a comma-separated list of member fields. Defaults: avatarHash, fullName, initials, username
	MembersVoted      *bool   `json:"membersVoted,omitempty"`       // Whether to return member objects for members who voted on the card. Default: false
	MemberVotedFields *string `json:"memberVoted_fields,omitempty"` // All or a comma-separated list of member fields. Defaults: avatarHash, fullName, initials, username
	CheckItemStates   *bool   `json:"checkItemStates,omitempty"`    // Whether to return check item states. Default: false
	Checklists        *string `json:"checklists,omitempty"`         // Whether to return the checklists on the card. Valid values: all or none. Default: none
	ChecklistFields   *string `json:"checklist_fields,omitempty"`   // All or a comma-separated list of: idBoard, idCard, name, pos. Default: all
	Board             *bool   `json:"board,omitempty"`              // Whether to return the board object the card is on. Default: false
	BoardFields       *string `json:"board_fields,omitempty"`       // All or a comma-separated list of board fields. Defaults: name, desc, descData, closed, idOrganization, pinned, url, prefs
	List              *bool   `json:"list,omitempty"`               // See the Lists Nested Resource. Default: false
	PluginData        *bool   `json:"pluginData,omitempty"`         // Whether to include pluginData on the card with the response. Default: false
	Stickers          *bool   `json:"stickers,omitempty"`           // Whether to include sticker models with the response. Default: false
	StickerFields     *string `json:"sticker_fields,omitempty"`     // All or a comma-separated list of sticker fields. Default: all
	CustomFieldItems  *bool   `json:"customFieldItems,omitempty"`   // Whether to include the customFieldItems. Default: false
}

type UpdateCardParams struct {
	ID                string     `json:"id"`                          // Required. The ID of the card
	Name              *string    `json:"name,omitempty"`              // The new name for the card
	Desc              *string    `json:"desc,omitempty"`              // The new description for the card
	Closed            *bool      `json:"closed,omitempty"`            // Whether the card should be archived (closed: true)
	IdMembers         *string    `json:"idMembers,omitempty"`         // Comma-separated list of member IDs. Pattern: ^[0-9a-fA-F]{24}$
	IdAttachmentCover *string    `json:"idAttachmentCover,omitempty"` // The ID of the image attachment the card should use as its cover, or null for none. Pattern: ^[0-9a-fA-F]{24}$
	IdList            *string    `json:"idList,omitempty"`            // The ID of the list the card should be in. Pattern: ^[0-9a-fA-F]{24}$
	IdLabels          *string    `json:"idLabels,omitempty"`          // Comma-separated list of label IDs. Pattern: ^[0-9a-fA-F]{24}$
	IdBoard           *string    `json:"idBoard,omitempty"`           // The ID of the board the card should be on. Pattern: ^[0-9a-fA-F]{24}$
	Pos               *string    `json:"pos,omitempty"`               // The position of the card in its list. Valid values: top, bottom, or a positive float
	Due               *string    `json:"due,omitempty"`               // When the card is due, or null. Format: date
	Start             *string    `json:"start,omitempty"`             // The start date of a card, or null. Format: date
	DueComplete       *bool      `json:"dueComplete,omitempty"`       // Whether the status of the card is complete
	Subscribed        *bool      `json:"subscribed,omitempty"`        // Whether the member should be subscribed to the card
	Address           *string    `json:"address,omitempty"`           // For use with/by the Map View
	LocationName      *string    `json:"locationName,omitempty"`      // For use with/by the Map View
	Coordinates       *string    `json:"coordinates,omitempty"`       // For use with/by the Map View. Should be latitude,longitude
	Cover             *CardCover `json:"cover,omitempty"`             // Updates the card's cover
}

type AddCardLabelParams struct {
	ID      string  `json:"id"`
	Name    *string `json:"name,omitempty"`
	LabelID string  `json:"value"`
}

type DeleteCardLabelParams struct {
	ID      string  `json:"id"`
	Name    *string `json:"name,omitempty"`
	LabelID string  `json:"value"`
}
