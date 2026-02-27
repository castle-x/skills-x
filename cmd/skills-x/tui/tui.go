// Package tui provides terminal interactive UI components
package tui

import (
	"fmt"
	"path/filepath"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/pkg/products"
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
	installSkills, uninstallSkills, updateSkills, err := RunSkillsSelect(product, allSkills, opts.Version, targetDir)
	if err != nil {
		if err.Error() == "quit" {
			return nil
		}
		if err.Error() == "go back" {
			return runTUIFlow(opts)
		}
		return fmt.Errorf("%s: %w", i18n.T("error_skills_select"), err)
	}

	if len(installSkills) == 0 && len(uninstallSkills) == 0 && len(updateSkills) == 0 {
		return nil
	}

	// =========================================================================
	// Level 4: Install/Update/Uninstall Skills
	// =========================================================================
	fmt.Print(ClearScreen)
	completed, failed, err := RunInstaller(installSkills, uninstallSkills, updateSkills, targetDir)
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

