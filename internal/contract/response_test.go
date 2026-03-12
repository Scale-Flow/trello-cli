package contract_test

import (
	"bytes"
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

func TestErrorEnvelope(t *testing.T) {
	result, err := contract.ErrorEnvelope(contract.NotFound, "board not found")
	if err != nil {
		t.Fatalf("ErrorEnvelope() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("ErrorEnvelope() produced invalid JSON: %v", err)
	}

	if envelope["ok"] != false {
		t.Errorf("ok = %v, want false", envelope["ok"])
	}

	errObj, exists := envelope["error"]
	if !exists {
		t.Fatal("envelope missing 'error' field")
	}

	errMap, ok := errObj.(map[string]any)
	if !ok {
		t.Fatal("error is not an object")
	}

	if errMap["code"] != "NOT_FOUND" {
		t.Errorf("error.code = %v, want NOT_FOUND", errMap["code"])
	}
	if errMap["message"] != "board not found" {
		t.Errorf("error.message = %v, want 'board not found'", errMap["message"])
	}
}

func TestErrorEnvelopeHasNoDataField(t *testing.T) {
	result, err := contract.ErrorEnvelope(contract.AuthRequired, "not logged in")
	if err != nil {
		t.Fatalf("ErrorEnvelope() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if _, exists := envelope["data"]; exists {
		t.Error("error envelope should not contain 'data' field")
	}
}

func TestErrorEnvelopeFromContractError(t *testing.T) {
	ce := &contract.ContractError{Code: contract.ValidationError, Message: "missing --board"}
	result, err := contract.ErrorFromContractError(ce)
	if err != nil {
		t.Fatalf("ErrorFromContractError() returned error: %v", err)
	}

	var envelope map[string]any
	if err := json.Unmarshal(result, &envelope); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	errMap := envelope["error"].(map[string]any)
	if errMap["code"] != "VALIDATION_ERROR" {
		t.Errorf("error.code = %v, want VALIDATION_ERROR", errMap["code"])
	}
}

func TestRenderCompact(t *testing.T) {
	data, _ := contract.Success(map[string]string{"id": "abc"})
	var buf bytes.Buffer
	err := contract.Render(&buf, data, false)
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	output := buf.String()
	// Compact: should be a single line ending with newline
	if output[len(output)-1] != '\n' {
		t.Error("Render() output should end with newline")
	}

	// Should be valid JSON
	var envelope map[string]any
	if err := json.Unmarshal([]byte(output), &envelope); err != nil {
		t.Fatalf("Render() produced invalid JSON: %v", err)
	}

	// Should not contain any indentation
	if bytes.Contains([]byte(output), []byte("  ")) {
		t.Error("compact output should not contain indentation")
	}
}

func TestRenderPretty(t *testing.T) {
	data, _ := contract.Success(map[string]string{"id": "abc"})
	var buf bytes.Buffer
	err := contract.Render(&buf, data, true)
	if err != nil {
		t.Fatalf("Render() returned error: %v", err)
	}

	output := buf.String()
	// Pretty: should contain indentation
	if !bytes.Contains([]byte(output), []byte("  ")) {
		t.Error("pretty output should contain indentation")
	}

	// Should be valid JSON
	var envelope map[string]any
	if err := json.Unmarshal([]byte(output), &envelope); err != nil {
		t.Fatalf("Render() produced invalid JSON: %v", err)
	}

	// Should end with newline
	if output[len(output)-1] != '\n' {
		t.Error("Render() output should end with newline")
	}
}
