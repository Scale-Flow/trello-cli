package trello_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/trello"
)

func TestListCustomFieldsByBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/boards/b1/customFields" {
			t.Errorf("path = %s, want /1/boards/b1/customFields", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{
				"id":        "cf1",
				"name":      "Phase",
				"idModel":   "b1",
				"modelType": "board",
				"type":      "list",
				"display":   map[string]any{"cardFront": true},
			},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	fields, err := client.ListCustomFieldsByBoard(context.Background(), "b1")
	if err != nil {
		t.Fatalf("ListCustomFieldsByBoard() error: %v", err)
	}
	if len(fields) != 1 {
		t.Fatalf("len = %d, want 1", len(fields))
	}
	if field := fields[0]; field.ID != "cf1" || field.Name != "Phase" {
		t.Fatalf("field[0] = %+v", field)
	} else {
		if field.ModelType != "board" {
			t.Fatalf("modelType = %q", field.ModelType)
		}
		if field.Type != "list" {
			t.Fatalf("type = %q", field.Type)
		}
		if !field.Display.CardFront {
			t.Fatalf("display.cardFront = %v", field.Display.CardFront)
		}
	}
}

func TestGetCustomField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/customFields/cf1" {
			t.Errorf("path = %s, want /1/customFields/cf1", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id":        "cf1",
			"name":      "Priority",
			"idModel":   "b1",
			"modelType": "board",
			"type":      "list",
			"display":   map[string]any{"cardFront": true},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	field, err := client.GetCustomField(context.Background(), "cf1")
	if err != nil {
		t.Fatalf("GetCustomField() error: %v", err)
	}
	if field.ID != "cf1" {
		t.Fatalf("ID = %q", field.ID)
	}
	if field.ModelType != "board" {
		t.Fatalf("modelType = %q", field.ModelType)
	}
	if field.Display.CardFront != true {
		t.Fatalf("display.cardFront = %v", field.Display.CardFront)
	}
}

func TestCreateCustomField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/1/customFields" {
			t.Errorf("path = %s, want /1/customFields", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode() error: %v", err)
		}
		if payload["idModel"] != "b1" {
			t.Fatalf("idModel = %v", payload["idModel"])
		}
		if payload["modelType"] != "board" {
			t.Fatalf("modelType = %v", payload["modelType"])
		}
		if payload["name"] != "Status" {
			t.Fatalf("name = %v", payload["name"])
		}
		if payload["type"] != "list" {
			t.Fatalf("type = %v", payload["type"])
		}
		display, _ := payload["display"].(map[string]any)
		if display == nil || display["cardFront"] != true {
			t.Fatalf("display.cardFront = %v", display["cardFront"])
		}
		options, _ := payload["options"].([]any)
		if len(options) != 1 {
			t.Fatalf("options = %+v", options)
		}
		option, _ := options[0].(map[string]any)
		value, _ := option["value"].(map[string]any)
		if value["text"] != "To Do" {
			t.Fatalf("option value = %+v", value)
		}
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id":        "cf1",
			"name":      "Status",
			"idModel":   "b1",
			"modelType": "board",
			"type":      "list",
			"display":   map[string]any{"cardFront": true},
			"options": []map[string]any{
				{
					"id":    "opt1",
					"value": map[string]any{"text": "To Do"},
					"color": "green",
				},
			},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	field, err := client.CreateCustomField(context.Background(), trello.CreateCustomFieldParams{
		IDModel:   "b1",
		ModelType: "board",
		Name:      "Status",
		Type:      "list",
		Display: trello.CustomFieldDisplay{
			CardFront: true,
		},
		Options: []trello.CustomFieldOption{
			{
				Value: trello.CustomFieldOptionValue{Text: "To Do"},
			},
		},
	})
	if err != nil {
		t.Fatalf("CreateCustomField() error: %v", err)
	}
	if field.ID != "cf1" {
		t.Fatalf("ID = %q", field.ID)
	}
	if field.IDModel != "b1" {
		t.Fatalf("idModel = %q", field.IDModel)
	}
	if !field.Display.CardFront {
		t.Fatalf("display.cardFront = %v", field.Display.CardFront)
	}
	if len(field.Options) != 1 {
		t.Fatalf("options = %+v", field.Options)
	}
	opt := field.Options[0]
	if opt.Value.Text != "To Do" {
		t.Fatalf("option value.text = %q", opt.Value.Text)
	}
	if opt.Color != "green" {
		t.Fatalf("option color = %q", opt.Color)
	}
}

func TestUpdateCustomField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/1/customFields/cf1" {
			t.Errorf("path = %s, want /1/customFields/cf1", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode() error: %v", err)
		}
		if payload["name"] != "Status Updated" {
			t.Fatalf("name = %v", payload["name"])
		}
		// Trello expects flat key "display/cardFront", not nested object
		if payload["display/cardFront"] != false {
			t.Fatalf("display/cardFront = %v", payload["display/cardFront"])
		}
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id":      "cf1",
			"name":    "Status Updated",
			"display": map[string]any{"cardFront": false},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	name := "Status Updated"
	cardFront := false
	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	field, err := client.UpdateCustomField(context.Background(), "cf1", trello.UpdateCustomFieldParams{
		Name: &name,
		Display: &trello.CustomFieldDisplay{
			CardFront: cardFront,
		},
	})
	if err != nil {
		t.Fatalf("UpdateCustomField() error: %v", err)
	}
	if field.Name != "Status Updated" {
		t.Fatalf("Name = %q", field.Name)
	}
	if field.Display.CardFront {
		t.Fatalf("display.cardFront = %v", field.Display.CardFront)
	}
}

func TestDeleteCustomField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/1/customFields/cf1" {
			t.Errorf("path = %s, want /1/customFields/cf1", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.DeleteCustomField(context.Background(), "cf1"); err != nil {
		t.Fatalf("DeleteCustomField() error: %v", err)
	}
}

func TestListCustomFieldOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/customFields/cf1/options" {
			t.Errorf("path = %s, want /1/customFields/cf1/options", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{
				"id":    "opt1",
				"value": map[string]any{"text": "One"},
				"color": "blue",
			},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	options, err := client.ListCustomFieldOptions(context.Background(), "cf1")
	if err != nil {
		t.Fatalf("ListCustomFieldOptions() error: %v", err)
	}
	if len(options) != 1 {
		t.Fatalf("len = %d, want 1", len(options))
	}
	opt := options[0]
	if opt.ID != "opt1" {
		t.Fatalf("option id = %q", opt.ID)
	}
	if opt.Value.Text != "One" {
		t.Fatalf("value.text = %q", opt.Value.Text)
	}
	if opt.Color != "blue" {
		t.Fatalf("color = %q", opt.Color)
	}
}

func TestCreateCustomFieldOption(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/1/customFields/cf1/options" {
			t.Errorf("path = %s, want /1/customFields/cf1/options", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode() error: %v", err)
		}
		if payload["color"] != "green" {
			t.Fatalf("color = %v", payload["color"])
		}
		value, _ := payload["value"].(map[string]any)
		if value["text"] != "Go" {
			t.Fatalf("value = %+v", value)
		}
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id":    "opt1",
			"color": "green",
			"value": map[string]any{"text": "Go"},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	option, err := client.CreateCustomFieldOption(context.Background(), "cf1", trello.CreateCustomFieldOptionParams{
		Value: trello.CustomFieldOptionValue{Text: "Go"},
		Color: "green",
	})
	if err != nil {
		t.Fatalf("CreateCustomFieldOption() error: %v", err)
	}
	if option.ID != "opt1" {
		t.Fatalf("ID = %q", option.ID)
	}
	if option.Color != "green" {
		t.Fatalf("color = %q", option.Color)
	}
	if option.Value.Text != "Go" {
		t.Fatalf("value.text = %q", option.Value.Text)
	}
}

func TestUpdateCustomFieldOption(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/1/customFields/cf1/options/opt1" {
			t.Errorf("path = %s, want /1/customFields/cf1/options/opt1", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode() error: %v", err)
		}
		value, _ := payload["value"].(map[string]any)
		if value["text"] != "Done" {
			t.Fatalf("value = %+v", value)
		}
		if err := json.NewEncoder(w).Encode(map[string]any{
			"id":    "opt1",
			"value": value,
			"color": "yellow",
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	valueText := "Done"
	option, err := client.UpdateCustomFieldOption(context.Background(), "cf1", "opt1", trello.UpdateCustomFieldOptionParams{
		Value: &trello.CustomFieldOptionValue{Text: valueText},
	})
	if err != nil {
		t.Fatalf("UpdateCustomFieldOption() error: %v", err)
	}
	if option.Value.Text != valueText {
		t.Fatalf("value = %+v", option.Value)
	}
	if option.Color != "yellow" {
		t.Fatalf("color = %q", option.Color)
	}
}

func TestDeleteCustomFieldOption(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		if r.URL.Path != "/1/customFields/cf1/options/opt1" {
			t.Errorf("path = %s, want /1/customFields/cf1/options/opt1", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.DeleteCustomFieldOption(context.Background(), "cf1", "opt1"); err != nil {
		t.Fatalf("DeleteCustomFieldOption() error: %v", err)
	}
}

func TestListCardCustomFieldItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/1/cards/c1/customFieldItems" {
			t.Errorf("path = %s, want /1/cards/c1/customFieldItems", r.URL.Path)
		}
		if err := json.NewEncoder(w).Encode([]map[string]any{
			{
				"id":            "item1",
				"idCustomField": "cf1",
				"idValue":       "opt1",
			},
			{
				"id":            "item2",
				"idCustomField": "cf2",
				"value": map[string]any{
					"text": "Follow-up",
				},
			},
		}); err != nil {
			t.Fatalf("Encode() error: %v", err)
		}
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	items, err := client.ListCardCustomFieldItems(context.Background(), "c1")
	if err != nil {
		t.Fatalf("ListCardCustomFieldItems() error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("len = %d, want 2", len(items))
	}
	listItem := items[0]
	if listItem.ID != "item1" {
		t.Fatalf("list item id = %q", listItem.ID)
	}
	if listItem.IDValue != "opt1" {
		t.Fatalf("list item idValue = %q", listItem.IDValue)
	}
	textItem := items[1]
	if textItem.ID != "item2" {
		t.Fatalf("text item id = %q", textItem.ID)
	}
	if textItem.Value.Text != "Follow-up" {
		t.Fatalf("text item value = %q", textItem.Value.Text)
	}
}

func TestSetCardCustomFieldItem(t *testing.T) {
	cases := []struct {
		name            string
		params          trello.SetCardCustomFieldItemParams
		assertRequest   func(t *testing.T, payload map[string]any)
		responseIDValue string
		responseValue   map[string]any
		assertResponse  func(t *testing.T, item trello.CardCustomFieldItem)
	}{
		{
			name: "list option value",
			params: trello.SetCardCustomFieldItemParams{
				Value: trello.CardCustomFieldItemValue{IDValue: "opt1"},
			},
			assertRequest: func(t *testing.T, payload map[string]any) {
				if idValue, _ := payload["idValue"].(string); idValue != "opt1" {
					t.Fatalf("idValue = %v", payload)
				}
				if _, exists := payload["value"]; exists {
					t.Fatalf("value should not be present for list option")
				}
			},
			responseIDValue: "opt1",
			assertResponse: func(t *testing.T, item trello.CardCustomFieldItem) {
				if item.Value.IDValue != "opt1" {
					t.Fatalf("value.idValue = %q", item.Value.IDValue)
				}
				if item.IDValue != "opt1" {
					t.Fatalf("idValue = %q", item.IDValue)
				}
			},
		},
		{
			name: "text value",
			params: trello.SetCardCustomFieldItemParams{
				Value: trello.CardCustomFieldItemValue{Text: "Need Info"},
			},
			assertRequest: func(t *testing.T, payload map[string]any) {
				value, ok := payload["value"].(map[string]any)
				if !ok {
					t.Fatalf("value missing: %+v", payload)
				}
				if text, _ := value["text"].(string); text != "Need Info" {
					t.Fatalf("text = %v", value)
				}
				if _, exists := payload["idValue"]; exists {
					t.Fatalf("idValue should be absent for text value")
				}
			},
			responseValue: map[string]any{"text": "Need Info"},
			assertResponse: func(t *testing.T, item trello.CardCustomFieldItem) {
				if item.Value.Text != "Need Info" {
					t.Fatalf("value.text = %q", item.Value.Text)
				}
				if item.Value.IDValue != "" {
					t.Fatalf("idValue should be empty for text value")
				}
			},
		},
		{
			name: "number/date/checked value",
			params: trello.SetCardCustomFieldItemParams{
				Value: trello.CardCustomFieldItemValue{
					Number:  "42",
					Date:    "2026-03-13T00:00:00Z",
					Checked: "true",
				},
			},
			assertRequest: func(t *testing.T, payload map[string]any) {
				value, ok := payload["value"].(map[string]any)
				if !ok {
					t.Fatalf("value missing: %+v", payload)
				}
				if payload["idValue"] != nil {
					t.Fatalf("idValue should be absent for non-list value")
				}
				if number, _ := value["number"].(string); number != "42" {
					t.Fatalf("number = %v", number)
				}
				if date, _ := value["date"].(string); date != "2026-03-13T00:00:00Z" {
					t.Fatalf("date = %v", date)
				}
				if checked, _ := value["checked"].(string); checked != "true" {
					t.Fatalf("checked = %v", checked)
				}
			},
			responseValue: map[string]any{
				"number":  "42",
				"date":    "2026-03-13T00:00:00Z",
				"checked": "true",
			},
			assertResponse: func(t *testing.T, item trello.CardCustomFieldItem) {
				if item.Value.Number != "42" {
					t.Fatalf("value.number = %q", item.Value.Number)
				}
				if item.Value.Date != "2026-03-13T00:00:00Z" {
					t.Fatalf("value.date = %q", item.Value.Date)
				}
				if item.Value.Checked != "true" {
					t.Fatalf("value.checked = %q", item.Value.Checked)
				}
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("method = %s, want PUT", r.Method)
				}
				if r.URL.Path != "/1/cards/c1/customField/cf1/item" {
					t.Errorf("path = %s, want /1/cards/c1/customField/cf1/item", r.URL.Path)
				}
				var payload map[string]any
				if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
					t.Fatalf("Decode() error: %v", err)
				}
				tc.assertRequest(t, payload)
				resp := map[string]any{
					"id":            "item1",
					"idCustomField": "cf1",
				}
				if tc.responseIDValue != "" {
					resp["idValue"] = tc.responseIDValue
				} else {
					resp["value"] = tc.responseValue
				}
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					t.Fatalf("Encode() error: %v", err)
				}
			}))
			defer server.Close()

			client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
			item, err := client.SetCardCustomFieldItem(context.Background(), "c1", "cf1", tc.params)
			if err != nil {
				t.Fatalf("SetCardCustomFieldItem() error: %v", err)
			}
			if item.ID != "item1" {
				t.Fatalf("ID = %q", item.ID)
			}
			if item.IDCustomField != "cf1" {
				t.Fatalf("idCustomField = %q", item.IDCustomField)
			}
			tc.assertResponse(t, item)
		})
	}
}

func TestClearCardCustomFieldItem(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/1/cards/c1/customField/cf1/item" {
			t.Errorf("path = %s, want /1/cards/c1/customField/cf1/item", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("Decode() error: %v", err)
		}
		if payload["value"] != "" {
			t.Fatalf("value = %v, want empty string", payload["value"])
		}
		if payload["idValue"] != "" {
			t.Fatalf("idValue = %v, want empty string", payload["idValue"])
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())
	if err := client.ClearCardCustomFieldItem(context.Background(), "c1", "cf1"); err != nil {
		t.Fatalf("ClearCardCustomFieldItem() error: %v", err)
	}
}
