// Package tui provides terminal interactive UI components
package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
	"github.com/castle-x/skills-x/pkg/products"
	tea "github.com/charmbracelet/bubbletea"
)

// ============================================================================
// Product Select Model - Level 1: Select AI Tool
// ============================================================================

// ProductModel for selecting AI tool
type ProductModel struct {
	products  []products.Product
	cursor    int
	quitting  bool
	goBack    bool
	version   string
	targetDir string // current working directory for project-level check
}

// NewProductModel creates a new product selection model
func NewProductModel(version string, targetDir string) ProductModel {
	return ProductModel{
		products:  products.AllProducts,
		cursor:    0,
		version:   version,
		targetDir: targetDir,
	}
}

// countSkillsInDir counts valid skills in a directory.
// Uses isSkillDir from skills_list.go for consistent detection across all pages.
func countSkillsInDir(dir string) int {
	if dir == "" {
		return 0
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())
		if isSkillDir(path) {
			count++
		}
	}
	return count
}

func (m ProductModel) Init() tea.Cmd {
	return nil
}

func (m ProductModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.products) - 1
			}
		case "down":
			if m.cursor < len(m.products)-1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter", " ":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m ProductModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// 1. Logo 区域
	b.WriteString(RenderLogo(m.version))

	// 2. 表头: AI Tool | Global | Project
	headerName := padRight(i18n.T("tui_col_ai_tool"), 22)
	b.WriteString(fmt.Sprintf("  %s  %s  %s\n",
		titleStyle.Render(headerName),
		titleStyle.Render(padRight(i18n.T("tui_global"), 4)),
		titleStyle.Render(padRight(i18n.T("tui_project"), 4))))
	b.WriteString(separatorStyle.Render(strings.Repeat("─", SeparatorWidth)))
	b.WriteString("\n")

	for i, p := range m.products {
		prefix := "  "
		style := selectableStyle

		if m.cursor == i {
			prefix = cursorStyle.Render("❯ ")
			style = selectedStyle
		}

		// 计算全局和项目的技能数量
		globalCount := countSkillsInDir(p.GlobalPath())
		eachProductProjectPath := filepath.Join(m.targetDir, p.ProjectSkills)
		projectCount := countSkillsInDir(eachProductProjectPath)

		// 数量样式: 有技能时高亮，无技能时灰色
		globalCountStyle := normalStyle
		projectCountStyle := normalStyle
		if globalCount > 0 {
			globalCountStyle = selectedStyle
		}
		if projectCount > 0 {
			projectCountStyle = selectedStyle
		}

		// 先对齐纯文本，再套样式
		name := padRight(p.Name, 22)
		globalStr := fmt.Sprintf("%4d", globalCount)
		projectStr := fmt.Sprintf("%4d", projectCount)

		b.WriteString(fmt.Sprintf("%s%s  %s  %s\n",
			prefix,
			style.Render(name),
			globalCountStyle.Render(globalStr),
			projectCountStyle.Render(projectStr)))
	}

	// 4. 提示区域
	b.WriteString(RenderHint(i18n.T("tui_hint_select")))

	return b.String()
}

// SelectedProduct returns the selected product, nil if cancelled
func (m ProductModel) SelectedProduct() *products.Product {
	if m.quitting {
		return nil
	}
	if m.cursor >= 0 && m.cursor < len(m.products) {
		return &m.products[m.cursor]
	}
	return nil
}

func (m ProductModel) IsQuitting() bool {
	return m.quitting
}

// RunProductSelect runs the product selection page
func RunProductSelect(version string, targetDir string) (*products.Product, error) {
	m := NewProductModel(version, targetDir)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	result := finalModel.(ProductModel)
	return result.SelectedProduct(), nil
}

// ============================================================================
// Install Target Select Model - Level 3: Select Global or Project
// ============================================================================

// InstallTargetModel for selecting install target (global or project)
type InstallTargetModel struct {
	product    *products.Product
	projectDir string
	cursor     int
	quitting   bool
}

// NewInstallTargetModel creates a new install target selection model
func NewInstallTargetModel(product *products.Product, projectDir string) InstallTargetModel {
	return InstallTargetModel{
		product:    product,
		projectDir: projectDir,
		cursor:     0,
	}
}

func (m InstallTargetModel) Init() tea.Cmd {
	return nil
}

func (m InstallTargetModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = 1
			}
		case "down":
			if m.cursor < 1 {
				m.cursor++
			} else {
				m.cursor = 0
			}
		case "enter", " ":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m InstallTargetModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Logo 区域
	b.WriteString(RenderLogo(""))
	b.WriteString("\n")

	// 标题
	b.WriteString(titleStyle.Render(i18n.T("tui_install_target_title")))
	b.WriteString("\n\n")

	// 选项
	targets := []string{i18n.T("tui_global"), i18n.T("tui_project")}
	for i, target := range targets {
		prefix := "  "
		style := selectableStyle

		if m.cursor == i {
			prefix = cursorStyle.Render("❯ ")
			style = selectedStyle
		}

		desc := ""
		if i == 0 {
			desc = hintStyle.Render("(" + m.product.GlobalPath() + ")")
		} else {
			// 项目路径 = 当前目录 + 工具的项目skills目录
			projectSkillPath := filepath.Join(m.projectDir, m.product.ProjectSkills)
			desc = hintStyle.Render("(" + projectSkillPath + ")")
		}

		b.WriteString(fmt.Sprintf("%s%s %s\n", prefix, style.Render(target), desc))
	}

	// 提示
	b.WriteString(RenderHint(i18n.T("tui_hint_select")))

	return b.String()
}

func (m InstallTargetModel) SelectedTarget() string {
	if m.quitting {
		return ""
	}
	if m.cursor == 0 {
		return "global"
	}
	return "project"
}

func (m InstallTargetModel) IsQuitting() bool {
	return m.quitting
}

// RunInstallTargetSelect runs the install target selection page
func RunInstallTargetSelect(product *products.Product, projectDir string) (string, error) {
	m := NewInstallTargetModel(product, projectDir)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	result := finalModel.(InstallTargetModel)
	return result.SelectedTarget(), nil
}
