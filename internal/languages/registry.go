package languages

import (
	"fmt"
	"os"
	"strings"

	"github.com/thesouldev/goboxd/config"
	"gopkg.in/yaml.v3"
)

const StrategyFromRequest = "from_request"

type Config struct {
	Languages []Language `yaml:"languages"`
}

type Language struct {
	ID                       string   `yaml:"id"`
	Name                     string   `yaml:"name"`
	SourceFilename           string   `yaml:"source_filename"`
	SourceFilenameStrategy   string   `yaml:"source_filename_strategy"`
	Artifact                 string   `yaml:"artifact"`
	ArtifactFilenameStrategy string   `yaml:"artifact_filename_strategy"`
	Build                    *Command `yaml:"build"`
	Run                      *Command `yaml:"run"`
}

type Command struct {
	Cmd           string   `yaml:"cmd"`
	Args          []string `yaml:"args"`
	Limits        Limits   `yaml:"limits"`
	FlagAllowlist []string `yaml:"flag_allowlist"`
}

type Limits struct {
	WallTimeSeconds int `yaml:"wall_time_s"`
	MemoryKB        int `yaml:"memory_kb"`
	MaxProcesses    int `yaml:"max_processes"`
}

type Registry struct {
	byID map[string]Language
}

func LoadDefaultRegistry() (*Registry, error) {
	return LoadRegistryBytes(config.LanguagesYAML)
}

func LoadRegistry(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadRegistryBytes(data)
}

func LoadRegistryBytes(data []byte) (*Registry, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return NewRegistry(cfg)
}

func NewRegistry(cfg Config) (*Registry, error) {
	reg := &Registry{byID: make(map[string]Language, len(cfg.Languages))}
	for _, lang := range cfg.Languages {
		if err := validateLanguage(lang); err != nil {
			return nil, err
		}
		if _, exists := reg.byID[lang.ID]; exists {
			return nil, fmt.Errorf("duplicate language id %q", lang.ID)
		}
		reg.byID[lang.ID] = lang
	}
	if len(reg.byID) == 0 {
		return nil, fmt.Errorf("at least one language must be configured")
	}
	return reg, nil
}

func (r *Registry) Get(id string) (Language, bool) {
	if r == nil {
		return Language{}, false
	}
	lang, ok := r.byID[id]
	return lang, ok
}

func (l Language) RequiresSourceFilename() bool {
	return l.SourceFilenameStrategy == StrategyFromRequest
}

func (l Language) RequiresArtifactFilename() bool {
	return l.ArtifactFilenameStrategy == StrategyFromRequest
}

func (c Command) AllowsFlag(flag string) bool {
	for _, allowed := range c.FlagAllowlist {
		if allowed == flag {
			return true
		}
		if strings.HasSuffix(allowed, "*") && strings.HasPrefix(flag, strings.TrimSuffix(allowed, "*")) {
			return true
		}
	}
	return false
}

func validateLanguage(lang Language) error {
	if strings.TrimSpace(lang.ID) == "" {
		return fmt.Errorf("language id is required")
	}
	if err := validateStrategy(lang.ID, "source_filename_strategy", lang.SourceFilenameStrategy); err != nil {
		return err
	}
	if err := validateStrategy(lang.ID, "artifact_filename_strategy", lang.ArtifactFilenameStrategy); err != nil {
		return err
	}
	if lang.SourceFilename == "" && lang.SourceFilenameStrategy != StrategyFromRequest {
		return fmt.Errorf("language %q must define source_filename or source_filename_strategy", lang.ID)
	}
	if lang.Build != nil && strings.TrimSpace(lang.Build.Cmd) == "" {
		return fmt.Errorf("language %q build.cmd is required when build is configured", lang.ID)
	}
	if lang.Run == nil || strings.TrimSpace(lang.Run.Cmd) == "" {
		return fmt.Errorf("language %q run.cmd is required", lang.ID)
	}
	return nil
}

func validateStrategy(languageID, field, strategy string) error {
	if strategy == "" || strategy == StrategyFromRequest {
		return nil
	}
	return fmt.Errorf("language %q has unsupported %s %q", languageID, field, strategy)
}
