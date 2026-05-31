package languages

import "testing"

func TestLoadDefaultRegistry(t *testing.T) {
	registry, err := LoadDefaultRegistry()
	if err != nil {
		t.Fatalf("LoadDefaultRegistry() error = %v", err)
	}

	if _, ok := registry.Get("py3"); !ok {
		t.Fatal("expected py3 language to be configured")
	}
	if _, ok := registry.Get("cpp"); !ok {
		t.Fatal("expected cpp language to be configured")
	}
}

func TestNewRegistryRejectsDuplicateIDs(t *testing.T) {
	_, err := NewRegistry(Config{Languages: []Language{
		{ID: "py3", SourceFilename: "solution.py", Run: &Command{Cmd: "/usr/bin/python3"}},
		{ID: "py3", SourceFilename: "main.py", Run: &Command{Cmd: "/usr/bin/python3"}},
	}})
	if err == nil {
		t.Fatal("expected duplicate language id error")
	}
}

func TestNewRegistryRejectsUnknownStrategy(t *testing.T) {
	_, err := NewRegistry(Config{Languages: []Language{
		{
			ID:                     "request-file",
			SourceFilenameStrategy: "from_path",
			Run:                    &Command{Cmd: "/bin/true"},
		},
	}})
	if err == nil {
		t.Fatal("expected unsupported strategy error")
	}
}

func TestCommandAllowsFlag(t *testing.T) {
	command := Command{FlagAllowlist: []string{"-Wall", "-std=*"}}

	tests := []struct {
		flag string
		want bool
	}{
		{flag: "-Wall", want: true},
		{flag: "-std=c++17", want: true},
		{flag: "-pipe", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.flag, func(t *testing.T) {
			if got := command.AllowsFlag(tt.flag); got != tt.want {
				t.Fatalf("AllowsFlag(%q) = %v, want %v", tt.flag, got, tt.want)
			}
		})
	}
}
