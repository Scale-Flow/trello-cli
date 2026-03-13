package trello

import (
	"context"
	"fmt"
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
