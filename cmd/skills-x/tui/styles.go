// Package tui provides terminal interactive UI components
package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// ============================================================================
// Style Definitions
// ============================================================================

var (
	// Colors
	primaryColor   = lipgloss.Color("#00D4AA") // 青绿色
	secondaryColor = lipgloss.Color("#FF6B6B") // 红色
	accentColor    = lipgloss.Color("#4ECDC4") // 浅青色
	dimColor       = lipgloss.Color("#666666") // 灰色
	whiteColor     = lipgloss.Color("#FFFFFF")
	yellowColor    = lipgloss.Color("#FFCC00")
	blueColor      = lipgloss.Color("#5599FF")
	cyanColor      = lipgloss.Color("#5A9FB8") // 可选项颜色 (柔和青色)

	// Logo 样式
	logoStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// 标题样式
	titleStyle = lipgloss.NewStyle().
			Foreground(whiteColor).
			Bold(true)

	// 选中项样式
	selectedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// 未选中项样式
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA"))

	// 可选项样式 (未选中时)
	selectableStyle = lipgloss.NewStyle().
			Foreground(cyanColor)

	// 提示样式
	hintStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	// 光标样式
	cursorStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	// 成功样式
	successStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))

	// 错误样式
	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF4444"))

	// 值样式
	valueStyle = lipgloss.NewStyle().
			Foreground(accentColor)

	// 分隔线
	separatorStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	// 已安装标记样式
	installedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	// 搜索框样式
	searchStyle = lipgloss.NewStyle().
			Foreground(whiteColor).
			Background(lipgloss.Color("#333333"))
)

// ASCII Logo for Skills-X
const logo = `
   ███████╗██╗  ██╗██╗██╗     ██╗     ███████╗     ██╗  ██╗
   ██╔════╝██║ ██╔╝██║██║     ██║     ██╔════╝     ╚██╗██╔╝
   ███████╗█████╔╝ ██║██║     ██║     ███████╗█████╗╚███╔╝
   ╚════██║██╔═██╗ ██║██║     ██║     ╚════██║╚════╝██╔██╗
   ███████║██║  ██╗██║███████╗███████╗███████║     ██╔╝ ██╗
   ╚══════╝╚═╝  ╚═╝╚═╝╚══════╝╚══════╝╚══════╝     ╚═╝  ╚═╝
`

// 分隔线宽度
const SeparatorWidth = 60

// termWidth returns the visual width of a string in terminal columns.
// CJK characters occupy 2 columns, ASCII occupies 1.
func termWidth(s string) int {
	w := 0
	for _, r := range s {
		if r >= 0x1100 &&
			(r <= 0x115F || r == 0x2329 || r == 0x232A ||
				(r >= 0x2E80 && r <= 0x303E) ||
				(r >= 0x3040 && r <= 0x33BF) ||
				(r >= 0x3400 && r <= 0x4DBF) ||
				(r >= 0x4E00 && r <= 0xA4CF) ||
				(r >= 0xA960 && r <= 0xA97C) ||
				(r >= 0xAC00 && r <= 0xD7A3) ||
				(r >= 0xF900 && r <= 0xFAFF) ||
				(r >= 0xFE10 && r <= 0xFE6F) ||
				(r >= 0xFF01 && r <= 0xFF60) ||
				(r >= 0xFFE0 && r <= 0xFFE6) ||
				(r >= 0x1F000 && r <= 0x1FFFF) ||
				(r >= 0x20000 && r <= 0x2FA1F)) {
			w += 2
		} else {
			w++
		}
	}
	return w
}

// padRight pads a string with spaces to reach the target visual width.
func padRight(s string, targetWidth int) string {
	w := termWidth(s)
	for w < targetWidth {
		s += " "
		w++
	}
	return s
}

// Color reset
const colorReset = "\033[0m"

// CurrentWorkDir 当前工作目录（供 RenderLogo 使用）
var CurrentWorkDir string

// ============================================================================
// Basic Page Components
// ============================================================================

// RenderLogo 渲染 Logo 区域 (标题+版本在上，分隔线，Logo 在下)
func RenderLogo(version string) string {
	var b strings.Builder

	// 第一行: 标题 + 版本 + 工作目录
	titleText := "🚀 Skills-X - AI Agent Skills Manager"
	b.WriteString(lipgloss.NewStyle().Foreground(accentColor).Render(titleText))
	if version != "" {
		displayVersion := strings.TrimSuffix(version, "-dirty")
		b.WriteString("  ")
		b.WriteString(hintStyle.Render(displayVersion))
	}
	if CurrentWorkDir != "" {
		b.WriteString("  ")
		b.WriteString(lipgloss.NewStyle().Foreground(dimColor).Render("["+CurrentWorkDir+"]"))
	}
	b.WriteString("\n")

	// Logo 上方分隔线
	b.WriteString(separatorStyle.Render(strings.Repeat("─", SeparatorWidth)))
	b.WriteString("\n")

	// Logo
	b.WriteString(logoStyle.Render(logo))
	b.WriteString("\n")

	// Logo 下方分隔线
	b.WriteString(separatorStyle.Render(strings.Repeat("─", SeparatorWidth)))
	b.WriteString("\n")

	return b.String()
}

// RenderSeparator 渲染分隔线
func RenderSeparator() string {
	return separatorStyle.Render(strings.Repeat("─", SeparatorWidth)) + "\n"
}

// RenderHint 渲染提示区域
func RenderHint(hint string) string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(RenderSeparator())
	b.WriteString(hintStyle.Render(hint))
	return b.String()
}

// RenderStatusBar 渲染底部状态栏（显示总数、已安装数量、新增数量、反选数量）
func RenderStatusBar(total, installed, newSelected, deselected int) string {
	var b strings.Builder
	b.WriteString("\n")
	statusText := fmt.Sprintf("(共 %d 个，已安装 %d 个，将安装 %d 个，将卸载 %d 个)", total, installed, newSelected, deselected)
	b.WriteString(hintStyle.Render(statusText))
	return b.String()
}

// RenderInstallProgress 渲染安装进度（底部动态显示）
func RenderInstallProgress(current, total int, skillName string) string {
	var b strings.Builder
	b.WriteString("\n")
	progressText := fmt.Sprintf("正在安装 (%d/%d): %s...", current, total, skillName)
	b.WriteString(successStyle.Render(progressText))
	return b.String()
}

// ============================================================================
// Terminal Control Sequences
// ============================================================================

// ClearScreen 清屏 ANSI 序列
const ClearScreen = "\033[2J\033[H"

// EnterAltScreen 进入备用屏幕
const EnterAltScreen = "\033[?1049h"

// ExitAltScreen 退出备用屏幕
const ExitAltScreen = "\033[?1049l"

// HideCursor 隐藏光标
const HideCursor = "\033[?25l"

// ShowCursor 显示光标
const ShowCursor = "\033[?25h"

// clearTerminal 清理终端
func clearTerminal() {
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		fmt.Print("\033[H") // 移动到左上角
		for i := 0; i < h; i++ {
			fmt.Print("\033[K\n") // 清除当前行
		}
		fmt.Print("\033[H") // 回到左上角
	} else {
		fmt.Print(ClearScreen)
	}
}

// PrintWelcome 打印欢迎界面
func PrintWelcome(version string) {
	fmt.Print(ClearScreen)
	fmt.Print(RenderLogo(version))
	fmt.Println()
}
