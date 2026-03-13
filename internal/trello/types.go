package trello

import "encoding/json"

// Board represents a Trello board.
type Board struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Closed bool   `json:"closed"`
	URL    string `json:"url"`
}

// CreateBoardParams holds fields for board creation.
type CreateBoardParams struct {
	Name           string  `json:"name"`
	Desc           *string `json:"desc,omitempty"`
	DefaultLists   *bool   `json:"defaultLists,omitempty"`
	DefaultLabels  *bool   `json:"defaultLabels,omitempty"`
	IDOrganization *string `json:"idOrganization,omitempty"`
	IDBoardSource  *string `json:"idBoardSource,omitempty"`
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
	ID       string `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Bytes    int    `json:"bytes"`
	MimeType string `json:"mimeType"`
	Date     string `json:"date"`
	IsUpload bool   `json:"isUpload"`
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

// CustomField represents a Trello custom field.
type CustomField struct {
	ID        string              `json:"id"`
	IDModel   string              `json:"idModel"`
	ModelType string              `json:"modelType"`
	Name      string              `json:"name"`
	Type      string              `json:"type"`
	Display   CustomFieldDisplay  `json:"display"`
	Options   []CustomFieldOption `json:"options,omitempty"`
}

// CustomFieldDisplay holds render preferences for a custom field.
type CustomFieldDisplay struct {
	CardFront bool `json:"cardFront"`
	CardBack  bool `json:"cardBack"`
}

// CustomFieldOption describes one of the selectable values for list-type fields.
type CustomFieldOption struct {
	ID            string                 `json:"id,omitempty"`
	IDCustomField string                 `json:"idCustomField,omitempty"`
	Color         string                 `json:"color,omitempty"`
	Value         CustomFieldOptionValue `json:"value"`
}

// CustomFieldOptionValue contains the user-visible text for an option.
type CustomFieldOptionValue struct {
	Text string `json:"text"`
}

// CustomFieldItem represents a card's current entry for a custom field.
type CustomFieldItem struct {
	ID            string               `json:"id"`
	IDCustomField string               `json:"idCustomField"`
	IDValue       string               `json:"idValue,omitempty"`
	Value         CustomFieldItemValue `json:"value,omitempty"`
}

// CustomFieldItemValue models the various payloads Trello accepts and returns for an item.
type CustomFieldItemValue struct {
	IDValue string `json:"idValue,omitempty"`
	Text    string `json:"text,omitempty"`
	Number  string `json:"number,omitempty"`
	Date    string `json:"date,omitempty"`
	Checked string `json:"checked,omitempty"`
}

// CardCustomFieldItem is an alias for the card-specific custom field item shape.
type CardCustomFieldItem = CustomFieldItem

// CardCustomFieldItemValue mirrors CustomFieldItemValue for card interactions.
type CardCustomFieldItemValue = CustomFieldItemValue

// CreateCustomFieldParams holds request fields when creating a custom field.
type CreateCustomFieldParams struct {
	IDModel   string              `json:"idModel"`
	ModelType string              `json:"modelType"`
	Name      string              `json:"name"`
	Type      string              `json:"type"`
	Display   CustomFieldDisplay  `json:"display"`
	Options   []CustomFieldOption `json:"options,omitempty"`
}

// UpdateCustomFieldParams contains fields that can be updated on a custom field.
type UpdateCustomFieldParams struct {
	Name    *string             `json:"name,omitempty"`
	Display *CustomFieldDisplay `json:"display,omitempty"`
}

// CreateCustomFieldOptionParams contains fields for creating a new option.
type CreateCustomFieldOptionParams struct {
	Value CustomFieldOptionValue `json:"value"`
	Color string                 `json:"color,omitempty"`
}

// UpdateCustomFieldOptionParams contains updatable option fields.
type UpdateCustomFieldOptionParams struct {
	Value *CustomFieldOptionValue `json:"value,omitempty"`
	Color *string                 `json:"color,omitempty"`
}

// SetCardCustomFieldItemParams holds the payload for setting a card custom field value.
type SetCardCustomFieldItemParams struct {
	Value   CardCustomFieldItemValue `json:"-"`
	IDValue string                   `json:"-"`
}

// MarshalJSON ensures only the correct custom field payload is emitted.
func (p SetCardCustomFieldItemParams) MarshalJSON() ([]byte, error) {
	if p.IDValue != "" {
		return json.Marshal(struct {
			IDValue string `json:"idValue"`
		}{p.IDValue})
	}
	if p.Value.isIDValueOnly() {
		return json.Marshal(struct {
			IDValue string `json:"idValue"`
		}{p.Value.IDValue})
	}
	if !p.Value.isZero() {
		return json.Marshal(struct {
			Value CardCustomFieldItemValue `json:"value"`
		}{p.Value})
	}
	return []byte(`{}`), nil
}

func (v CardCustomFieldItemValue) isZero() bool {
	return v.IDValue == "" && v.Text == "" && v.Number == "" && v.Date == "" && v.Checked == ""
}

func (v CardCustomFieldItemValue) isIDValueOnly() bool {
	return v.IDValue != "" && v.Text == "" && v.Number == "" && v.Date == "" && v.Checked == ""
}
