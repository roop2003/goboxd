package runs

import "unicode/utf8"

const DefaultMaxOutputBytes = 64 * 1024

func TruncateOutput(output string, maxBytes int) (string, bool) {
	if maxBytes < 0 {
		maxBytes = 0
	}
	if len(output) <= maxBytes {
		return output, false
	}

	truncated := output[:maxBytes]
	for !utf8.ValidString(truncated) && len(truncated) > 0 {
		truncated = truncated[:len(truncated)-1]
	}
	return truncated, true
}
