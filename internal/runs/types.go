package runs

type Request struct {
	Language         *string  `json:"language"`
	Source           *string  `json:"source"`
	SourceFilename   string   `json:"source_filename,omitempty"`
	ArtifactFilename string   `json:"artifact_filename,omitempty"`
	Build            *Options `json:"build,omitempty"`
	Run              *Options `json:"run,omitempty"`
	Tests            []Test   `json:"tests"`
}

type Options struct {
	Limits *Limits  `json:"limits,omitempty"`
	Flags  []string `json:"flags,omitempty"`
}

type Limits struct {
	WallTimeSeconds *int `json:"wall_time_s,omitempty"`
	MemoryKB        *int `json:"memory_kb,omitempty"`
	MaxProcesses    *int `json:"max_processes,omitempty"`
}

type Test struct {
	Stdin          *string `json:"stdin"`
	ExpectedStdout *string `json:"expected_stdout"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Response struct {
	Status Status        `json:"status"`
	Build  *BuildResult  `json:"build"`
	Tests  []TestResult  `json:"tests"`
}

type BuildResult struct {
	Status     Status `json:"status"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	DurationMS int64  `json:"duration_ms"`
}

type TestResult struct {
	Status       Status `json:"status"`
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	DurationMS   int64  `json:"duration_ms"`
	MemoryPeakKB int64  `json:"memory_peak_kb"`
}
