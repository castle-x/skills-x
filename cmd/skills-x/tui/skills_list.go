// Package tui provides terminal interactive UI components
package tui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
)

// LoadSkillsFromRegistry loads skills from registry and checks installed status
// targetDir: the directory to check for installation status
func LoadSkillsFromRegistry(targetDir string) ([]SkillItem, error) {
	reg, err := loadMergedRegistry()
	if err != nil {
		return nil, err
	}

	starredSet := LoadStarred()

	var skills []SkillItem

	// Load skills from each source
	for _, source := range reg.Sources {
		for _, skill := range source.Skills {
			fullName := source.Name + "/" + skill.Name

			skillDir := filepath.Join(targetDir, skill.Name)
			installed := targetDir != "" && isSkillDir(skillDir)

			description := skill.GetDescription(i18n.GetLanguage())

			item := SkillItem{
				Name:        skill.Name,
				FullName:    fullName,
				Source:      source.Repo,
				SourceName:  source.Name,
				Description: description,
				Tags:        skill.Tags,
				Installed:   installed,
				Starred:     starredSet[fullName],
			}
			if installed {
				item.Meta, _ = ReadSkillMeta(skillDir)
			}
			skills = append(skills, item)
		}
	}

	// Sort: starred first, then alphabetical by FullName
	sort.Slice(skills, func(i, j int) bool {
		if skills[i].Starred != skills[j].Starred {
			return skills[i].Starred
		}
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

// SortSkills sorts skills: starred first, then alphabetical by FullName
func SortSkills(skills []SkillItem) {
	sort.Slice(skills, func(i, j int) bool {
		if skills[i].Starred != skills[j].Starred {
			return skills[i].Starred
		}
		return strings.ToLower(skills[i].FullName) < strings.ToLower(skills[j].FullName)
	})
}

// FilterSkills filters skills based on search query
// Supports #tag prefix for tag-based filtering
func FilterSkills(skills []SkillItem, query string) []SkillItem {
	if query == "" {
		return skills
	}

	query = strings.TrimSpace(query)
	var filtered []SkillItem

	if strings.HasPrefix(query, "#") {
		tagQuery := strings.ToLower(query[1:])
		if tagQuery == "" {
			return skills
		}
		for _, s := range skills {
			for _, t := range s.Tags {
				if strings.ToLower(t) == tagQuery {
					filtered = append(filtered, s)
					break
				}
			}
		}
	} else {
		queryLower := strings.ToLower(query)
		for _, s := range skills {
			if strings.Contains(strings.ToLower(s.FullName), queryLower) ||
				strings.Contains(strings.ToLower(s.SourceName), queryLower) ||
				strings.Contains(strings.ToLower(s.Name), queryLower) ||
				strings.Contains(strings.ToLower(s.Description), queryLower) {
				filtered = append(filtered, s)
			}
		}
	}

	return filtered
}
