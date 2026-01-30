// Package skills provides x (self-developed) skills data
package skills

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
)

// go:embed directive embeds the x/ directory (synced by Makefile before build)
// Note: embed.FS always uses "/" as path separator, regardless of OS (Windows compatible)
//
//go:embed all:x
var xFS embed.FS

// SkillInfo holds metadata about a skill
type SkillInfo struct {
	Name        string
	Description string
	IsX         bool   // true if from x (self-developed)
	Path        string // path in embed.FS (always uses "/" separator)
}

// xSkillDescriptions provides descriptions for x self-developed skills
// (kept for backward compatibility, but descriptions now come from i18n)
var xSkillDescriptions = map[string]string{
	"skills-x": "Contribute skills to skills-x collection", // Deprecated: use i18n instead
}

// GetXFS returns the embedded filesystem for x skills
func GetXFS() embed.FS {
	return xFS
}

// ListXSkills returns all x (self-developed) skills with metadata
func ListXSkills() ([]SkillInfo, error) {
	var skills []SkillInfo

	// The embedded filesystem has "x" as root directory
	// Structure: x/skill-name/SKILL.md
	// Note: embed.FS always uses "/" as separator, regardless of OS
	const xRoot = "x"

	// Walk through x skills directory
	err := fs.WalkDir(xFS, xRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip root "x" directory and non-directories
		if path == xRoot || !d.IsDir() {
			return nil
		}

		// Only process direct children of "x" (skill folders)
		// Skip subdirectories within skills
		// Note: must use "/" for embed.FS, not filepath.Join (which uses "\" on Windows)
		rel := strings.TrimPrefix(path, xRoot+"/")
		if strings.Contains(rel, "/") {
			return fs.SkipDir
		}

		name := d.Name()

		// Check if SKILL.md exists
		// Note: must use "/" for embed.FS paths, not filepath.Join
		skillMdPath := path + "/SKILL.md"
		if _, err := xFS.Open(skillMdPath); err != nil {
			return nil // Skip directories without SKILL.md
		}

		// Get description from i18n system
		descKey := "skill_" + name
		description := i18n.T(descKey)

		info := SkillInfo{
			Name:        name,
			Description: description,
			IsX:         true,
			Path:        path, // Full path including "x/" prefix, uses "/" separator
		}

		skills = append(skills, info)
		return nil
	})

	return skills, err
}

// GetXSkill returns a specific x skill by name
func GetXSkill(name string) (*SkillInfo, error) {
	skills, err := ListXSkills()
	if err != nil {
		return nil, err
	}

	for _, s := range skills {
		if s.Name == name {
			return &s, nil
		}
	}

	return nil, nil
}

// XSkillExists checks if an x skill exists
func XSkillExists(name string) bool {
	skill, _ := GetXSkill(name)
	return skill != nil
}

// GetXSkillFS returns the filesystem and path for an x skill
func GetXSkillFS(name string) (embed.FS, string, bool) {
	skill, err := GetXSkill(name)
	if err != nil || skill == nil {
		return xFS, "", false
	}
	return xFS, skill.Path, true
}
