package runs

import "testing"

func TestBuildStatus(t *testing.T) {
	tests := []struct {
		name    string
		exit    int
		timeout bool
		want    Status
	}{
		{name: "ok", want: StatusOK},
		{name: "failed", exit: 1, want: StatusBuildFailed},
		{name: "timeout", timeout: true, want: StatusTimeLimitExceeded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildStatus(tt.exit, tt.timeout); got != tt.want {
				t.Fatalf("BuildStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTestStatus(t *testing.T) {
	tests := []struct {
		name           string
		stdout         string
		expectedStdout string
		exit           int
		timeout        bool
		memoryExceeded bool
		want           Status
	}{
		{name: "accepted", stdout: "OK\n", expectedStdout: "OK\n", want: StatusAccepted},
		{name: "wrong output", stdout: "NO\n", expectedStdout: "OK\n", want: StatusWrongOutput},
		{name: "runtime error", exit: 2, want: StatusRuntimeError},
		{name: "timeout", timeout: true, want: StatusTimeLimitExceeded},
		{name: "memory exceeded", memoryExceeded: true, want: StatusMemoryLimitExceeded},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TestStatus(tt.stdout, tt.expectedStdout, tt.exit, tt.timeout, tt.memoryExceeded)
			if got != tt.want {
				t.Fatalf("TestStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOverallStatus(t *testing.T) {
	tests := []struct {
		name  string
		build Status
		tests []Status
		want  Status
	}{
		{name: "all accepted", build: StatusOK, tests: []Status{StatusAccepted}, want: StatusAccepted},
		{name: "build failed", build: StatusBuildFailed, tests: []Status{StatusAccepted}, want: StatusBuildFailed},
		{name: "first failing test", build: StatusOK, tests: []Status{StatusAccepted, StatusWrongOutput}, want: StatusWrongOutput},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OverallStatus(tt.build, tt.tests); got != tt.want {
				t.Fatalf("OverallStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}
