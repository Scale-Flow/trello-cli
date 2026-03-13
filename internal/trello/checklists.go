package trello

import (
	"context"
	"fmt"
)

func (c *Client) ListChecklists(ctx context.Context, cardID string) ([]Checklist, error) {
	var checklists []Checklist
	err := c.Get(ctx, fmt.Sprintf("/1/cards/%s/checklists", cardID), nil, &checklists)
	return checklists, err
}

func (c *Client) CreateChecklist(ctx context.Context, cardID, name string) (Checklist, error) {
	var checklist Checklist
	err := c.Post(ctx, "/1/checklists", map[string]string{
		"idCard": cardID,
		"name":   name,
	}, &checklist)
	return checklist, err
}

func (c *Client) DeleteChecklist(ctx context.Context, checklistID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/checklists/%s", checklistID), nil)
}

func (c *Client) AddCheckItem(ctx context.Context, checklistID, name string) (CheckItem, error) {
	var item CheckItem
	err := c.Post(ctx, fmt.Sprintf("/1/checklists/%s/checkItems", checklistID), map[string]string{
		"name": name,
	}, &item)
	return item, err
}

func (c *Client) UpdateCheckItem(ctx context.Context, cardID, itemID, state string) (CheckItem, error) {
	var item CheckItem
	err := c.Put(ctx, fmt.Sprintf("/1/cards/%s/checkItem/%s", cardID, itemID), map[string]string{
		"state": state,
	}, &item)
	return item, err
}

func (c *Client) DeleteCheckItem(ctx context.Context, checklistID, itemID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/checklists/%s/checkItems/%s", checklistID, itemID), nil)
}
