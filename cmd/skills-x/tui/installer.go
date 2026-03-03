// Package tui provides terminal interactive UI components
package tui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/pkg/discover"
	"github.com/castle-x/skills-x/pkg/gitutil"
	"github.com/castle-x/skills-x/pkg/registry"
	tea "github.com/charmbracelet/bubbletea"
)

// InstallerModel represents the installation state
type InstallerModel struct {
	installSkills    []SkillItem
	updateSkills     []SkillItem
	uninstallSkills  []SkillItem
	targetDir        string
	currentIdx       int
	completed        int
	failed           int
	quitting         bool
	finished         bool
	err              error
	progressMsg      string
	phase            string   // "install", "update", or "uninstall"
	installResults   []string // per-skill result: "ok", "fail"
	updateResults    []string
	uninstallResults []string
}

// NewInstallerModel creates a new installer model with three operation lists
func NewInstallerModel(installSkills, uninstallSkills, updateSkills []SkillItem, targetDir string) InstallerModel {
	phase := "install"
	if len(installSkills) == 0 {
		phase = "update"
		if len(updateSkills) == 0 {
			phase = "uninstall"
		}
	}
	return InstallerModel{
		installSkills:    installSkills,
		updateSkills:     updateSkills,
		uninstallSkills:  uninstallSkills,
		targetDir:        targetDir,
		currentIdx:       0,
		completed:        0,
		failed:           0,
		phase:            phase,
		installResults:   make([]string, len(installSkills)),
		updateResults:    make([]string, len(updateSkills)),
		uninstallResults: make([]string, len(uninstallSkills)),
	}
}

func (m InstallerModel) Init() tea.Cmd {
	return m.installNext
}

// installProgressMsg is a message sent during installation progress
type installProgressMsg struct {
	current      int
	completedAdd int
	failedAdd    int
	skill        string
	progress     string
	result       string // "ok" or "fail"
}

// installNext executes the current operation and returns a progress message
func (m *InstallerModel) installNext() tea.Msg {
	switch m.phase {
	case "install":
		if m.currentIdx >= len(m.installSkills) {
			return nil
		}
		skill := m.installSkills[m.currentIdx]
		var progress, result string
		var completed, failed int

		tempDir, err := m.installSkill(skill)
		if err != nil {
			failed = 1
			result = "fail"
			progress = i18n.Tf("tui_installer_fail_install", m.currentIdx+1, len(m.installSkills), skill.FullName, err)
		} else {
			completed = 1
			result = "ok"
			progress = i18n.Tf("tui_installer_progress_install", m.currentIdx+1, len(m.installSkills), skill.FullName)
			m.writeMetaForSkill(skill, tempDir)
		}
		return installProgressMsg{
			current: m.currentIdx + 1, completedAdd: completed, failedAdd: failed,
			skill: skill.FullName, progress: progress, result: result,
		}

	case "update":
		if m.currentIdx >= len(m.updateSkills) {
			return nil
		}
		skill := m.updateSkills[m.currentIdx]
		var progress, result string
		var completed, failed int

		tempDir, err := m.updateSkill(skill)
		if err != nil {
			failed = 1
			result = "fail"
			progress = i18n.Tf("tui_installer_fail_update", m.currentIdx+1, len(m.updateSkills), skill.FullName, err)
		} else {
			completed = 1
			result = "ok"
			progress = i18n.Tf("tui_installer_progress_update", m.currentIdx+1, len(m.updateSkills), skill.FullName)
			m.writeMetaForSkill(skill, tempDir)
		}
		return installProgressMsg{
			current: m.currentIdx + 1, completedAdd: completed, failedAdd: failed,
			skill: skill.FullName, progress: progress, result: result,
		}

	case "uninstall":
		if m.currentIdx >= len(m.uninstallSkills) {
			return nil
		}
		skill := m.uninstallSkills[m.currentIdx]
		var progress, result string
		var completed, failed int

		err := m.uninstallSkill(skill)
		if err != nil {
			failed = 1
			result = "fail"
			progress = i18n.Tf("tui_installer_fail_uninstall", m.currentIdx+1, len(m.uninstallSkills), skill.FullName, err)
		} else {
			completed = 1
			result = "ok"
			progress = i18n.Tf("tui_installer_progress_uninstall", m.currentIdx+1, len(m.uninstallSkills), skill.FullName)
		}
		return installProgressMsg{
			current: m.currentIdx + 1, completedAdd: completed, failedAdd: failed,
			skill: skill.FullName, progress: progress, result: result,
		}
	}
	return nil
}

func (m InstallerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case installProgressMsg:
		idx := msg.current - 1
		switch m.phase {
		case "install":
			if idx >= 0 && idx < len(m.installResults) {
				m.installResults[idx] = msg.result
			}
		case "update":
			if idx >= 0 && idx < len(m.updateResults) {
				m.updateResults[idx] = msg.result
			}
		case "uninstall":
			if idx >= 0 && idx < len(m.uninstallResults) {
				m.uninstallResults[idx] = msg.result
			}
		}

		m.currentIdx = msg.current
		m.completed += msg.completedAdd
		m.failed += msg.failedAdd
		m.progressMsg = msg.progress

		// Continue current phase if not done
		switch m.phase {
		case "install":
			if m.currentIdx < len(m.installSkills) {
				return m, m.installNext
			}
			// Transition to update phase
			if len(m.updateSkills) > 0 {
				m.phase = "update"
				m.currentIdx = 0
				return m, m.installNext
			}
			// Transition to uninstall phase
			if len(m.uninstallSkills) > 0 {
				m.phase = "uninstall"
				m.currentIdx = 0
				return m, m.installNext
			}
		case "update":
			if m.currentIdx < len(m.updateSkills) {
				return m, m.installNext
			}
			// Transition to uninstall phase
			if len(m.uninstallSkills) > 0 {
				m.phase = "uninstall"
				m.currentIdx = 0
				return m, m.installNext
			}
		case "uninstall":
			if m.currentIdx < len(m.uninstallSkills) {
				return m, m.installNext
			}
		}

		m.finished = true
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		if m.finished {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m InstallerModel) View() string {
	var b strings.Builder

	totalItems := len(m.installSkills) + len(m.updateSkills) + len(m.uninstallSkills)

	b.WriteString(titleStyle.Render(i18n.T("tui_installer_title")))
	b.WriteString("\n\n")

	// Progress bar
	barWidth := 40
	currentTotal := m.currentIdx
	switch m.phase {
	case "update":
		currentTotal = len(m.installSkills) + m.currentIdx
	case "uninstall":
		currentTotal = len(m.installSkills) + len(m.updateSkills) + m.currentIdx
	}
	if m.finished {
		currentTotal = totalItems
	}

	progressPercent := 100
	filled := barWidth
	if totalItems > 0 {
		progressPercent = (currentTotal * 100) / totalItems
		filled = (currentTotal * barWidth) / totalItems
	}

	b.WriteString(i18n.T("tui_installer_progress_label") + " [")
	for i := 0; i < barWidth; i++ {
		if i < filled {
			b.WriteString("=")
		} else if i == filled && !m.finished {
			b.WriteString(">")
		} else {
			b.WriteString(" ")
		}
	}
	b.WriteString(fmt.Sprintf("] %d%% (%d/%d)\n", progressPercent, currentTotal, totalItems))
	b.WriteString(i18n.Tf("tui_installer_summary", m.completed, m.failed) + "\n")

	// Install list
	if len(m.installSkills) > 0 {
		b.WriteString("\n")
		for i, skill := range m.installSkills {
			status := "  "
			if i < len(m.installResults) && m.installResults[i] != "" {
				switch m.installResults[i] {
				case "ok":
					status = successStyle.Render("✓ ")
				case "fail":
					status = errorStyle.Render("✗ ")
				}
			} else if m.phase == "install" && i == m.currentIdx {
				status = hintStyle.Render("▸ ")
			}
			b.WriteString(fmt.Sprintf("  %s%s\n", status, skill.FullName))
		}
	}

	// Update list
	if len(m.updateSkills) > 0 {
		b.WriteString("\n")
		b.WriteString(updateStyle.Render(i18n.T("tui_installer_updating")))
		b.WriteString("\n")
		for i, skill := range m.updateSkills {
			status := "  "
			if i < len(m.updateResults) && m.updateResults[i] != "" {
				switch m.updateResults[i] {
				case "ok":
					status = successStyle.Render("✓ ")
				case "fail":
					status = errorStyle.Render("✗ ")
				}
			} else if m.phase == "update" && i == m.currentIdx {
				status = hintStyle.Render("▸ ")
			}
			b.WriteString(fmt.Sprintf("  %s%s\n", status, skill.FullName))
		}
	}

	// Uninstall list
	if len(m.uninstallSkills) > 0 {
		b.WriteString("\n")
		b.WriteString(hintStyle.Render(i18n.T("tui_installer_uninstalling")))
		b.WriteString("\n")
		for i, skill := range m.uninstallSkills {
			status := "  "
			if i < len(m.uninstallResults) && m.uninstallResults[i] != "" {
				switch m.uninstallResults[i] {
				case "ok":
					status = successStyle.Render("✓ ")
				case "fail":
					status = errorStyle.Render("✗ ")
				}
			} else if m.phase == "uninstall" && i == m.currentIdx {
				status = hintStyle.Render("▸ ")
			}
			b.WriteString(fmt.Sprintf("  %s%s\n", status, skill.FullName))
		}
	}

	// Bottom hint
	b.WriteString("\n")
	if m.finished {
		if m.failed == 0 {
			b.WriteString(RenderHint(i18n.Tf("tui_installer_done", m.completed)))
		} else {
			b.WriteString(RenderHint(i18n.Tf("tui_installer_done_failed", m.completed, m.failed)))
		}
	} else {
		b.WriteString(RenderHint(i18n.T("tui_installer_cancel")))
	}

	return b.String()
}

func (m InstallerModel) IsQuitting() bool {
	return m.quitting
}

func (m InstallerModel) IsFinished() bool {
	return m.finished
}

func (m InstallerModel) Completed() int {
	return m.completed
}

func (m InstallerModel) Failed() int {
	return m.failed
}

func (m InstallerModel) Error() error {
	return m.err
}

// writeMetaForSkill writes .skills-x-meta.json after a successful install/update
func (m *InstallerModel) writeMetaForSkill(item SkillItem, tempDir string) {
	dstPath := filepath.Join(m.targetDir, item.Name)

	commit := ""
	if tempDir != "" {
		commit, _ = gitutil.GetRepoHeadCommit(tempDir)
	}

	meta := SkillMeta{
		Skill:       item.Name,
		Source:      item.SourceName,
		Repo:        item.Source,
		Commit:      commit,
		InstalledAt: time.Now().UTC().Format(time.RFC3339),
	}
	_ = WriteSkillMeta(dstPath, meta)
}

// installSkill installs a single skill, returns the temp/cache dir for meta writing
func (m *InstallerModel) installSkill(item SkillItem) (string, error) {
	targetDir := m.targetDir
	if targetDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		targetDir = cwd
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	return m.installRegistrySkill(item, targetDir)
}

// updateSkill updates a single skill using refresh=true
func (m *InstallerModel) updateSkill(item SkillItem) (string, error) {
	targetDir := m.targetDir
	if targetDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		targetDir = cwd
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	return m.installRegistrySkillWithRefresh(item, targetDir, true)
}

// uninstallSkill removes a skill from the target directory
func (m *InstallerModel) uninstallSkill(item SkillItem) error {
	targetDir := m.targetDir
	if targetDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		targetDir = cwd
	}

	skillPath := filepath.Join(targetDir, item.Name)
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(skillPath); err != nil {
		return fmt.Errorf("failed to remove skill directory: %w", err)
	}
	return nil
}

// installRegistrySkill installs a registry skill (refresh=false)
func (m *InstallerModel) installRegistrySkill(item SkillItem, targetDir string) (string, error) {
	return m.installRegistrySkillWithRefresh(item, targetDir, false)
}

// installRegistrySkillWithRefresh installs/updates a registry skill
func (m *InstallerModel) installRegistrySkillWithRefresh(item SkillItem, targetDir string, refresh bool) (string, error) {
	reg, err := loadMergedRegistry()
	if err != nil {
		return "", fmt.Errorf("failed to load registry: %w", err)
	}

	matches := reg.FindSkillsWithConflict(item.Name)
	if len(matches) == 0 {
		return "", fmt.Errorf("skill not found in registry: %s", item.Name)
	}

	var skill *registry.Skill
	var source *registry.Source
	for _, match := range matches {
		if match.Source.Name == item.SourceName {
			skill = match.Skill
			source = match.Source
			break
		}
	}
	if skill == nil || source == nil {
		skill = matches[0].Skill
		source = matches[0].Source
	}

	var result *gitutil.CloneResult
	if source.SkipFetch && skill.Path != "" {
		result, err = gitutil.SparseCloneRepo(source.GetGitURL(), source.Repo, []string{skill.Path})
	} else {
		result, err = gitutil.CloneRepoWithRefresh(source.GetGitURL(), source.Repo, refresh)
	}
	if err != nil {
		return "", fmt.Errorf("clone failed: %w", err)
	}

	var skillPath string
	if skill.Path != "" {
		skillPath = filepath.Join(result.TempDir, skill.Path)
	} else {
		discovered, err := discover.DiscoverSkillByPath(result.TempDir, skill.Name)
		if err != nil || discovered == nil {
			discovered, _ = findSkillInRepo(result.TempDir, skill.Name)
		}
		if discovered != nil {
			skillPath = discovered.Path
		}
	}

	if skillPath == "" || !dirExists(skillPath) {
		return "", fmt.Errorf("skill path not found: %s", skill.Name)
	}

	dstPath := filepath.Join(targetDir, skill.Name)
	os.RemoveAll(dstPath)

	if err := copyDir(skillPath, dstPath); err != nil {
		os.RemoveAll(dstPath)
		return "", fmt.Errorf("copy failed: %w", err)
	}

	return result.TempDir, nil
}

// findSkillInRepo searches for a skill by name in common locations
func findSkillInRepo(repoPath string, skillName string) (*discover.DiscoveredSkill, error) {
	commonPaths := []string{
		filepath.Join("skills", skillName),
		filepath.Join("packages", "docs", "skills", skillName),
		skillName,
	}

	for _, path := range commonPaths {
		fullPath := filepath.Join(repoPath, path)
		if dirExists(fullPath) && fileExists(filepath.Join(fullPath, "SKILL.md")) {
			return &discover.DiscoveredSkill{
				Name: skillName,
				Path: fullPath,
			}, nil
		}
	}

	discovered, err := discover.DiscoverSkills(repoPath, nil)
	if err != nil {
		return nil, err
	}

	for _, d := range discovered {
		if strings.EqualFold(d.Name, skillName) {
			return &d, nil
		}
	}

	return nil, nil
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// copyDir copies a directory from source to destination
func copyDir(srcPath string, dstPath string) error {
	os.RemoveAll(dstPath)

	return filepath.WalkDir(srcPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dstPath, relPath)

		if d.Type()&os.ModeSymlink != 0 {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return nil
			}
			info, err := os.Stat(realPath)
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return copyDir(realPath, targetPath)
			}
			return copyFile(realPath, targetPath)
		}

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		return copyFile(path, targetPath)
	})
}

// copyFile copies a single file
func copyFile(srcPath string, dstPath string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// RunInstaller runs the installer UI with three operation lists
func RunInstaller(installSkills, uninstallSkills, updateSkills []SkillItem, targetDir string) (completed, failed int, err error) {
	m := NewInstallerModel(installSkills, uninstallSkills, updateSkills, targetDir)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return 0, 0, err
	}

	result := finalModel.(InstallerModel)
	return result.Completed(), result.Failed(), result.Error()
}
