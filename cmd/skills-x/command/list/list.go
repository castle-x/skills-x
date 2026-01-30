// Package list implements the list command
package list

import (
	"fmt"
	"sort"
	"strings"

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
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
	colorBlue   = "\033[34m"
	colorDim    = "\033[2m"
)

var (
	flagVerbose bool
	flagFetch   bool // changed: default to NOT fetch, use --fetch to enable
)

// NewCommand creates the list command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   i18n.T("cmd_list_short"),
		Long:    i18n.T("cmd_list_long"),
		RunE:    runList,
	}

	cmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, i18n.T("cmd_list_flag_verbose"))
	cmd.Flags().BoolVar(&flagFetch, "fetch", false, i18n.T("cmd_list_flag_fetch"))

	return cmd
}

// sourceSkills holds skills grouped by source
type sourceSkills struct {
	Source *registry.Source
	Skills []skillDisplay
}

// skillDisplay holds skill display information
type skillDisplay struct {
	Name        string
	Description string
	Version     string
	FromRepo    bool // true if dynamically fetched from repo
}

func runList(cmd *cobra.Command, args []string) error {
	// Load registry
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Title removed as requested - no need to display "Available Skills from Registry" or "注册表中的 Skills"

	// Get all sources
	sources := reg.GetAllSources()
	
	// Sort sources by name
	sort.Slice(sources, func(i, j int) bool {
		// Put "anthropic" first, then alphabetically
		if sources[i].Name == "anthropic" {
			return true
		}
		if sources[j].Name == "anthropic" {
			return false
		}
		return sources[i].Name < sources[j].Name
	})

	totalSkills := 0
	totalSources := 0

	for _, source := range sources {
		skills, err := getSkillsForSource(source, flagFetch)
		if err != nil {
			// Print error but continue
			fmt.Printf("%s⚠ %s: %v%s\n\n", colorYellow, source.Repo, err, colorReset)
			continue
		}

		if len(skills) == 0 {
			continue
		}

		totalSources++
		totalSkills += len(skills)

		// Print source header
		printSourceHeader(source)

		// Print skills
		for _, skill := range skills {
			printSkill(skill)
		}
		fmt.Println()
	}

	// Print x (self-developed) skills if available
	xSkills := getXSkills()
	if len(xSkills) > 0 {
		totalSources++
		totalSkills += len(xSkills)
		
		fmt.Printf("%s📦 %sskills-x%s %s(Original)%s\n",
			colorBold, colorCyan, colorReset, colorGray, colorReset)
		
		for _, skill := range xSkills {
			printSkill(skill)
		}
		fmt.Println()
	}

	// Print summary
	fmt.Printf("%s%s%s\n", colorGray, i18n.Tf("list_summary", totalSkills, totalSources), colorReset)

	return nil
}

// getSkillsForSource gets skills for a source
// If fetch is true, it will clone the repo and discover skills dynamically
// Otherwise, it uses the registry data
func getSkillsForSource(source *registry.Source, fetch bool) ([]skillDisplay, error) {
	if !fetch {
		// Use registry data only
		return getSkillsFromRegistry(source), nil
	}

	// Try to fetch from repo
	skills, err := fetchSkillsFromRepo(source)
	if err != nil {
		// Fallback to registry data
		if flagVerbose {
			fmt.Printf("%s  (fallback to registry: %v)%s\n", colorDim, err, colorReset)
		}
		return getSkillsFromRegistry(source), nil
	}

	return skills, nil
}

// getSkillsFromRegistry returns skills from registry definition
func getSkillsFromRegistry(source *registry.Source) []skillDisplay {
	lang := i18n.GetLanguage()
	skills := make([]skillDisplay, 0, len(source.Skills))
	for _, s := range source.Skills {
		skills = append(skills, skillDisplay{
			Name:        s.Name,
			Description: s.GetDescription(lang),
			Version:     s.Version,
			FromRepo:    false,
		})
	}
	
	// Sort by name
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})
	
	return skills
}

// fetchSkillsFromRepo clones the repo and discovers skills
func fetchSkillsFromRepo(source *registry.Source) ([]skillDisplay, error) {
	gitURL := source.GetGitURL()

	// Show progress
	fmt.Printf("%s  %s %s...%s", colorDim, i18n.T("list_fetching"), source.GetRepoShortName(), colorReset)

	// Clone or use cached
	result, err := gitutil.CloneRepo(gitURL, source.Repo)
	if err != nil {
		fmt.Printf("\r%s  %s %s ✗%s\n", colorYellow, i18n.T("list_fetch_failed"), source.GetRepoShortName(), colorReset)
		return nil, err
	}

	// Clear progress line
	fmt.Printf("\r%s\r", strings.Repeat(" ", 60))

	// Discover skills
	discovered, err := discover.DiscoverSkills(result.TempDir, nil)
	if err != nil {
		return nil, err
	}

	// Convert to display format
	skills := make([]skillDisplay, 0, len(discovered))
	for _, d := range discovered {
		skills = append(skills, skillDisplay{
			Name:        d.Name,
			Description: d.Description,
			Version:     d.Version,
			FromRepo:    true,
		})
	}

	// Sort by name
	sort.Slice(skills, func(i, j int) bool {
		return skills[i].Name < skills[j].Name
	})

	return skills, nil
}

// getXSkills returns x (self-developed) skills
func getXSkills() []skillDisplay {
	// Load from embedded x/ directory via skills package
	xSkillsList, err := skills.ListXSkills()
	if err != nil {
		return nil
	}

	result := make([]skillDisplay, 0, len(xSkillsList))
	for _, s := range xSkillsList {
		result = append(result, skillDisplay{
			Name:        s.Name,
			Description: s.Description,
			FromRepo:    false,
		})
	}
	return result
}

// printSourceHeader prints the source header
func printSourceHeader(source *registry.Source) {
	license := ""
	if source.License != "" {
		license = fmt.Sprintf(" %s(%s)%s", colorGray, source.License, colorReset)
	}

	fmt.Printf("%s📦 %s%s%s%s\n",
		colorBold, colorCyan, source.Repo, colorReset, license)
}

// printSkill prints a skill entry
func printSkill(skill skillDisplay) {
	version := ""
	if skill.Version != "" {
		version = fmt.Sprintf(" %s(%s)%s", colorGreen, skill.Version, colorReset)
	}

	desc := skill.Description
	if desc == "" {
		desc = "-"
	}

	// Truncate description if too long
	maxDescLen := 50
	if len(desc) > maxDescLen {
		desc = desc[:maxDescLen-3] + "..."
	}

	fmt.Printf("   %s%-35s%s%s %s%s%s\n",
		colorCyan, skill.Name, colorReset,
		version,
		colorGray, desc, colorReset)
}
