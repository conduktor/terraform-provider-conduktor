package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetConsoleLicensePlan_NoDoubleApiPrefix(t *testing.T) {
	orgsHit := false
	licenseHit := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/organizations":
			orgsHit = true
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode([]map[string]any{{"slug": "my-org"}})
		case "/api/organizations/my-org/platform-license":
			licenseHit = true
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"plan": "enterprise"})
		default:
			// Return HTML like Console SPA does for unmatched routes — this is what triggers the bug
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte("<html>not found</html>"))
		}
	}))
	defer ts.Close()

	// Use Make so uniformizeBaseUrl is in the callstack — it appends /api to the base URL,
	// which is the real production path. Passing ts.URL (no /api) verifies the fix end-to-end.
	c, err := Make(context.Background(), CONSOLE, ApiParameter{
		BaseUrl: ts.URL,
		ApiKey:  "test-key",
	}, "test")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	plan, err := c.GetConsoleLicensePlan(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plan != "enterprise" {
		t.Errorf("expected plan 'enterprise', got %q", plan)
	}
	if !orgsHit {
		t.Error("GET /api/organizations was never called — double /api prefix likely")
	}
	if !licenseHit {
		t.Error("GET /api/organizations/my-org/platform-license was never called — double /api prefix likely")
	}
}

func TestMakeAuthMethod(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Console credential auth logs in against /login before any resource call.
		if r.URL.Path == "/api/login" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "minted-token"})
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	apiKeyClient, err := Make(context.Background(), CONSOLE, ApiParameter{BaseUrl: ts.URL, ApiKey: "test-key"}, "test")
	if err != nil {
		t.Fatalf("failed to create api-key client: %v", err)
	}
	if apiKeyClient.AuthMethod != AuthMethodApiKey {
		t.Errorf("expected AuthMethodApiKey, got %q", apiKeyClient.AuthMethod)
	}

	credClient, err := Make(context.Background(), CONSOLE, ApiParameter{BaseUrl: ts.URL, CdkUser: "admin", CdkPassword: "secret"}, "test")
	if err != nil {
		t.Fatalf("failed to create credential client: %v", err)
	}
	if credClient.AuthMethod != AuthMethodCredentials {
		t.Errorf("expected AuthMethodCredentials, got %q", credClient.AuthMethod)
	}
}

func TestUniformizeBaseUrl(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"https://console.example.com", "https://console.example.com/api"},
		{"https://console.example.com/", "https://console.example.com/api"},
		{"https://console.example.com/api", "https://console.example.com/api"},
		{"https://console.example.com/api/", "https://console.example.com/api"},
	}
	for _, tc := range cases {
		got := uniformizeBaseUrl(tc.input)
		if got != tc.expected {
			t.Errorf("uniformizeBaseUrl(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
