// Package initcmd implements the init command
package initcmd

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/castle-x/skills-x/cmd/skills-x/errmsg"
	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/cmd/skills-x/skills"
	"github.com/castle-x/skills-x/pkg/discover"
	"github.com/castle-x/skills-x/pkg/gitutil"
	"github.com/castle-x/skills-x/pkg/registry"
	"github.com/spf13/cobra"
)

// ANSI colors
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

var (
	flagAll    bool
	flagTarget string
	flagForce  bool
)

// NewCommand creates the init command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [skill_name]",
		Short: i18n.T("cmd_init_short"),
		Long:  i18n.T("cmd_init_long"),
		RunE:  runInit,
	}

	cmd.Flags().BoolVar(&flagAll, "all", false, i18n.T("cmd_init_flag_all"))
	cmd.Flags().StringVarP(&flagTarget, "target", "t", "", i18n.T("cmd_init_flag_target"))
	cmd.Flags().BoolVarP(&flagForce, "force", "f", false, i18n.T("cmd_init_flag_force"))

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := flagTarget
	if targetDir == "" {
		// Default to current directory
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		targetDir = cwd
	}

	// Create target directory if not exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return errmsg.TargetDirCreateError(targetDir)
	}

	fmt.Printf("%s%s%s\n\n", colorCyan, i18n.Tf("init_target_dir", targetDir), colorReset)

	// Load registry
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	if flagAll {
		return initAll(reg, targetDir)
	}

	if len(args) == 0 {
		return errmsg.MissingArgument("skill_name")
	}

	return initSkill(reg, args[0], targetDir)
}

func initSkill(reg *registry.Registry, name string, targetDir string) error {
	// First check if it's an X (self-developed) skill
	if skills.XSkillExists(name) {
		return initXSkill(name, targetDir)
	}

	// Otherwise, treat as a registry skill
	return initRegistrySkill(reg, name, targetDir)
}

func initXSkill(name string, targetDir string) error {
	// Get the embedded filesystem for this skill
	embedFS, skillPath, ok := skills.GetXSkillFS(name)
	if !ok || skillPath == "" {
		return fmt.Errorf("%s: %s", i18n.T("init_skill_path_not_found"), name)
	}

	dstPath := filepath.Join(targetDir, name)

	// Check if already exists
	if dirExists(dstPath) {
		if !flagForce {
			if !confirmOverwrite(name) {
				fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_skipped", name), colorReset)
				return nil
			}
		}
		fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_overwrite", name), colorReset)
	} else {
		fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_downloading", name), colorReset)
	}

	// Copy skill from embedded filesystem to target
	if err := copyFromEmbedFS(embedFS, skillPath, dstPath); err != nil {
		return errmsg.CopyFailed(name)
	}

	fmt.Printf("%s✓ %s%s\n", colorGreen, i18n.Tf("init_success", name), colorReset)
	fmt.Printf("  %s%s%s\n", colorGray, i18n.T("init_from_embedded"), colorReset)

	return nil
}

func initRegistrySkill(reg *registry.Registry, name string, targetDir string) error {
	// Find skill in registry
	matches := reg.FindSkillsWithConflict(name)

	if len(matches) == 0 {
		return errmsg.SkillNotFound(name)
	}

	var skill *registry.Skill
	var source *registry.Source

	if len(matches) > 1 {
		// Multiple skills with same name - prompt user to choose
		fmt.Printf("%s%s%s\n\n", colorYellow, i18n.Tf("init_conflict_found", name, len(matches)), colorReset)
		
		for i, m := range matches {
			fmt.Printf("  %d. %s%s%s from %s%s%s\n",
				i+1,
				colorCyan, m.Skill.Name, colorReset,
				colorGray, m.Source.Repo, colorReset)
		}
		fmt.Println()

		choice := promptChoice(len(matches))
		if choice < 0 {
			fmt.Printf("%s%s%s\n", colorYellow, i18n.T("init_cancelled"), colorReset)
			return nil
		}

		skill = matches[choice].Skill
		source = matches[choice].Source
	} else {
		skill = matches[0].Skill
		source = matches[0].Source
	}

	// Clone the repository
	fmt.Printf("%s%s %s...%s\n", colorGray, i18n.T("init_cloning"), source.GetRepoShortName(), colorReset)

	var result *gitutil.CloneResult
	var err error

	// For large repos (marked with skip_fetch), use sparse checkout
	if source.SkipFetch && skill.Path != "" {
		// Use sparse checkout for large repos - only fetch the specific skill path
		result, err = gitutil.SparseCloneRepo(source.GetGitURL(), source.Repo, []string{skill.Path})
	} else {
		result, err = gitutil.CloneRepo(source.GetGitURL(), source.Repo)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("init_clone_failed"), err)
	}

	// Find the skill in the cloned repo
	var skillPath string
	if skill.Path != "" {
		skillPath = filepath.Join(result.TempDir, skill.Path)
	} else {
		// Try to discover the skill
		discovered, err := discover.DiscoverSkillByPath(result.TempDir, skill.Name)
		if err != nil || discovered == nil {
			// Fallback: search in common locations
			discovered, _ = findSkillInRepo(result.TempDir, skill.Name)
		}
		if discovered != nil {
			skillPath = discovered.Path
		}
	}

	if skillPath == "" || !dirExists(skillPath) {
		return fmt.Errorf("%s: %s", i18n.T("init_skill_path_not_found"), skill.Name)
	}

	dstPath := filepath.Join(targetDir, skill.Name)

	// Check if already exists
	if dirExists(dstPath) {
		if !flagForce {
			if !confirmOverwrite(skill.Name) {
				fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_skipped", skill.Name), colorReset)
				return nil
			}
		}
		fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_overwrite", skill.Name), colorReset)
	} else {
		fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_downloading", skill.Name), colorReset)
	}

	// Copy skill to target
	if err := copyDir(skillPath, dstPath); err != nil {
		return errmsg.CopyFailed(skill.Name)
	}

	fmt.Printf("%s✓ %s%s\n", colorGreen, i18n.Tf("init_success", skill.Name), colorReset)
	fmt.Printf("  %s%s%s\n", colorGray, i18n.Tf("init_from_source", source.Repo), colorReset)

	return nil
}

func initAll(reg *registry.Registry, targetDir string) error {
	sources := reg.GetAllSources()
	
	count := 0
	skipped := 0
	errors := 0

	for _, source := range sources {
		fmt.Printf("\n%s📦 %s%s\n", colorBold, source.Repo, colorReset)

		// For large repos, use sparse checkout for each skill
		if source.SkipFetch {
			// Install each skill with sparse checkout
			for _, skill := range source.Skills {
				if skill.Path == "" {
					fmt.Printf("%s  ⚠ %s: %s%s\n", colorYellow, skill.Name, i18n.T("init_skill_path_not_found"), colorReset)
					errors++
					continue
				}

				result, err := gitutil.SparseCloneRepo(source.GetGitURL(), source.Repo, []string{skill.Path})
				if err != nil {
					fmt.Printf("%s  ⚠ %s: %v%s\n", colorYellow, skill.Name, err, colorReset)
					errors++
					continue
				}

				skillPath := filepath.Join(result.TempDir, skill.Path)
				dstPath := filepath.Join(targetDir, skill.Name)

				if dirExists(dstPath) && !flagForce {
					fmt.Printf("%s  - %s%s\n", colorGray, i18n.Tf("init_skipped", skill.Name), colorReset)
					skipped++
					continue
				}

				if !dirExists(skillPath) {
					fmt.Printf("%s  ⚠ %s: %s%s\n", colorYellow, skill.Name, i18n.T("init_skill_path_not_found"), colorReset)
					errors++
					continue
				}

				if err := copyDir(skillPath, dstPath); err != nil {
					fmt.Printf("%s  ✗ %s: %v%s\n", colorRed, skill.Name, err, colorReset)
					errors++
					continue
				}

				fmt.Printf("%s  ✓ %s%s\n", colorGreen, skill.Name, colorReset)
				count++
			}
			continue
		}

		// Clone repository for normal repos
		result, err := gitutil.CloneRepo(source.GetGitURL(), source.Repo)
		if err != nil {
			fmt.Printf("%s  ⚠ %s: %v%s\n", colorYellow, i18n.T("init_clone_failed"), err, colorReset)
			errors++
			continue
		}

		// Install each skill
		for _, skill := range source.Skills {
			var skillPath string
			if skill.Path != "" {
				skillPath = filepath.Join(result.TempDir, skill.Path)
			} else {
				discovered, _ := findSkillInRepo(result.TempDir, skill.Name)
				if discovered != nil {
					skillPath = discovered.Path
				}
			}

			if skillPath == "" || !dirExists(skillPath) {
				fmt.Printf("%s  ⚠ %s: %s%s\n", colorYellow, skill.Name, i18n.T("init_skill_path_not_found"), colorReset)
				errors++
				continue
			}

			dstPath := filepath.Join(targetDir, skill.Name)

			// Check if exists
			if dirExists(dstPath) && !flagForce {
				fmt.Printf("%s  - %s%s\n", colorGray, i18n.Tf("init_skipped", skill.Name), colorReset)
				skipped++
				continue
			}

			if err := copyDir(skillPath, dstPath); err != nil {
				fmt.Printf("%s  ✗ %s: %v%s\n", colorRed, skill.Name, err, colorReset)
				errors++
				continue
			}

			fmt.Printf("%s  ✓ %s%s\n", colorGreen, skill.Name, colorReset)
			count++
		}
	}

	fmt.Printf("\n%s%s%s\n", colorGreen, i18n.Tf("init_all_success", count), colorReset)
	if skipped > 0 {
		fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_all_skipped", skipped), colorReset)
	}
	if errors > 0 {
		fmt.Printf("%s%s%s\n", colorRed, i18n.Tf("init_all_errors", errors), colorReset)
	}

	return nil
}

// findSkillInRepo searches for a skill by name in common locations
func findSkillInRepo(repoPath string, skillName string) (*discover.DiscoveredSkill, error) {
	// First, try exact path matches
	commonPaths := []string{
		filepath.Join("skills", skillName),
		filepath.Join("packages", "docs", "skills", skillName),
		skillName,
	}

	for _, path := range commonPaths {
		fullPath := filepath.Join(repoPath, path)
		if dirExists(fullPath) && fileExists(filepath.Join(fullPath, "SKILL.md")) {
			return &discover.DiscoveredSkill{
				Name: skillName,
				Path: fullPath,
			}, nil
		}
	}

	// Fall back to discovery
	discovered, err := discover.DiscoverSkills(repoPath, nil)
	if err != nil {
		return nil, err
	}

	for _, d := range discovered {
		if strings.EqualFold(d.Name, skillName) {
			return &d, nil
		}
	}

	return nil, nil
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// confirmOverwrite prompts user to confirm overwrite for a single skill
func confirmOverwrite(name string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s%s [y/N]: %s", colorYellow, i18n.Tf("init_confirm_overwrite", name), colorReset)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// promptChoice prompts user to choose from options
func promptChoice(max int) int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s%s [1-%d]: %s", colorCyan, i18n.T("init_choose_source"), max, colorReset)
	response, err := reader.ReadString('\n')
	if err != nil {
		return -1
	}
	response = strings.TrimSpace(response)
	
	var choice int
	if _, err := fmt.Sscanf(response, "%d", &choice); err != nil {
		return -1
	}
	
	if choice < 1 || choice > max {
		return -1
	}
	
	return choice - 1
}

// copyDir copies a directory from source to destination
// Follows symlinks to copy actual content
func copyDir(srcPath string, dstPath string) error {
	// Remove existing directory
	os.RemoveAll(dstPath)

	return filepath.WalkDir(srcPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dstPath, relPath)

		// Check if it's a symlink
		if d.Type()&os.ModeSymlink != 0 {
			// Resolve symlink to get actual path
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return nil // Skip broken symlinks
			}

			info, err := os.Stat(realPath)
			if err != nil {
				return nil // Skip if can't stat
			}

			if info.IsDir() {
				// Recursively copy symlinked directory
				return copyDir(realPath, targetPath)
			}
			// Copy symlinked file
			return copyFile(realPath, targetPath)
		}

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Copy file
		return copyFile(path, targetPath)
	})
}

// copyFile copies a single file
func copyFile(srcPath string, dstPath string) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	// Open source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy content
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// copyFromEmbedFS copies a directory from embedded filesystem to target path
func copyFromEmbedFS(embedFS embed.FS, srcPath string, dstPath string) error {
	// Remove existing directory
	os.RemoveAll(dstPath)

	// Walk through embedded filesystem
	// Note: embed.FS always uses "/" as separator, regardless of OS
	return fs.WalkDir(embedFS, srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path using string operations
		// Note: embed.FS paths always use "/", not filepath functions (which use "\" on Windows)
		relPath := path
		if len(path) > len(srcPath) {
			relPath = strings.TrimPrefix(path, srcPath+"/")
		} else if path == srcPath {
			relPath = "."
		}

		// Convert to OS-specific path for target filesystem
		targetPath := filepath.Join(dstPath, filepath.FromSlash(relPath))

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Read file from embedded FS
		data, err := embedFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// Write file to target
		return os.WriteFile(targetPath, data, 0644)
	})
}
