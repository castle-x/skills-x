package registry

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/pkg/registry"
	"github.com/castle-x/skills-x/pkg/skillvalidator"
	"github.com/castle-x/skills-x/pkg/userregistry"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	var descFlag string
	var descZhFlag string
	var forceFlag bool
	var allFlag bool

	cmd := &cobra.Command{
		Use:   "add <owner/repo[/skill]> | <local-path> [skill-path]",
		Short: i18n.T("cmd_registry_add_short"),
		Long:  i18n.T("cmd_registry_add_long"),
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Legacy 2-arg mode.
			if len(args) == 2 {
				return runAddSingle(args[0], args[1], descFlag, descZhFlag, forceFlag)
			}

			parsed := skillvalidator.ParseInput(args[0])

			switch parsed.Kind {
			case skillvalidator.InputKindRepoScan:
				return runAddDiscover(parsed.Repo, allFlag, descFlag, descZhFlag, forceFlag)
			case skillvalidator.InputKindSingleSkill:
				if parsed.SkillHint != "" {
					return runAddFind(parsed.Repo, parsed.SkillHint, descFlag, descZhFlag, forceFlag)
				}
				return runAddDiscover(parsed.Repo, allFlag, descFlag, descZhFlag, forceFlag)
			default:
				return runAddSingle(parsed.Repo, "", descFlag, descZhFlag, forceFlag)
			}
		},
	}

	cmd.Flags().StringVar(&descFlag, "desc", "", i18n.T("flag_registry_desc"))
	cmd.Flags().StringVar(&descZhFlag, "desc-zh", "", i18n.T("flag_registry_desc_zh"))
	cmd.Flags().BoolVar(&forceFlag, "force", false, i18n.T("flag_registry_force"))
	cmd.Flags().BoolVar(&allFlag, "all", false, i18n.T("flag_registry_all"))

	return cmd
}

// runAddDiscover scans a repo and lets the user pick skills to add.
func runAddDiscover(repo string, addAll bool, descFlag, descZhFlag string, force bool) error {
	fmt.Printf("%s %s ...\n", i18n.T("registry_scanning"), repoShortName(repo))

	skills, err := skillvalidator.Discover(repo)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_scan_failed"), err)
	}

	if len(skills) == 0 {
		fmt.Println(i18n.Tf("registry_scan_empty", repoShortName(repo)))
		return nil
	}

	builtinNames := loadBuiltinNames()

	// Print discovered list.
	fmt.Printf("\n%s %d %s\n\n", i18n.T("registry_scan_found"), len(skills), i18n.T("registry_scan_skills"))
	for idx, s := range skills {
		desc := s.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		marker := "✓"
		if !s.Valid {
			marker = "✗"
		}
		suffix := ""
		if builtinNames != nil {
			if _, exists := builtinNames[strings.ToLower(s.Name)]; exists {
				suffix = " " + i18n.T("registry_scan_builtin")
			}
		}
		fmt.Printf("  %2d  %-28s  %-50s  %s%s\n", idx+1, s.Name, desc, marker, suffix)
	}
	fmt.Println()

	if addAll {
		return addMultipleSkills(repo, skills, builtinNames, force)
	}

	// Interactive: ask user which to add.
	fmt.Printf("%s\n", i18n.T("registry_add_prompt"))
	fmt.Printf("  %s\n", i18n.T("registry_add_prompt_hint"))
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	if line == "" || strings.ToLower(line) == "q" {
		fmt.Println(i18n.T("registry_add_cancelled"))
		return nil
	}

	if strings.ToLower(line) == "all" || line == "*" {
		return addMultipleSkills(repo, skills, builtinNames, force)
	}

	// Parse numbers.
	indices := parseNumberList(line, len(skills))
	if len(indices) == 0 {
		fmt.Println(i18n.T("registry_add_no_selection"))
		return nil
	}

	var selected []skillvalidator.DiscoveredSkill
	for _, idx := range indices {
		selected = append(selected, skills[idx])
	}
	return addMultipleSkills(repo, selected, builtinNames, force)
}

// runAddFind searches for a specific skill and asks to add it.
func runAddFind(repo, skillHint, descFlag, descZhFlag string, force bool) error {
	fmt.Printf("%s %s ...\n", i18n.T("registry_scanning"), repoShortName(repo))

	ds, err := skillvalidator.FindSkill(repo, skillHint)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_scan_failed"), err)
	}

	if ds == nil {
		// Fallback: prompt user to input path manually.
		fmt.Printf("\n%s\n", i18n.Tf("registry_skill_not_found", skillHint, repoShortName(repo)))
		fmt.Printf("%s\n", i18n.T("registry_input_path_hint"))
		fmt.Print("> ")

		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		manualPath := strings.TrimSpace(line)

		if manualPath == "" {
			return nil
		}

		// Retry with the manual path as an explicit path.
		return runAddSingle("github.com/"+repoShortName(repo), manualPath, descFlag, descZhFlag, force)
	}

	printDiscoveredSkill(ds)

	if !ds.Valid && !force {
		return fmt.Errorf("%s", i18n.T("registry_add_aborted_invalid"))
	}

	fmt.Printf("\n%s (y/n) ", i18n.T("registry_add_confirm"))
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "y") {
		fmt.Println(i18n.T("registry_add_cancelled"))
		return nil
	}

	return addSingleDiscoveredSkill(repo, ds, descFlag, descZhFlag)
}

// runAddSingle is the original direct-path add.
func runAddSingle(repo, path, descFlag, descZhFlag string, force bool) error {
	req := skillvalidator.ValidateRequest{Repo: repo, Path: path}
	fmt.Println(i18n.Tf("registry_checking", req.Repo, req.Path))

	result, err := skillvalidator.Validate(req)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_check_failed"), err)
	}

	printValidateResult(result)

	if !result.Valid && !force {
		return fmt.Errorf("%s", i18n.T("registry_add_aborted_invalid"))
	}

	desc := result.Description
	if descFlag != "" {
		desc = descFlag
	}
	descZh := result.DescriptionZh
	if descZhFlag != "" {
		descZh = descZhFlag
	}

	builtinReg, _ := registry.Load()
	var builtinNames map[string][]string
	if builtinReg != nil {
		builtinNames = builtinReg.BuiltinSkillNameMap()
	}

	ur, err := userregistry.Load()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_load_user_failed"), err)
	}

	addResult, err := ur.Add(req.Repo, req.Path, result.SkillName, desc, descZh, result.License, builtinNames)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_add_failed"), err)
	}

	for _, src := range addResult.ConflictSources {
		fmt.Printf("⚠ %s\n", i18n.Tf("registry_conflict_warn", result.SkillName, src))
	}

	fmt.Printf("✓ %s\n", i18n.Tf("registry_add_success", result.SkillName, addResult.SourceName))
	return nil
}

// addSingleDiscoveredSkill adds one DiscoveredSkill to user registry.
func addSingleDiscoveredSkill(repo string, ds *skillvalidator.DiscoveredSkill, descFlag, descZhFlag string) error {
	desc := ds.Description
	if descFlag != "" {
		desc = descFlag
	}
	descZh := ""
	if descZhFlag != "" {
		descZh = descZhFlag
	}

	builtinNames := loadBuiltinNames()
	ur, err := userregistry.Load()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_load_user_failed"), err)
	}

	addResult, err := ur.Add(repo, ds.Path, ds.Name, desc, descZh, ds.License, builtinNames)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_add_failed"), err)
	}

	for _, src := range addResult.ConflictSources {
		fmt.Printf("⚠ %s\n", i18n.Tf("registry_conflict_warn", ds.Name, src))
	}

	fmt.Printf("✓ %s\n", i18n.Tf("registry_add_success", ds.Name, addResult.SourceName))
	return nil
}

// addMultipleSkills adds a batch of discovered skills.
func addMultipleSkills(repo string, skills []skillvalidator.DiscoveredSkill, builtinNames map[string][]string, force bool) error {
	ur, err := userregistry.Load()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_load_user_failed"), err)
	}

	added := 0
	skipped := 0
	for _, s := range skills {
		if !s.Valid && !force {
			fmt.Printf("  ✗ %s (%s)\n", s.Name, i18n.T("registry_check_failed"))
			skipped++
			continue
		}

		addResult, err := ur.Add(repo, s.Path, s.Name, s.Description, "", s.License, builtinNames)
		if err != nil {
			fmt.Printf("  ✗ %s: %v\n", s.Name, err)
			skipped++
			continue
		}

		for _, src := range addResult.ConflictSources {
			fmt.Printf("  ⚠ %s\n", i18n.Tf("registry_conflict_warn", s.Name, src))
		}

		fmt.Printf("  ✓ %s\n", i18n.Tf("registry_add_success", s.Name, addResult.SourceName))
		added++
	}

	fmt.Printf("\n%s\n", i18n.Tf("registry_add_batch_summary", added, skipped))
	return nil
}

// parseNumberList parses a comma/space separated list of 1-based numbers.
func parseNumberList(input string, max int) []int {
	input = strings.ReplaceAll(input, ",", " ")
	parts := strings.Fields(input)
	var result []int
	seen := make(map[int]bool)
	for _, p := range parts {
		var n int
		if _, err := fmt.Sscanf(p, "%d", &n); err == nil {
			idx := n - 1
			if idx >= 0 && idx < max && !seen[idx] {
				result = append(result, idx)
				seen[idx] = true
			}
		}
	}
	return result
}
