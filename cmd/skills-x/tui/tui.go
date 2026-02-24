// Package tui provides terminal interactive UI components
package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/pkg/products"
	"github.com/spf13/cobra"
)

// ============================================================================
// TUI Main Entry - Page Flow Orchestration
// Level 1: Select Product -> Level 2: Select Skills -> Install
// ============================================================================

// TUIOptions contains options for running the TUI
type TUIOptions struct {
	Version   string
	TargetDir string
}

// RunTUI runs the complete TUI flow: Product Select -> Skills Select -> Install
func RunTUI(opts TUIOptions) error {
	// Enter alt screen once for the entire TUI session
	fmt.Print(EnterAltScreen)
	fmt.Print(HideCursor)
	defer func() {
		fmt.Print(ShowCursor)
		fmt.Print(ExitAltScreen)
	}()

	return runTUIFlow(opts)
}

// runTUIFlow runs the TUI page flow (can be called recursively for "go back")
func runTUIFlow(opts TUIOptions) error {
	// =========================================================================
	// Level 1: Select Product
	// =========================================================================
	fmt.Print(ClearScreen)
	product, err := RunProductSelect(opts.Version, opts.TargetDir)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("error_product_select"), err)
	}
	if product == nil {
		return nil
	}

	// =========================================================================
	// Level 2: Select Install Target (Global or Project)
	// =========================================================================
	fmt.Print(ClearScreen)
	installTarget, err := RunInstallTargetSelect(product, opts.TargetDir)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("error_install_target_select"), err)
	}
	if installTarget == "" {
		return nil
	}

	var targetDir string
	if installTarget == "global" {
		targetDir = product.GlobalPath()
	} else {
		targetDir = filepath.Join(opts.TargetDir, product.ProjectSkills)
	}

	// =========================================================================
	// Level 3: Select Skills
	// =========================================================================
	allSkills, err := LoadSkillsForProduct(product, targetDir)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("error_load_skills"), err)
	}

	fmt.Print(ClearScreen)
	selectedSkills, deselectedSkills, err := RunSkillsSelect(product, allSkills, opts.Version, targetDir)
	if err != nil {
		if err.Error() == "quit" {
			return nil
		}
		if err.Error() == "go back" {
			return runTUIFlow(opts)
		}
		return fmt.Errorf("%s: %w", i18n.T("error_skills_select"), err)
	}

	if len(selectedSkills) == 0 && len(deselectedSkills) == 0 {
		return nil
	}

	// =========================================================================
	// Level 4: Install/Uninstall Skills
	// =========================================================================
	fmt.Print(ClearScreen)
	completed, failed, err := RunInstaller(selectedSkills, deselectedSkills, targetDir)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("error_installation"), err)
	}
	_ = completed
	_ = failed

	return nil
}

// LoadSkillsForProduct loads skills for a specific product
func LoadSkillsForProduct(product *products.Product, targetDir string) ([]SkillItem, error) {
	// Load skills from registry
	skills, err := LoadSkillsFromRegistry(targetDir)
	if err != nil {
		return nil, err
	}

	// Filter skills for the selected product if needed
	// For now, we return all skills - the user can filter in the UI

	return skills, nil
}

// TUICommand returns a cobra command for the TUI
func TUICommand(version string) *cobra.Command {
	var targetDir string

	cmd := &cobra.Command{
		Use:   "tui",
		Short: i18n.T("tui_desc"),
		Long:  i18n.T("tui_long_desc"),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get target directory
			if targetDir == "" {
				// Use current working directory
				cwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("%s: %w", i18n.T("error_get_cwd"), err)
				}
				targetDir = cwd
			}

			// Run TUI
			opts := TUIOptions{
				Version:   version,
				TargetDir: targetDir,
			}

			if err := RunTUI(opts); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&targetDir, "dir", "d", "", i18n.T("tui_flag_dir"))

	return cmd
}
