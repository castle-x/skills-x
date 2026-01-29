package main

import (
	"os"

	"github.com/castle-x/skills-x/cmd/skills-x/command/initcmd"
	"github.com/castle-x/skills-x/cmd/skills-x/command/list"
	"github.com/castle-x/skills-x/cmd/skills-x/errmsg"
	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/spf13/cobra"
)

// Version and build info (set by ldflags)
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// Initialize i18n
	i18n.MustInit()

	rootCmd := &cobra.Command{
		Use:     "skills-x",
		Short:   i18n.T("app_desc"),
		Long:    i18n.T("app_long_desc"),
		Version: Version,
	}

	// Register subcommands
	rootCmd.AddCommand(list.NewCommand())    // list
	rootCmd.AddCommand(initcmd.NewCommand()) // init

	// Disable cobra's default error output
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	// Set version template
	rootCmd.SetVersionTemplate("skills-x version {{.Version}}\n")

	if err := rootCmd.Execute(); err != nil {
		// Use custom error format
		errmsg.PrintError(err)
		os.Exit(1)
	}
}
