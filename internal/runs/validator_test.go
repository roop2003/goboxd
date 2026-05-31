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
	requestFilenameRegistry, err := languages.NewRegistry(languages.Config{Languages: []languages.Language{
		{
			ID:                       "request-file",
			SourceFilenameStrategy:   languages.StrategyFromRequest,
			ArtifactFilenameStrategy: languages.StrategyFromRequest,
			Run:                      &languages.Command{Cmd: "/bin/true"},
		},
	}})
	if err != nil {
		t.Fatalf("build request filename registry: %v", err)
	}

	stdin := ""
	expected := "OK\n"
	source := "print('OK')"
	langPy := "python"
	langCPP := "cpp"
	langRequestFile := "request-file"
	badLang := "ruby"

	tests := []struct {
		name     string
		req      Request
		registry *languages.Registry
		code     string
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
			name: "configured language requires request filenames",
			req: Request{
				Language: &langRequestFile,
				Source:   &source,
				Tests:    []Test{{Stdin: &stdin, ExpectedStdout: &expected}},
			},
			registry: requestFilenameRegistry,
			code:     "missing_source_filename",
		},
		{
			name: "reject path traversal filename",
			req: Request{
				Language:       &langRequestFile,
				Source:         &source,
				SourceFilename: "../solution.txt",
				Tests:          []Test{{Stdin: &stdin, ExpectedStdout: &expected}},
			},
			registry: requestFilenameRegistry,
			code:     "invalid_filename",
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
			activeRegistry := registry
			if tt.registry != nil {
				activeRegistry = tt.registry
			}
			err := ValidateRequest(tt.req, activeRegistry)
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
