// Package tui provides terminal interactive UI components
package tui

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/castle-x/skills-x/pkg/gitutil"
	"github.com/castle-x/skills-x/pkg/products"
	"github.com/castle-x/skills-x/pkg/registry"
	tea "github.com/charmbracelet/bubbletea"
)

// SkillAction represents the intended operation for a skill
type SkillAction string

const (
	ActionNone      SkillAction = "none"
	ActionInstall   SkillAction = "install"
	ActionUninstall SkillAction = "uninstall"
	ActionUpdate    SkillAction = "update"
)

// SkillItem represents a skill in the list
type SkillItem struct {
	Name        string
	FullName    string // "source/skill-name"
	Source      string
	SourceName  string
	Description string
	Installed   bool        // installed in target directory
	IsX         bool        // true if from x (self-developed)
	Action      SkillAction // intended operation
	Meta        *SkillMeta  // nil if not installed or no meta
	Checking    bool        // true while u-key check is in progress
	HasUpdate   *bool       // nil=unknown, true=has update, false=no update
}

// checkUpdateResultMsg is returned by the async update check command
type checkUpdateResultMsg struct {
	skillFullName string
	hasUpdate     bool
	localCommit   string
	remoteCommit  string
	err           error
}

// ============================================================================
// Skills Model - Level 3: Select Skills (with search and multi-select)
// ============================================================================

// SkillsModel for selecting skills with search and multi-select
type SkillsModel struct {
	product        *products.Product
	allSkills      []SkillItem
	filtered       []SkillItem
	cursor         int
	offset         int
	search         string
	searching      bool // 是否处于搜索模式
	quitting       bool
	goBack         bool
	version        string
	pageSize       int
	targetDir      string // current working directory for project-level skills
	errMsg         string // error message to display
	selectAllState int    // 0=none, 1=install/update, 2=none/uninstall
}

// NewSkillsModel creates a new skills selection model
func NewSkillsModel(product *products.Product, allSkills []SkillItem, version string, targetDir string) SkillsModel {
	for i := range allSkills {
		allSkills[i].Action = ActionNone
	}
	return SkillsModel{
		product:   product,
		allSkills: allSkills,
		filtered:  allSkills,
		cursor:    0,
		offset:    0,
		search:    "",
		version:   version,
		pageSize:  10,
		targetDir: targetDir,
	}
}

func (m SkillsModel) Init() tea.Cmd {
	return nil
}

// filterSkills filters skills based on search query
func (m *SkillsModel) filterSkills() {
	query := strings.ToLower(strings.TrimSpace(m.search))

	if query == "" {
		m.filtered = m.allSkills
	} else {
		m.filtered = make([]SkillItem, 0)
		for _, s := range m.allSkills {
			if strings.Contains(strings.ToLower(s.FullName), query) ||
				strings.Contains(strings.ToLower(s.Source), query) ||
				strings.Contains(strings.ToLower(s.Name), query) {
				m.filtered = append(m.filtered, s)
			}
		}
	}

	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.offset = 0
}

// syncToAllSkills syncs a filtered item's Action back to allSkills
func (m *SkillsModel) syncToAllSkills(fullName string, action SkillAction) {
	for i := range m.allSkills {
		if m.allSkills[i].FullName == fullName {
			m.allSkills[i].Action = action
			break
		}
	}
}

// syncFilteredFromAll copies state from allSkills to the matching filtered items
func (m *SkillsModel) syncFilteredFromAll() {
	lookup := make(map[string]*SkillItem, len(m.allSkills))
	for i := range m.allSkills {
		lookup[m.allSkills[i].FullName] = &m.allSkills[i]
	}
	for i := range m.filtered {
		if src, ok := lookup[m.filtered[i].FullName]; ok {
			m.filtered[i] = *src
		}
	}
}

// toggleSelect implements the space key state machine
func (m *SkillsModel) toggleSelect() {
	if m.cursor < 0 || m.cursor >= len(m.filtered) {
		return
	}

	item := &m.filtered[m.cursor]
	if !item.Installed {
		if item.Action == ActionNone {
			item.Action = ActionInstall
		} else {
			item.Action = ActionNone
		}
	} else {
		switch item.Action {
		case ActionNone:
			item.Action = ActionUninstall
		case ActionUninstall, ActionUpdate:
			item.Action = ActionNone
		}
	}
	m.syncToAllSkills(item.FullName, item.Action)
}

// selectAll implements the A key 3-state cycle
func (m *SkillsModel) selectAll() {
	m.selectAllState = (m.selectAllState + 1) % 3

	switch m.selectAllState {
	case 1: // Not installed -> Install; Installed -> Update
		for i := range m.allSkills {
			if !m.allSkills[i].Installed {
				m.allSkills[i].Action = ActionInstall
			} else {
				m.allSkills[i].Action = ActionUpdate
			}
		}
	case 2: // Not installed -> None; Installed -> Uninstall
		for i := range m.allSkills {
			if !m.allSkills[i].Installed {
				m.allSkills[i].Action = ActionNone
			} else {
				m.allSkills[i].Action = ActionUninstall
			}
		}
	case 0: // All -> None
		for i := range m.allSkills {
			m.allSkills[i].Action = ActionNone
		}
	}

	m.syncFilteredFromAll()
}

// updateOffset updates the scroll offset based on cursor
func (m *SkillsModel) updateOffset() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+m.pageSize {
		m.offset = m.cursor - m.pageSize + 1
	}
}

// checkUpdateForSkill creates a command that checks for updates
func checkUpdateForSkill(item SkillItem, targetDir string) tea.Cmd {
	return func() tea.Msg {
		skillDir := filepath.Join(targetDir, item.Name)
		meta, _ := ReadSkillMeta(skillDir)

		if item.IsX {
			return checkUpdateResultMsg{
				skillFullName: item.FullName,
				err:           fmt.Errorf("内置 skill 随程序更新"),
			}
		}

		reg, err := registry.Load()
		if err != nil {
			return checkUpdateResultMsg{
				skillFullName: item.FullName,
				err:           fmt.Errorf("加载注册表失败: %w", err),
			}
		}

		matches := reg.FindSkillsWithConflict(item.Name)
		if len(matches) == 0 {
			return checkUpdateResultMsg{
				skillFullName: item.FullName,
				err:           fmt.Errorf("注册表中未找到此 skill"),
			}
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
			result, err = gitutil.CloneRepoWithRefresh(source.GetGitURL(), source.Repo, false)
		}
		if err != nil {
			return checkUpdateResultMsg{
				skillFullName: item.FullName,
				err:           fmt.Errorf("获取仓库失败: %w", err),
			}
		}

		remoteCommit, err := gitutil.GetRepoHeadCommit(result.TempDir)
		if err != nil {
			return checkUpdateResultMsg{
				skillFullName: item.FullName,
				err:           fmt.Errorf("获取提交信息失败: %w", err),
			}
		}

		localCommit := ""
		if meta != nil {
			localCommit = meta.Commit
		}

		hasUpdate := localCommit == "" || localCommit != remoteCommit

		return checkUpdateResultMsg{
			skillFullName: item.FullName,
			hasUpdate:     hasUpdate,
			localCommit:   localCommit,
			remoteCommit:  remoteCommit,
		}
	}
}

func (m SkillsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case checkUpdateResultMsg:
		for i := range m.allSkills {
			if m.allSkills[i].FullName == msg.skillFullName {
				m.allSkills[i].Checking = false
				if msg.err != nil {
					m.errMsg = fmt.Sprintf("检测失败: %v", msg.err)
				} else {
					m.allSkills[i].HasUpdate = &msg.hasUpdate
					if msg.hasUpdate {
						m.allSkills[i].Action = ActionUpdate
						if msg.localCommit != "" {
							m.errMsg = fmt.Sprintf("✓ %s 有新版可用 (%s → %s)", m.allSkills[i].Name, msg.localCommit, msg.remoteCommit)
						} else {
							m.errMsg = fmt.Sprintf("✓ %s 有新版可用 (→ %s)", m.allSkills[i].Name, msg.remoteCommit)
						}
					} else {
						m.errMsg = m.allSkills[i].Name + " 已是最新 (" + msg.localCommit + ")"
					}
				}
				break
			}
		}
		m.syncFilteredFromAll()
		return m, nil

	case tea.KeyMsg:
		m.errMsg = ""

		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "q":
			if !m.searching {
				m.quitting = true
				return m, tea.Quit
			}
			m.search += "q"
			m.filterSkills()
			return m, nil
		case "esc":
			if m.searching {
				m.searching = false
				m.search = ""
				m.filterSkills()
				return m, nil
			}
			m.goBack = true
			return m, tea.Quit
		case "b":
			if !m.searching {
				m.goBack = true
				return m, tea.Quit
			}
			m.search += "b"
			m.filterSkills()
			return m, nil
		case "/", "ctrl+f":
			m.searching = true
			return m, nil
		case "up":
			if m.cursor > 0 {
				m.cursor--
				m.updateOffset()
			} else {
				m.cursor = len(m.filtered) - 1
				m.updateOffset()
			}
			return m, nil
		case "down":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				m.updateOffset()
			} else {
				m.cursor = 0
				m.offset = 0
			}
			return m, nil
		case "pgup":
			m.cursor -= m.pageSize
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.updateOffset()
		case "pgdown":
			m.cursor += m.pageSize
			if m.cursor >= len(m.filtered) {
				m.cursor = len(m.filtered) - 1
			}
			m.updateOffset()
		case "home":
			m.cursor = 0
			m.offset = 0
		case "end":
			m.cursor = len(m.filtered) - 1
			m.updateOffset()
		case " ":
			mm := &m
			mm.toggleSelect()
		case "u":
			if m.searching {
				m.search += "u"
				m.filterSkills()
				return m, nil
			}
			if m.cursor >= 0 && m.cursor < len(m.filtered) {
				item := &m.filtered[m.cursor]
				if !item.Installed {
					m.errMsg = "仅已安装 skill 可检测更新"
					return m, nil
				}
				if item.Action == ActionUpdate {
					item.Action = ActionNone
					m.syncToAllSkills(item.FullName, ActionNone)
					return m, nil
				}
				item.Checking = true
				for i := range m.allSkills {
					if m.allSkills[i].FullName == item.FullName {
						m.allSkills[i].Checking = true
						break
					}
				}
				return m, checkUpdateForSkill(*item, m.targetDir)
			}
		case "a", "A":
			if !m.searching {
				m.selectAll()
				return m, nil
			}
			m.search += msg.String()
			m.filterSkills()
		case "enter":
			installCount := 0
			uninstallCount := 0
			updateCount := 0
			for _, s := range m.allSkills {
				switch s.Action {
				case ActionInstall:
					installCount++
				case ActionUninstall:
					uninstallCount++
				case ActionUpdate:
					updateCount++
				}
			}
			if installCount == 0 && uninstallCount == 0 && updateCount == 0 {
				m.errMsg = "请用空格选择要操作的技能，或按 Q 退出"
				return m, nil
			}
			m.errMsg = ""
			return m, tea.Quit
		case "backspace":
			if len(m.search) > 0 {
				m.search = m.search[:len(m.search)-1]
				m.filterSkills()
			} else if m.searching {
				m.searching = false
			}
		default:
			if m.searching && len(msg.String()) == 1 {
				m.search += msg.String()
				m.filterSkills()
				return m, nil
			}
		}
	}
	return m, nil
}

func (m SkillsModel) View() string {
	if m.quitting || m.goBack {
		return ""
	}

	var b strings.Builder

	// 1. Logo
	b.WriteString(RenderLogo(m.version))
	b.WriteString("\n")

	// 2. Title
	b.WriteString(titleStyle.Render("选择要安装的 Skills"))
	b.WriteString(" ")
	b.WriteString(hintStyle.Render(m.product.Name))
	b.WriteString("\n")

	// 2.1 Project directory
	if m.targetDir != "" {
		b.WriteString(hintStyle.Render("📁 " + m.targetDir))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// 3. Search box
	searchPlaceholder := "搜索技能..."
	if m.search != "" {
		searchPlaceholder = m.search
	}
	if m.searching {
		b.WriteString(searchStyle.Render(" / " + searchPlaceholder + " "))
	} else {
		b.WriteString(searchStyle.Render(" 🔍 " + searchPlaceholder))
	}
	b.WriteString("\n")

	// Legend line
	b.WriteString(hintStyle.Render("[ ]未安装  [●]已安装  [+]安装  [-]卸载  [↑]更新  [↻]检测中"))
	b.WriteString("\n")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", SeparatorWidth)))
	b.WriteString("\n")

	// 4. Skill list
	start := m.offset
	end := m.offset + m.pageSize
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := start; i < end; i++ {
		s := m.filtered[i]

		// Cursor prefix
		prefix := "  "
		if m.cursor == i {
			prefix = cursorStyle.Render("❯ ")
		}

		// Action marker
		var marker string
		nameStyle := selectableStyle

		if s.Checking {
			marker = hintStyle.Render("[↻]")
			nameStyle = normalStyle
		} else {
			switch s.Action {
			case ActionNone:
				if s.Installed {
					marker = normalStyle.Render("[●]")
					nameStyle = normalStyle
				} else {
					marker = hintStyle.Render("[ ]")
					nameStyle = selectableStyle
				}
			case ActionInstall:
				marker = successStyle.Render("[+]")
				nameStyle = successStyle
			case ActionUninstall:
				marker = errorStyle.Render("[-]")
				nameStyle = errorStyle
			case ActionUpdate:
				marker = updateStyle.Render("[↑]")
				nameStyle = updateStyle
			}
		}

		if m.cursor == i {
			nameStyle = selectedStyle
		}

		displayName := s.FullName
		if len(displayName) > 35 {
			displayName = displayName[:32] + "..."
		}
		displayName = padRight(displayName, 35)

		// Installation date from meta
		dateStr := ""
		if s.Installed && s.Meta != nil && s.Meta.InstalledAt != "" {
			if t, err := time.Parse(time.RFC3339, s.Meta.InstalledAt); err == nil {
				dateStr = "  " + hintStyle.Render(t.Format("2006-01-02"))
			}
		}

		// Update available indicator
		updateHint := ""
		if s.HasUpdate != nil && *s.HasUpdate {
			updateHint = " " + warningStyle.Render("⚠ 有新版")
		}

		b.WriteString(fmt.Sprintf("%s%s %s%s%s\n", prefix, marker, nameStyle.Render(displayName), dateStr, updateHint))
	}

	// Padding for stable layout
	for i := end - start; i < m.pageSize; i++ {
		b.WriteString("\n")
	}

	// 5. Scroll indicator
	if len(m.filtered) > m.pageSize {
		scrollInfo := fmt.Sprintf("%d/%d", m.cursor+1, len(m.filtered))
		if m.offset > 0 {
			scrollInfo = "↑ " + scrollInfo
		}
		if m.offset+m.pageSize < len(m.filtered) {
			scrollInfo = scrollInfo + " ↓"
		}
		b.WriteString("\n")
		b.WriteString(hintStyle.Render(scrollInfo))
	}

	// 6. Status bar
	installCount := 0
	updateCount := 0
	uninstallCount := 0
	for _, s := range m.allSkills {
		switch s.Action {
		case ActionInstall:
			installCount++
		case ActionUpdate:
			updateCount++
		case ActionUninstall:
			uninstallCount++
		}
	}
	b.WriteString("\n")
	statusText := fmt.Sprintf("安装: %d | 更新: %d | 卸载: %d", installCount, updateCount, uninstallCount)
	b.WriteString(hintStyle.Render(statusText))

	// 7. Error message
	if m.errMsg != "" {
		b.WriteString("\n")
		if strings.HasPrefix(m.errMsg, "✓") {
			b.WriteString(successStyle.Render(m.errMsg))
		} else {
			b.WriteString(cursorStyle.Render("⚠ " + m.errMsg))
		}
	}

	// 8. Bottom hints
	if m.searching {
		b.WriteString(RenderHint("输入搜索 | Esc 退出搜索 | Enter 确认"))
	} else {
		b.WriteString(RenderHint("/ 搜索 | 空格 选择 | u 检测更新 | A 全选 | Enter 确认 | b 返回 | q 退出"))
	}

	return b.String()
}

// InstallSkills returns skills with Action == ActionInstall
func (m SkillsModel) InstallSkills() []SkillItem {
	var result []SkillItem
	for _, s := range m.allSkills {
		if s.Action == ActionInstall {
			result = append(result, s)
		}
	}
	return result
}

// UninstallSkills returns skills with Action == ActionUninstall
func (m SkillsModel) UninstallSkills() []SkillItem {
	var result []SkillItem
	for _, s := range m.allSkills {
		if s.Action == ActionUninstall {
			result = append(result, s)
		}
	}
	return result
}

// UpdateSkills returns skills with Action == ActionUpdate
func (m SkillsModel) UpdateSkills() []SkillItem {
	var result []SkillItem
	for _, s := range m.allSkills {
		if s.Action == ActionUpdate {
			result = append(result, s)
		}
	}
	return result
}

func (m SkillsModel) IsQuitting() bool {
	return m.quitting
}

func (m SkillsModel) IsGoBack() bool {
	return m.goBack
}

// RunSkillsSelect runs the skills selection page, returns three action lists
func RunSkillsSelect(product *products.Product, allSkills []SkillItem, version string, targetDir string) (install, uninstall, update []SkillItem, err error) {
	sort.Slice(allSkills, func(i, j int) bool {
		return allSkills[i].FullName < allSkills[j].FullName
	})

	m := NewSkillsModel(product, allSkills, version, targetDir)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return nil, nil, nil, err
	}

	result := finalModel.(SkillsModel)
	if result.IsQuitting() {
		return nil, nil, nil, fmt.Errorf("quit")
	}
	if result.IsGoBack() {
		return nil, nil, nil, fmt.Errorf("go back")
	}

	return result.InstallSkills(), result.UninstallSkills(), result.UpdateSkills(), nil
}
