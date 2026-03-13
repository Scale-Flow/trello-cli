package main

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/credentials"
	"github.com/brettmcdowell/trello-cli/internal/trello"
)

func (m *mockAPI) ListComments(ctx context.Context, cardID string) ([]trello.Comment, error) {
	if m.listCommentsFn != nil {
		return m.listCommentsFn(ctx, cardID)
	}
	return nil, nil
}

func (m *mockAPI) AddComment(ctx context.Context, cardID, text string) (trello.Comment, error) {
	if m.addCommentFn != nil {
		return m.addCommentFn(ctx, cardID, text)
	}
	return trello.Comment{}, nil
}

func (m *mockAPI) UpdateComment(ctx context.Context, actionID, text string) (trello.Comment, error) {
	if m.updateCommentFn != nil {
		return m.updateCommentFn(ctx, actionID, text)
	}
	return trello.Comment{}, nil
}

func (m *mockAPI) DeleteComment(ctx context.Context, actionID string) error {
	if m.deleteCommentFn != nil {
		return m.deleteCommentFn(ctx, actionID)
	}
	return nil
}

func TestCommentsListCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		listCommentsFn: func(ctx context.Context, cardID string) ([]trello.Comment, error) {
			if cardID != "c1" {
				t.Errorf("card ID = %q, want c1", cardID)
			}
			return []trello.Comment{{ID: "a1", Data: trello.CommentData{Text: "hello"}}}, nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"comments", "list", "--card", "c1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("comments list failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].([]any)
	if len(data) != 1 {
		t.Fatalf("len(data) = %d, want 1", len(data))
	}
}

func TestCommentsAddCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		addCommentFn: func(ctx context.Context, cardID, text string) (trello.Comment, error) {
			if cardID != "c1" || text != "hello" {
				t.Fatalf("card/text = %q/%q", cardID, text)
			}
			return trello.Comment{ID: "a1", Data: trello.CommentData{Text: text}}, nil
		},
	}

	if err := executeRootArgs("comments", "add", "--card", "c1", "--text", "hello"); err != nil {
		t.Fatalf("comments add failed: %v", err)
	}
}

func TestCommentsUpdateCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		updateCommentFn: func(ctx context.Context, actionID, text string) (trello.Comment, error) {
			if actionID != "a1" || text != "updated" {
				t.Fatalf("action/text = %q/%q", actionID, text)
			}
			return trello.Comment{ID: actionID, Data: trello.CommentData{Text: text}}, nil
		},
	}

	if err := executeRootArgs("comments", "update", "--action", "a1", "--text", "updated"); err != nil {
		t.Fatalf("comments update failed: %v", err)
	}
}

func TestCommentsDeleteCommand(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})
	apiClient = &mockAPI{
		deleteCommentFn: func(ctx context.Context, actionID string) error {
			if actionID != "a1" {
				t.Errorf("action ID = %q, want a1", actionID)
			}
			return nil
		},
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"comments", "delete", "--action", "a1"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("comments delete failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	data := envelope["data"].(map[string]any)
	if data["deleted"] != true || data["id"] != "a1" {
		t.Fatalf("data = %+v", data)
	}
}

func TestCommentsValidation(t *testing.T) {
	setupTestAuth(t)
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	assertContractCode(t, executeRootArgs("comments", "list"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("comments", "add", "--card", "c1"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("comments", "update", "--action", "a1"), "VALIDATION_ERROR")
	assertContractCode(t, executeRootArgs("comments", "delete"), "VALIDATION_ERROR")
}
