// Package tui provides terminal interactive UI components
package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/castle-x/skills-x/pkg/products"
	tea "github.com/charmbracelet/bubbletea"
)

// SkillItem represents a skill in the list
type SkillItem struct {
	Name        string
	FullName    string // "source/skill-name"
	Source      string
	SourceName  string
	Description string
	Installed   bool // installed in target directory
	Selected    bool
	IsX         bool // true if from x (self-developed)
}

// ============================================================================
// Skills Model - Level 2: Select Skills (with search and multi-select)
// ============================================================================

// SkillsModel for selecting skills with search and multi-select
type SkillsModel struct {
	product    *products.Product
	allSkills []SkillItem
	filtered  []SkillItem
	cursor    int
	offset    int
	search    string
	searching bool // 是否处于搜索模式
	quitting  bool
	goBack    bool
	version   string
	pageSize  int
	targetDir string // current working directory for project-level skills
	errMsg    string // error message to display
}

// NewSkillsModel creates a new skills selection model
func NewSkillsModel(product *products.Product, allSkills []SkillItem, version string, targetDir string) SkillsModel {
	// Auto-select skills that are already installed (no need to re-select)
	// Users can toggle selection with space
	for i := range allSkills {
		if allSkills[i].Installed {
			allSkills[i].Selected = true
		}
	}
	return SkillsModel{
		product:    product,
		allSkills: allSkills,
		filtered:  allSkills,
		cursor:    0,
		offset:    0,
		search:    "",
		version:   version,
		pageSize:  10, // 默认显示 10 个
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
			// 支持按来源搜索 (anthropic/) 或模糊搜索 (react)
			if strings.Contains(strings.ToLower(s.FullName), query) ||
				strings.Contains(strings.ToLower(s.Source), query) ||
				strings.Contains(strings.ToLower(s.Name), query) {
				m.filtered = append(m.filtered, s)
			}
		}
	}

	// 保持 cursor 在有效范围
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}

	// 重置偏移
	m.offset = 0
}

// toggleSelect toggles selection at cursor
func (m *SkillsModel) toggleSelect() {
	if m.cursor < 0 || m.cursor >= len(m.filtered) {
		return
	}

	idx := m.cursor
	m.filtered[idx].Selected = !m.filtered[idx].Selected

	// 同步到 allSkills
	fullName := m.filtered[idx].FullName
	for i := range m.allSkills {
		if m.allSkills[i].FullName == fullName {
			m.allSkills[i].Selected = m.filtered[idx].Selected
			break
		}
	}
}

// selectAll selects all visible skills
func (m *SkillsModel) selectAll() {
	for i := range m.filtered {
		m.filtered[i].Selected = true
	}
	// 同步到 allSkills
	for i := range m.allSkills {
		m.allSkills[i].Selected = true
	}
}

// updateOffset updates the scroll offset based on cursor
func (m *SkillsModel) updateOffset() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+m.pageSize {
		m.offset = m.cursor - m.pageSize + 1
	}
}

func (m SkillsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			// 如果在搜索模式，退出搜索模式
			if m.searching {
				m.searching = false
				m.search = ""
				m.filterSkills()
				return m, nil
			}
			// 否则返回上一级
			m.goBack = true
			return m, tea.Quit
		case "b":
			// 非搜索模式时 b 返回
			if !m.searching {
				m.goBack = true
				return m, tea.Quit
			}
			// 搜索模式时 b 作为搜索输入
			m.search += "b"
			m.filterSkills()
			return m, nil
		case "/", "ctrl+f":
			// 进入搜索模式
			m.searching = true
			return m, nil
		case "up":
			// 方向键始终用于导航
			if m.cursor > 0 {
				m.cursor--
				m.updateOffset()
			} else {
				// 循环到底部
				m.cursor = len(m.filtered) - 1
				m.updateOffset()
			}
			return m, nil
		case "down":
			// 方向键始终用于导航
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
				m.updateOffset()
			} else {
				// 循环到顶部
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
		case " ": // 空格键多选
			mm := &m
			mm.toggleSelect()
		case "a", "A": // 全选
			// 非搜索模式时 a/A 是全选
			if !m.searching {
				m.selectAll()
				return m, nil
			}
			// 搜索模式时 a/A 作为搜索输入
			m.search += msg.String()
			m.filterSkills()
		case "enter":
			// 检查是否有新增的技能
			hasNew := false
			for _, s := range m.allSkills {
				if s.Selected && !s.Installed {
					hasNew = true
					break
				}
			}
			// 检查是否有需要卸载的技能
			hasDeselected := false
			for _, s := range m.allSkills {
				if s.Installed && !s.Selected {
					hasDeselected = true
					break
				}
			}
			if !hasNew && !hasDeselected {
				// 没有新增也没有反选，显示错误消息
				m.errMsg = "请用空格选择要安装/卸载的技能，或按Q退出"
				return m, nil
			}
			m.errMsg = ""
			return m, tea.Quit
		case "backspace":
			if len(m.search) > 0 {
				m.search = m.search[:len(m.search)-1]
				m.filterSkills()
			} else if m.searching {
				// 如果搜索框为空且在搜索模式，退出搜索模式
				m.searching = false
			}
		default:
			// 搜索模式下：所有可打印字符都输入搜索框
			if m.searching && len(msg.String()) == 1 {
				m.search += msg.String()
				m.filterSkills()
				return m, nil
			}
			// 非搜索模式时忽略其他按键
		}
	}
	return m, nil
}

func (m SkillsModel) View() string {
	if m.quitting || m.goBack {
		return ""
	}

	var b strings.Builder

	// 1. Logo 区域
	b.WriteString(RenderLogo(m.version))
	b.WriteString("\n")

	// 2. 产品标题
	b.WriteString(titleStyle.Render("选择要安装的 Skills"))
	b.WriteString(" ")
	b.WriteString(hintStyle.Render(m.product.Name))
	b.WriteString("\n")

	// 2.1 当前项目目录
	if m.targetDir != "" {
		b.WriteString(hintStyle.Render("📁 " + m.targetDir))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// 3. 搜索框
	searchPlaceholder := "搜索技能..."
	if m.search != "" {
		searchPlaceholder = m.search
	}
	// 搜索模式指示器
	if m.searching {
		b.WriteString(searchStyle.Render(" / " + searchPlaceholder + " "))
	} else {
		b.WriteString(searchStyle.Render(" 🔍 " + searchPlaceholder))
	}
	b.WriteString("\n")

	// 表头 - 状态列
	b.WriteString(fmt.Sprintf("%s %s\n",
		strings.Repeat(" ", 42),
		titleStyle.Render("已安装")))
	b.WriteString(separatorStyle.Render(strings.Repeat("─", SeparatorWidth)))
	b.WriteString("\n")

	// 4. 技能列表
	start := m.offset
	end := m.offset + m.pageSize
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	for i := start; i < end; i++ {
		s := m.filtered[i]

		// 选择状态
		prefix := "  "
		style := selectableStyle
		checked := "☐"

		if m.cursor == i {
			prefix = cursorStyle.Render("❯ ")
			// 只有未选中时使用选中样式，光标行的已选中项使用不同的样式
			if !s.Selected {
				style = selectedStyle
			}
		}

		if s.Selected {
			checked = "☑"
			// 已选中项使用青蓝色样式 (不是绿色)
			style = selectedStyle
		}

		// 安装状态
		installedStatus := "✗"
		installedStatusStyle := normalStyle
		if s.Installed {
			installedStatus = "✓"
			installedStatusStyle = selectedStyle
		}

		installed := fmt.Sprintf(" %s", installedStatusStyle.Render(installedStatus))

		// 截断过长的名称
		displayName := s.FullName
		if len(displayName) > 40 {
			displayName = displayName[:37] + "..."
		}
		// 补齐空格对齐
		for len([]rune(displayName)) < 40 {
			displayName += " "
		}

		b.WriteString(fmt.Sprintf("%s%s %s%s%s\n", prefix, checked, style.Render(displayName), installed, colorReset))
	}

	// 填充空行保持页面稳定
	for i := end - start; i < m.pageSize; i++ {
		b.WriteString("\n")
	}

	// 5. 滚动提示
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

	// 6. 状态栏
	installedCount := 0
	newSelectedCount := 0
	deselectedCount := 0 // 反选 = 已安装但未选中（将卸载）
	for _, s := range m.allSkills {
		if s.Installed {
			installedCount++
		}
		// 新增 = 选中 且 未安装
		if s.Selected && !s.Installed {
			newSelectedCount++
		}
		// 反选 = 已安装但未选中（将卸载）
		if s.Installed && !s.Selected {
			deselectedCount++
		}
	}
	b.WriteString(RenderStatusBar(len(m.filtered), installedCount, newSelectedCount, deselectedCount))

	// 7. 错误消息
	if m.errMsg != "" {
		b.WriteString("\n")
		b.WriteString(cursorStyle.Render("⚠ " + m.errMsg))
		b.WriteString("\n")
	}

	// 8. 底部提示
	if m.searching {
		b.WriteString(RenderHint("输入搜索 | Esc 退出搜索 | Enter 确认"))
	} else {
		b.WriteString(RenderHint("/ 搜索 | 空格 选择/反选 | A 全选 | Enter 确认安装/卸载 | b 返回 | q 退出"))
	}

	return b.String()
}

// SelectedSkills returns all selected skills
func (m SkillsModel) SelectedSkills() []SkillItem {
	result := make([]SkillItem, 0)
	for _, s := range m.allSkills {
		if s.Selected {
			result = append(result, s)
		}
	}
	return result
}

// DeselectedSkills returns skills that were installed but are now unselected (to be uninstalled)
func (m SkillsModel) DeselectedSkills() []SkillItem {
	result := make([]SkillItem, 0)
	for _, s := range m.allSkills {
		if s.Installed && !s.Selected {
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

// RunSkillsSelect runs the skills selection page
func RunSkillsSelect(product *products.Product, allSkills []SkillItem, version string, targetDir string) ([]SkillItem, []SkillItem, error) {
	// 按 FullName 排序
	sort.Slice(allSkills, func(i, j int) bool {
		return allSkills[i].FullName < allSkills[j].FullName
	})

	m := NewSkillsModel(product, allSkills, version, targetDir)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return nil, nil, err
	}

	result := finalModel.(SkillsModel)
	if result.IsQuitting() {
		return nil, nil, fmt.Errorf("quit")
	}
	if result.IsGoBack() {
		return nil, nil, fmt.Errorf("go back")
	}

	return result.SelectedSkills(), result.DeselectedSkills(), nil
}
