package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestGetConsoleLicensePlan_NoDoubleApiPrefix(t *testing.T) {
	orgsHit := false
	licenseHit := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/organizations":
			orgsHit = true
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{{"slug": "my-org"}})
		case "/api/organizations/my-org/platform-license":
			licenseHit = true
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"plan": "enterprise"})
		default:
			// Return HTML like Console SPA does for unmatched routes — this is what triggers the bug
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<html>not found</html>"))
		}
	}))
	defer ts.Close()

	c := &Client{
		BaseUrl: ts.URL + "/api",
		Client:  resty.New(),
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
