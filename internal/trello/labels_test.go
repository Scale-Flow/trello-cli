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

func TestListLabels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/boards/b1/labels" {
			t.Errorf("path = %s, want /1/boards/b1/labels", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{"id": "lab1", "idBoard": "b1", "name": "Bug", "color": "red"},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	labels, err := client.ListLabels(context.Background(), "b1")
	if err != nil {
		t.Fatalf("ListLabels() error: %v", err)
	}
	if len(labels) != 1 || labels[0].ID != "lab1" {
		t.Fatalf("labels = %+v", labels)
	}
}

func TestCreateLabel(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/1/labels" {
			t.Errorf("path = %s, want /1/labels", r.URL.Path)
		}
		capturedQuery = r.URL.RawQuery
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id": "lab1", "idBoard": "b1", "name": "Bug", "color": "red",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	label, err := client.CreateLabel(context.Background(), "b1", "Bug", "red")
	if err != nil {
		t.Fatalf("CreateLabel() error: %v", err)
	}
	for _, want := range []string{"idBoard=b1", "name=Bug", "color=red"} {
		if !strings.Contains(capturedQuery, want) {
			t.Errorf("query missing %s: %s", want, capturedQuery)
		}
	}
	if label.ID != "lab1" {
		t.Errorf("ID = %q, want lab1", label.ID)
	}
}

func TestAddLabelToCard(t *testing.T) {
	var capturedMethod string
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedQuery = r.URL.RawQuery
		if r.URL.Path != "/1/cards/c1/idLabels" {
			t.Errorf("path = %s, want /1/cards/c1/idLabels", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.AddLabelToCard(context.Background(), "c1", "lab1"); err != nil {
		t.Fatalf("AddLabelToCard() error: %v", err)
	}
	if capturedMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", capturedMethod)
	}
	if !strings.Contains(capturedQuery, "value=lab1") {
		t.Errorf("query missing value=lab1: %s", capturedQuery)
	}
}

func TestRemoveLabelFromCard(t *testing.T) {
	var capturedMethod string
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.RemoveLabelFromCard(context.Background(), "c1", "lab1"); err != nil {
		t.Fatalf("RemoveLabelFromCard() error: %v", err)
	}
	if capturedMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", capturedMethod)
	}
	if capturedPath != "/1/cards/c1/idLabels/lab1" {
		t.Errorf("path = %s, want /1/cards/c1/idLabels/lab1", capturedPath)
	}
}
