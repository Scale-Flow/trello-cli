package trello

import (
	"context"
	"fmt"
	"strconv"
)

func (c *Client) ListLists(ctx context.Context, boardID string) ([]List, error) {
	var lists []List
	err := c.Get(ctx, fmt.Sprintf("/1/boards/%s/lists", boardID), nil, &lists)
	return lists, err
}

func (c *Client) CreateList(ctx context.Context, boardID, name string) (List, error) {
	var list List
	err := c.Post(ctx, "/1/lists", map[string]string{
		"idBoard": boardID,
		"name":    name,
	}, &list)
	return list, err
}

func (c *Client) UpdateList(ctx context.Context, listID string, params UpdateListParams) (List, error) {
	queryParams := map[string]string{}
	if params.Name != nil {
		queryParams["name"] = *params.Name
	}
	if params.Pos != nil {
		queryParams["pos"] = strconv.FormatFloat(*params.Pos, 'f', -1, 64)
	}

	var list List
	err := c.Put(ctx, fmt.Sprintf("/1/lists/%s", listID), queryParams, &list)
	return list, err
}

func (c *Client) ArchiveList(ctx context.Context, listID string) (List, error) {
	var list List
	err := c.Put(ctx, fmt.Sprintf("/1/lists/%s/closed", listID), map[string]string{
		"value": "true",
	}, &list)
	return list, err
}

func (c *Client) MoveList(ctx context.Context, listID, boardID string, pos *float64) (List, error) {
	queryParams := map[string]string{"idBoard": boardID}
	if pos != nil {
		queryParams["pos"] = strconv.FormatFloat(*pos, 'f', -1, 64)
	}

	var list List
	err := c.Put(ctx, fmt.Sprintf("/1/lists/%s", listID), queryParams, &list)
	return list, err
}
