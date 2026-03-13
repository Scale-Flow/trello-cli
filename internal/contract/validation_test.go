package contract_test

import (
    "testing"

    "github.com/brettmcdowell/trello-cli/internal/contract"
)

func TestRequireExactlyOneWithOne(t *testing.T) {
    err := contract.RequireExactlyOne(map[string]string{
        "board": "abc",
        "list":  "",
    })
    if err != nil {
        t.Errorf("RequireExactlyOne() with one value should return nil, got %v", err)
    }
}

func TestRequireExactlyOneWithNone(t *testing.T) {
    err := contract.RequireExactlyOne(map[string]string{
        "board": "",
        "list":  "",
    })
    if err == nil {
        t.Fatal("RequireExactlyOne() with no values should return error")
    }
    ce := err.(*contract.ContractError)
    if ce.Code != contract.ValidationError {
        t.Errorf("Code = %q, want %q", ce.Code, contract.ValidationError)
    }
}

func TestRequireExactlyOneWithTwo(t *testing.T) {
    err := contract.RequireExactlyOne(map[string]string{
        "board": "abc",
        "list":  "def",
    })
    if err == nil {
        t.Fatal("RequireExactlyOne() with two values should return error")
    }
    ce := err.(*contract.ContractError)
    if ce.Code != contract.ValidationError {
        t.Errorf("Code = %q, want %q", ce.Code, contract.ValidationError)
    }
}

func TestRequireAtLeastOneWithOne(t *testing.T) {
    err := contract.RequireAtLeastOne(map[string]string{
        "name": "new name",
        "pos":  "",
    })
    if err != nil {
        t.Errorf("RequireAtLeastOne() with one value should return nil, got %v", err)
    }
}

func TestRequireAtLeastOneWithNone(t *testing.T) {
    err := contract.RequireAtLeastOne(map[string]string{
        "name": "",
        "pos":  "",
    })
    if err == nil {
        t.Fatal("RequireAtLeastOne() with no values should return error")
    }
    ce := err.(*contract.ContractError)
    if ce.Code != contract.ValidationError {
        t.Errorf("Code = %q, want %q", ce.Code, contract.ValidationError)
    }
}

func TestRequireAtLeastOneWithAll(t *testing.T) {
    err := contract.RequireAtLeastOne(map[string]string{
        "name": "new name",
        "pos":  "top",
    })
    if err != nil {
        t.Errorf("RequireAtLeastOne() with all values should return nil, got %v", err)
    }
}
