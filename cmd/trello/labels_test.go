package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/credentials"
	"github.com/brettmcdowell/trello-cli/internal/trello"
)

func (m *mockAPI) ListLabels(ctx context.Context, boardID string) ([]trello.Label, error) {
	if m.listLabelsFn != nil {
		return m.listLabelsFn(ctx, boardID)
	}
	return nil, nil
}

func (m *mockAPI) CreateLabel(ctx context.Context, boardID, name, color string) (trello.Label, error) {
	if m.createLabelFn != nil {
		return m.createLabelFn(ctx, boardID, name, color)
	}
	return trello.Label{}, nil
}

func (m *mockAPI) AddLabelToCard(ctx context.Context, cardID, labelID string) error {
	if m.addLabelToCardFn != nil {
		return m.addLabelToCardFn(ctx, cardID, labelID)
	}
	return nil
}

func (m *mockAPI) RemoveLabelFromCard(ctx context.Context, cardID, labelID string) error {
	if m.removeLabelFromCardFn != nil {
		return m.removeLabelFromCardFn(ctx, cardID, labelID)
	}
	return nil
}

func TestLabelsListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listLabelsFn: func(ctx context.Context, boardID string) ([]trello.Label, error) {
			return []trello.Label{{ID: "lab1", Name: "Bug"}}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"labels", "list", "--board", "b1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("labels list failed: %v", err)
	}
}

func TestLabelsCreateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		createLabelFn: func(ctx context.Context, boardID, name, color string) (trello.Label, error) {
			if boardID != "b1" || name != "Bug" || color != "red" {
				t.Fatalf("board/name/color = %q/%q/%q", boardID, name, color)
			}
			return trello.Label{ID: "lab1", Name: name, Color: color}, nil
		},
	}

	if err := executeRootArgs("labels", "create", "--board", "b1", "--name", "Bug", "--color", "red"); err != nil {
		t.Fatalf("labels create failed: %v", err)
	}
}

func TestLabelsAddCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		addLabelToCardFn: func(ctx context.Context, cardID, labelID string) error {
			if cardID != "c1" || labelID != "lab1" {
				t.Fatalf("card/label = %q/%q", cardID, labelID)
			}
			return nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"labels", "add", "--card", "c1", "--label", "lab1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("labels add failed: %v", err)
	}
	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["success"] != true || data["id"] != "lab1" {
		t.Fatalf("data = %+v", data)
	}
}

func TestLabelsRemoveCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		removeLabelFromCardFn: func(ctx context.Context, cardID, labelID string) error {
			if cardID != "c1" || labelID != "lab1" {
				t.Fatalf("card/label = %q/%q", cardID, labelID)
			}
			return nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"labels", "remove", "--card", "c1", "--label", "lab1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("labels remove failed: %v", err)
	}
	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["success"] != true || data["id"] != "lab1" {
		t.Fatalf("data = %+v", data)
	}
}
