package trello

import (
	"context"
	"fmt"
	"strconv"
)

func (c *Client) ListBoards(ctx context.Context) ([]Board, error) {
	var boards []Board
	err := c.Get(ctx, "/1/members/me/boards", nil, &boards)
	return boards, err
}

func (c *Client) GetBoard(ctx context.Context, boardID string) (Board, error) {
	var board Board
	err := c.Get(ctx, fmt.Sprintf("/1/boards/%s", boardID), nil, &board)
	return board, err
}

func (c *Client) CreateBoard(ctx context.Context, params CreateBoardParams) (Board, error) {
	queryParams := map[string]string{
		"name": params.Name,
	}
	if params.Desc != nil {
		queryParams["desc"] = *params.Desc
	}
	if params.DefaultLists != nil {
		queryParams["defaultLists"] = strconv.FormatBool(*params.DefaultLists)
	}
	if params.DefaultLabels != nil {
		queryParams["defaultLabels"] = strconv.FormatBool(*params.DefaultLabels)
	}
	if params.IDOrganization != nil {
		queryParams["idOrganization"] = *params.IDOrganization
	}
	if params.IDBoardSource != nil {
		queryParams["idBoardSource"] = *params.IDBoardSource
	}

	var board Board
	err := c.Post(ctx, "/1/boards", queryParams, &board)
	return board, err
}
