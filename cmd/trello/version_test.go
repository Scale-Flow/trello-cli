package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// Set test build metadata
	version = "1.0.0-test"
	commit = "abc1234"
	date = "2026-03-12T00:00:00Z"

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(buf.Bytes(), &envelope); err != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", err, buf.String())
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	data, ok := envelope["data"].(map[string]any)
	if !ok {
		t.Fatal("data is not an object")
	}

	if data["version"] != "1.0.0-test" {
		t.Errorf("version = %v, want '1.0.0-test'", data["version"])
	}
	if data["commit"] != "abc1234" {
		t.Errorf("commit = %v, want 'abc1234'", data["commit"])
	}
	if data["date"] != "2026-03-12T00:00:00Z" {
		t.Errorf("date = %v, want '2026-03-12T00:00:00Z'", data["date"])
	}
}
