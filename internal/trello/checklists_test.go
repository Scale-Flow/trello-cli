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

func TestListChecklists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/cards/c1/checklists" {
			t.Errorf("path = %s, want /1/cards/c1/checklists", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{"id": "cl1", "name": "Checklist", "idCard": "c1", "checkItems": []map[string]any{{"id": "i1", "name": "Item", "state": "incomplete"}}},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	checklists, err := client.ListChecklists(context.Background(), "c1")
	if err != nil {
		t.Fatalf("ListChecklists() error: %v", err)
	}
	if len(checklists) != 1 || checklists[0].ID != "cl1" {
		t.Fatalf("checklists = %+v", checklists)
	}
}

func TestCreateChecklist(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/1/checklists" {
			t.Errorf("path = %s, want /1/checklists", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "cl1", "name": "Checklist", "idCard": "c1", "checkItems": []any{},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	checklist, err := client.CreateChecklist(context.Background(), "c1", "Checklist")
	if err != nil {
		t.Fatalf("CreateChecklist() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "idCard=c1") {
		t.Errorf("query missing idCard: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "name=Checklist") {
		t.Errorf("query missing name: %s", capturedQuery)
	}
	if checklist.ID != "cl1" {
		t.Errorf("ID = %q, want cl1", checklist.ID)
	}
}

func TestDeleteChecklist(t *testing.T) {
	var capturedMethod string
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.DeleteChecklist(context.Background(), "cl1"); err != nil {
		t.Fatalf("DeleteChecklist() error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", capturedMethod)
	}
	if capturedPath != "/1/checklists/cl1" {
		t.Errorf("path = %s, want /1/checklists/cl1", capturedPath)
	}
}

func TestAddCheckItem(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/1/checklists/cl1/checkItems" {
			t.Errorf("path = %s, want /1/checklists/cl1/checkItems", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "i1", "name": "Item", "state": "incomplete",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	item, err := client.AddCheckItem(context.Background(), "cl1", "Item")
	if err != nil {
		t.Fatalf("AddCheckItem() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "name=Item") {
		t.Errorf("query missing name: %s", capturedQuery)
	}
	if item.ID != "i1" {
		t.Errorf("ID = %q, want i1", item.ID)
	}
}

func TestUpdateCheckItem(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/1/cards/c1/checkItem/i1" {
			t.Errorf("path = %s, want /1/cards/c1/checkItem/i1", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "i1", "name": "Item", "state": "complete",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	item, err := client.UpdateCheckItem(context.Background(), "c1", "i1", "complete")
	if err != nil {
		t.Fatalf("UpdateCheckItem() error: %v", err)
	}
	if !strings.Contains(capturedQuery, "state=complete") {
		t.Errorf("query missing state=complete: %s", capturedQuery)
	}
	if item.State != "complete" {
		t.Errorf("State = %q, want complete", item.State)
	}
}

func TestDeleteCheckItem(t *testing.T) {
	var capturedMethod string
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.DeleteCheckItem(context.Background(), "cl1", "i1"); err != nil {
		t.Fatalf("DeleteCheckItem() error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", capturedMethod)
	}
	if capturedPath != "/1/checklists/cl1/checkItems/i1" {
		t.Errorf("path = %s, want /1/checklists/cl1/checkItems/i1", capturedPath)
	}
}
