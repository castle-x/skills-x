package initcmd

import (
	"path/filepath"
	"testing"

	"github.com/castle-x/skills-x/cmd/skills-x/skills"
	"github.com/castle-x/skills-x/pkg/registry"
)

func TestInitAllIncludeXInstallsXSkills(t *testing.T) {
	cmd := NewCommand()
	if err := cmd.Flags().Set("include-x", "true"); err != nil {
		t.Fatalf("expected include-x flag to be supported: %v", err)
	}

	tmpDir := t.TempDir()
	reg := &registry.Registry{Sources: map[string]*registry.Source{}}

	if err := initAll(reg, tmpDir); err != nil {
		t.Fatalf("initAll failed: %v", err)
	}

	xSkills, err := skills.ListXSkills()
	if err != nil {
		t.Fatalf("ListXSkills failed: %v", err)
	}

	for _, s := range xSkills {
		path := filepath.Join(tmpDir, s.Name)
		if !dirExists(path) {
			t.Fatalf("expected x skill to be installed: %s", s.Name)
		}
	}
}
