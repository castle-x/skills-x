package skillvalidator

import (
	"os"
	"path/filepath"
	"testing"
)

// ---------------------------------------------------------------------------
// ParseInput — table-driven, pure unit tests (no network)
// ---------------------------------------------------------------------------

func TestParseInput(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantKind  InputKind
		wantRepo  string
		wantHint  string
	}{
		{
			name:     "owner/repo → repo scan",
			input:    "affaan-m/everything-claude-code",
			wantKind: InputKindRepoScan,
			wantRepo: "github.com/affaan-m/everything-claude-code",
		},
		{
			name:     "owner/repo/skill → single skill",
			input:    "affaan-m/everything-claude-code/golang-testing",
			wantKind: InputKindSingleSkill,
			wantRepo: "github.com/affaan-m/everything-claude-code",
			wantHint: "golang-testing",
		},
		{
			name:     "github.com prefix compat",
			input:    "github.com/affaan-m/everything-claude-code",
			wantKind: InputKindRepoScan,
			wantRepo: "github.com/affaan-m/everything-claude-code",
		},
		{
			name:     "github.com prefix with skill path",
			input:    "github.com/affaan-m/everything-claude-code/skills/golang-testing",
			wantKind: InputKindSingleSkill,
			wantRepo: "github.com/affaan-m/everything-claude-code",
			wantHint: "skills/golang-testing",
		},
		{
			name:     "https URL stripped",
			input:    "https://github.com/affaan-m/everything-claude-code",
			wantKind: InputKindRepoScan,
			wantRepo: "github.com/affaan-m/everything-claude-code",
		},
		{
			name:     "https URL with tree/main path",
			input:    "https://github.com/affaan-m/everything-claude-code/tree/main/skills/golang-testing",
			wantKind: InputKindSingleSkill,
			wantRepo: "github.com/affaan-m/everything-claude-code",
			wantHint: "skills/golang-testing",
		},
		{
			name:     "absolute local path",
			input:    "/home/user/skills/my-skill",
			wantKind: InputKindLocal,
			wantRepo: "/home/user/skills/my-skill",
		},
		{
			name:     "relative local path",
			input:    "./skills/go-i18n",
			wantKind: InputKindLocal,
			wantRepo: "./skills/go-i18n",
		},
		{
			name:     "home-relative path",
			input:    "~/my-skills/test",
			wantKind: InputKindLocal,
			wantRepo: "~/my-skills/test",
		},
		{
			name:     "whitespace trimmed",
			input:    "  affaan-m/everything-claude-code  ",
			wantKind: InputKindRepoScan,
			wantRepo: "github.com/affaan-m/everything-claude-code",
		},
		{
			name:     "trailing slash stripped from hint",
			input:    "affaan-m/everything-claude-code/golang-testing/",
			wantKind: InputKindSingleSkill,
			wantRepo: "github.com/affaan-m/everything-claude-code",
			wantHint: "golang-testing",
		},
		{
			name:     "http prefix stripped",
			input:    "http://github.com/affaan-m/everything-claude-code",
			wantKind: InputKindRepoScan,
			wantRepo: "github.com/affaan-m/everything-claude-code",
		},
		{
			name:     "blob/main stripped from URL",
			input:    "https://github.com/affaan-m/everything-claude-code/blob/main/skills/golang-testing",
			wantKind: InputKindSingleSkill,
			wantRepo: "github.com/affaan-m/everything-claude-code",
			wantHint: "skills/golang-testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseInput(tt.input)

			if got.Kind != tt.wantKind {
				t.Errorf("Kind = %d; want %d", got.Kind, tt.wantKind)
			}
			if got.Repo != tt.wantRepo {
				t.Errorf("Repo = %q; want %q", got.Repo, tt.wantRepo)
			}
			if got.SkillHint != tt.wantHint {
				t.Errorf("SkillHint = %q; want %q", got.SkillHint, tt.wantHint)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Validate — local path validation (no network)
// ---------------------------------------------------------------------------

func TestValidate_LocalSkill(t *testing.T) {
	// Create a valid skill directory in a temp dir.
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, "my-skill")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("failed to create skill dir: %v", err)
	}

	skillMD := `---
name: my-skill
description: A test skill for validation
license: MIT
---

# My Skill

Test content.
`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillMD), 0644); err != nil {
		t.Fatalf("failed to write SKILL.md: %v", err)
	}

	t.Run("valid skill passes", func(t *testing.T) {
		result, err := Validate(ValidateRequest{Repo: skillDir})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Valid {
			t.Errorf("expected valid, got errors: %v", result.Errors)
		}
		if result.SkillName != "my-skill" {
			t.Errorf("SkillName = %q; want %q", result.SkillName, "my-skill")
		}
		if result.Description != "A test skill for validation" {
			t.Errorf("Description = %q; want %q", result.Description, "A test skill for validation")
		}
		if result.License != "MIT" {
			t.Errorf("License = %q; want %q", result.License, "MIT")
		}
	})

	t.Run("warns about missing LICENSE.txt", func(t *testing.T) {
		result, _ := Validate(ValidateRequest{Repo: skillDir})
		if len(result.Warnings) == 0 {
			t.Error("expected warning about missing LICENSE.txt")
		}
	})

	t.Run("LICENSE.txt suppresses warning", func(t *testing.T) {
		licensePath := filepath.Join(skillDir, "LICENSE.txt")
		os.WriteFile(licensePath, []byte("MIT"), 0644)
		defer os.Remove(licensePath)

		result, _ := Validate(ValidateRequest{Repo: skillDir})
		for _, w := range result.Warnings {
			if w == "LICENSE.txt not found" {
				t.Error("LICENSE.txt warning should not appear")
			}
		}
	})
}

func TestValidate_InvalidSkill(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("missing SKILL.md", func(t *testing.T) {
		emptyDir := filepath.Join(tmpDir, "empty")
		os.MkdirAll(emptyDir, 0755)

		result, _ := Validate(ValidateRequest{Repo: emptyDir})
		if result.Valid {
			t.Error("expected invalid for missing SKILL.md")
		}
	})

	t.Run("no frontmatter", func(t *testing.T) {
		noFM := filepath.Join(tmpDir, "no-fm")
		os.MkdirAll(noFM, 0755)
		os.WriteFile(filepath.Join(noFM, "SKILL.md"), []byte("# Just a heading\nNo frontmatter here."), 0644)

		result, _ := Validate(ValidateRequest{Repo: noFM})
		if result.Valid {
			t.Error("expected invalid for missing frontmatter")
		}
	})

	t.Run("invalid name format", func(t *testing.T) {
		badName := filepath.Join(tmpDir, "bad-name")
		os.MkdirAll(badName, 0755)
		os.WriteFile(filepath.Join(badName, "SKILL.md"), []byte("---\nname: INVALID_NAME\ndescription: test\n---\n"), 0644)

		result, _ := Validate(ValidateRequest{Repo: badName})
		if result.Valid {
			t.Error("expected invalid for bad name format")
		}
	})

	t.Run("empty description", func(t *testing.T) {
		noDesc := filepath.Join(tmpDir, "no-desc")
		os.MkdirAll(noDesc, 0755)
		os.WriteFile(filepath.Join(noDesc, "SKILL.md"), []byte("---\nname: valid-name\ndescription: \"\"\n---\n"), 0644)

		result, _ := Validate(ValidateRequest{Repo: noDesc})
		if result.Valid {
			t.Error("expected invalid for empty description")
		}
	})

	t.Run("nonexistent path", func(t *testing.T) {
		result, _ := Validate(ValidateRequest{Repo: "/tmp/this-does-not-exist-12345"})
		if result.Valid {
			t.Error("expected invalid for nonexistent path")
		}
	})
}

// ---------------------------------------------------------------------------
// Discover — integration test (requires network, uses real repo)
// ---------------------------------------------------------------------------

func TestDiscover_RealRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	skills, err := Discover("github.com/affaan-m/everything-claude-code")
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	if len(skills) == 0 {
		t.Fatal("expected to discover at least one skill")
	}

	// This repo is known to have golang-testing.
	found := false
	for _, s := range skills {
		if s.Name == "golang-testing" {
			found = true
			if s.Description == "" {
				t.Error("golang-testing should have a description")
			}
			if !s.Valid {
				t.Errorf("golang-testing should be valid, got errors: %v", s.Errors)
			}
			break
		}
	}
	if !found {
		t.Error("expected to find golang-testing in discovered skills")
	}

	// Verify all discovered skills have at least a name.
	for _, s := range skills {
		if s.Name == "" {
			t.Errorf("skill at path %q has empty name", s.Path)
		}
	}
}

// ---------------------------------------------------------------------------
// FindSkill — integration test (requires network)
// ---------------------------------------------------------------------------

func TestFindSkill_RealRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	repo := "github.com/affaan-m/everything-claude-code"

	t.Run("find by name", func(t *testing.T) {
		ds, err := FindSkill(repo, "golang-testing")
		if err != nil {
			t.Fatalf("FindSkill failed: %v", err)
		}
		if ds == nil {
			t.Fatal("expected to find golang-testing, got nil")
		}
		if ds.Name != "golang-testing" {
			t.Errorf("Name = %q; want %q", ds.Name, "golang-testing")
		}
		if ds.Description == "" {
			t.Error("expected non-empty description")
		}
		if !ds.Valid {
			t.Errorf("expected valid, got errors: %v", ds.Errors)
		}
	})

	t.Run("find by explicit path", func(t *testing.T) {
		ds, err := FindSkill(repo, "skills/golang-testing")
		if err != nil {
			t.Fatalf("FindSkill failed: %v", err)
		}
		if ds == nil {
			t.Fatal("expected to find skill at explicit path")
		}
		if ds.Name != "golang-testing" {
			t.Errorf("Name = %q; want %q", ds.Name, "golang-testing")
		}
	})

	t.Run("nonexistent skill returns nil", func(t *testing.T) {
		ds, err := FindSkill(repo, "this-skill-does-not-exist-xyz")
		if err != nil {
			t.Fatalf("FindSkill failed: %v", err)
		}
		if ds != nil {
			t.Errorf("expected nil for nonexistent skill, got %+v", ds)
		}
	})
}

// ---------------------------------------------------------------------------
// parseFrontmatter — unit tests
// ---------------------------------------------------------------------------

func TestParseFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		wantName    string
		wantDesc    string
		wantLicense string
	}{
		{
			name:        "standard frontmatter",
			content:     "---\nname: test-skill\ndescription: A test\nlicense: MIT\n---\n# Content",
			wantName:    "test-skill",
			wantDesc:    "A test",
			wantLicense: "MIT",
		},
		{
			name:     "frontmatter with leading blank lines",
			content:  "\n\n---\nname: leading-blanks\ndescription: Has blanks\n---\n",
			wantName: "leading-blanks",
			wantDesc: "Has blanks",
		},
		{
			name:     "no frontmatter at all",
			content:  "# Just a heading\nSome content.",
			wantName: "",
			wantDesc: "",
		},
		{
			name:     "empty file",
			content:  "",
			wantName: "",
			wantDesc: "",
		},
		{
			name:     "multiline description",
			content:  "---\nname: multi\ndescription: \"Line one. Line two.\"\n---\n",
			wantName: "multi",
			wantDesc: "Line one. Line two.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name+".md")
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			fm, err := parseFrontmatter(path)
			if err != nil {
				t.Fatalf("parseFrontmatter error: %v", err)
			}
			if fm.Name != tt.wantName {
				t.Errorf("Name = %q; want %q", fm.Name, tt.wantName)
			}
			if fm.Description != tt.wantDesc {
				t.Errorf("Description = %q; want %q", fm.Description, tt.wantDesc)
			}
			if fm.License != tt.wantLicense {
				t.Errorf("License = %q; want %q", fm.License, tt.wantLicense)
			}
		})
	}
}
