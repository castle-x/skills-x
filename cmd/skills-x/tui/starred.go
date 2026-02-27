// Package tui provides terminal interactive UI components
package tui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

// starredFilePath returns the path to the starred skills config file.
func starredFilePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = os.Getenv("HOME")
		if dir == "" {
			dir = "."
		}
		dir = filepath.Join(dir, ".config")
	}
	return filepath.Join(dir, "skills-x", "starred.json")
}

// LoadStarred reads the persisted starred set from disk.
// Returns an empty map if the file does not exist or cannot be parsed.
func LoadStarred() map[string]bool {
	data, err := os.ReadFile(starredFilePath())
	if err != nil {
		return map[string]bool{}
	}
	var list []string
	if err := json.Unmarshal(data, &list); err != nil {
		return map[string]bool{}
	}
	set := make(map[string]bool, len(list))
	for _, name := range list {
		set[name] = true
	}
	return set
}

// SaveStarred writes the starred set to disk as a sorted JSON array.
func SaveStarred(set map[string]bool) error {
	list := make([]string, 0, len(set))
	for name := range set {
		list = append(list, name)
	}
	sort.Strings(list)

	data, err := json.Marshal(list)
	if err != nil {
		return err
	}

	path := starredFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
