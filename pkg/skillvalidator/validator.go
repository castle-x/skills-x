// Package skillvalidator validates skill directories against the Agent Skills specification.
// It is designed to be used by both the CLI registry command and the TUI.
package skillvalidator

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/castle-x/skills-x/pkg/gitutil"
	"gopkg.in/yaml.v3"
)

// SourceType indicates where the skill comes from.
type SourceType string

const (
	SourceTypeGitHub SourceType = "github"
	SourceTypeLocal  SourceType = "local"
)

// ValidateRequest describes what to validate.
type ValidateRequest struct {
	// Repo is one of:
	//   - "github.com/owner/repo"  → GitHub sparse clone
	//   - "/abs/path"              → local absolute path to skill dir
	//   - "./rel/path"             → local relative path (resolved against cwd)
	Repo string
	// Path is the skill sub-path inside the repository.
	// Empty when Repo is a local path pointing directly at the skill dir.
	Path string
}

// ValidateResult is the outcome of a validation run.
// Both CLI and TUI consume this struct directly.
type ValidateResult struct {
	// Parsed from SKILL.md frontmatter
	SkillName     string
	Description   string
	DescriptionZh string
	License       string

	SourceType   SourceType
	ResolvedPath string // local directory that was actually inspected

	Valid    bool
	Errors   []string
	Warnings []string
}

// skillFrontmatter is the minimal YAML structure we look for in SKILL.md.
type skillFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	License     string `yaml:"license"`
}

var validNameRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// InputKind classifies the user input.
type InputKind int

const (
	InputKindSingleSkill InputKind = iota // owner/repo/skill-name or explicit path
	InputKindRepoScan                     // owner/repo → discover all skills
	InputKindLocal                        // /abs/path or ./rel/path or ~/path
)

// ParsedInput is the result of ParseInput.
type ParsedInput struct {
	Kind      InputKind
	Repo      string // "github.com/owner/repo"
	SkillHint string // skill name to search for (empty in RepoScan mode)
	// For local paths, Repo is the path itself and SkillHint is empty.
}

// ParseInput interprets flexible user input into a structured form.
//
// Supported formats:
//
//	owner/repo               → RepoScan
//	owner/repo/skill-name    → SingleSkill (search by name in repo)
//	github.com/owner/repo    → RepoScan (compat)
//	github.com/owner/repo/x  → SingleSkill (compat)
//	https://github.com/...   → strip prefix, then parse
//	/abs/path                → Local
//	./rel/path               → Local
//	~/path                   → Local
func ParseInput(input string) ParsedInput {
	s := strings.TrimSpace(input)

	// Strip URL prefixes.
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")

	// Local paths.
	if strings.HasPrefix(s, "/") || strings.HasPrefix(s, "./") || strings.HasPrefix(s, "~/") {
		return ParsedInput{Kind: InputKindLocal, Repo: s}
	}

	// Strip "github.com/" if already present.
	withoutGH := strings.TrimPrefix(s, "github.com/")

	parts := strings.SplitN(withoutGH, "/", 3)
	if len(parts) < 2 || parts[0] == "" || parts[1] == "" {
		// Fallback: treat as local path.
		return ParsedInput{Kind: InputKindLocal, Repo: s}
	}

	repo := "github.com/" + parts[0] + "/" + parts[1]

	if len(parts) == 2 {
		return ParsedInput{Kind: InputKindRepoScan, Repo: repo}
	}

	// 3+ segments → the remainder is the skill hint.
	skillHint := parts[2]
	// Remove leading "tree/main/" or "tree/master/" from GitHub URLs.
	for _, prefix := range []string{"tree/main/", "tree/master/", "blob/main/", "blob/master/"} {
		skillHint = strings.TrimPrefix(skillHint, prefix)
	}
	// Remove trailing slashes.
	skillHint = strings.TrimRight(skillHint, "/")

	if skillHint == "" {
		return ParsedInput{Kind: InputKindRepoScan, Repo: repo}
	}
	return ParsedInput{Kind: InputKindSingleSkill, Repo: repo, SkillHint: skillHint}
}

// DiscoveredSkill is one skill found during a Discover scan.
type DiscoveredSkill struct {
	Name        string
	Path        string // relative path inside the repo (e.g. "skills/pdf")
	Description string
	License     string
	Valid       bool
	Errors      []string
}

// Discover clones a GitHub repo and finds all directories containing SKILL.md.
// Returns the list of skills found, sorted by path.
func Discover(repo string) ([]DiscoveredSkill, error) {
	gitURL := "https://" + repo + ".git"
	repoName := strings.TrimPrefix(repo, "github.com/")

	cloneResult, err := gitutil.CloneRepo(gitURL, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to clone %s: %w", repo, err)
	}

	var skills []DiscoveredSkill
	_ = filepath.Walk(cloneResult.TempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() != "SKILL.md" {
			return nil
		}

		dir := filepath.Dir(path)
		relPath, _ := filepath.Rel(cloneResult.TempDir, dir)
		if relPath == "." {
			return nil
		}

		fm, fmErr := parseFrontmatter(path)
		ds := DiscoveredSkill{
			Path: relPath,
		}
		if fmErr != nil {
			ds.Errors = append(ds.Errors, fmt.Sprintf("frontmatter parse error: %v", fmErr))
		} else {
			ds.Name = fm.Name
			ds.Description = fm.Description
			ds.License = fm.License
			if fm.Name == "" {
				ds.Errors = append(ds.Errors, "missing name")
			}
			if fm.Description == "" {
				ds.Errors = append(ds.Errors, "missing description")
			}
		}
		ds.Valid = len(ds.Errors) == 0
		if ds.Name == "" {
			ds.Name = filepath.Base(relPath)
		}
		skills = append(skills, ds)
		return nil
	})

	return skills, nil
}

// FindSkill clones a repo and searches for a skill by name hint.
// It tries common locations first, then falls back to a full walk.
// Returns the discovered skill and its relative path, or nil if not found.
func FindSkill(repo, skillHint string) (*DiscoveredSkill, error) {
	gitURL := "https://" + repo + ".git"
	repoName := strings.TrimPrefix(repo, "github.com/")

	cloneResult, err := gitutil.CloneRepo(gitURL, repoName)
	if err != nil {
		return nil, fmt.Errorf("failed to clone %s: %w", repo, err)
	}

	// If skillHint looks like a path (contains "/"), try it directly.
	if strings.Contains(skillHint, "/") {
		skillDir := filepath.Join(cloneResult.TempDir, skillHint)
		skillMD := filepath.Join(skillDir, "SKILL.md")
		if _, err := os.Stat(skillMD); err == nil {
			return buildDiscoveredSkill(skillMD, skillHint), nil
		}
	}

	// Try common locations.
	baseName := filepath.Base(skillHint)
	candidates := []string{
		"skills/" + baseName,
		baseName,
		"packages/" + baseName,
		"plugins/" + baseName,
		"packages/skills/" + baseName,
		"plugins/skills/" + baseName,
	}

	for _, cand := range candidates {
		skillMD := filepath.Join(cloneResult.TempDir, cand, "SKILL.md")
		if _, err := os.Stat(skillMD); err == nil {
			return buildDiscoveredSkill(skillMD, cand), nil
		}
	}

	// Full walk: find a directory whose basename matches.
	var found *DiscoveredSkill
	_ = filepath.Walk(cloneResult.TempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || found != nil {
			return nil
		}
		if !info.IsDir() || info.Name() != baseName {
			return nil
		}
		skillMD := filepath.Join(path, "SKILL.md")
		if _, err := os.Stat(skillMD); err == nil {
			relPath, _ := filepath.Rel(cloneResult.TempDir, path)
			found = buildDiscoveredSkill(skillMD, relPath)
		}
		return nil
	})

	return found, nil
}

func buildDiscoveredSkill(skillMDPath, relPath string) *DiscoveredSkill {
	fm, err := parseFrontmatter(skillMDPath)
	ds := &DiscoveredSkill{Path: relPath}
	if err != nil {
		ds.Errors = append(ds.Errors, fmt.Sprintf("frontmatter parse error: %v", err))
	} else {
		ds.Name = fm.Name
		ds.Description = fm.Description
		ds.License = fm.License
		if fm.Name == "" {
			ds.Errors = append(ds.Errors, "missing name")
		}
		if fm.Description == "" {
			ds.Errors = append(ds.Errors, "missing description")
		}
	}
	ds.Valid = len(ds.Errors) == 0
	if ds.Name == "" {
		ds.Name = filepath.Base(relPath)
	}
	return ds
}

// Validate runs all checks for the given request and returns a result.
// It never returns a non-nil error for validation failures — those go into
// result.Errors. A non-nil error means the operation itself failed (e.g.
// network, permission).
func Validate(req ValidateRequest) (*ValidateResult, error) {
	result := &ValidateResult{}

	// Detect source type and resolve to a local directory.
	localDir, sourceType, err := resolve(req)
	if err != nil {
		result.Errors = append(result.Errors, err.Error())
		return result, nil
	}
	result.SourceType = sourceType
	result.ResolvedPath = localDir

	// Check SKILL.md presence.
	skillMDPath := filepath.Join(localDir, "SKILL.md")
	if _, err := os.Stat(skillMDPath); os.IsNotExist(err) {
		result.Errors = append(result.Errors, fmt.Sprintf("SKILL.md not found at %s", skillMDPath))
		return result, nil
	}

	// Parse frontmatter.
	fm, err := parseFrontmatter(skillMDPath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("failed to parse SKILL.md frontmatter: %v", err))
		return result, nil
	}

	result.SkillName = fm.Name
	result.Description = fm.Description
	result.License = fm.License

	// Validate name.
	if fm.Name == "" {
		result.Errors = append(result.Errors, "SKILL.md frontmatter missing required field: name")
	} else if !validNameRe.MatchString(fm.Name) {
		result.Errors = append(result.Errors,
			fmt.Sprintf("invalid skill name %q: must be lowercase letters, numbers, and hyphens only (no leading/trailing/consecutive hyphens)", fm.Name))
	} else if len(fm.Name) > 64 {
		result.Errors = append(result.Errors, fmt.Sprintf("skill name %q exceeds 64 character limit", fm.Name))
	}

	// Validate description.
	if fm.Description == "" {
		result.Errors = append(result.Errors, "SKILL.md frontmatter missing required field: description")
	} else if len(fm.Description) > 1024 {
		result.Errors = append(result.Errors,
			fmt.Sprintf("description exceeds 1024 characters (%d chars)", len(fm.Description)))
	}

	// Warnings.
	if _, err := os.Stat(filepath.Join(localDir, "LICENSE.txt")); os.IsNotExist(err) {
		result.Warnings = append(result.Warnings, "LICENSE.txt not found")
	}

	result.Valid = len(result.Errors) == 0
	return result, nil
}

// resolve determines the source type and clones/resolves the repo to a local dir.
func resolve(req ValidateRequest) (string, SourceType, error) {
	repo := strings.TrimSpace(req.Repo)

	// GitHub
	if strings.HasPrefix(repo, "github.com/") {
		return resolveGitHub(repo, req.Path)
	}

	// Local path (absolute, relative, or ~/)
	return resolveLocal(repo, req.Path)
}

func resolveGitHub(repo, skillPath string) (string, SourceType, error) {
	gitURL := "https://" + repo + ".git"
	repoName := strings.TrimPrefix(repo, "github.com/")

	var cloneResult *gitutil.CloneResult
	var err error

	if skillPath != "" {
		cloneResult, err = gitutil.SparseCloneRepo(gitURL, repoName, []string{skillPath})
	} else {
		cloneResult, err = gitutil.CloneRepo(gitURL, repoName)
	}
	if err != nil {
		return "", SourceTypeGitHub, fmt.Errorf("failed to clone %s: %w", repo, err)
	}

	skillDir := filepath.Join(cloneResult.TempDir, skillPath)
	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		return "", SourceTypeGitHub,
			fmt.Errorf("skill path %q not found in repository %s", skillPath, repo)
	}

	return skillDir, SourceTypeGitHub, nil
}

func resolveLocal(rawPath, subPath string) (string, SourceType, error) {
	// Expand ~/
	if strings.HasPrefix(rawPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", SourceTypeLocal, fmt.Errorf("cannot resolve home directory: %w", err)
		}
		rawPath = filepath.Join(home, rawPath[2:])
	}

	// Make absolute.
	absPath, err := filepath.Abs(rawPath)
	if err != nil {
		return "", SourceTypeLocal, fmt.Errorf("invalid path %q: %w", rawPath, err)
	}

	skillDir := absPath
	if subPath != "" {
		skillDir = filepath.Join(absPath, subPath)
	}

	if _, err := os.Stat(skillDir); os.IsNotExist(err) {
		return "", SourceTypeLocal, fmt.Errorf("path %q does not exist", skillDir)
	}

	return skillDir, SourceTypeLocal, nil
}

// parseFrontmatter extracts the YAML front matter block from a SKILL.md file.
// Front matter is delimited by "---" lines at the start of the file.
func parseFrontmatter(path string) (*skillFrontmatter, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	// First non-empty line must be "---"
	started := false
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if !started {
			if strings.TrimSpace(line) == "" {
				continue
			}
			if strings.TrimSpace(line) != "---" {
				// No frontmatter — return empty struct (will trigger validation errors)
				return &skillFrontmatter{}, nil
			}
			started = true
			continue
		}
		if strings.TrimSpace(line) == "---" {
			break
		}
		lines = append(lines, line)
	}

	if !started || len(lines) == 0 {
		return &skillFrontmatter{}, nil
	}

	var fm skillFrontmatter
	if err := yaml.Unmarshal([]byte(strings.Join(lines, "\n")), &fm); err != nil {
		return nil, err
	}
	return &fm, nil
}
