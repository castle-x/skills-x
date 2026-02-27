package tui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const metaFileName = ".skills-x-meta.json"

// SkillMeta stores metadata about an installed skill
type SkillMeta struct {
	Skill       string `json:"skill"`
	Source      string `json:"source"`
	Repo        string `json:"repo"`
	Commit      string `json:"commit"`
	InstalledAt string `json:"installed_at"`
}

// WriteSkillMeta writes meta to .skills-x-meta.json inside the skill directory
func WriteSkillMeta(skillDir string, meta SkillMeta) error {
	if meta.InstalledAt == "" {
		meta.InstalledAt = time.Now().UTC().Format(time.RFC3339)
	}
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(skillDir, metaFileName), data, 0644)
}

// ReadSkillMeta reads meta from .skills-x-meta.json inside the skill directory
func ReadSkillMeta(skillDir string) (*SkillMeta, error) {
	data, err := os.ReadFile(filepath.Join(skillDir, metaFileName))
	if err != nil {
		return nil, err
	}
	var meta SkillMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}
	return &meta, nil
}
