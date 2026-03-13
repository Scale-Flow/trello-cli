package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/credentials"
	"github.com/brettmcdowell/trello-cli/internal/trello"
)

// mockAPI implements trello.API for command testing.
type mockAPI struct {
	trello.API
	listBoardsFn func(ctx context.Context) ([]trello.Board, error)
	getBoardFn   func(ctx context.Context, id string) (trello.Board, error)
}

func (m *mockAPI) ListBoards(ctx context.Context) ([]trello.Board, error) {
	if m.listBoardsFn != nil {
		return m.listBoardsFn(ctx)
	}
	return nil, nil
}

func (m *mockAPI) GetBoard(ctx context.Context, id string) (trello.Board, error) {
	if m.getBoardFn != nil {
		return m.getBoardFn(ctx, id)
	}
	return trello.Board{}, nil
}

func TestBoardsListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listBoardsFn: func(ctx context.Context) ([]trello.Board, error) {
			return []trello.Board{
				{ID: "b1", Name: "Board One"},
				{ID: "b2", Name: "Board Two"},
			}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"boards", "list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards list failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}
	data, ok := envelope["data"].([]any)
	if !ok {
		t.Fatal("data should be an array")
	}
	if len(data) != 2 {
		t.Errorf("len(data) = %d, want 2", len(data))
	}
}

func TestBoardsListRequiresAuth(t *testing.T) {
	setupTestAuth(t)

	rootCmd.SetArgs([]string{"boards", "list"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("boards list should fail without auth")
	}
}

func TestBoardsGetCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		getBoardFn: func(ctx context.Context, id string) (trello.Board, error) {
			if id != "b1" {
				t.Errorf("board ID = %q, want %q", id, "b1")
			}
			return trello.Board{ID: "b1", Name: "My Board"}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"boards", "get", "--board", "b1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("boards get failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	data := envelope["data"].(map[string]any)
	if data["id"] != "b1" {
		t.Errorf("data.id = %v, want b1", data["id"])
	}
}

func TestBoardsGetMissingFlag(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	rootCmd.SetArgs([]string{"boards", "get"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("boards get should fail without --board")
	}
}
