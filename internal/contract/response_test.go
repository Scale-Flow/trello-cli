package contract_test

import (
	"encoding/json"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/contract"
)

func TestSuccessEnvelope(t *testing.T) {
	data := map[string]string{"name": "My Board"}
	result, err := contract.Success(data)
	if err != nil {
		t.Fatalf("Success() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("Success() produced invalid JSON: %v", err)
	}

	ok, exists := envelope["ok"]
	if !exists {
		t.Fatal("envelope missing 'ok' field")
	}
	if ok != true {
		t.Errorf("ok = %v, want true", ok)
	}

	d, exists := envelope["data"]
	if !exists {
		t.Fatal("envelope missing 'data' field")
	}

	dataMap, ok2 := d.(map[string]any)
	if !ok2 {
		t.Fatal("data is not an object")
	}
	if dataMap["name"] != "My Board" {
		t.Errorf("data.name = %v, want 'My Board'", dataMap["name"])
	}
}

func TestSuccessEnvelopeWithSlice(t *testing.T) {
	data := []map[string]string{{"id": "1"}, {"id": "2"}}
	result, err := contract.Success(data)
	if err != nil {
		t.Fatalf("Success() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("Success() produced invalid JSON: %v", err)
	}

	if envelope["ok"] != true {
		t.Errorf("ok = %v, want true", envelope["ok"])
	}

	arr, ok := envelope["data"].([]any)
	if !ok {
		t.Fatal("data is not an array")
	}
	if len(arr) != 2 {
		t.Errorf("data length = %d, want 2", len(arr))
	}
}

func TestSuccessEnvelopeHasNoErrorField(t *testing.T) {
	result, err := contract.Success("hello")
	if err != nil {
		t.Fatalf("Success() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if _, exists := envelope["error"]; exists {
		t.Error("success envelope should not contain 'error' field")
	}
}
