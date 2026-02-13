package command

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAccountDirectlySuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v2/account/create" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}

		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("invalid request payload: %v", err)
		}
		if payload["data"] == nil {
			t.Fatalf("expected data payload, got %v", payload)
		}

		w.WriteHeader(http.StatusCreated)
		response := map[string]any{
			"data": map[string]any{
				"id": "acc-1",
				"attributes": map[string]any{
					"email":    "test@example.com",
					"username": "testuser",
				},
			},
		}

		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	id, email, username, err := createAccountDirectly(server.URL, "test@example.com", "testuser", "secret123")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if id != "acc-1" || email != "test@example.com" || username != "testuser" {
		t.Fatalf("unexpected response values: id=%q email=%q username=%q", id, email, username)
	}
}

func TestCreateAccountDirectlyNonCreatedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid request"}`))
	}))
	defer server.Close()

	_, _, _, err := createAccountDirectly(server.URL, "test@example.com", "testuser", "secret123")
	if err == nil {
		t.Fatal("expected error for non-created status")
	}
	if !strings.Contains(err.Error(), "API error (status 400)") {
		t.Fatalf("expected API error, got %v", err)
	}
}

func TestCreateAccountDirectlyInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("not-json"))
	}))
	defer server.Close()

	_, _, _, err := createAccountDirectly(server.URL, "test@example.com", "testuser", "secret123")
	if err == nil {
		t.Fatal("expected JSON parse error")
	}
}

func TestCreateAccountDirectlyInvalidAddress(t *testing.T) {
	_, _, _, err := createAccountDirectly(":%", "test@example.com", "testuser", "secret123")
	if err == nil {
		t.Fatal("expected request creation error")
	}
}
