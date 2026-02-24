// Package registry provides registry.yaml parsing and management
package registry

import (
	"embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed registry.yaml
var registryFS embed.FS

// Source represents a skill source (repository)
type Source struct {
	Name      string  // Source identifier (e.g., "anthropic", "vercel")
	Repo      string  `yaml:"repo"`       // Repository URL (e.g., "github.com/anthropics/skills")
	License   string  `yaml:"license"`    // License type
	SkipFetch bool    `yaml:"skip_fetch"` // Skip dynamic fetching (for large repos)
	Skills    []Skill `yaml:"skills"`     // Skills in this source
}

// Skill represents a skill entry in the registry
type Skill struct {
	Name          string `yaml:"name"`           // Skill name
	Path          string `yaml:"path"`           // Path in repository
	Description   string `yaml:"description"`    // Short description (English)
	DescriptionZh string `yaml:"description_zh"` // Short description (Chinese)
	Version       string `yaml:"version"`        // Version (optional)
}

// GetDescription returns the description based on language
// lang should be "zh" for Chinese, otherwise English
func (s *Skill) GetDescription(lang string) string {
	if lang == "zh" && s.DescriptionZh != "" {
		return s.DescriptionZh
	}
	return s.Description
}

// Registry holds all sources from registry.yaml
type Registry struct {
	Sources map[string]*Source
}

// registryYAML is the raw YAML structure
type registryYAML map[string]struct {
	Repo      string `yaml:"repo"`
	License   string `yaml:"license"`
	SkipFetch bool   `yaml:"skip_fetch"`
	Skills    []struct {
		Name          string `yaml:"name"`
		Path          string `yaml:"path"`
		Description   string `yaml:"description"`
		DescriptionZh string `yaml:"description_zh"`
		Version       string `yaml:"version"`
	} `yaml:"skills"`
}

// Load loads the embedded registry.yaml
func Load() (*Registry, error) {
	data, err := registryFS.ReadFile("registry.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read registry.yaml: %w", err)
	}

	return Parse(data)
}

// Parse parses registry.yaml content
func Parse(data []byte) (*Registry, error) {
	var raw registryYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse registry.yaml: %w", err)
	}

	registry := &Registry{
		Sources: make(map[string]*Source),
	}

	for name, src := range raw {
		source := &Source{
			Name:      name,
			Repo:      src.Repo,
			License:   src.License,
			SkipFetch: src.SkipFetch,
			Skills:    make([]Skill, 0, len(src.Skills)),
		}

		for _, s := range src.Skills {
			skill := Skill{
				Name:          s.Name,
				Path:          s.Path,
				Description:   s.Description,
				DescriptionZh: s.DescriptionZh,
				Version:       s.Version,
			}
			source.Skills = append(source.Skills, skill)
		}

		registry.Sources[name] = source
	}

	return registry, nil
}

// GetAllSources returns all sources in the registry
func (r *Registry) GetAllSources() []*Source {
	sources := make([]*Source, 0, len(r.Sources))
	for _, src := range r.Sources {
		sources = append(sources, src)
	}
	return sources
}

// GetSource returns a source by name
func (r *Registry) GetSource(name string) *Source {
	return r.Sources[name]
}

// FindSkill finds a skill by name across all sources
// Returns the skill and its source, or nil if not found
func (r *Registry) FindSkill(skillName string) (*Skill, *Source) {
	skillNameLower := strings.ToLower(skillName)
	for _, src := range r.Sources {
		for i := range src.Skills {
			if strings.ToLower(src.Skills[i].Name) == skillNameLower {
				return &src.Skills[i], src
			}
		}
	}
	return nil, nil
}

// FindSkillsWithConflict finds all skills matching a name (for conflict detection)
func (r *Registry) FindSkillsWithConflict(skillName string) []struct {
	Skill  *Skill
	Source *Source
} {
	var matches []struct {
		Skill  *Skill
		Source *Source
	}

	skillNameLower := strings.ToLower(skillName)
	for _, src := range r.Sources {
		for i := range src.Skills {
			if strings.ToLower(src.Skills[i].Name) == skillNameLower {
				matches = append(matches, struct {
					Skill  *Skill
					Source *Source
				}{&src.Skills[i], src})
			}
		}
	}
	return matches
}

// TotalSkillCount returns the total number of skills in the registry
func (r *Registry) TotalSkillCount() int {
	count := 0
	for _, src := range r.Sources {
		count += len(src.Skills)
	}
	return count
}

// GetGitURL returns the git clone URL for a source
func (s *Source) GetGitURL() string {
	// Convert github.com/owner/repo to https://github.com/owner/repo.git
	if strings.HasPrefix(s.Repo, "github.com/") {
		return "https://" + s.Repo + ".git"
	}
	return s.Repo
}

// GetRepoShortName returns a short display name for the repo
func (s *Source) GetRepoShortName() string {
	// github.com/owner/repo -> owner/repo
	if strings.HasPrefix(s.Repo, "github.com/") {
		return strings.TrimPrefix(s.Repo, "github.com/")
	}
	return s.Repo
}
