// Package products provides AI tool product configuration
package products

import (
	"os"
	"path/filepath"
	"strings"
)

// Product represents an AI tool that supports skills
type Product struct {
	Name         string // Display name
	GlobalSkills string // Global skills directory (~ expansion)
	ProjectSkills string // Project skills directory
}

// AllProducts returns all supported AI tools
var AllProducts = []Product{
	{
		Name:         "Claude Code",
		GlobalSkills: "~/.claude/skills/",
		ProjectSkills: ".claude/skills/",
	},
	{
		Name:         "Cursor",
		GlobalSkills: "~/.cursor/skills/",
		ProjectSkills: ".cursor/skills/",
	},
	{
		Name:         "Windsurf",
		GlobalSkills: "~/.windsurf/skills/",
		ProjectSkills: ".windsurf/skills/",
	},
	{
		Name:         "Trae",
		GlobalSkills: "~/.trae/skills/",
		ProjectSkills: ".trae/skills/",
	},
	{
		Name:         "Qoder",
		GlobalSkills: "~/.qoder/skills/",
		ProjectSkills: ".qoder/skills/",
	},
	{
		Name:         "CodeX",
		GlobalSkills: "~/.codex/skills/",
		ProjectSkills: ".codex/skills/",
	},
	{
		Name:         "Kimi",
		GlobalSkills: "~/.kimi/skills/",
		ProjectSkills: ".kimi/skills/",
	},
	{
		Name:         "CodeBuddy",
		GlobalSkills: "~/.codebuddy/skills/",
		ProjectSkills: ".codebuddy/skills/",
	},
	{
		Name:         "Roo Code",
		GlobalSkills: "~/.roocode/skills/",
		ProjectSkills: ".roocode/skills/",
	},
	{
		Name:         "Opencode",
		GlobalSkills: "~/.config/opencode/skills/",
		ProjectSkills: ".agents/skills/",
	},
	{
		Name:         "Aider",
		GlobalSkills: "~/.aider/skills/",
		ProjectSkills: ".aider/skills/",
	},
}

// ExpandPath expands ~ to user's home directory
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}
	return path
}

// GlobalPath returns the expanded global skills path
func (p *Product) GlobalPath() string {
	return ExpandPath(p.GlobalSkills)
}

// ProjectPath returns the project skills path
func (p *Product) ProjectPath() string {
	return p.ProjectSkills
}

// IsInstalled checks if the product has any skills installed in the global directory
func (p *Product) IsInstalled() bool {
	path := p.GlobalPath()
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetProductByName finds a product by name (case-insensitive)
func GetProductByName(name string) *Product {
	lowerName := strings.ToLower(name)
	for i := range AllProducts {
		if strings.ToLower(AllProducts[i].Name) == lowerName {
			return &AllProducts[i]
		}
	}
	return nil
}

// GetProductCount returns the number of supported products
func GetProductCount() int {
	return len(AllProducts)
}
