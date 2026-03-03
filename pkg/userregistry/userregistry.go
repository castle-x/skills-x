// Package userregistry manages the user-local skill registry stored at
// ~/.config/skills-x/user-registry.yaml.
//
// The format is intentionally identical to pkg/registry/registry.yaml so
// that users can share entries, and the CLI/TUI can load both with the same
// parser.
package userregistry

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// SkillEntry is a single skill inside a source.
type SkillEntry struct {
	Name          string   `yaml:"name"`
	Path          string   `yaml:"path"`
	Tags          []string `yaml:"tags,omitempty"`
	Description   string   `yaml:"description"`
	DescriptionZh string   `yaml:"description_zh,omitempty"`
	License       string   `yaml:"license,omitempty"`
}

// SourceEntry is a source (repository or local dir) containing skills.
type SourceEntry struct {
	Repo    string       `yaml:"repo"`
	License string       `yaml:"license,omitempty"`
	Skills  []SkillEntry `yaml:"skills"`
}

// UserRegistry holds the contents of user-registry.yaml.
type UserRegistry struct {
	// Sources maps source-name → SourceEntry, preserving YAML order.
	Sources map[string]*SourceEntry
	// orderedKeys keeps insertion order for deterministic YAML output.
	orderedKeys []string
}

// FilePath returns the canonical path for the user registry file.
func FilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(configDir, "skills-x", "user-registry.yaml")
}

// Load reads (or initialises) the user registry from disk.
// Returns an empty registry (not an error) when the file does not exist yet.
func Load() (*UserRegistry, error) {
	path := FilePath()
	ur := &UserRegistry{Sources: make(map[string]*SourceEntry)}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ur, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading user registry: %w", err)
	}

	// Decode into ordered map to preserve key order on save.
	var raw yaml.Node
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing user registry: %w", err)
	}

	if raw.Kind == yaml.DocumentNode && len(raw.Content) > 0 {
		mapNode := raw.Content[0]
		if mapNode.Kind == yaml.MappingNode {
			for i := 0; i+1 < len(mapNode.Content); i += 2 {
				key := mapNode.Content[i].Value
				var entry SourceEntry
				if err := mapNode.Content[i+1].Decode(&entry); err != nil {
					return nil, fmt.Errorf("decoding source %q: %w", key, err)
				}
				ur.Sources[key] = &entry
				ur.orderedKeys = append(ur.orderedKeys, key)
			}
		}
	}

	return ur, nil
}

// Save writes the registry back to disk, creating directories if needed.
func (ur *UserRegistry) Save() error {
	path := FilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	// Serialise in insertion order.
	var sb strings.Builder
	for _, key := range ur.orderedKeys {
		src := ur.Sources[key]
		data, err := yaml.Marshal(map[string]*SourceEntry{key: src})
		if err != nil {
			return fmt.Errorf("encoding source %q: %w", key, err)
		}
		sb.Write(data)
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

// AddResult is returned by Add to inform callers about conflicts.
type AddResult struct {
	// ConflictSources lists built-in source names that have the same skill name.
	ConflictSources []string
	// SourceName is the source key used/created in the user registry.
	SourceName string
}

// Add appends a skill to the user registry.
// builtinSkillNames is the set of skill names already present in the built-in
// registry (used for conflict detection). Pass nil to skip conflict checks.
//
// The skill is placed under a source key derived from the repo:
//   - "github.com/owner/repo" → "owner-repo"
//   - "/local/path"           → "local"
func (ur *UserRegistry) Add(
	repo, path, skillName, description, descriptionZh, license string,
	builtinSkillNames map[string][]string, // name → []sourceName
) (*AddResult, error) {
	result := &AddResult{}

	// Conflict detection.
	if builtinSkillNames != nil {
		if sources, ok := builtinSkillNames[strings.ToLower(skillName)]; ok {
			result.ConflictSources = sources
		}
	}

	// Also check within the user registry itself.
	for _, src := range ur.Sources {
		for _, s := range src.Skills {
			if strings.EqualFold(s.Name, skillName) {
				return nil, fmt.Errorf("skill %q already exists in user registry", skillName)
			}
		}
	}

	sourceName := deriveSourceName(repo)
	result.SourceName = sourceName

	entry := SkillEntry{
		Name:          skillName,
		Path:          path,
		Description:   description,
		DescriptionZh: descriptionZh,
		License:       license,
	}

	if src, exists := ur.Sources[sourceName]; exists {
		src.Skills = append(src.Skills, entry)
	} else {
		ur.Sources[sourceName] = &SourceEntry{
			Repo:    repo,
			License: license,
			Skills:  []SkillEntry{entry},
		}
		ur.orderedKeys = append(ur.orderedKeys, sourceName)
	}

	return result, ur.Save()
}

// Remove deletes a skill entry by name. Returns an error if not found.
func (ur *UserRegistry) Remove(skillName string) error {
	nameLower := strings.ToLower(skillName)
	for srcKey, src := range ur.Sources {
		for i, s := range src.Skills {
			if strings.ToLower(s.Name) == nameLower {
				src.Skills = append(src.Skills[:i], src.Skills[i+1:]...)
				// Remove empty sources.
				if len(src.Skills) == 0 {
					delete(ur.Sources, srcKey)
					ur.orderedKeys = removeKey(ur.orderedKeys, srcKey)
				}
				return ur.Save()
			}
		}
	}
	return fmt.Errorf("skill %q not found in user registry", skillName)
}

// ListAll returns all skills as a flat list with their source name attached.
type ListedSkill struct {
	SourceName string
	Repo       string
	SkillEntry
}

func (ur *UserRegistry) ListAll() []ListedSkill {
	var out []ListedSkill
	for _, key := range ur.orderedKeys {
		src := ur.Sources[key]
		for _, s := range src.Skills {
			out = append(out, ListedSkill{
				SourceName: key,
				Repo:       src.Repo,
				SkillEntry: s,
			})
		}
	}
	return out
}

// IsEmpty returns true when the user registry has no skills.
func (ur *UserRegistry) IsEmpty() bool {
	for _, src := range ur.Sources {
		if len(src.Skills) > 0 {
			return false
		}
	}
	return true
}

// TotalSkillCount returns the total number of skills in the user registry.
func (ur *UserRegistry) TotalSkillCount() int {
	count := 0
	for _, src := range ur.Sources {
		count += len(src.Skills)
	}
	return count
}

// deriveSourceName creates a short, filesystem-safe source key from a repo string.
func deriveSourceName(repo string) string {
	// github.com/owner/repo-name → owner-repo-name
	if strings.HasPrefix(repo, "github.com/") {
		parts := strings.SplitN(strings.TrimPrefix(repo, "github.com/"), "/", 2)
		if len(parts) == 2 {
			return strings.ReplaceAll(parts[0]+"-"+parts[1], "/", "-")
		}
	}
	// local path → "local"
	if strings.HasPrefix(repo, "/") || strings.HasPrefix(repo, "./") || strings.HasPrefix(repo, "~/") {
		return "local"
	}
	// fallback: sanitise
	safe := strings.NewReplacer("/", "-", ".", "-", ":", "-").Replace(repo)
	return safe
}

func removeKey(keys []string, target string) []string {
	out := make([]string, 0, len(keys)-1)
	for _, k := range keys {
		if k != target {
			out = append(out, k)
		}
	}
	return out
}
