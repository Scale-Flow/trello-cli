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

func TestListCardsByBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/boards/b1/cards" {
			t.Errorf("path = %s, want /1/boards/b1/cards", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{"id": "c1", "name": "Card One", "idBoard": "b1", "idList": "l1"},
			{"id": "c2", "name": "Card Two", "idBoard": "b1", "idList": "l2"},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	cards, err := client.ListCardsByBoard(context.Background(), "b1")
	if err != nil {
		t.Fatalf("ListCardsByBoard() error: %v", err)
	}
	if len(cards) != 2 {
		t.Fatalf("len = %d, want 2", len(cards))
	}
}

func TestListCardsByList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1/lists/l1/cards" {
			t.Errorf("path = %s, want /1/lists/l1/cards", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{"id": "c1", "name": "Card One", "idBoard": "b1", "idList": "l1"},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	cards, err := client.ListCardsByList(context.Background(), "l1")
	if err != nil {
		t.Fatalf("ListCardsByList() error: %v", err)
	}
	if len(cards) != 1 || cards[0].ID != "c1" {
		t.Fatalf("cards = %+v, want one c1", cards)
	}
}

func TestGetCard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1/cards/c1" {
			t.Errorf("path = %s, want /1/cards/c1", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "c1", "name": "Card One", "idBoard": "b1", "idList": "l1",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	card, err := client.GetCard(context.Background(), "c1")
	if err != nil {
		t.Fatalf("GetCard() error: %v", err)
	}
	if card.ID != "c1" || card.Name != "Card One" {
		t.Fatalf("card = %+v", card)
	}
}

func TestCreateCard(t *testing.T) {
	var capturedQuery string
	desc := "Details"
	due := "2026-03-12T12:00:00Z"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/1/cards" {
			t.Errorf("path = %s, want /1/cards", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "c1", "name": "Card One", "desc": desc, "due": due, "idBoard": "b1", "idList": "l1",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	card, err := client.CreateCard(context.Background(), trello.CreateCardParams{
		IDList: "l1",
		Name:   "Card One",
		Desc:   &desc,
		Due:    &due,
	})
	if err != nil {
		t.Fatalf("CreateCard() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "idList=l1") {
		t.Errorf("query missing idList: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "name=Card+One") && !strings.Contains(capturedQuery, "name=Card%20One") {
		t.Errorf("query missing name: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "desc=Details") {
		t.Errorf("query missing desc: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "due=2026-03-12T12%3A00%3A00Z") {
		t.Errorf("query missing due: %s", capturedQuery)
	}
	if card.ID != "c1" {
		t.Errorf("ID = %q, want c1", card.ID)
	}
}

func TestUpdateCard(t *testing.T) {
	var capturedQuery string
	name := "Updated"
	desc := "Changed"
	due := "2026-03-13T12:00:00Z"
	labels := "lab1,lab2"
	members := "mem1,mem2"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/1/cards/c1" {
			t.Errorf("path = %s, want /1/cards/c1", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "c1", "name": name, "desc": desc, "due": due, "idBoard": "b1", "idList": "l1",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	card, err := client.UpdateCard(context.Background(), "c1", trello.UpdateCardParams{
		Name:    &name,
		Desc:    &desc,
		Due:     &due,
		Labels:  &labels,
		Members: &members,
	})
	if err != nil {
		t.Fatalf("UpdateCard() error: %v", err)
	}
	for _, want := range []string{
		"name=Updated",
		"desc=Changed",
		"due=2026-03-13T12%3A00%3A00Z",
		"idLabels=lab1",
		"idMembers=mem1",
	} {
		if !strings.Contains(capturedQuery, want) {
			t.Errorf("query missing %s: %s", want, capturedQuery)
		}
	}
	if card.Name != name {
		t.Errorf("Name = %q, want %q", card.Name, name)
	}
}

func TestMoveCard(t *testing.T) {
	var capturedQuery string
	pos := 2.5
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1/cards/c1" {
			t.Errorf("path = %s, want /1/cards/c1", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "c1", "name": "Moved", "idBoard": "b1", "idList": "l2",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	card, err := client.MoveCard(context.Background(), "c1", "l2", &pos)
	if err != nil {
		t.Fatalf("MoveCard() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "idList=l2") {
		t.Errorf("query missing idList: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "pos=2.5") {
		t.Errorf("query missing pos: %s", capturedQuery)
	}
	if card.IDList != "l2" {
		t.Errorf("IDList = %q, want l2", card.IDList)
	}
}

func TestArchiveCard(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1/cards/c1/closed" {
			t.Errorf("path = %s, want /1/cards/c1/closed", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "c1", "name": "Archived", "closed": true, "idBoard": "b1", "idList": "l1",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	card, err := client.ArchiveCard(context.Background(), "c1")
	if err != nil {
		t.Fatalf("ArchiveCard() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "value=true") {
		t.Errorf("query missing value=true: %s", capturedQuery)
	}
	if !card.Closed {
		t.Errorf("Closed = %v, want true", card.Closed)
	}
}

func TestDeleteCard(t *testing.T) {
	var capturedMethod string
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.DeleteCard(context.Background(), "c1"); err != nil {
		t.Fatalf("DeleteCard() error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", capturedMethod)
	}
	if capturedPath != "/1/cards/c1" {
		t.Errorf("path = %s, want /1/cards/c1", capturedPath)
	}
}
