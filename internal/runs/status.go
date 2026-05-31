package runs

type Status string

const (
	StatusOK                  Status = "ok"
	StatusAccepted            Status = "accepted"
	StatusWrongOutput         Status = "wrong_output"
	StatusBuildFailed         Status = "build_failed"
	StatusRuntimeError        Status = "runtime_error"
	StatusTimeLimitExceeded   Status = "time_limit_exceeded"
	StatusMemoryLimitExceeded Status = "memory_limit_exceeded"
)

func BuildStatus(exitCode int, timedOut bool) Status {
	if timedOut {
		return StatusTimeLimitExceeded
	}
	if exitCode != 0 {
		return StatusBuildFailed
	}
	return StatusOK
}

func TestStatus(stdout, expectedStdout string, exitCode int, timedOut bool, memoryExceeded bool) Status {
	switch {
	case timedOut:
		return StatusTimeLimitExceeded
	case memoryExceeded:
		return StatusMemoryLimitExceeded
	case exitCode != 0:
		return StatusRuntimeError
	case stdout != expectedStdout:
		return StatusWrongOutput
	default:
		return StatusAccepted
	}
}

func OverallStatus(build Status, tests []Status) Status {
	if build != "" && build != StatusOK {
		return build
	}
	for _, status := range tests {
		if status != StatusAccepted && status != StatusOK {
			return status
		}
	}
	return StatusAccepted
}
