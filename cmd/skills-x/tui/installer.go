// Package tui provides terminal interactive UI components
package tui

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/cmd/skills-x/skills"
	"github.com/castle-x/skills-x/pkg/discover"
	"github.com/castle-x/skills-x/pkg/gitutil"
	"github.com/castle-x/skills-x/pkg/registry"
	tea "github.com/charmbracelet/bubbletea"
)

// InstallerModel represents the installation state
type InstallerModel struct {
	skills           []SkillItem
	deselectedSkills []SkillItem // skills to uninstall
	targetDir        string
	currentIdx       int
	completed        int
	failed           int
	quitting         bool
	finished         bool
	err              error
	progressMsg      string
	phase            string // "install" or "uninstall"
	installResults   []string // per-skill result: "ok", "fail", "skip"
	uninstallResults []string // per-skill result for uninstalls
}

// NewInstallerModel creates a new installer model
func NewInstallerModel(skills []SkillItem, deselectedSkills []SkillItem, targetDir string) InstallerModel {
	phase := "install"
	if len(skills) == 0 && len(deselectedSkills) > 0 {
		phase = "uninstall"
	}
	return InstallerModel{
		skills:           skills,
		deselectedSkills: deselectedSkills,
		targetDir:        targetDir,
		currentIdx:       0,
		completed:        0,
		failed:           0,
		phase:            phase,
		installResults:   make([]string, len(skills)),
		uninstallResults: make([]string, len(deselectedSkills)),
	}
}

func (m InstallerModel) Init() tea.Cmd {
	return m.installNext
}

// installNext executes the current install/uninstall operation and returns a progress message.
// All state transitions (phase switching, finished) are handled by Update, not here.
func (m *InstallerModel) installNext() tea.Msg {
	if m.phase == "install" {
		if m.currentIdx >= len(m.skills) {
			return nil
		}

		skill := m.skills[m.currentIdx]

		var progress, result string
		var completed, failed int
		err := m.installSkill(skill)
		if err != nil {
			failed = 1
			result = "fail"
			progress = fmt.Sprintf("Failed (%d/%d): %s - %v", m.currentIdx+1, len(m.skills), skill.FullName, err)
		} else {
			completed = 1
			result = "ok"
			progress = fmt.Sprintf("Installed (%d/%d): %s", m.currentIdx+1, len(m.skills), skill.FullName)
		}

		return installProgressMsg{
			current:      m.currentIdx + 1,
			completedAdd: completed,
			failedAdd:    failed,
			skill:        skill.FullName,
			progress:     progress,
			result:       result,
		}
	}

	if m.phase == "uninstall" {
		if m.currentIdx >= len(m.deselectedSkills) {
			return nil
		}

		skill := m.deselectedSkills[m.currentIdx]

		var progress, result string
		var completed, failed int
		err := m.uninstallSkill(skill)
		if err != nil {
			failed = 1
			result = "fail"
			progress = fmt.Sprintf("Uninstall Failed (%d/%d): %s - %v", m.currentIdx+1, len(m.deselectedSkills), skill.FullName, err)
		} else {
			completed = 1
			result = "ok"
			progress = fmt.Sprintf("Uninstalled (%d/%d): %s", m.currentIdx+1, len(m.deselectedSkills), skill.FullName)
		}

		return installProgressMsg{
			current:      m.currentIdx + 1,
			completedAdd: completed,
			failedAdd:    failed,
			skill:        skill.FullName,
			progress:     progress,
			result:       result,
		}
	}

	return nil
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

func (m InstallerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case installProgressMsg:
		// Record per-skill result
		idx := msg.current - 1
		if m.phase == "install" && idx >= 0 && idx < len(m.installResults) {
			m.installResults[idx] = msg.result
		} else if m.phase == "uninstall" && idx >= 0 && idx < len(m.uninstallResults) {
			m.uninstallResults[idx] = msg.result
		}

		m.currentIdx = msg.current
		m.completed += msg.completedAdd
		m.failed += msg.failedAdd
		m.progressMsg = msg.progress

		if m.phase == "install" && m.currentIdx < len(m.skills) {
			return m, m.installNext
		}
		if m.phase == "uninstall" && m.currentIdx < len(m.deselectedSkills) {
			return m, m.installNext
		}

		if m.phase == "install" && len(m.deselectedSkills) > 0 {
			m.phase = "uninstall"
			m.currentIdx = 0
			return m, m.installNext
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

	totalItems := len(m.skills) + len(m.deselectedSkills)

	b.WriteString(titleStyle.Render("Installing Skills"))
	b.WriteString("\n\n")

	// Progress bar
	barWidth := 40
	currentTotal := m.currentIdx
	if m.phase == "uninstall" {
		currentTotal = len(m.skills) + m.currentIdx
	}
	if m.finished {
		currentTotal = totalItems
	}

	progressPercent := 0
	if totalItems > 0 {
		progressPercent = (currentTotal * 100) / totalItems
	} else {
		progressPercent = 100
	}

	filled := 0
	if totalItems > 0 {
		filled = (currentTotal * barWidth) / totalItems
	} else {
		filled = barWidth
	}

	b.WriteString("Progress: [")
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
	b.WriteString(fmt.Sprintf("Completed: %d | Failed: %d\n", m.completed, m.failed))

	// Install list
	if len(m.skills) > 0 {
		b.WriteString("\n")
		for i, skill := range m.skills {
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

	// Uninstall list
	if len(m.deselectedSkills) > 0 {
		b.WriteString("\n")
		b.WriteString(hintStyle.Render("Uninstalling:"))
		b.WriteString("\n")
		for i, skill := range m.deselectedSkills {
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
			b.WriteString(RenderHint(fmt.Sprintf("Done! %d completed. Press any key to exit.", m.completed)))
		} else {
			b.WriteString(RenderHint(fmt.Sprintf("Done! %d completed, %d failed. Press any key to exit.", m.completed, m.failed)))
		}
	} else {
		b.WriteString(RenderHint("Press q to cancel"))
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

// installSkill installs a single skill
func (m *InstallerModel) installSkill(item SkillItem) error {
	// If target directory is empty, use current directory
	targetDir := m.targetDir
	if targetDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		targetDir = cwd
	}

	// Create target directory if not exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	// Check if it's an X (self-developed) skill
	if item.IsX || skills.XSkillExists(item.Name) {
		return m.installXSkill(item.Name, targetDir)
	}

	// Otherwise, treat as a registry skill
	return m.installRegistrySkill(item, targetDir)
}

// installXSkill installs an x (self-developed) skill from embedded filesystem
func (m *InstallerModel) installXSkill(name, targetDir string) error {
	embedFS, skillPath, ok := skills.GetXSkillFS(name)
	if !ok || skillPath == "" {
		return fmt.Errorf("%s: %s", i18n.T("init_skill_path_not_found"), name)
	}

	dstPath := filepath.Join(targetDir, name)

	os.RemoveAll(dstPath)

	if err := copyFromEmbedFS(embedFS, skillPath, dstPath); err != nil {
		os.RemoveAll(dstPath)
		return err
	}
	return nil
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

// installRegistrySkill installs a skill from the registry
func (m *InstallerModel) installRegistrySkill(item SkillItem, targetDir string) error {
	// Load registry
	reg, err := registry.Load()
	if err != nil {
		return fmt.Errorf("failed to load registry: %w", err)
	}

	// Find skill in registry
	matches := reg.FindSkillsWithConflict(item.Name)
	if len(matches) == 0 {
		return fmt.Errorf("skill not found in registry: %s", item.Name)
	}

	// Find the matching skill
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
		// Fall back to first match
		skill = matches[0].Skill
		source = matches[0].Source
	}

	// Clone the repository
	var result *gitutil.CloneResult

	// For large repos (marked with skip_fetch), use sparse checkout
	if source.SkipFetch && skill.Path != "" {
		result, err = gitutil.SparseCloneRepo(source.GetGitURL(), source.Repo, []string{skill.Path})
	} else {
		result, err = gitutil.CloneRepoWithRefresh(source.GetGitURL(), source.Repo, false)
	}
	if err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	// Find the skill in the cloned repo
	var skillPath string
	if skill.Path != "" {
		skillPath = filepath.Join(result.TempDir, skill.Path)
	} else {
		// Try to discover the skill
		discovered, err := discover.DiscoverSkillByPath(result.TempDir, skill.Name)
		if err != nil || discovered == nil {
			// Fallback: search in common locations
			discovered, _ = findSkillInRepo(result.TempDir, skill.Name)
		}
		if discovered != nil {
			skillPath = discovered.Path
		}
	}

	if skillPath == "" || !dirExists(skillPath) {
		return fmt.Errorf("skill path not found: %s", skill.Name)
	}

	dstPath := filepath.Join(targetDir, skill.Name)

	os.RemoveAll(dstPath)

	if err := copyDir(skillPath, dstPath); err != nil {
		os.RemoveAll(dstPath)
		return fmt.Errorf("copy failed: %w", err)
	}

	return nil
}

// findSkillInRepo searches for a skill by name in common locations
func findSkillInRepo(repoPath string, skillName string) (*discover.DiscoveredSkill, error) {
	// First, try exact path matches
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

	// Fall back to discovery
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

		// Handle symlinks
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

// copyFromEmbedFS copies a directory from embedded filesystem to target path
func copyFromEmbedFS(embedFS fs.FS, srcPath string, dstPath string) error {
	// Remove existing directory
	os.RemoveAll(dstPath)

	// Walk through embedded filesystem
	// Note: embed.FS always uses "/" as separator, regardless of OS
	return fs.WalkDir(embedFS, srcPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path using string operations
		// Note: embed.FS paths always use "/", not filepath functions (which use "\" on Windows)
		relPath := path
		if len(path) > len(srcPath) {
			relPath = strings.TrimPrefix(path, srcPath+"/")
		} else if path == srcPath {
			relPath = "."
		}

		// Convert to OS-specific path for target filesystem
		targetPath := filepath.Join(dstPath, filepath.FromSlash(relPath))

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Read file from embedded filesystem using fs.ReadFile
		data, err := fs.ReadFile(embedFS, path)
		if err != nil {
			return err
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// Write file to target
		return os.WriteFile(targetPath, data, 0644)
	})
}

// RunInstaller runs the installer UI
func RunInstaller(skills []SkillItem, deselectedSkills []SkillItem, targetDir string) (completed, failed int, err error) {
	m := NewInstallerModel(skills, deselectedSkills, targetDir)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return 0, 0, err
	}

	result := finalModel.(InstallerModel)
	return result.Completed(), result.Failed(), result.Error()
}
