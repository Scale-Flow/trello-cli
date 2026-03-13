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

func TestListLists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/boards/b1/lists" {
			t.Errorf("path = %s, want /1/boards/b1/lists", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{"id": "l1", "name": "Todo", "closed": false, "idBoard": "b1", "pos": 1.0},
			{"id": "l2", "name": "Done", "closed": true, "idBoard": "b1", "pos": 2.0},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	lists, err := client.ListLists(context.Background(), "b1")
	if err != nil {
		t.Fatalf("ListLists() error: %v", err)
	}
	if len(lists) != 2 {
		t.Fatalf("len = %d, want 2", len(lists))
	}
	if lists[0].ID != "l1" || lists[0].Name != "Todo" {
		t.Errorf("lists[0] = %+v", lists[0])
	}
	if !lists[1].Closed {
		t.Errorf("lists[1].Closed = %v, want true", lists[1].Closed)
	}
}

func TestCreateList(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/1/lists" {
			t.Errorf("path = %s, want /1/lists", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "l1", "name": "New List", "closed": false, "idBoard": "b1", "pos": 1.0,
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	list, err := client.CreateList(context.Background(), "b1", "New List")
	if err != nil {
		t.Fatalf("CreateList() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "idBoard=b1") {
		t.Errorf("query missing idBoard: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "name=New+List") && !strings.Contains(capturedQuery, "name=New%20List") {
		t.Errorf("query missing name: %s", capturedQuery)
	}
	if list.ID != "l1" {
		t.Errorf("ID = %q, want %q", list.ID, "l1")
	}
}

func TestUpdateList(t *testing.T) {
	var capturedQuery string
	newName := "Doing"
	newPos := 3.5
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/1/lists/l1" {
			t.Errorf("path = %s, want /1/lists/l1", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "l1", "name": "Doing", "closed": false, "idBoard": "b1", "pos": 3.5,
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	list, err := client.UpdateList(context.Background(), "l1", trello.UpdateListParams{
		Name: &newName,
		Pos:  &newPos,
	})
	if err != nil {
		t.Fatalf("UpdateList() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "name=Doing") {
		t.Errorf("query missing name: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "pos=3.5") {
		t.Errorf("query missing pos: %s", capturedQuery)
	}
	if list.Pos != 3.5 {
		t.Errorf("Pos = %v, want 3.5", list.Pos)
	}
}

func TestArchiveList(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/1/lists/l1/closed" {
			t.Errorf("path = %s, want /1/lists/l1/closed", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "l1", "name": "Done", "closed": true, "idBoard": "b1", "pos": 2.0,
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	list, err := client.ArchiveList(context.Background(), "l1")
	if err != nil {
		t.Fatalf("ArchiveList() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "value=true") {
		t.Errorf("query missing value=true: %s", capturedQuery)
	}
	if !list.Closed {
		t.Errorf("Closed = %v, want true", list.Closed)
	}
}

func TestMoveList(t *testing.T) {
	var capturedQuery string
	pos := 2.5
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/1/lists/l1" {
			t.Errorf("path = %s, want /1/lists/l1", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "l1", "name": "Doing", "closed": false, "idBoard": "b2", "pos": 2.5,
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	list, err := client.MoveList(context.Background(), "l1", "b2", &pos)
	if err != nil {
		t.Fatalf("MoveList() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "idBoard=b2") {
		t.Errorf("query missing idBoard: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "pos=2.5") {
		t.Errorf("query missing pos: %s", capturedQuery)
	}
	if list.IDBoard != "b2" {
		t.Errorf("IDBoard = %q, want %q", list.IDBoard, "b2")
	}
}
