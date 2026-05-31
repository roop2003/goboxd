package runs

import (
	"reflect"
	"testing"

	"github.com/thesouldev/goboxd/internal/languages"
)

func TestExpandArgs(t *testing.T) {
	got := expandArgs(
		[]string{"{{flags}}", "-o", "{{artifact}}", "{{source}}"},
		[]string{"-std=c++17", "-Wall"},
		"solution.cpp",
		"solution",
	)
	want := []string{"-std=c++17", "-Wall", "-o", "solution", "solution.cpp"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expandArgs() = %#v, want %#v", got, want)
	}
}

func TestMergeLimits(t *testing.T) {
	wallTime := 5
	defaults := languages.Limits{
		WallTimeSeconds: 3,
		MemoryKB:        1024,
		MaxProcesses:    8,
	}

	got := mergeLimits(defaults, &Options{
		Limits: &Limits{WallTimeSeconds: &wallTime},
	})

	if got.WallTimeSeconds != wallTime {
		t.Fatalf("WallTimeSeconds = %d, want %d", got.WallTimeSeconds, wallTime)
	}
	if got.MemoryKB != defaults.MemoryKB {
		t.Fatalf("MemoryKB = %d, want %d", got.MemoryKB, defaults.MemoryKB)
	}
	if got.MaxProcesses != defaults.MaxProcesses {
		t.Fatalf("MaxProcesses = %d, want %d", got.MaxProcesses, defaults.MaxProcesses)
	}
}
