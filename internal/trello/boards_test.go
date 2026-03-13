package trello_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/trello"
)

func TestListBoards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/members/me/boards" {
			t.Errorf("path = %s, want /1/members/me/boards", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{"id": "b1", "name": "Board One", "desc": "First", "closed": false, "url": "https://trello.com/b/b1"},
			{"id": "b2", "name": "Board Two", "desc": "", "closed": true, "url": "https://trello.com/b/b2"},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	boards, err := client.ListBoards(context.Background())
	if err != nil {
		t.Fatalf("ListBoards() error: %v", err)
	}
	if len(boards) != 2 {
		t.Fatalf("len = %d, want 2", len(boards))
	}
	if boards[0].ID != "b1" || boards[0].Name != "Board One" {
		t.Errorf("boards[0] = %+v", boards[0])
	}
	if boards[1].Closed != true {
		t.Errorf("boards[1].Closed = %v, want true", boards[1].Closed)
	}
}

func TestGetBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1/boards/b1" {
			t.Errorf("path = %s, want /1/boards/b1", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "b1", "name": "My Board", "desc": "A board", "closed": false, "url": "https://trello.com/b/b1",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	board, err := client.GetBoard(context.Background(), "b1")
	if err != nil {
		t.Fatalf("GetBoard() error: %v", err)
	}
	if board.ID != "b1" {
		t.Errorf("ID = %q, want %q", board.ID, "b1")
	}
	if board.Name != "My Board" {
		t.Errorf("Name = %q, want %q", board.Name, "My Board")
	}
}

func TestCreateBoard(t *testing.T) {
	var capturedQuery string
	desc := "Board description"
	defaultLists := true
	defaultLabels := true
	orgID := "org1"
	sourceBoardID := "template1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/1/boards" {
			t.Errorf("path = %s, want /1/boards", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id":     "b3",
			"name":   "Project Board",
			"desc":   desc,
			"closed": false,
			"url":    "https://trello.com/b/b3",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	board, err := client.CreateBoard(context.Background(), trello.CreateBoardParams{
		Name:           "Project Board",
		Desc:           &desc,
		DefaultLists:   &defaultLists,
		DefaultLabels:  &defaultLabels,
		IDOrganization: &orgID,
		IDBoardSource:  &sourceBoardID,
	})
	if err != nil {
		t.Fatalf("CreateBoard() error: %v", err)
	}
	for _, want := range []string{
		"name=Project+Board",
		"desc=Board+description",
		"defaultLists=true",
		"defaultLabels=true",
		"idOrganization=org1",
		"idBoardSource=template1",
	} {
		if !strings.Contains(capturedQuery, want) {
			t.Errorf("query missing %s: %s", want, capturedQuery)
		}
	}
	if board.ID != "b3" || board.Name != "Project Board" {
		t.Errorf("board = %+v", board)
	}
}
