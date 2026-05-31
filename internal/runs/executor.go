package runs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/thesouldev/goboxd/internal/languages"
)

type Executor struct {
	Registry *languages.Registry
}

func (e Executor) Execute(ctx context.Context, req Request) (Response, error) {
	lang, ok := e.Registry.Get(*req.Language)
	if !ok {
		return Response{}, fmt.Errorf("language %q disappeared after validation", *req.Language)
	}

	workDir, err := os.MkdirTemp("", "goboxd-run-*")
	if err != nil {
		return Response{}, err
	}
	defer os.RemoveAll(workDir)

	sourceName := filenameFor(lang.SourceFilename, req.SourceFilename, lang.RequiresSourceFilename())
	artifactName := filenameFor(lang.Artifact, req.ArtifactFilename, lang.RequiresArtifactFilename())

	if err := os.WriteFile(filepath.Join(workDir, sourceName), []byte(*req.Source), 0600); err != nil {
		return Response{}, err
	}

	var build *BuildResult
	var buildStatus Status
	if lang.Build != nil {
		result, err := runPhase(ctx, workDir, *lang.Build, req.Build, sourceName, artifactName, "")
		if err != nil {
			return Response{}, err
		}
		status := BuildStatus(result.exitCode, result.timedOut)
		buildStatus = status
		build = &BuildResult{
			Status:     status,
			Stdout:     result.stdout,
			Stderr:     result.stderr,
			DurationMS: result.duration.Milliseconds(),
		}
		if status != StatusOK {
			return Response{Status: status, Build: build, Tests: []TestResult{}}, nil
		}
	}

	testResults := make([]TestResult, 0, len(req.Tests))
	testStatuses := make([]Status, 0, len(req.Tests))
	for _, test := range req.Tests {
		result, err := runPhase(ctx, workDir, *lang.Run, req.Run, sourceName, artifactName, *test.Stdin)
		if err != nil {
			return Response{}, err
		}
		status := TestStatus(result.stdout, *test.ExpectedStdout, result.exitCode, result.timedOut, false)
		testStatuses = append(testStatuses, status)
		testResults = append(testResults, TestResult{
			Status:       status,
			Stdout:       result.stdout,
			Stderr:       result.stderr,
			DurationMS:   result.duration.Milliseconds(),
			MemoryPeakKB: 0,
		})
	}

	return Response{
		Status: OverallStatus(buildStatus, testStatuses),
		Build:  build,
		Tests:  testResults,
	}, nil
}

type phaseResult struct {
	stdout   string
	stderr   string
	duration time.Duration
	exitCode int
	timedOut bool
}

func runPhase(ctx context.Context, workDir string, command languages.Command, overrides *Options, sourceName, artifactName, stdin string) (phaseResult, error) {
	limits := mergeLimits(command.Limits, overrides)
	timeout := time.Duration(limits.WallTimeSeconds) * time.Second
	if timeout <= 0 {
		timeout = time.Second
	}

	phaseCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmdName := replacePlaceholders(command.Cmd, sourceName, artifactName)
	args := expandArgs(command.Args, flagsFor(overrides), sourceName, artifactName)
	cmd := exec.CommandContext(phaseCtx, cmdName, args...)
	cmd.Dir = workDir
	cmd.Stdin = strings.NewReader(stdin)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	started := time.Now()
	err := cmd.Run()
	duration := time.Since(started)

	out, _ := TruncateOutput(stdout.String(), DefaultMaxOutputBytes)
	errOut, _ := TruncateOutput(stderr.String(), DefaultMaxOutputBytes)
	result := phaseResult{stdout: out, stderr: errOut, duration: duration}

	if errors.Is(phaseCtx.Err(), context.DeadlineExceeded) {
		result.timedOut = true
		return result, nil
	}
	if err == nil {
		return result, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.exitCode = exitErr.ExitCode()
		return result, nil
	}
	return result, err
}

func filenameFor(configured string, requested string, fromRequest bool) string {
	if fromRequest {
		return requested
	}
	return configured
}

func flagsFor(options *Options) []string {
	if options == nil {
		return nil
	}
	return options.Flags
}

func mergeLimits(defaults languages.Limits, options *Options) languages.Limits {
	limits := defaults
	if options == nil || options.Limits == nil {
		return limits
	}
	if options.Limits.WallTimeSeconds != nil {
		limits.WallTimeSeconds = *options.Limits.WallTimeSeconds
	}
	if options.Limits.MemoryKB != nil {
		limits.MemoryKB = *options.Limits.MemoryKB
	}
	if options.Limits.MaxProcesses != nil {
		limits.MaxProcesses = *options.Limits.MaxProcesses
	}
	return limits
}

func expandArgs(args []string, flags []string, sourceName string, artifactName string) []string {
	expanded := make([]string, 0, len(args)+len(flags))
	for _, arg := range args {
		if arg == "{{flags}}" {
			expanded = append(expanded, flags...)
			continue
		}
		expanded = append(expanded, replacePlaceholders(arg, sourceName, artifactName))
	}
	return expanded
}

func replacePlaceholders(value string, sourceName string, artifactName string) string {
	value = strings.ReplaceAll(value, "{{source}}", sourceName)
	value = strings.ReplaceAll(value, "{{artifact}}", artifactName)
	return value
}
