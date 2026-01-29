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

// Category display names (i18n keys)
var categoryNames = map[string]string{
	"creative":    "🎨 Creative & Design",
	"document":    "📄 Document Processing",
	"devtools":    "🛠️  Development Tools",
	"workflow":    "🔄 Workflows",
	"git":         "📝 Git & Code Review",
	"writing":     "✍️  Writing",
	"integration": "🔗 Integrations",
	"business":    "📊 Business & Analytics",
	"files":       "🗂️  File Management",
	"utility":     "🎲 Utilities",
	"skilldev":    "🧰 Skills Development",
	"other":       "📦 Other",
	"castle-x":    "🏰 Castle-X (Original)",
}

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

		catName := categoryNames[cat]
		fmt.Printf("%s%s%s\n", colorBold, catName, colorReset)

		for _, s := range catSkills {
			tag := ""
			if s.IsCastleX {
				tag = fmt.Sprintf(" %s⭐ 作者自研%s", colorYellow, colorReset)
			}

			desc := s.Description
			if desc == "" {
				desc = "-"
			}

			// Special note for skills-x (meta skill)
			if s.Name == "skills-x" {
				desc = "🔄 套娃! Contribution guide (not for regular use)"
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
