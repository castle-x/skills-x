package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/castle-x/skills-x/cmd/skills-x/command/initcmd"
	"github.com/castle-x/skills-x/cmd/skills-x/command/list"
	"github.com/castle-x/skills-x/cmd/skills-x/errmsg"
	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/cmd/skills-x/tui"
	"github.com/castle-x/skills-x/pkg/versioncheck"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			// Run TUI by default when no arguments provided
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("%s: %w", i18n.T("error_get_cwd"), err)
			}
			opts := tui.TUIOptions{
				Version:   Version,
				TargetDir: cwd,
			}
			return tui.RunTUI(opts)
		},
	}

	// Register subcommands
	rootCmd.AddCommand(list.NewCommand())    // list
	rootCmd.AddCommand(initcmd.NewCommand()) // init
	rootCmd.AddCommand(tui.TUICommand(Version)) // tui

	// Disable cobra's default error output
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	// Set version template
	rootCmd.SetVersionTemplate("skills-x version {{.Version}}\n")

	os.Exit(run(rootCmd, func() { checkForUpdate(Version) }))
}

func checkForUpdate(currentVersion string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	latest, err := versioncheck.FetchLatestVersion(ctx, "skills-x")
	if err != nil {
		return
	}
	if versioncheck.ShouldNotify(currentVersion, latest) {
		fmt.Println(i18n.Tf("update_available", latest, currentVersion))
		fmt.Println(i18n.T("update_command"))
	}
}

func run(rootCmd *cobra.Command, postRun func()) int {
	err := rootCmd.Execute()
	if err != nil {
		// Use custom error format
		errmsg.PrintError(err)
	}
	if postRun != nil {
		postRun()
	}
	if err != nil {
		return 1
	}
	return 0
}
