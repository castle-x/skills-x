package registry

import (
	"fmt"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/pkg/userregistry"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: i18n.T("cmd_registry_list_short"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ur, err := userregistry.Load()
			if err != nil {
				return fmt.Errorf("%s: %w", i18n.T("registry_load_user_failed"), err)
			}

			if ur.IsEmpty() {
				fmt.Println(i18n.T("registry_list_empty"))
				fmt.Printf("  %s\n", i18n.T("registry_list_empty_hint"))
				return nil
			}

			fmt.Printf("%s (%d)\n", i18n.T("registry_list_header"), ur.TotalSkillCount())
			fmt.Println("────────────────────────────────────────────────")

			skills := ur.ListAll()
			for _, s := range skills {
				desc := s.Description
				if len(desc) > 50 {
					desc = desc[:47] + "..."
				}
				fmt.Printf("  %-30s  %s\n", s.Name, desc)
				fmt.Printf("  %s%-30s  %s: %s\n", "  ", "", i18n.T("registry_list_source"), s.Repo)
			}
			return nil
		},
	}
}
