// Package initcmd implements the init command
package initcmd

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthropics/skills-x/cmd/skills-x/errmsg"
	"github.com/anthropics/skills-x/cmd/skills-x/i18n"
	"github.com/anthropics/skills-x/cmd/skills-x/skills"
	"github.com/spf13/cobra"
)

// ANSI colors
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
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
	cmd.Flags().StringVar(&flagTarget, "target", "", i18n.T("cmd_init_flag_target"))
	cmd.Flags().BoolVarP(&flagForce, "force", "f", false, i18n.T("cmd_init_flag_force"))

	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine target directory
	targetDir := flagTarget
	if targetDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		targetDir = filepath.Join(home, ".claude", "skills")
	}

	// Create target directory if not exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return errmsg.TargetDirCreateError(targetDir)
	}

	fmt.Printf("%s%s%s\n\n", colorCyan, i18n.Tf("init_target_dir", targetDir), colorReset)

	if flagAll {
		return initAll(targetDir)
	}

	if len(args) == 0 {
		return errmsg.MissingArgument("skill_name")
	}

	return initSkill(args[0], targetDir)
}

func initSkill(name string, targetDir string) error {
	// Check if skill exists
	if !skills.SkillExists(name) {
		return errmsg.SkillNotFound(name)
	}

	// Copy skill to target
	srcPath := filepath.Join("data", name)
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

	if err := copyDir(skills.GetFS(), srcPath, dstPath); err != nil {
		return errmsg.CopyFailed(name)
	}

	fmt.Printf("%s✓ %s%s\n", colorGreen, i18n.Tf("init_success", name), colorReset)
	return nil
}

func initAll(targetDir string) error {
	skillList, err := skills.ListSkills()
	if err != nil {
		return err
	}

	// Check for existing skills
	existingSkills := []string{}
	for _, s := range skillList {
		dstPath := filepath.Join(targetDir, s.Name)
		if dirExists(dstPath) {
			existingSkills = append(existingSkills, s.Name)
		}
	}

	// Prompt for overwrite if there are existing skills and not forced
	overwriteAll := flagForce
	if len(existingSkills) > 0 && !flagForce {
		fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_existing_count", len(existingSkills)), colorReset)
		overwriteAll = confirmOverwriteAll()
	}

	count := 0
	skipped := 0
	for _, s := range skillList {
		srcPath := filepath.Join("data", s.Name)
		dstPath := filepath.Join(targetDir, s.Name)

		// Check if exists
		exists := dirExists(dstPath)
		if exists && !overwriteAll {
			fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_skipped", s.Name), colorReset)
			skipped++
			continue
		}

		if exists {
			fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_overwrite", s.Name), colorReset)
		} else {
			fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_downloading", s.Name), colorReset)
		}

		if err := copyDir(skills.GetFS(), srcPath, dstPath); err != nil {
			fmt.Printf("%s✗ %s: %v%s\n", colorYellow, s.Name, err, colorReset)
			continue
		}

		fmt.Printf("%s✓ %s%s\n", colorGreen, i18n.Tf("init_success", s.Name), colorReset)
		count++
	}

	fmt.Printf("\n%s%s%s\n", colorGreen, i18n.Tf("init_all_success", count), colorReset)
	if skipped > 0 {
		fmt.Printf("%s%s%s\n", colorYellow, i18n.Tf("init_all_skipped", skipped), colorReset)
	}
	return nil
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
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

// confirmOverwriteAll prompts user to confirm overwrite for all existing skills
func confirmOverwriteAll() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s%s [y/N]: %s", colorYellow, i18n.T("init_confirm_overwrite_all"), colorReset)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// copyDir copies a directory from embed.FS to local filesystem
func copyDir(srcFS fs.FS, srcPath string, dstPath string) error {
	// Remove existing directory
	os.RemoveAll(dstPath)

	return fs.WalkDir(srcFS, srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dstPath, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Copy file
		return copyFile(srcFS, path, targetPath)
	})
}

// copyFile copies a single file from embed.FS to local filesystem
func copyFile(srcFS fs.FS, srcPath string, dstPath string) error {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	// Open source file
	srcFile, err := srcFS.Open(srcPath)
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
