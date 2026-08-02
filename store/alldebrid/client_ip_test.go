// Package alldebrid verifies client IP forwarding behavior.
package alldebrid

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/MunifTanjim/stremthru/internal/request"
)

func TestWithClientIP(t *testing.T) {
	originalQuery := url.Values{
		"existing": {"value"},
		"ip":       {"untrusted"},
	}
	ctx := request.Ctx{Query: &originalQuery}
	ctx = withClientIP(ctx, "203.0.113.7")

	client := NewAPIClient(&APIClientConfig{})
	req, err := ctx.NewRequest(client.BaseURL, http.MethodGet, "/v4/user", client.reqHeader, client.reqQuery)
	if err != nil {
		t.Fatal(err)
	}

	query := req.URL.Query()
	if got, want := query.Get("ip"), "203.0.113.7"; got != want {
		t.Fatalf("ip query parameter = %q, want %q", got, want)
	}
	if got, want := query.Get("agent"), "stremthru"; got != want {
		t.Fatalf("agent query parameter = %q, want %q", got, want)
	}
	if got, want := query.Get("existing"), "value"; got != want {
		t.Fatalf("existing query parameter = %q, want %q", got, want)
	}
	if got, want := originalQuery.Get("ip"), "untrusted"; got != want {
		t.Fatalf("original query was mutated: ip = %q, want %q", got, want)
	}
}

func TestWithClientIPEmpty(t *testing.T) {
	ctx := withClientIP(request.Ctx{}, "")
	client := NewAPIClient(&APIClientConfig{})
	req, err := ctx.NewRequest(client.BaseURL, http.MethodGet, "/v4/user", client.reqHeader, client.reqQuery)
	if err != nil {
		t.Fatal(err)
	}

	if got := req.URL.Query().Get("ip"); got != "" {
		t.Fatalf("ip query parameter = %q, want empty", got)
	}
}
