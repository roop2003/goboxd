package runs

import (
	"testing"

	"github.com/thesouldev/goboxd/internal/languages"
)

func TestValidateRequest(t *testing.T) {
	registry, err := languages.LoadDefaultRegistry()
	if err != nil {
		t.Fatalf("load default registry: %v", err)
	}

	stdin := ""
	expected := "OK\n"
	source := "print('OK')"
	langPy := "py3"
	langCPP := "cpp"
	langJava := "java"
	badLang := "ruby"

	tests := []struct {
		name string
		req  Request
		code string
	}{
		{
			name: "valid interpreted language",
			req: Request{
				Language: &langPy,
				Source:   &source,
				Tests:    []Test{{Stdin: &stdin, ExpectedStdout: &expected}},
			},
		},
		{
			name: "valid wildcard build flag",
			req: Request{
				Language: &langCPP,
				Source:   &source,
				Build:    &Options{Flags: []string{"-std=c++17"}},
				Tests:    []Test{{Stdin: &stdin, ExpectedStdout: &expected}},
			},
		},
		{
			name: "unknown language",
			req: Request{
				Language: &badLang,
				Source:   &source,
				Tests:    []Test{{Stdin: &stdin, ExpectedStdout: &expected}},
			},
			code: "unknown_language",
		},
		{
			name: "java requires request filenames",
			req: Request{
				Language: &langJava,
				Source:   &source,
				Tests:    []Test{{Stdin: &stdin, ExpectedStdout: &expected}},
			},
			code: "missing_source_filename",
		},
		{
			name: "reject path traversal filename",
			req: Request{
				Language:       &langJava,
				Source:         &source,
				SourceFilename: "../Main.java",
				Tests:          []Test{{Stdin: &stdin, ExpectedStdout: &expected}},
			},
			code: "invalid_filename",
		},
		{
			name: "reject disallowed flag",
			req: Request{
				Language: &langCPP,
				Source:   &source,
				Build:    &Options{Flags: []string{"-pipe"}},
				Tests:    []Test{{Stdin: &stdin, ExpectedStdout: &expected}},
			},
			code: "disallowed_flag",
		},
		{
			name: "missing expected stdout",
			req: Request{
				Language: &langPy,
				Source:   &source,
				Tests:    []Test{{Stdin: &stdin}},
			},
			code: "missing_expected_stdout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRequest(tt.req, registry)
			if tt.code == "" {
				if err != nil {
					t.Fatalf("expected valid request, got %s: %s", err.Code, err.Message)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error code %q", tt.code)
			}
			if err.Code != tt.code {
				t.Fatalf("expected error code %q, got %q", tt.code, err.Code)
			}
		})
	}
}
