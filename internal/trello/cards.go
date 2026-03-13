package trello

import (
	"context"
	"fmt"
	"strconv"
)

func (c *Client) ListCardsByBoard(ctx context.Context, boardID string) ([]Card, error) {
	var cards []Card
	err := c.Get(ctx, fmt.Sprintf("/1/boards/%s/cards", boardID), nil, &cards)
	return cards, err
}

func (c *Client) ListCardsByList(ctx context.Context, listID string) ([]Card, error) {
	var cards []Card
	err := c.Get(ctx, fmt.Sprintf("/1/lists/%s/cards", listID), nil, &cards)
	return cards, err
}

func (c *Client) GetCard(ctx context.Context, cardID string) (Card, error) {
	var card Card
	err := c.Get(ctx, fmt.Sprintf("/1/cards/%s", cardID), nil, &card)
	return card, err
}

func (c *Client) CreateCard(ctx context.Context, params CreateCardParams) (Card, error) {
	queryParams := map[string]string{
		"idList": params.IDList,
		"name":   params.Name,
	}
	if params.Desc != nil {
		queryParams["desc"] = *params.Desc
	}
	if params.Due != nil {
		queryParams["due"] = *params.Due
	}

	var card Card
	err := c.Post(ctx, "/1/cards", queryParams, &card)
	return card, err
}

func (c *Client) UpdateCard(ctx context.Context, cardID string, params UpdateCardParams) (Card, error) {
	queryParams := map[string]string{}
	if params.Name != nil {
		queryParams["name"] = *params.Name
	}
	if params.Desc != nil {
		queryParams["desc"] = *params.Desc
	}
	if params.Due != nil {
		queryParams["due"] = *params.Due
	}
	if params.Labels != nil {
		queryParams["idLabels"] = *params.Labels
	}
	if params.Members != nil {
		queryParams["idMembers"] = *params.Members
	}

	var card Card
	err := c.Put(ctx, fmt.Sprintf("/1/cards/%s", cardID), queryParams, &card)
	return card, err
}

func (c *Client) MoveCard(ctx context.Context, cardID, listID string, pos *float64) (Card, error) {
	queryParams := map[string]string{"idList": listID}
	if pos != nil {
		queryParams["pos"] = strconv.FormatFloat(*pos, 'f', -1, 64)
	}

	var card Card
	err := c.Put(ctx, fmt.Sprintf("/1/cards/%s", cardID), queryParams, &card)
	return card, err
}

func (c *Client) ArchiveCard(ctx context.Context, cardID string) (Card, error) {
	var card Card
	err := c.Put(ctx, fmt.Sprintf("/1/cards/%s/closed", cardID), map[string]string{
		"value": "true",
	}, &card)
	return card, err
}

func (c *Client) DeleteCard(ctx context.Context, cardID string) error {
	return c.Delete(ctx, fmt.Sprintf("/1/cards/%s", cardID), nil)
}
