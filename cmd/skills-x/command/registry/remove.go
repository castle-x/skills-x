package registry

import (
	"fmt"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/pkg/userregistry"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <skill-name>",
		Aliases: []string{"rm"},
		Short:   i18n.T("cmd_registry_remove_short"),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			skillName := args[0]

			ur, err := userregistry.Load()
			if err != nil {
				return fmt.Errorf("%s: %w", i18n.T("registry_load_user_failed"), err)
			}

			if err := ur.Remove(skillName); err != nil {
				return fmt.Errorf("%s: %w", i18n.T("registry_remove_failed"), err)
			}

			fmt.Printf("✓ %s\n", i18n.Tf("registry_remove_success", skillName))
			return nil
		},
	}
}
