// Package tui provides terminal interactive UI components
package tui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	xskills "github.com/castle-x/skills-x/cmd/skills-x/skills"
	"github.com/castle-x/skills-x/pkg/registry"
)

// LoadSkillsFromRegistry loads skills from registry and checks installed status
// targetDir: the directory to check for installation status
func LoadSkillsFromRegistry(targetDir string) ([]SkillItem, error) {
	reg, err := registry.Load()
	if err != nil {
		return nil, err
	}

	var skills []SkillItem

	// Load skills from each source
	for _, source := range reg.Sources {
		for _, skill := range source.Skills {
			fullName := source.Name + "/" + skill.Name

			installed := targetDir != "" && isSkillDir(filepath.Join(targetDir, skill.Name))

			description := skill.GetDescription("")

			skills = append(skills, SkillItem{
				Name:        skill.Name,
				FullName:    fullName,
				Source:      source.Repo,
				SourceName:  source.Name,
				Description: description,
				Installed:   installed,
			})
		}
	}

	// Load X (self-developed) skills
	xSkills, err := xskills.ListXSkills()
	if err == nil {
		for _, xSkill := range xSkills {
			installed := targetDir != "" && isSkillDir(filepath.Join(targetDir, xSkill.Name))

			skills = append(skills, SkillItem{
				Name:        xSkill.Name,
				FullName:    "skills-x/" + xSkill.Name,
				Source:      "skills-x",
				SourceName:  "skills-x",
				Description: xSkill.Description,
				Installed:   installed,
				IsX:         true,
			})
		}
	}

	// Sort by FullName (source/skill-name format)
	sort.Slice(skills, func(i, j int) bool {
		return strings.ToLower(skills[i].FullName) < strings.ToLower(skills[j].FullName)
	})

	return skills, nil
}

// isSkillDir checks if a path is a valid skill directory.
// Follows symlinks and requires SKILL.md to exist inside.
func isSkillDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	if _, err := os.Stat(filepath.Join(path, "SKILL.md")); err == nil {
		return true
	}
	return false
}

// CheckSkillInstalled checks if a skill is installed in the target directory
func CheckSkillInstalled(skillName, targetDir string) bool {
	if targetDir == "" {
		return false
	}
	return isSkillDir(filepath.Join(targetDir, skillName))
}

// GetSkillByFullName finds a skill by its full name (source/skill-name)
func GetSkillByFullName(skills []SkillItem, fullName string) *SkillItem {
	for i := range skills {
		if skills[i].FullName == fullName {
			return &skills[i]
		}
	}
	return nil
}

// SortSkills sorts skills by FullName
func SortSkills(skills []SkillItem) {
	sort.Slice(skills, func(i, j int) bool {
		return strings.ToLower(skills[i].FullName) < strings.ToLower(skills[j].FullName)
	})
}

// FilterSkills filters skills based on search query
func FilterSkills(skills []SkillItem, query string) []SkillItem {
	if query == "" {
		return skills
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var filtered []SkillItem

	for _, s := range skills {
		if strings.Contains(strings.ToLower(s.FullName), query) ||
			strings.Contains(strings.ToLower(s.SourceName), query) ||
			strings.Contains(strings.ToLower(s.Name), query) ||
			strings.Contains(strings.ToLower(s.Description), query) {
			filtered = append(filtered, s)
		}
	}

	return filtered
}
