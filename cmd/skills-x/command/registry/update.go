package registry

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	pkgregistry "github.com/castle-x/skills-x/pkg/registry"
	"github.com/spf13/cobra"
)

const remoteRegistryURL = "https://raw.githubusercontent.com/castle-x/skills-x/main/pkg/registry/registry.yaml"

func newUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: i18n.T("cmd_registry_update_short"),
		Long:  i18n.T("cmd_registry_update_long"),
		Args:  cobra.NoArgs,
		RunE:  runUpdate,
	}
}

func runUpdate(_ *cobra.Command, _ []string) error {
	cachePath, err := pkgregistry.CachedRegistryPath()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_update_path_error"), err)
	}

	fmt.Printf("%s\n", i18n.T("registry_update_fetching"))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(remoteRegistryURL)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_update_fetch_error"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: HTTP %d", i18n.T("registry_update_fetch_error"), resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_update_fetch_error"), err)
	}

	// Validate it parses correctly before saving
	reg, err := pkgregistry.Parse(data)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_update_parse_error"), err)
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_update_save_error"), err)
	}

	if err := os.WriteFile(cachePath, data, 0o644); err != nil {
		return fmt.Errorf("%s: %w", i18n.T("registry_update_save_error"), err)
	}

	fmt.Printf("✓ %s\n", i18n.Tf("registry_update_success", reg.TotalSkillCount(), cachePath))
	return nil
}
