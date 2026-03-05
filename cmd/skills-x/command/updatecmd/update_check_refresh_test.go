package updatecmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/castle-x/skills-x/pkg/gitutil"
)

func TestRunUpdate_CheckModeRefreshesRepoCache(t *testing.T) {
	targetDir := t.TempDir()
	skillDir := filepath.Join(targetDir, "gve")
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# test\n"), 0644); err != nil {
		t.Fatalf("write SKILL.md: %v", err)
	}
	meta := []byte(`{"skill":"gve","source":"castle-x-gve","repo":"github.com/castle-x/gve","commit":"old1234","installed_at":"2026-03-05T00:00:00Z"}`)
	if err := os.WriteFile(filepath.Join(skillDir, ".skills-x-meta.json"), meta, 0644); err != nil {
		t.Fatalf("write meta: %v", err)
	}

	origAll, origCheck, origTarget := flagAll, flagCheck, flagTarget
	origClone := cloneRepoWithRefresh
	origHead := getRepoHeadCommit
	defer func() {
		flagAll, flagCheck, flagTarget = origAll, origCheck, origTarget
		cloneRepoWithRefresh = origClone
		getRepoHeadCommit = origHead
	}()

	flagAll = false
	flagCheck = true
	flagTarget = targetDir

	refreshCalled := false
	cloneRepoWithRefresh = func(gitURL string, repoName string, refresh bool) (*gitutil.CloneResult, error) {
		if repoName == "github.com/castle-x/gve" {
			refreshCalled = refresh
		}
		return &gitutil.CloneResult{TempDir: t.TempDir(), Repo: repoName}, nil
	}
	getRepoHeadCommit = func(repoDir string) (string, error) {
		return "new9999", nil
	}

	if err := runUpdate(nil, []string{"gve"}); err != nil {
		t.Fatalf("runUpdate returned error: %v", err)
	}

	if !refreshCalled {
		t.Fatalf("expected check mode to refresh repository cache")
	}
}

