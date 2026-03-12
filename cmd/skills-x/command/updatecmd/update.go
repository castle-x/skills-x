// Package updatecmd implements the update command for skills-x
package updatecmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/cmd/skills-x/tui"
	"github.com/castle-x/skills-x/pkg/discover"
	"github.com/castle-x/skills-x/pkg/gitutil"
	"github.com/castle-x/skills-x/pkg/registry"
	"github.com/spf13/cobra"
)

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
	flagCheck  bool
	flagTarget string
)

var (
	cloneRepoWithRefresh = gitutil.CloneRepoWithRefresh
	sparseCloneRepo      = gitutil.SparseCloneRepo
	getRepoHeadCommit    = gitutil.GetRepoHeadCommit
)

// NewCommand creates the update command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [skill_name...]",
		Short: i18n.T("cmd_update_short"),
		Long:  i18n.T("cmd_update_long"),
		RunE:  runUpdate,
	}

	cmd.Flags().BoolVarP(&flagAll, "all", "a", false, i18n.T("cmd_update_flag_all"))
	cmd.Flags().BoolVarP(&flagCheck, "check", "c", false, i18n.T("cmd_update_flag_check"))
	cmd.Flags().StringVarP(&flagTarget, "target", "t", "", i18n.T("cmd_update_flag_target"))

	return cmd
}

type skillCheckResult struct {
	name         string
	status       string // "up_to_date", "update_available", "no_meta", "error"
	localCommit  string
	remoteCommit string
	err          error
}

func runUpdate(cmd *cobra.Command, args []string) error {
	targetDir := flagTarget
	if targetDir == "" {
		// Default: ~/.claude/skills (or similar)
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot determine home directory: %w", err)
		}
		targetDir = filepath.Join(home, ".claude", "skills")
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("target directory does not exist: %s", targetDir)
	}

	reg, warnings, err := registry.LoadWithUser()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "⚠ %s\n", w)
	}

	// Find installed skills in target directory
	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	type installedSkill struct {
		name   string
		meta   *tui.SkillMeta
		source *registry.Source
		skill  *registry.Skill
	}

	var installed []installedSkill

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skillDir := filepath.Join(targetDir, entry.Name())
		// Must be a valid skill directory
		if _, err := os.Stat(filepath.Join(skillDir, "SKILL.md")); err != nil {
			continue
		}

		meta, _ := tui.ReadSkillMeta(skillDir)

		// Find in registry
		matches := reg.FindSkillsWithConflict(entry.Name())
		var source *registry.Source
		var skill *registry.Skill
		if len(matches) > 0 {
			if meta != nil && meta.Source != "" {
				for _, m := range matches {
					if m.Source.Name == meta.Source {
						source = m.Source
						skill = m.Skill
						break
					}
				}
			}
			if source == nil {
				source = matches[0].Source
				skill = matches[0].Skill
			}
		}

		// Filter by args if specific names given
		if !flagAll && len(args) > 0 {
			found := false
			for _, arg := range args {
				if strings.EqualFold(entry.Name(), arg) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		installed = append(installed, installedSkill{
			name:   entry.Name(),
			meta:   meta,
			source: source,
			skill:  skill,
		})
	}

	if !flagAll && len(args) == 0 {
		return fmt.Errorf("specify skill names or use --all to update all installed skills")
	}

	if len(installed) == 0 {
		fmt.Println("No installed skills found to update.")
		return nil
	}

	fmt.Printf("Checking for updates (%s)...\n\n", targetDir)

	var results []skillCheckResult
	updateAvailable := 0

	for _, is := range installed {
		if is.source == nil || is.skill == nil {
			results = append(results, skillCheckResult{
				name:   is.name,
				status: "no_meta",
			})
			continue
		}

		// Get cached/fresh repo
		var cloneResult *gitutil.CloneResult
			// Check mode must also refresh, otherwise stale cache can hide updates.
			refresh := true
			if is.source.SkipFetch && is.skill.Path != "" {
				cloneResult, err = sparseCloneRepo(is.source.GetGitURL(), is.source.Repo, is.source.Branch, []string{is.skill.Path})
			} else {
				cloneResult, err = cloneRepoWithRefresh(is.source.GetGitURL(), is.source.Repo, is.source.Branch, refresh)
			}
		if err != nil {
			results = append(results, skillCheckResult{
				name:   is.name,
				status: "error",
				err:    err,
			})
			continue
		}

			remoteCommit, err := getRepoHeadCommit(cloneResult.TempDir)
		if err != nil {
			results = append(results, skillCheckResult{
				name:   is.name,
				status: "error",
				err:    err,
			})
			continue
		}

		localCommit := ""
		if is.meta != nil {
			localCommit = is.meta.Commit
		}

		if localCommit == "" {
			results = append(results, skillCheckResult{
				name:         is.name,
				status:       "no_meta",
				remoteCommit: remoteCommit,
			})
			updateAvailable++
		} else if localCommit == remoteCommit {
			results = append(results, skillCheckResult{
				name:         is.name,
				status:       "up_to_date",
				localCommit:  localCommit,
				remoteCommit: remoteCommit,
			})
		} else {
			results = append(results, skillCheckResult{
				name:         is.name,
				status:       "update_available",
				localCommit:  localCommit,
				remoteCommit: remoteCommit,
			})
			updateAvailable++
		}

		// If not check-only, perform the update
		if !flagCheck && (localCommit != remoteCommit || localCommit == "") {
			var skillPath string
			if is.skill.Path != "" {
				skillPath = filepath.Join(cloneResult.TempDir, is.skill.Path)
			} else {
				discovered, err := discover.DiscoverSkillByPath(cloneResult.TempDir, is.skill.Name)
				if err != nil || discovered == nil {
					discovered, _ = findSkillInRepo(cloneResult.TempDir, is.skill.Name)
				}
				if discovered != nil {
					skillPath = discovered.Path
				}
			}

			if skillPath == "" {
				results[len(results)-1].status = "error"
				results[len(results)-1].err = fmt.Errorf("skill path not found")
				continue
			}

			dstPath := filepath.Join(targetDir, is.skill.Name)
			os.RemoveAll(dstPath)

			if err := copyDir(skillPath, dstPath); err != nil {
				results[len(results)-1].status = "error"
				results[len(results)-1].err = fmt.Errorf("copy failed: %w", err)
				continue
			}

			// Write meta
			meta := tui.SkillMeta{
				Skill:  is.skill.Name,
				Source: is.source.Name,
				Repo:   is.source.Repo,
				Commit: remoteCommit,
			}
			_ = tui.WriteSkillMeta(dstPath, meta)
		}
	}

	// Print results
	for _, r := range results {
		name := padRight(r.name, 25)
		switch r.status {
		case "up_to_date":
			fmt.Printf("  %s✓%s %s %sup to date%s (%s)\n", colorGreen, colorReset, name, colorGray, colorReset, r.localCommit)
		case "update_available":
			if flagCheck {
				fmt.Printf("  %s↑%s %s %supdate available%s (%s → %s)\n", colorYellow, colorReset, name, colorYellow, colorReset, r.localCommit, r.remoteCommit)
			} else {
				fmt.Printf("  %s✓%s %s %supdated%s (%s → %s)\n", colorGreen, colorReset, name, colorGreen, colorReset, r.localCommit, r.remoteCommit)
			}
		case "no_meta":
			if flagCheck {
				fmt.Printf("  %s-%s %s %sno metadata, skip%s\n", colorGray, colorReset, name, colorGray, colorReset)
			} else {
				fmt.Printf("  %s-%s %s %sno metadata, reinstalling%s\n", colorGray, colorReset, name, colorGray, colorReset)
			}
		case "error":
			fmt.Printf("  %s✗%s %s %serror: %v%s\n", colorRed, colorReset, name, colorRed, r.err, colorReset)
		}
	}

	fmt.Println()

	if flagCheck && updateAvailable > 0 {
		fmt.Printf("%d skill(s) can be updated. Run: %sskills-x update --all%s\n", updateAvailable, colorBold, colorReset)
	} else if !flagCheck {
		fmt.Printf("Update complete.\n")
	} else {
		fmt.Printf("All skills are up to date.\n")
	}

	return nil
}

func padRight(s string, width int) string {
	for len(s) < width {
		s += " "
	}
	return s
}

func findSkillInRepo(repoPath string, skillName string) (*discover.DiscoveredSkill, error) {
	commonPaths := []string{
		filepath.Join("skills", skillName),
		filepath.Join("packages", "docs", "skills", skillName),
		skillName,
	}

	for _, path := range commonPaths {
		fullPath := filepath.Join(repoPath, path)
		info, err := os.Stat(fullPath)
		if err != nil || !info.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(fullPath, "SKILL.md")); err == nil {
			return &discover.DiscoveredSkill{
				Name: skillName,
				Path: fullPath,
			}, nil
		}
	}

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

// copyDir copies a directory recursively
func copyDir(srcPath string, dstPath string) error {
	os.RemoveAll(dstPath)

	return filepath.WalkDir(srcPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dstPath, relPath)

		if d.Type()&os.ModeSymlink != 0 {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return nil
			}
			info, err := os.Stat(realPath)
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return copyDir(realPath, targetPath)
			}
			return copyFile(realPath, targetPath)
		}

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return copyFile(path, targetPath)
	})
}

// copyFile copies a single file
func copyFile(srcPath string, dstPath string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}
	return os.WriteFile(dstPath, data, 0644)
}
