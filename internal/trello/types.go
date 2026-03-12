package trello

// Board represents a Trello board.
type Board struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Closed bool   `json:"closed"`
	URL    string `json:"url"`
}

// List represents a Trello list.
type List struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Closed  bool    `json:"closed"`
	IDBoard string  `json:"idBoard"`
	Pos     float64 `json:"pos"`
}

// Card represents a Trello card.
type Card struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Desc    string  `json:"desc"`
	Closed  bool    `json:"closed"`
	IDBoard string  `json:"idBoard"`
	IDList  string  `json:"idList"`
	Due     *string `json:"due"`
	URL     string  `json:"url"`
}

// Comment represents a Trello comment action.
type Comment struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Date          string        `json:"date"`
	MemberCreator MemberCreator `json:"memberCreator"`
	Data          CommentData   `json:"data"`
}

// MemberCreator is the member who created an action.
type MemberCreator struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
}

// CommentData holds the text of a comment action.
type CommentData struct {
	Text string `json:"text"`
}

// Checklist represents a Trello checklist.
type Checklist struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	IDCard     string      `json:"idCard"`
	CheckItems []CheckItem `json:"checkItems"`
}

// CheckItem represents an item in a checklist.
type CheckItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

// Attachment represents a Trello card attachment.
type Attachment struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	URL      string  `json:"url"`
	Bytes    int     `json:"bytes"`
	MimeType string  `json:"mimeType"`
	Date     string  `json:"date"`
	IsUpload bool    `json:"isUpload"`
}

// Label represents a Trello label.
type Label struct {
	ID      string `json:"id"`
	IDBoard string `json:"idBoard"`
	Name    string `json:"name"`
	Color   string `json:"color"`
}

// Member represents a Trello member.
type Member struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
}

// CardSearchResult wraps search results for cards.
type CardSearchResult struct {
	Query string `json:"query"`
	Cards []Card `json:"cards"`
}

// BoardSearchResult wraps search results for boards.
type BoardSearchResult struct {
	Query  string  `json:"query"`
	Boards []Board `json:"boards"`
}

// UpdateListParams holds optional fields for list updates.
type UpdateListParams struct {
	Name *string  `json:"name,omitempty"`
	Pos  *float64 `json:"pos,omitempty"`
}

// CreateCardParams holds fields for card creation.
type CreateCardParams struct {
	IDList string  `json:"idList"`
	Name   string  `json:"name"`
	Desc   *string `json:"desc,omitempty"`
	Due    *string `json:"due,omitempty"`
}

// UpdateCardParams holds optional fields for card updates.
type UpdateCardParams struct {
	Name    *string `json:"name,omitempty"`
	Desc    *string `json:"desc,omitempty"`
	Due     *string `json:"due,omitempty"`
	Labels  *string `json:"idLabels,omitempty"`
	Members *string `json:"idMembers,omitempty"`
}

// DeleteResult is the response shape for delete operations.
type DeleteResult struct {
	Deleted bool   `json:"deleted"`
	ID      string `json:"id"`
}

// ActionResult is the response shape for void add/remove operations
// (labels add, labels remove, members add, members remove).
type ActionResult struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
