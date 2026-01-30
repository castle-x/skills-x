// Package discover provides skill discovery from cloned repositories
package discover

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const (
	// SkillFileName is the name of the skill definition file
	SkillFileName = "SKILL.md"
	// MaxSearchDepth is the maximum depth for recursive search
	MaxSearchDepth = 5
)

// SkipDirs are directories to skip during discovery
var SkipDirs = map[string]bool{
	"node_modules":   true,
	".git":           true,
	"dist":           true,
	"build":          true,
	"__pycache__":    true,
	".venv":          true,
	"venv":           true,
	".next":          true,
	".nuxt":          true,
	"coverage":       true,
	".turbo":         true,
	".cache":         true,
}

// PriorityDirs are directories to search first (common skill locations)
var PriorityDirs = []string{
	"skills",
	".agent/skills",
	".agents/skills",
	".cursor/skills",
	".windsurf/skills",
	".codebuddy/skills",
	"packages/docs/skills", // remotion
}

// DiscoveredSkill represents a discovered skill
type DiscoveredSkill struct {
	Name        string // Skill name (from frontmatter or directory name)
	Description string // Description (from frontmatter)
	Version     string // Version (from frontmatter, optional)
	Path        string // Absolute path to the skill directory
	SkillMdPath string // Absolute path to SKILL.md
}

// DiscoverOptions configures skill discovery
type DiscoverOptions struct {
	IncludeInternal bool // Include internal skills
	FullDepth       bool // Search all subdirectories even if found in root
}

// DiscoverSkills discovers all skills in a repository
func DiscoverSkills(basePath string, opts *DiscoverOptions) ([]DiscoveredSkill, error) {
	if opts == nil {
		opts = &DiscoverOptions{}
	}

	var skills []DiscoveredSkill
	seen := make(map[string]bool)

	// First, check if basePath itself is a skill
	if isSkillDir(basePath) {
		skill, err := parseSkillDir(basePath)
		if err == nil && skill != nil {
			skills = append(skills, *skill)
			seen[skill.Name] = true
			if !opts.FullDepth {
				return skills, nil
			}
		}
	}

	// Search priority directories first
	for _, dir := range PriorityDirs {
		priorityPath := filepath.Join(basePath, dir)
		if dirExists(priorityPath) {
			found, err := searchDir(priorityPath, 1, seen, opts)
			if err == nil {
				skills = append(skills, found...)
			}
		}
	}

	// If no skills found in priority dirs, do recursive search
	if len(skills) == 0 {
		found, err := searchDir(basePath, 0, seen, opts)
		if err == nil {
			skills = append(skills, found...)
		}
	}

	return skills, nil
}

// DiscoverSkillByPath discovers a skill at a specific path
func DiscoverSkillByPath(basePath string, skillPath string) (*DiscoveredSkill, error) {
	fullPath := filepath.Join(basePath, skillPath)
	
	if !isSkillDir(fullPath) {
		// Try direct path (skillPath might already be the full path)
		if isSkillDir(skillPath) {
			fullPath = skillPath
		} else {
			return nil, nil
		}
	}

	return parseSkillDir(fullPath)
}

// searchDir recursively searches for skills in a directory
func searchDir(dir string, depth int, seen map[string]bool, opts *DiscoverOptions) ([]DiscoveredSkill, error) {
	if depth > MaxSearchDepth {
		return nil, nil
	}

	var skills []DiscoveredSkill

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if SkipDirs[name] {
			continue
		}

		subdir := filepath.Join(dir, name)

		// Check if this directory is a skill
		if isSkillDir(subdir) {
			skill, err := parseSkillDir(subdir)
			if err == nil && skill != nil {
				if !seen[skill.Name] {
					// Filter internal skills if not requested
					if !opts.IncludeInternal || !isInternalSkill(skill) {
						skills = append(skills, *skill)
						seen[skill.Name] = true
					}
				}
			}
		} else {
			// Recurse into subdirectory
			found, err := searchDir(subdir, depth+1, seen, opts)
			if err == nil {
				skills = append(skills, found...)
			}
		}
	}

	return skills, nil
}

// isSkillDir checks if a directory contains SKILL.md
func isSkillDir(dir string) bool {
	skillMdPath := filepath.Join(dir, SkillFileName)
	info, err := os.Stat(skillMdPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// parseSkillDir parses a skill directory
func parseSkillDir(dir string) (*DiscoveredSkill, error) {
	skillMdPath := filepath.Join(dir, SkillFileName)

	// Read and parse SKILL.md
	file, err := os.Open(skillMdPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	frontmatter := parseFrontmatter(file)

	// Get name from frontmatter or directory name
	name := frontmatter["name"]
	if name == "" {
		name = filepath.Base(dir)
	}

	skill := &DiscoveredSkill{
		Name:        name,
		Description: frontmatter["description"],
		Version:     frontmatter["version"],
		Path:        dir,
		SkillMdPath: skillMdPath,
	}

	return skill, nil
}

// parseFrontmatter extracts YAML frontmatter from SKILL.md
func parseFrontmatter(file *os.File) map[string]string {
	result := make(map[string]string)
	scanner := bufio.NewScanner(file)

	// Check for opening ---
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return result
	}

	// Read frontmatter lines until closing ---
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}

		// Parse key: value
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			result[key] = value
		}
	}

	return result
}

// isInternalSkill checks if a skill is marked as internal
func isInternalSkill(skill *DiscoveredSkill) bool {
	// Read the file to check metadata.internal
	// For simplicity, we skip this for now
	return false
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
