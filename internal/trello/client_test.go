package trello_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/trello"
)

func TestClientInjectsAuthParams(t *testing.T) {
	var capturedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		json.NewEncoder(w).Encode(map[string]string{"id": "test"})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "my-key", "my-token", trello.DefaultClientOptions())

	var result map[string]string
	err := client.Get(context.Background(), "/1/members/me", nil, &result)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if capturedURL == "" {
		t.Fatal("no request was captured")
	}
	if !containsParam(capturedURL, "key=my-key") {
		t.Errorf("URL missing key param: %s", capturedURL)
	}
	if !containsParam(capturedURL, "token=my-token") {
		t.Errorf("URL missing token param: %s", capturedURL)
	}
}

func TestClientGetDecodesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"id":   "board1",
			"name": "Test Board",
		})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())

	var board trello.Board
	err := client.Get(context.Background(), "/1/boards/board1", nil, &board)
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}

	if board.ID != "board1" {
		t.Errorf("ID = %q, want %q", board.ID, "board1")
	}
	if board.Name != "Test Board" {
		t.Errorf("Name = %q, want %q", board.Name, "Test Board")
	}
}

func TestClientGetWithQueryParams(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode([]any{})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())

	params := map[string]string{"filter": "open", "fields": "id,name"}
	var result []any
	client.Get(context.Background(), "/1/boards/b1/lists", params, &result)

	if !containsParam(capturedQuery, "filter=open") {
		t.Errorf("query missing filter param: %s", capturedQuery)
	}
	if !containsParam(capturedQuery, "fields=id%2Cname") {
		t.Errorf("query missing fields param: %s", capturedQuery)
	}
}

func TestClientPostSendsQueryParams(t *testing.T) {
	var capturedQuery string
	var capturedMethod string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method

		capturedQuery = r.URL.RawQuery
		json.NewEncoder(w).Encode(map[string]string{"id": "new1"})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())

	params := map[string]string{"name": "New List", "idBoard": "b1"}
	var result map[string]string
	err := client.Post(context.Background(), "/1/lists", params, &result)
	if err != nil {
		t.Fatalf("Post() returned error: %v", err)
	}

	if capturedMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", capturedMethod)
	}
	if !strings.Contains(capturedQuery, "name=New+List") && !strings.Contains(capturedQuery, "name=New%20List") {
		t.Errorf("query missing name param: %s", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "idBoard=b1") {
		t.Errorf("query missing idBoard param: %s", capturedQuery)
	}
}

func TestClientPostJSONSendsBody(t *testing.T) {
	var capturedQuery string
	var capturedMethod string
	var contentType string
	var bodyBytes []byte
	var handlerBodyErr error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedQuery = r.URL.RawQuery
		contentType = r.Header.Get("Content-Type")
		bodyBytes, handlerBodyErr = io.ReadAll(r.Body)
		json.NewEncoder(w).Encode(map[string]string{"id": "new1"})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())

	params := map[string]string{"name": "New List", "idBoard": "b1"}
	var result map[string]string
	err := client.PostJSON(context.Background(), "/1/lists", params, &result)
	if err != nil {
		t.Fatalf("PostJSON() returned error: %v", err)
	}
	if handlerBodyErr != nil {
		t.Fatalf("failed to read body in handler: %v", handlerBodyErr)
	}

	if capturedMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", capturedMethod)
	}
	if !containsParam(capturedQuery, "key=k") {
		t.Errorf("query missing key param: %s", capturedQuery)
	}
	if !containsParam(capturedQuery, "token=t") {
		t.Errorf("query missing token param: %s", capturedQuery)
	}
	if !strings.HasPrefix(contentType, "application/json") {
		t.Errorf("content-type = %s, want application/json", contentType)
	}
	var body map[string]string
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["name"] != "New List" {
		t.Errorf("name = %s, want New List", body["name"])
	}
	if body["idBoard"] != "b1" {
		t.Errorf("idBoard = %s, want b1", body["idBoard"])
	}
}

func TestClientPutJSONSendsBody(t *testing.T) {
	var capturedQuery string
	var capturedMethod string
	var contentType string
	var bodyBytes []byte
	var handlerBodyErr error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		capturedQuery = r.URL.RawQuery
		contentType = r.Header.Get("Content-Type")
		bodyBytes, handlerBodyErr = io.ReadAll(r.Body)
		json.NewEncoder(w).Encode(map[string]string{"id": "updated1"})
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())

	params := map[string]string{"name": "Updated List", "idBoard": "b1"}
	var result map[string]string
	err := client.PutJSON(context.Background(), "/1/lists/l1", params, &result)
	if err != nil {
		t.Fatalf("PutJSON() returned error: %v", err)
	}
	if handlerBodyErr != nil {
		t.Fatalf("failed to read body in handler: %v", handlerBodyErr)
	}

	if capturedMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", capturedMethod)
	}
	if !containsParam(capturedQuery, "key=k") {
		t.Errorf("query missing key param: %s", capturedQuery)
	}
	if !containsParam(capturedQuery, "token=t") {
		t.Errorf("query missing token param: %s", capturedQuery)
	}
	if !strings.HasPrefix(contentType, "application/json") {
		t.Errorf("content-type = %s, want application/json", contentType)
	}
	var body map[string]string
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["name"] != "Updated List" {
		t.Errorf("name = %s, want Updated List", body["name"])
	}
	if body["idBoard"] != "b1" {
		t.Errorf("idBoard = %s, want b1", body["idBoard"])
	}
}

func TestClientContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	client := trello.NewClient(server.URL, "k", "t", trello.DefaultClientOptions())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var result map[string]any
	err := client.Get(ctx, "/1/members/me", nil, &result)
	if err == nil {
		t.Fatal("Get() should return error when context is cancelled")
	}
}

func containsParam(query, param string) bool {
	for _, part := range strings.Split(query, "&") {
		if part == param || strings.TrimPrefix(part, "?") == param {
			return true
		}
	}
	return false
}
