// Package registry implements the "skills-x registry" subcommand group.
package registry

import (
	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/spf13/cobra"
)

// NewCommand returns the root "registry" command with all subcommands attached.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: i18n.T("cmd_registry_short"),
		Long:  i18n.T("cmd_registry_long"),
	}

	cmd.AddCommand(newCheckCommand())
	cmd.AddCommand(newAddCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newRemoveCommand())
	cmd.AddCommand(newUpdateCommand())

	return cmd
}
