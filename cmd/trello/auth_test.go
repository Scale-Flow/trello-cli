package main

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/auth"
	"github.com/brettmcdowell/trello-cli/internal/credentials"
)

func setupTestAuth(t *testing.T) {
	t.Helper()
	// Use a memory store for tests
	credStore = credentials.NewMemoryStore()
	apiClient = nil
	runAuthLogin = auth.Login
	// Reset root command output state to avoid cross-test contamination
	rootCmd.SetOut(nil)
	rootCmd.SetArgs(nil)
	if err := authSetCmd.Flags().Set("api-key", ""); err != nil {
		t.Fatalf("failed to reset api-key flag: %v", err)
	}
	if err := authSetCmd.Flags().Set("token", ""); err != nil {
		t.Fatalf("failed to reset token flag: %v", err)
	}
	if err := authSetKeyCmd.Flags().Set("api-key", ""); err != nil {
		t.Fatalf("failed to reset set-key api-key flag: %v", err)
	}
	resetFlag := func(cmdName, flagName, value string) {
		t.Helper()
		cmd, _, err := rootCmd.Find([]string{cmdName})
		if err != nil {
			return
		}
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			return
		}
		if err := cmd.Flags().Set(flagName, value); err != nil {
			t.Fatalf("failed to reset %s --%s flag: %v", cmdName, flagName, err)
		}
		flag.Changed = false
	}
	resetSubFlag := func(parentName, childName, flagName, value string) {
		t.Helper()
		parent, _, err := rootCmd.Find([]string{parentName})
		if err != nil {
			return
		}
		cmd, _, err := parent.Find([]string{childName})
		if err != nil {
			return
		}
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			return
		}
		if err := cmd.Flags().Set(flagName, value); err != nil {
			t.Fatalf("failed to reset %s %s --%s flag: %v", parentName, childName, flagName, err)
		}
		flag.Changed = false
	}
	resetPathFlag := func(path []string, flagName, value string) {
		t.Helper()
		cmd, _, err := rootCmd.Find(path)
		if err != nil {
			return
		}
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			return
		}
		if err := cmd.Flags().Set(flagName, value); err != nil {
			t.Fatalf("failed to reset %v --%s flag: %v", path, flagName, err)
		}
		flag.Changed = false
	}
	resetSubFlag("boards", "get", "board", "")
	resetSubFlag("boards", "create", "name", "")
	resetSubFlag("boards", "create", "desc", "")
	resetSubFlag("boards", "create", "default-lists", "false")
	resetSubFlag("boards", "create", "default-labels", "false")
	resetSubFlag("boards", "create", "organization", "")
	resetSubFlag("boards", "create", "source-board", "")
	resetSubFlag("lists", "list", "board", "")
	resetSubFlag("lists", "create", "board", "")
	resetSubFlag("lists", "create", "name", "")
	resetSubFlag("lists", "update", "list", "")
	resetSubFlag("lists", "update", "name", "")
	resetSubFlag("lists", "update", "pos", "0")
	resetSubFlag("lists", "archive", "list", "")
	resetSubFlag("lists", "move", "list", "")
	resetSubFlag("lists", "move", "board", "")
	resetSubFlag("lists", "move", "pos", "0")
	resetSubFlag("cards", "list", "board", "")
	resetSubFlag("cards", "list", "list", "")
	resetSubFlag("cards", "get", "card", "")
	resetSubFlag("cards", "create", "list", "")
	resetSubFlag("cards", "create", "name", "")
	resetSubFlag("cards", "create", "desc", "")
	resetSubFlag("cards", "create", "due", "")
	resetSubFlag("cards", "create", "labels", "")
	resetSubFlag("cards", "create", "members", "")
	resetSubFlag("cards", "update", "card", "")
	resetSubFlag("cards", "update", "name", "")
	resetSubFlag("cards", "update", "desc", "")
	resetSubFlag("cards", "update", "due", "")
	resetSubFlag("cards", "update", "labels", "")
	resetSubFlag("cards", "update", "members", "")
	resetSubFlag("cards", "move", "card", "")
	resetSubFlag("cards", "move", "list", "")
	resetSubFlag("cards", "move", "pos", "0")
	resetSubFlag("cards", "archive", "card", "")
	resetSubFlag("cards", "delete", "card", "")
	resetSubFlag("comments", "list", "card", "")
	resetSubFlag("comments", "add", "card", "")
	resetSubFlag("comments", "add", "text", "")
	resetSubFlag("comments", "update", "action", "")
	resetSubFlag("comments", "update", "text", "")
	resetSubFlag("comments", "delete", "action", "")
	resetSubFlag("checklists", "list", "card", "")
	resetSubFlag("checklists", "create", "card", "")
	resetSubFlag("checklists", "create", "name", "")
	resetSubFlag("checklists", "delete", "checklist", "")
	resetPathFlag([]string{"checklists", "items", "add"}, "checklist", "")
	resetPathFlag([]string{"checklists", "items", "add"}, "name", "")
	resetPathFlag([]string{"checklists", "items", "update"}, "card", "")
	resetPathFlag([]string{"checklists", "items", "update"}, "item", "")
	resetPathFlag([]string{"checklists", "items", "update"}, "state", "")
	resetPathFlag([]string{"checklists", "items", "delete"}, "checklist", "")
	resetPathFlag([]string{"checklists", "items", "delete"}, "item", "")
	resetSubFlag("attachments", "list", "card", "")
	resetSubFlag("attachments", "add-file", "card", "")
	resetSubFlag("attachments", "add-file", "path", "")
	resetSubFlag("attachments", "add-file", "name", "")
	resetSubFlag("attachments", "add-url", "card", "")
	resetSubFlag("attachments", "add-url", "url", "")
	resetSubFlag("attachments", "add-url", "name", "")
	resetSubFlag("attachments", "delete", "card", "")
	resetSubFlag("attachments", "delete", "attachment", "")
	resetSubFlag("labels", "list", "board", "")
	resetSubFlag("labels", "create", "board", "")
	resetSubFlag("labels", "create", "name", "")
	resetSubFlag("labels", "create", "color", "")
	resetSubFlag("labels", "add", "card", "")
	resetSubFlag("labels", "add", "label", "")
	resetSubFlag("labels", "remove", "card", "")
	resetSubFlag("labels", "remove", "label", "")
	resetSubFlag("members", "list", "board", "")
	resetSubFlag("members", "add", "card", "")
	resetSubFlag("members", "add", "member", "")
	resetSubFlag("members", "remove", "card", "")
	resetSubFlag("members", "remove", "member", "")
	resetSubFlag("search", "cards", "query", "")
	resetSubFlag("search", "boards", "query", "")
	_ = resetFlag
}

func TestAuthSetCommand(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "set", "--api-key", "test-key", "--token", "test-token"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth set failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data := envelope["data"].(map[string]any)
	if data["configured"] != true {
		t.Errorf("configured = %v, want true", data["configured"])
	}
	if data["authMode"] != "manual" {
		t.Errorf("authMode = %v, want manual", data["authMode"])
	}
}

func TestAuthSetMissingFlags(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "set", "--api-key", "test-key"})
	// Missing --token

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("auth set should fail without --token")
	}
}

func TestAuthSetKeyCommand(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "set-key", "--api-key", "test-key"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth set-key failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data := envelope["data"].(map[string]any)
	if data["configured"] != false {
		t.Errorf("configured = %v, want false", data["configured"])
	}
	if data["authMode"] != "key_only" {
		t.Errorf("authMode = %v, want key_only", data["authMode"])
	}
}

func TestAuthSetKeyMissingFlag(t *testing.T) {
	setupTestAuth(t)

	rootCmd.SetArgs([]string{"auth", "set-key"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("auth set-key should fail without --api-key")
	}
}

func TestAuthClearCommand(t *testing.T) {
	setupTestAuth(t)
	// Pre-set credentials
	credStore.Set("default", credentials.Credentials{APIKey: "k", Token: "t", AuthMode: "manual"})

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "clear"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth clear failed: %v", err)
	}

	var envelope map[string]any
	json.Unmarshal(buf.Bytes(), &envelope)

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data := envelope["data"].(map[string]any)
	if data["configured"] != false {
		t.Errorf("configured = %v, want false", data["configured"])
	}
	if data["authMode"] != nil {
		t.Errorf("authMode = %v, want nil", data["authMode"])
	}
}

func TestAuthLoginCommand(t *testing.T) {
	setupTestAuth(t)

	runAuthLogin = func(_ loginCommandContext, _ credentials.Store, _ string, _ string, _ string, _ loginBrowserOpener, _ loginOutputWriter) (auth.LoginResult, error) {
		authMode := "interactive"
		return auth.LoginResult{
			Configured: true,
			AuthMode:   &authMode,
			Member: &auth.Member{
				ID:       "member123",
				Username: "brett",
				FullName: "Brett McDowell",
			},
		}, nil
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "login"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth login failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data := envelope["data"].(map[string]any)
	if data["configured"] != true {
		t.Errorf("configured = %v, want true", data["configured"])
	}
	if data["authMode"] != "interactive" {
		t.Errorf("authMode = %v, want interactive", data["authMode"])
	}
}

func TestAuthStatusNotConfigured(t *testing.T) {
	setupTestAuth(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"auth", "status"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth status failed: %v", err)
	}

	var envelope map[string]any
	json.Unmarshal(buf.Bytes(), &envelope)

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data := envelope["data"].(map[string]any)
	if data["configured"] != false {
		t.Errorf("configured = %v, want false", data["configured"])
	}
}
