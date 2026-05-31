package runs

import "testing"

func TestTruncateOutput(t *testing.T) {
	tests := []struct {
		name      string
		output    string
		maxBytes  int
		want      string
		truncated bool
	}{
		{name: "under limit", output: "hello", maxBytes: 10, want: "hello"},
		{name: "over limit", output: "hello", maxBytes: 2, want: "he", truncated: true},
		{name: "negative limit", output: "hello", maxBytes: -1, want: "", truncated: true},
		{name: "utf8 boundary", output: "a界", maxBytes: 2, want: "a", truncated: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, truncated := TruncateOutput(tt.output, tt.maxBytes)
			if got != tt.want {
				t.Fatalf("TruncateOutput() output = %q, want %q", got, tt.want)
			}
			if truncated != tt.truncated {
				t.Fatalf("TruncateOutput() truncated = %v, want %v", truncated, tt.truncated)
			}
		})
	}
}
