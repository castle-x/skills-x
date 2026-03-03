package registry

import (
	"fmt"
	"strings"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/pkg/registry"
	"github.com/castle-x/skills-x/pkg/skillvalidator"
	"github.com/spf13/cobra"
)

func newCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "check <owner/repo[/skill]> | <local-path> [skill-path]",
		Short: i18n.T("cmd_registry_check_short"),
		Long:  i18n.T("cmd_registry_check_long"),
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Legacy 2-arg mode: check github.com/owner/repo skills/x
			if len(args) == 2 {
				return runCheckSingle(args[0], args[1])
			}

			parsed := skillvalidator.ParseInput(args[0])

			switch parsed.Kind {
			case skillvalidator.InputKindRepoScan:
				return runCheckDiscover(parsed.Repo)
			case skillvalidator.InputKindSingleSkill:
				if parsed.SkillHint != "" {
					return runCheckFind(parsed.Repo, parsed.SkillHint)
				}
				return runCheckDiscover(parsed.Repo)
			default:
				return runCheckSingle(parsed.Repo, "")
			}
		},
	}
}

// runCheckDiscover clones the repo and lists all discovered skills.
func runCheckDiscover(repo string) error {
	fmt.Printf("%s %s ...\n", i18n.T("registry_scanning"), repoShortName(repo))

	skills, err := skillvalidator.Discover(repo)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_scan_failed"), err)
	}

	if len(skills) == 0 {
		fmt.Println(i18n.Tf("registry_scan_empty", repoShortName(repo)))
		return nil
	}

	// Load built-in for conflict hints.
	builtinNames := loadBuiltinNames()

	fmt.Printf("\n%s %d %s\n\n", i18n.T("registry_scan_found"), len(skills), i18n.T("registry_scan_skills"))
	fmt.Printf("  %-4s %-30s %-50s %s\n", "#", i18n.T("registry_field_name"), i18n.T("registry_field_desc"), i18n.T("registry_field_status"))
	fmt.Println("  " + strings.Repeat("─", 90))

	for idx, s := range skills {
		desc := s.Description
		if len(desc) > 47 {
			desc = desc[:44] + "..."
		}
		status := "✓"
		if !s.Valid {
			status = "✗"
		}
		if _, exists := builtinNames[strings.ToLower(s.Name)]; exists {
			status += " " + i18n.T("registry_scan_builtin")
		}
		fmt.Printf("  %-4d %-30s %-50s %s\n", idx+1, s.Name, desc, status)
	}
	fmt.Println()
	return nil
}

// runCheckFind searches for a specific skill in a repo.
func runCheckFind(repo, skillHint string) error {
	fmt.Printf("%s %s ...\n", i18n.T("registry_scanning"), repoShortName(repo))

	ds, err := skillvalidator.FindSkill(repo, skillHint)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_scan_failed"), err)
	}

	if ds == nil {
		fmt.Printf("\n%s\n", i18n.Tf("registry_skill_not_found", skillHint, repoShortName(repo)))
		return nil
	}

	printDiscoveredSkill(ds)
	return nil
}

// runCheckSingle is the original single-skill validation path.
func runCheckSingle(repo, path string) error {
	req := skillvalidator.ValidateRequest{Repo: repo, Path: path}
	fmt.Println(i18n.Tf("registry_checking", req.Repo, req.Path))

	result, err := skillvalidator.Validate(req)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_check_failed"), err)
	}

	printValidateResult(result)
	if !result.Valid {
		return fmt.Errorf("%s", i18n.T("registry_check_not_valid"))
	}
	return nil
}

// printValidateResult prints a human-readable summary of a ValidateResult.
func printValidateResult(r *skillvalidator.ValidateResult) {
	if r.Valid {
		fmt.Printf("✓ %s\n", i18n.T("registry_check_passed"))
	} else {
		fmt.Printf("✗ %s\n", i18n.T("registry_check_failed"))
	}

	if r.SkillName != "" {
		fmt.Printf("  %s: %s\n", i18n.T("registry_field_name"), r.SkillName)
	}
	if r.Description != "" {
		desc := r.Description
		if len(desc) > 80 {
			desc = desc[:77] + "..."
		}
		fmt.Printf("  %s: %s\n", i18n.T("registry_field_desc"), desc)
	}
	if r.License != "" {
		fmt.Printf("  %s: %s\n", i18n.T("registry_field_license"), r.License)
	}
	if r.ResolvedPath != "" {
		fmt.Printf("  %s: %s\n", i18n.T("registry_field_path"), r.ResolvedPath)
	}

	for _, e := range r.Errors {
		fmt.Printf("  ✗ %s: %s\n", i18n.T("registry_error"), e)
	}
	for _, w := range r.Warnings {
		fmt.Printf("  ⚠ %s: %s\n", i18n.T("registry_warning"), w)
	}
}

// printDiscoveredSkill prints details of a single found skill.
func printDiscoveredSkill(ds *skillvalidator.DiscoveredSkill) {
	fmt.Printf("\n%s %s:\n", i18n.T("registry_skill_found"), ds.Name)
	if ds.Description != "" {
		fmt.Printf("  %s: %s\n", i18n.T("registry_field_desc"), ds.Description)
	}
	fmt.Printf("  %s: %s\n", i18n.T("registry_field_path"), ds.Path)
	if ds.License != "" {
		fmt.Printf("  %s: %s\n", i18n.T("registry_field_license"), ds.License)
	}
	if ds.Valid {
		fmt.Printf("  ✓ %s\n", i18n.T("registry_check_passed"))
	} else {
		for _, e := range ds.Errors {
			fmt.Printf("  ✗ %s: %s\n", i18n.T("registry_error"), e)
		}
	}
}

func repoShortName(repo string) string {
	return strings.TrimPrefix(repo, "github.com/")
}

func loadBuiltinNames() map[string][]string {
	builtinReg, err := registry.Load()
	if err != nil {
		return nil
	}
	return builtinReg.BuiltinSkillNameMap()
}
