// Package tui provides terminal interactive UI components
package tui

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
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
	Tags        []string    // tags for filtering (e.g., featured, web-frontend)
	Installed   bool        // installed in target directory
	IsX         bool        // true if from x (self-developed)
	Action      SkillAction // intended operation
	Meta        *SkillMeta  // nil if not installed or no meta
	Checking    bool        // true while u-key check is in progress
	HasUpdate   *bool       // nil=unknown, true=has update, false=no update
	Starred     bool        // persisted in ~/.config/skills-x/starred.json
}

// checkUpdateResultMsg is returned by the async update check command
type checkUpdateResultMsg struct {
	skillFullName string
	hasUpdate     bool
	localCommit   string
	remoteCommit  string
	err           error
}

// tagAliases maps Chinese search terms to English tag identifiers
var tagAliases = map[string]string{
	"常用":     "featured",
	"ai效能":   "ai-efficiency",
	"规划":     "planning",
	"前端":     "web-frontend",
	"小程序":   "mobile",
	"后端":     "backend",
	"测试":     "testing",
	"审查":     "code-review",
	"文件":     "office",
	"设计":     "design",
	"写作":     "writing",
	"多媒体":   "media",
	"skills":  "skills-meta",
}

// allTagNames lists English tag names for direct match
var allTagNames = map[string]bool{
	"starred": true, "featured": true, "ai-efficiency": true, "planning": true,
	"web-frontend": true, "mobile": true, "backend": true,
	"testing": true, "code-review": true, "office": true,
	"design": true, "writing": true, "media": true,
	"skills-meta": true,
}

// getTagPickerList returns the ordered tag list with i18n labels at call time.
// values are English tag names matching allTagNames for direct filterSkills resolution.
func getTagPickerList() []struct{ label, value string } {
	return []struct{ label, value string }{
		{i18n.T("tui_tag_starred"), "starred"},
		{i18n.T("tui_tag_featured"), "featured"},
		{i18n.T("tui_tag_ai_efficiency"), "ai-efficiency"},
		{i18n.T("tui_tag_planning"), "planning"},
		{i18n.T("tui_tag_frontend"), "web-frontend"},
		{i18n.T("tui_tag_mobile"), "mobile"},
		{i18n.T("tui_tag_backend"), "backend"},
		{i18n.T("tui_tag_testing"), "testing"},
		{i18n.T("tui_tag_code_review"), "code-review"},
		{i18n.T("tui_tag_office"), "office"},
		{i18n.T("tui_tag_design"), "design"},
		{i18n.T("tui_tag_writing"), "writing"},
		{i18n.T("tui_tag_media"), "media"},
		{i18n.T("tui_tag_skills_meta"), "skills-meta"},
	}
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
	tagPicking     bool // true when # category picker is active
	tagCursor      int  // cursor position in tag picker list
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
// Supports #tag prefix for tag-based filtering (e.g., #前端, #featured)
func (m *SkillsModel) filterSkills() {
	query := strings.TrimSpace(m.search)

	if query == "" {
		m.filtered = m.allSkills
	} else if strings.HasPrefix(query, "#") {
		tagQuery := strings.ToLower(query[1:])
		if tagQuery == "" {
			m.filtered = m.allSkills
		} else {
			resolvedTag := tagQuery
			if mapped, ok := tagAliases[tagQuery]; ok {
				resolvedTag = mapped
			}
			m.filtered = make([]SkillItem, 0)
			if resolvedTag == "starred" {
				for _, s := range m.allSkills {
					if s.Starred {
						m.filtered = append(m.filtered, s)
					}
				}
			} else {
				for _, s := range m.allSkills {
					for _, t := range s.Tags {
						if strings.ToLower(t) == resolvedTag {
							m.filtered = append(m.filtered, s)
							break
						}
					}
				}
			}
		}
	} else {
		queryLower := strings.ToLower(query)
		m.filtered = make([]SkillItem, 0)
		for _, s := range m.allSkills {
			if strings.Contains(strings.ToLower(s.FullName), queryLower) ||
				strings.Contains(strings.ToLower(s.Source), queryLower) ||
				strings.Contains(strings.ToLower(s.Name), queryLower) {
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
// toggleStarred toggles the star/favorite status of the currently focused skill
// and persists the change immediately.
func (m *SkillsModel) toggleStarred() {
	if m.cursor < 0 || m.cursor >= len(m.filtered) {
		return
	}

	targetName := m.filtered[m.cursor].FullName
	newStarred := !m.filtered[m.cursor].Starred

	// Update in filtered list
	m.filtered[m.cursor].Starred = newStarred

	// Sync to allSkills
	for i := range m.allSkills {
		if m.allSkills[i].FullName == targetName {
			m.allSkills[i].Starred = newStarred
			break
		}
	}

	// Build updated set and persist
	set := make(map[string]bool)
	for _, s := range m.allSkills {
		if s.Starred {
			set[s.FullName] = true
		}
	}
	_ = SaveStarred(set)

	// Re-sort allSkills so starred ones move to front
	SortSkills(m.allSkills)
	// Re-apply filter to reflect new order
	m.filterSkills()
	// Find new cursor position for the just-toggled skill
	for i, s := range m.filtered {
		if s.FullName == targetName {
			m.cursor = i
			break
		}
	}
}

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
				err:           fmt.Errorf("%s", i18n.T("tui_builtin_skill_msg")),
			}
		}

		reg, err := registry.Load()
		if err != nil {
			return checkUpdateResultMsg{
				skillFullName: item.FullName,
				err:           fmt.Errorf(i18n.T("tui_err_fetch_repo"), err),
			}
		}

		matches := reg.FindSkillsWithConflict(item.Name)
		if len(matches) == 0 {
			return checkUpdateResultMsg{
				skillFullName: item.FullName,
				err:           fmt.Errorf("%s", i18n.T("tui_err_not_in_registry")),
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
				err:           fmt.Errorf(i18n.T("tui_err_fetch_repo"), err),
			}
		}

		remoteCommit, err := gitutil.GetRepoHeadCommit(result.TempDir)
		if err != nil {
			return checkUpdateResultMsg{
				skillFullName: item.FullName,
				err:           fmt.Errorf(i18n.T("tui_err_get_commit"), err),
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
					m.errMsg = i18n.Tf("tui_check_failed", msg.err)
				} else {
					m.allSkills[i].HasUpdate = &msg.hasUpdate
					if msg.hasUpdate {
						m.allSkills[i].Action = ActionUpdate
						if msg.localCommit != "" {
							m.errMsg = i18n.Tf("tui_update_available_fmt", m.allSkills[i].Name, msg.localCommit, msg.remoteCommit)
						} else {
							m.errMsg = i18n.Tf("tui_update_available_new", m.allSkills[i].Name, msg.remoteCommit)
						}
					} else {
						m.errMsg = i18n.Tf("tui_update_up_to_date", m.allSkills[i].Name, msg.localCommit)
					}
				}
				break
			}
		}
		m.syncFilteredFromAll()
		return m, nil

	case tea.KeyMsg:
		m.errMsg = ""

		// Tag picker mode: intercept all navigation before the main switch
		if m.tagPicking {
			tagList := getTagPickerList()
			switch msg.String() {
			case "up":
				if m.tagCursor > 0 {
					m.tagCursor--
				}
			case "down":
				if m.tagCursor < len(tagList)-1 {
					m.tagCursor++
				}
			case "enter":
				selected := tagList[m.tagCursor]
				m.search = "#" + selected.value
				m.tagPicking = false
				m.filterSkills()
			case "esc", "backspace":
				m.tagPicking = false
				m.search = ""
				m.filterSkills()
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}
			return m, nil
		}

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
		case "/", "ctrl+f", "、": // 、is the Chinese input method equivalent of /
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
				m.errMsg = i18n.T("tui_only_installed_check")
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
		case "f":
			if m.searching {
				m.search += "f"
				m.filterSkills()
				return m, nil
			}
			m.toggleStarred()
			return m, nil
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
			m.errMsg = i18n.T("tui_select_required")
			return m, nil
		}
			m.errMsg = ""
			return m, tea.Quit
		case "backspace":
			if len(m.search) > 0 {
				runes := []rune(m.search)
				m.search = string(runes[:len(runes)-1])
				m.filterSkills()
			} else if m.searching {
				m.searching = false
			}
		default:
			if m.searching {
				s := msg.String()
				if s == "#" {
					m.tagPicking = true
					m.tagCursor = 0
					return m, nil
				}
				if len(s) > 0 && s != " " {
					m.search += s
					m.filterSkills()
					return m, nil
				}
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

	// 2. Title + path on one line
	titleLine := titleStyle.Render(i18n.T("tui_skills_for") + " " + m.product.Name)
	if m.targetDir != "" {
		titleLine += "  " + hintStyle.Render(m.targetDir)
	}
	b.WriteString(titleLine)
	b.WriteString("\n")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", SeparatorWidth)))
	b.WriteString("\n")

	// 3. Search box
	searchPlaceholder := i18n.T("tui_search_placeholder")
	if m.search != "" {
		searchPlaceholder = m.search
	}
	if m.searching || m.tagPicking {
		tagSuffix := ""
		if m.tagPicking {
			tagSuffix = "#"
		}
		b.WriteString(searchStyle.Render(" / " + searchPlaceholder + tagSuffix + " "))
	} else {
		b.WriteString(searchStyle.Render(i18n.T("tui_search_idle")))
	}
	b.WriteString("\n")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", SeparatorWidth)))
	b.WriteString("\n")

	// Legend line
	b.WriteString(hintStyle.Render(i18n.T("tui_legend")))
	b.WriteString("\n")
	b.WriteString(separatorStyle.Render(strings.Repeat("─", SeparatorWidth)))
	b.WriteString("\n")

	// 4a. Tag picker (replaces skill list when active)
	if m.tagPicking {
		tagList := getTagPickerList()
		for i, tag := range tagList {
			prefix := "  "
			style := hintStyle
			if i == m.tagCursor {
				prefix = cursorStyle.Render("❯ ")
				style = selectedStyle
			}
			b.WriteString(fmt.Sprintf("%s%s\n", prefix, style.Render(tag.label)))
		}
		// Padding to fill remaining pageSize rows
		for i := len(tagList); i < m.pageSize; i++ {
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(hintStyle.Render(fmt.Sprintf("%d/%d", m.tagCursor+1, len(tagList))))
		b.WriteString("\n")
		b.WriteString(hintStyle.Render(i18n.Tf("tui_status_ops", 0, 0, 0)))
		b.WriteString(RenderHint(i18n.T("tui_tag_picker_hint")))
		return b.String()
	}

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
			updateHint = " " + warningStyle.Render(i18n.T("tui_update_badge"))
		}

		// Star indicator
		starHint := ""
		if s.Starred {
			starHint = " " + warningStyle.Render("★")
		}

		b.WriteString(fmt.Sprintf("%s%s %s%s%s%s\n", prefix, marker, nameStyle.Render(displayName), dateStr, updateHint, starHint))
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
	statusText := i18n.Tf("tui_status_ops", installCount, updateCount, uninstallCount)
	b.WriteString(hintStyle.Render(statusText))

	// 7. Info / error message area
	b.WriteString("\n")
	if m.errMsg != "" {
		if strings.HasPrefix(m.errMsg, "✓") {
			b.WriteString(successStyle.Render(m.errMsg))
		} else {
			b.WriteString(cursorStyle.Render("⚠ " + m.errMsg))
		}
	} else if m.cursor >= 0 && m.cursor < len(m.filtered) {
		desc := m.filtered[m.cursor].Description
		if desc != "" {
			b.WriteString(RenderDescriptionGradient(desc))
		}
	}

	// 8. Bottom hints
	if m.searching {
		b.WriteString("\n")
		b.WriteString(hintStyle.Render(i18n.T("tui_tag_search_hint")))
		b.WriteString(RenderHint(i18n.T("tui_hint_searching")))
	} else {
		b.WriteString(RenderHint(i18n.T("tui_hint_main")))
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
	// Sort with starred skills first, then alphabetical
	SortSkills(allSkills)

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
