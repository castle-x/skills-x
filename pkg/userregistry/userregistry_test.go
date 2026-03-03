package userregistry

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setTempConfigDir redirects FilePath() to a temp directory for test isolation.
// Returns a cleanup function.
func setTempConfigDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	return tmpDir
}

// ---------------------------------------------------------------------------
// FilePath — sanity check
// ---------------------------------------------------------------------------

func TestFilePath(t *testing.T) {
	tmpDir := setTempConfigDir(t)
	got := FilePath()
	want := filepath.Join(tmpDir, "skills-x", "user-registry.yaml")
	if got != want {
		t.Errorf("FilePath() = %q; want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// Load — empty / round-trip
// ---------------------------------------------------------------------------

func TestLoad_EmptyReturnsEmptyRegistry(t *testing.T) {
	setTempConfigDir(t)

	ur, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !ur.IsEmpty() {
		t.Error("new registry should be empty")
	}
	if ur.TotalSkillCount() != 0 {
		t.Errorf("TotalSkillCount = %d; want 0", ur.TotalSkillCount())
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	setTempConfigDir(t)

	ur, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	_, err = ur.Add(
		"github.com/affaan-m/everything-claude-code",
		"skills/golang-testing",
		"golang-testing",
		"Go testing patterns",
		"Go 测试模式",
		"MIT",
		nil,
	)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	ur2, err := Load()
	if err != nil {
		t.Fatalf("second Load: %v", err)
	}

	if ur2.IsEmpty() {
		t.Fatal("registry should not be empty after load")
	}
	if ur2.TotalSkillCount() != 1 {
		t.Errorf("TotalSkillCount = %d; want 1", ur2.TotalSkillCount())
	}

	listed := ur2.ListAll()
	if len(listed) != 1 {
		t.Fatalf("ListAll returned %d items; want 1", len(listed))
	}
	if listed[0].Name != "golang-testing" {
		t.Errorf("Name = %q; want %q", listed[0].Name, "golang-testing")
	}
	if listed[0].Description != "Go testing patterns" {
		t.Errorf("Description = %q; want %q", listed[0].Description, "Go testing patterns")
	}
}

// ---------------------------------------------------------------------------
// Add — source naming, multi-skill, conflict detection
// ---------------------------------------------------------------------------

func TestAdd_DeriveSourceName(t *testing.T) {
	setTempConfigDir(t)

	tests := []struct {
		name       string
		repo       string
		wantSource string
	}{
		{
			name:       "github repo",
			repo:       "github.com/affaan-m/everything-claude-code",
			wantSource: "affaan-m-everything-claude-code",
		},
		{
			name:       "local path",
			repo:       "/home/user/my-skills",
			wantSource: "local",
		},
		{
			name:       "relative path",
			repo:       "./my-skills",
			wantSource: "local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTempConfigDir(t) // fresh dir for each subtest
			ur, _ := Load()
			result, err := ur.Add(tt.repo, "test/path", "test-skill-"+strings.ReplaceAll(tt.name, " ", "-"),
				"desc", "", "", nil)
			if err != nil {
				t.Fatalf("Add: %v", err)
			}
			if result.SourceName != tt.wantSource {
				t.Errorf("SourceName = %q; want %q", result.SourceName, tt.wantSource)
			}
		})
	}
}

func TestAdd_MultipleSkillsSameSource(t *testing.T) {
	setTempConfigDir(t)

	ur, _ := Load()
	repo := "github.com/affaan-m/everything-claude-code"

	_, err := ur.Add(repo, "skills/golang-testing", "golang-testing", "Go tests", "", "MIT", nil)
	if err != nil {
		t.Fatalf("first Add: %v", err)
	}

	_, err = ur.Add(repo, "skills/react-best-practices", "react-best-practices", "React patterns", "", "MIT", nil)
	if err != nil {
		t.Fatalf("second Add: %v", err)
	}

	if ur.TotalSkillCount() != 2 {
		t.Errorf("TotalSkillCount = %d; want 2", ur.TotalSkillCount())
	}

	// Verify both skills are under the same source.
	src, ok := ur.Sources["affaan-m-everything-claude-code"]
	if !ok {
		t.Fatal("expected source 'affaan-m-everything-claude-code' to exist")
	}
	if len(src.Skills) != 2 {
		t.Errorf("source has %d skills; want 2", len(src.Skills))
	}
}

func TestAdd_DetectsDuplicate(t *testing.T) {
	setTempConfigDir(t)

	ur, _ := Load()
	repo := "github.com/affaan-m/everything-claude-code"

	_, err := ur.Add(repo, "skills/golang-testing", "golang-testing", "Go tests", "", "", nil)
	if err != nil {
		t.Fatalf("first Add: %v", err)
	}

	_, err = ur.Add(repo, "skills/golang-testing", "golang-testing", "Go tests dup", "", "", nil)
	if err == nil {
		t.Error("expected error for duplicate skill name")
	}
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error = %q; should mention 'already exists'", err.Error())
	}
}

func TestAdd_DuplicateCaseInsensitive(t *testing.T) {
	setTempConfigDir(t)

	ur, _ := Load()
	repo := "github.com/affaan-m/everything-claude-code"

	_, _ = ur.Add(repo, "skills/golang-testing", "golang-testing", "desc", "", "", nil)

	_, err := ur.Add(repo, "skills/foo", "Golang-Testing", "desc2", "", "", nil)
	if err == nil {
		t.Error("expected error for case-insensitive duplicate")
	}
}

func TestAdd_ConflictWithBuiltin(t *testing.T) {
	setTempConfigDir(t)

	ur, _ := Load()
	builtinNames := map[string][]string{
		"golang-testing": {"anthropic"},
	}

	result, err := ur.Add(
		"github.com/affaan-m/everything-claude-code",
		"skills/golang-testing",
		"golang-testing",
		"desc", "", "",
		builtinNames,
	)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if len(result.ConflictSources) == 0 {
		t.Error("expected conflict warning with builtin registry")
	}
	if result.ConflictSources[0] != "anthropic" {
		t.Errorf("ConflictSources = %v; want [anthropic]", result.ConflictSources)
	}

	// Skill should still be added (user takes precedence).
	if ur.TotalSkillCount() != 1 {
		t.Error("skill should be added despite conflict")
	}
}

// ---------------------------------------------------------------------------
// Remove — basic / not found / empty source cleanup
// ---------------------------------------------------------------------------

func TestRemove_Basic(t *testing.T) {
	setTempConfigDir(t)

	ur, _ := Load()
	repo := "github.com/affaan-m/everything-claude-code"
	ur.Add(repo, "skills/golang-testing", "golang-testing", "desc", "", "", nil)

	if err := ur.Remove("golang-testing"); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	if !ur.IsEmpty() {
		t.Error("registry should be empty after removing the only skill")
	}

	// Verify source key is also removed.
	if _, ok := ur.Sources["affaan-m-everything-claude-code"]; ok {
		t.Error("empty source should be removed from Sources map")
	}
}

func TestRemove_CaseInsensitive(t *testing.T) {
	setTempConfigDir(t)

	ur, _ := Load()
	ur.Add("github.com/affaan-m/everything-claude-code", "skills/golang-testing",
		"golang-testing", "desc", "", "", nil)

	if err := ur.Remove("Golang-Testing"); err != nil {
		t.Fatalf("Remove (case-insensitive): %v", err)
	}

	if !ur.IsEmpty() {
		t.Error("registry should be empty after case-insensitive remove")
	}
}

func TestRemove_NotFound(t *testing.T) {
	setTempConfigDir(t)

	ur, _ := Load()
	err := ur.Remove("nonexistent-skill")
	if err == nil {
		t.Error("expected error for removing nonexistent skill")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q; should mention 'not found'", err.Error())
	}
}

func TestRemove_KeepsOtherSkills(t *testing.T) {
	setTempConfigDir(t)

	ur, _ := Load()
	repo := "github.com/affaan-m/everything-claude-code"
	ur.Add(repo, "skills/golang-testing", "golang-testing", "desc1", "", "", nil)
	ur.Add(repo, "skills/react-patterns", "react-patterns", "desc2", "", "", nil)

	if err := ur.Remove("golang-testing"); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	if ur.TotalSkillCount() != 1 {
		t.Errorf("TotalSkillCount = %d; want 1", ur.TotalSkillCount())
	}

	listed := ur.ListAll()
	if listed[0].Name != "react-patterns" {
		t.Errorf("remaining skill = %q; want %q", listed[0].Name, "react-patterns")
	}

	// Source should still exist (still has skills).
	if _, ok := ur.Sources["affaan-m-everything-claude-code"]; !ok {
		t.Error("source should still exist with remaining skills")
	}
}

// ---------------------------------------------------------------------------
// ListAll — ordering preservation
// ---------------------------------------------------------------------------

func TestListAll_PreservesInsertionOrder(t *testing.T) {
	setTempConfigDir(t)

	ur, _ := Load()
	repos := []struct {
		repo  string
		name  string
		skill string
	}{
		{"github.com/alice/repo-a", "skills/alpha", "alpha"},
		{"github.com/bob/repo-b", "skills/beta", "beta"},
		{"github.com/carol/repo-c", "skills/gamma", "gamma"},
	}

	for _, r := range repos {
		ur.Add(r.repo, r.name, r.skill, "desc", "", "", nil)
	}

	listed := ur.ListAll()
	if len(listed) != 3 {
		t.Fatalf("ListAll returned %d items; want 3", len(listed))
	}

	for i, r := range repos {
		if listed[i].Name != r.skill {
			t.Errorf("listed[%d].Name = %q; want %q", i, listed[i].Name, r.skill)
		}
	}
}

// ---------------------------------------------------------------------------
// IsEmpty / TotalSkillCount — edge cases
// ---------------------------------------------------------------------------

func TestIsEmpty_WithEmptySource(t *testing.T) {
	ur := &UserRegistry{
		Sources: map[string]*SourceEntry{
			"test": {Repo: "github.com/test/repo", Skills: []SkillEntry{}},
		},
		orderedKeys: []string{"test"},
	}
	if !ur.IsEmpty() {
		t.Error("registry with zero-skill source should be empty")
	}
	if ur.TotalSkillCount() != 0 {
		t.Errorf("TotalSkillCount = %d; want 0", ur.TotalSkillCount())
	}
}

// ---------------------------------------------------------------------------
// deriveSourceName — unit tests
// ---------------------------------------------------------------------------

func TestDeriveSourceName(t *testing.T) {
	tests := []struct {
		repo string
		want string
	}{
		{"github.com/affaan-m/everything-claude-code", "affaan-m-everything-claude-code"},
		{"github.com/anthropics/courses", "anthropics-courses"},
		{"/home/user/skills", "local"},
		{"./my-skills", "local"},
		{"~/my-skills", "local"},
		{"custom.gitlab.com/foo/bar", "custom-gitlab-com-foo-bar"},
	}

	for _, tt := range tests {
		t.Run(tt.repo, func(t *testing.T) {
			got := deriveSourceName(tt.repo)
			if got != tt.want {
				t.Errorf("deriveSourceName(%q) = %q; want %q", tt.repo, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Persistence — YAML round-trip integrity
// ---------------------------------------------------------------------------

func TestPersistence_YAMLIntegrity(t *testing.T) {
	setTempConfigDir(t)

	ur, _ := Load()
	repo := "github.com/affaan-m/everything-claude-code"
	ur.Add(repo, "skills/golang-testing", "golang-testing", "Go tests", "Go 测试", "MIT", nil)
	ur.Add(repo, "skills/react-patterns", "react-patterns", "React", "", "Apache-2.0", nil)

	// Read raw YAML and verify structure.
	data, err := os.ReadFile(FilePath())
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "golang-testing") {
		t.Error("YAML should contain golang-testing")
	}
	if !strings.Contains(content, "react-patterns") {
		t.Error("YAML should contain react-patterns")
	}
	if !strings.Contains(content, "affaan-m-everything-claude-code") {
		t.Error("YAML should contain derived source key")
	}
	if !strings.Contains(content, "Go 测试") {
		t.Error("YAML should preserve Chinese description")
	}

	// Reload and verify full round-trip.
	ur2, err := Load()
	if err != nil {
		t.Fatalf("Load after save: %v", err)
	}
	if ur2.TotalSkillCount() != 2 {
		t.Errorf("round-trip TotalSkillCount = %d; want 2", ur2.TotalSkillCount())
	}

	listed := ur2.ListAll()
	for _, s := range listed {
		if s.Name == "golang-testing" {
			if s.DescriptionZh != "Go 测试" {
				t.Errorf("DescriptionZh = %q; want %q", s.DescriptionZh, "Go 测试")
			}
			if s.License != "MIT" {
				t.Errorf("License = %q; want %q", s.License, "MIT")
			}
		}
	}
}
