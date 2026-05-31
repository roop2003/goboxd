//go:build integration

package integration_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/thesouldev/goboxd/server"
)

func TestHealthz(t *testing.T) {
	testServer := httptest.NewServer(server.NewMux())
	defer testServer.Close()

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(testServer.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if string(body) != "OK" {
		t.Fatalf("body = %q, want %q", string(body), "OK")
	}
}

func TestRunRejectsUnknownLanguage(t *testing.T) {
	testServer := httptest.NewServer(server.NewMux())
	defer testServer.Close()

	body := `{"language":"ruby","source":"puts 'hi'","tests":[{"stdin":"","expected_stdout":"hi\n"}]}`
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Post(testServer.URL+"/run", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST /run: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}
