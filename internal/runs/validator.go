package runs

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/thesouldev/goboxd/internal/languages"
)

const (
	MaxSourceBytes  = 256 * 1024
	MaxFilenameByte = 128
)

func ValidateRequest(req Request, registry *languages.Registry) *Error {
	if req.Language == nil || strings.TrimSpace(*req.Language) == "" {
		return badRequest("missing_language", "language is required")
	}

	lang, ok := registry.Get(*req.Language)
	if !ok {
		return badRequest("unknown_language", "language is not configured")
	}

	if req.Source == nil {
		return badRequest("missing_source", "source is required")
	}
	if !utf8.ValidString(*req.Source) {
		return badRequest("invalid_source_encoding", "source must be UTF-8")
	}
	if len([]byte(*req.Source)) > MaxSourceBytes {
		return badRequest("source_too_large", "source exceeds maximum size")
	}

	if err := validateFilename("source_filename", req.SourceFilename, lang.RequiresSourceFilename()); err != nil {
		return err
	}
	if err := validateFilename("artifact_filename", req.ArtifactFilename, lang.RequiresArtifactFilename()); err != nil {
		return err
	}

	if err := validateOptions("build", req.Build, lang.Build); err != nil {
		return err
	}
	if err := validateOptions("run", req.Run, lang.Run); err != nil {
		return err
	}

	if len(req.Tests) == 0 {
		return badRequest("missing_tests", "tests must contain at least one entry")
	}
	for i, test := range req.Tests {
		if test.Stdin == nil {
			return badRequest("missing_test_stdin", fmt.Sprintf("tests[%d].stdin is required", i))
		}
		if test.ExpectedStdout == nil {
			return badRequest("missing_expected_stdout", fmt.Sprintf("tests[%d].expected_stdout is required", i))
		}
	}

	return nil
}

func validateFilename(field, name string, required bool) *Error {
	if name == "" {
		if required {
			return badRequest("missing_"+field, field+" is required for this language")
		}
		return nil
	}
	if len([]byte(name)) > MaxFilenameByte || strings.HasPrefix(name, ".") || strings.ContainsAny(name, `/\`) {
		return badRequest("invalid_filename", field+" must be a single path component")
	}
	return nil
}

func validateOptions(field string, requested *Options, configured *languages.Command) *Error {
	if requested == nil {
		return nil
	}
	if configured == nil {
		return badRequest("unsupported_"+field, field+" options are not supported for this language")
	}
	if err := validateLimits(field, requested.Limits); err != nil {
		return err
	}
	for _, flag := range requested.Flags {
		if flag == "" {
			return badRequest("invalid_flag", field+" flags cannot contain empty strings")
		}
		if !configured.AllowsFlag(flag) {
			return badRequest("disallowed_flag", field+" flag is not allowed: "+flag)
		}
	}
	return nil
}

func validateLimits(field string, limits *Limits) *Error {
	if limits == nil {
		return nil
	}
	if limits.WallTimeSeconds != nil && *limits.WallTimeSeconds <= 0 {
		return badRequest("invalid_limits", field+".limits.wall_time_s must be greater than zero")
	}
	if limits.MemoryKB != nil && *limits.MemoryKB <= 0 {
		return badRequest("invalid_limits", field+".limits.memory_kb must be greater than zero")
	}
	if limits.MaxProcesses != nil && *limits.MaxProcesses <= 0 {
		return badRequest("invalid_limits", field+".limits.max_processes must be greater than zero")
	}
	return nil
}

func badRequest(code, message string) *Error {
	return &Error{Code: code, Message: message}
}
