package registry

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadWithUserMergesSources(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)

	userFile := filepath.Join(configDir, "skills-x", "user-registry.yaml")
	if err := os.MkdirAll(filepath.Dir(userFile), 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	userYAML := `
my-team:
  repo: github.com/my-org/skills
  license: MIT
  skills:
    - name: my-skill
      path: skills/my-skill
      description: my skill
`
	if err := os.WriteFile(userFile, []byte(strings.TrimSpace(userYAML)+"\n"), 0644); err != nil {
		t.Fatalf("write user registry failed: %v", err)
	}

	reg, warnings, err := LoadWithUser()
	if err != nil {
		t.Fatalf("LoadWithUser failed: %v", err)
	}

	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", warnings)
	}

	src := reg.GetSource("user:my-team")
	if src == nil {
		t.Fatalf("expected merged source user:my-team")
	}
	if !src.IsUserSource() {
		t.Fatalf("expected source to be marked as user source")
	}
	if len(src.Skills) != 1 || src.Skills[0].Name != "my-skill" {
		t.Fatalf("unexpected user source skills: %+v", src.Skills)
	}
}

func TestLoadWithUserOverridesBuiltinsOnFindSkill(t *testing.T) {
	configDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configDir)

	userFile := filepath.Join(configDir, "skills-x", "user-registry.yaml")
	if err := os.MkdirAll(filepath.Dir(userFile), 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	// brainstorming exists in built-in registry and is used here for conflict testing.
	userYAML := `
override-source:
  repo: github.com/my-org/skills
  license: MIT
  skills:
    - name: brainstorming
      path: skills/brainstorming
      description: overridden brainstorming
`
	if err := os.WriteFile(userFile, []byte(strings.TrimSpace(userYAML)+"\n"), 0644); err != nil {
		t.Fatalf("write user registry failed: %v", err)
	}

	reg, warnings, err := LoadWithUser()
	if err != nil {
		t.Fatalf("LoadWithUser failed: %v", err)
	}
	if len(warnings) == 0 {
		t.Fatalf("expected conflict warning when overriding built-in skill")
	}

	skill, source := reg.FindSkill("brainstorming")
	if skill == nil || source == nil {
		t.Fatalf("expected to find brainstorming")
	}
	if !source.IsUserSource() {
		t.Fatalf("expected user source to override built-in source, got %q", source.Name)
	}
	if skill.Description != "overridden brainstorming" {
		t.Fatalf("expected overridden description, got %q", skill.Description)
	}
}
