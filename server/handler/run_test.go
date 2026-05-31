package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunHandlerRejectsUnknownLanguage(t *testing.T) {
	body := `{
		"language": "ruby",
		"source": "puts 'hello'",
		"tests": [{"stdin": "", "expected_stdout": "hello\n"}]
	}`
	req := httptest.NewRequest(http.MethodPost, "/run", strings.NewReader(body))
	rec := httptest.NewRecorder()

	RunHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), `"code":"unknown_language"`) {
		t.Fatalf("body = %s, want unknown_language error", rec.Body.String())
	}
}

func TestRunHandlerRejectsBadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/run", bytes.NewBufferString(`{`))
	rec := httptest.NewRecorder()

	RunHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), `"code":"bad_json"`) {
		t.Fatalf("body = %s, want bad_json error", rec.Body.String())
	}
}
