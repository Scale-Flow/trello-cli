package trello

import (
	"context"
	"fmt"
)

func (c *Client) ListLabels(ctx context.Context, boardID string) ([]Label, error) {
	var labels []Label
	err := c.Get(ctx, fmt.Sprintf("/1/boards/%s/labels", boardID), nil, &labels)
	return labels, err
}

func (c *Client) CreateLabel(ctx context.Context, boardID, name, color string) (Label, error) {
	var label Label
	err := c.Post(ctx, "/1/labels", map[string]string{
		"idBoard": boardID,
		"name":    name,
		"color":   color,
	}, &label)
	return label, err
}

func (c *Client) AddLabelToCard(ctx context.Context, cardID, labelID string) error {
	return c.Post(ctx, fmt.Sprintf("/1/cards/%s/idLabels", cardID), map[string]string{
		"value": labelID,
	}, nil)
}

func (c *Client) RemoveLabelFromCard(ctx context.Context, cardID, labelID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/cards/%s/idLabels/%s", cardID, labelID), nil)
}
