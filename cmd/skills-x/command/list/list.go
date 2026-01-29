// Package list implements the list command
package list

import (
	"fmt"
	"sort"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/cmd/skills-x/skills"
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
)

// Category order for display
var categoryOrder = []string{
	"castle-x", // Castle-X skills first
	"creative", "document", "devtools", "workflow", "git",
	"writing", "integration", "business", "files", "utility", "skilldev", "other",
}

// NewCommand creates the list command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: i18n.T("cmd_list_short"),
		Long:  i18n.T("cmd_list_long"),
		RunE:  runList,
	}
	return cmd
}

// getCategoryName returns the localized category name
func getCategoryName(cat string) string {
	key := "cat_" + cat
	name := i18n.T(key)
	if name == key {
		// Fallback if not found
		return "📦 " + cat
	}
	return name
}

// getSkillDescription returns the localized skill description
func getSkillDescription(skillName string, fallback string) string {
	key := "skill_" + skillName
	desc := i18n.T(key)
	if desc == key {
		// Fallback to original description if not translated
		if fallback != "" {
			return fallback
		}
		return "-"
	}
	return desc
}

func runList(cmd *cobra.Command, args []string) error {
	skillList, err := skills.ListSkills()
	if err != nil {
		return err
	}

	// Group by category
	byCategory := make(map[string][]skills.SkillInfo)
	for _, s := range skillList {
		byCategory[s.Category] = append(byCategory[s.Category], s)
	}

	// Sort skills within each category
	for cat := range byCategory {
		sort.Slice(byCategory[cat], func(i, j int) bool {
			// castle-x skills first
			if byCategory[cat][i].IsCastleX != byCategory[cat][j].IsCastleX {
				return byCategory[cat][i].IsCastleX
			}
			return byCategory[cat][i].Name < byCategory[cat][j].Name
		})
	}

	// Print header
	fmt.Printf("\n%s%s%s %s%s%s\n\n", colorBold, i18n.T("list_header"), colorReset, colorGray, i18n.Tf("list_total", len(skillList)), colorReset)

	// Print by category
	for _, cat := range categoryOrder {
		catSkills, ok := byCategory[cat]
		if !ok || len(catSkills) == 0 {
			continue
		}

		catName := getCategoryName(cat)
		fmt.Printf("%s%s%s\n", colorBold, catName, colorReset)

		for _, s := range catSkills {
			tag := ""
			if s.IsCastleX {
				tag = fmt.Sprintf(" %s%s%s", colorYellow, i18n.T("list_castlex_tag"), colorReset)
			}

			// Get localized description
			desc := getSkillDescription(s.Name, s.Description)

			// Special note for skills-x (meta skill)
			if s.Name == "skills-x" {
				desc = i18n.T("list_skillsx_desc")
			}

			fmt.Printf("  %s%-30s%s %s%s%s%s\n",
				colorCyan, s.Name, colorReset,
				colorGray, desc, colorReset,
				tag)
		}
		fmt.Println()
	}

	return nil
}
