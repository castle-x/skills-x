// Package skills provides embedded skills data and metadata
package skills

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed all:data
var skillsFS embed.FS

// SkillInfo holds metadata about a skill
type SkillInfo struct {
	Name        string
	Category    string
	Description string
	IsCastleX   bool // true if from castle-x (self-developed)
	Path        string
}

// skillCategories defines the category for each skill
var skillCategories = map[string]string{
	// Creative & Design
	"ui-ux-pro-max":    "creative",
	"algorithmic-art":  "creative",
	"canvas-design":    "creative",
	"brand-guidelines": "creative",
	"theme-factory":    "creative",
	"frontend-design":  "creative",
	"image-enhancer":   "creative",

	// Document Processing
	"pdf":             "document",
	"docx":            "document",
	"pptx":            "document",
	"xlsx":            "document",
	"document-skills": "document",
	"doc-coauthoring": "document",

	// Development Tools
	"mcp-builder":          "devtools",
	"artifacts-builder":    "devtools",
	"web-artifacts-builder": "devtools",
	"webapp-testing":       "devtools",
	"langsmith-fetch":      "devtools",
	"changelog-generator":  "devtools",

	// Workflows
	"brainstorming":                "workflow",
	"writing-plans":                "workflow",
	"executing-plans":              "workflow",
	"systematic-debugging":         "workflow",
	"test-driven-development":      "workflow",
	"verification-before-completion": "workflow",
	"subagent-driven-development":  "workflow",
	"dispatching-parallel-agents":  "workflow",

	// Git & Code Review
	"requesting-code-review":        "git",
	"receiving-code-review":         "git",
	"finishing-a-development-branch": "git",
	"using-git-worktrees":           "git",

	// Writing
	"content-research-writer":   "writing",
	"internal-comms":            "writing",
	"tailored-resume-generator": "writing",

	// Integrations
	"connect":             "integration",
	"connect-apps":        "integration",
	"connect-apps-plugin": "integration",
	"slack-gif-creator":   "integration",

	// Business & Analytics
	"competitive-ads-extractor":   "business",
	"developer-growth-analysis":   "business",
	"lead-research-assistant":     "business",
	"meeting-insights-analyzer":   "business",
	"twitter-algorithm-optimizer": "business",

	// File Management
	"file-organizer":    "files",
	"invoice-organizer": "files",

	// Utilities
	"video-downloader":        "utility",
	"domain-name-brainstormer": "utility",
	"raffle-winner-picker":    "utility",

	// Skills Development
	"skill-creator":    "skilldev",
	"writing-skills":   "skilldev",
	"skill-share":      "skilldev",
	"template-skill":   "skilldev",
	"using-superpowers": "skilldev",
}

// skillDescriptions provides short descriptions for each skill
var skillDescriptions = map[string]string{
	"ui-ux-pro-max":                   "UI/UX design intelligence with 67 styles, 96 palettes",
	"algorithmic-art":                 "p5.js generative art, flow fields, particles",
	"canvas-design":                   "Posters, visual art (.png/.pdf)",
	"brand-guidelines":                "Anthropic brand styling",
	"theme-factory":                   "Artifact theme switching (10 presets)",
	"frontend-design":                 "Frontend design assistance",
	"image-enhancer":                  "Image upscaling, sharpening, cleanup",
	"pdf":                             "PDF extract/fill/merge",
	"docx":                            "Word document processing",
	"pptx":                            "PowerPoint presentations",
	"xlsx":                            "Excel sheets/formulas/charts",
	"document-skills":                 "Comprehensive document processing",
	"doc-coauthoring":                 "Document collaboration",
	"mcp-builder":                     "Build MCP servers",
	"artifacts-builder":               "React+Tailwind+shadcn artifacts",
	"web-artifacts-builder":           "Complex HTML artifacts",
	"webapp-testing":                  "Playwright testing",
	"langsmith-fetch":                 "LangSmith debug tracing",
	"changelog-generator":             "Generate changelog from git commits",
	"brainstorming":                   "Brainstorm before creative work",
	"writing-plans":                   "Write task plans",
	"executing-plans":                 "Execute plans",
	"systematic-debugging":            "Systematic debugging methodology",
	"test-driven-development":         "TDD workflow",
	"verification-before-completion":  "Verify before completion",
	"subagent-driven-development":     "Subagent-driven development",
	"dispatching-parallel-agents":     "Parallel agent dispatching",
	"requesting-code-review":          "Request code review",
	"receiving-code-review":           "Handle CR feedback",
	"finishing-a-development-branch":  "Complete development branch",
	"using-git-worktrees":             "Git Worktree isolation",
	"content-research-writer":         "Content research and writing",
	"internal-comms":                  "Internal communications/reports",
	"tailored-resume-generator":       "Custom resume generation",
	"connect":                         "Connect 1000+ services",
	"connect-apps":                    "Gmail/Slack/GitHub integration",
	"connect-apps-plugin":             "App connection plugin",
	"slack-gif-creator":               "Slack GIF creation",
	"competitive-ads-extractor":       "Competitor ad analysis",
	"developer-growth-analysis":       "Developer growth analytics",
	"lead-research-assistant":         "Lead research assistant",
	"meeting-insights-analyzer":       "Meeting analysis",
	"twitter-algorithm-optimizer":     "Tweet optimization",
	"file-organizer":                  "File organization",
	"invoice-organizer":               "Invoice organization/tax prep",
	"video-downloader":                "YouTube download",
	"domain-name-brainstormer":        "Domain name ideas",
	"raffle-winner-picker":            "Raffle picker",
	"skill-creator":                   "Create new skills",
	"writing-skills":                  "Write/validate skills",
	"skill-share":                     "Share skills",
	"template-skill":                  "Skill template",
	"using-superpowers":               "How to use skills",
}

// castleXSkills lists skills from castle-x (self-developed)
var castleXSkills = map[string]bool{
	// Add castle-x skills here when available
}

// GetFS returns the embedded filesystem
func GetFS() embed.FS {
	return skillsFS
}

// ListSkills returns all available skills with metadata
func ListSkills() ([]SkillInfo, error) {
	var skills []SkillInfo

	// Walk through skills/data directory
	err := fs.WalkDir(skillsFS, "data", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip root and non-directories
		if path == "data" || !d.IsDir() {
			return nil
		}

		// Only process top-level directories (skill folders)
		rel, _ := filepath.Rel("data", path)
		if strings.Contains(rel, string(filepath.Separator)) {
			return fs.SkipDir
		}

		name := d.Name()
		
		// Check if SKILL.md exists
		skillMdPath := filepath.Join(path, "SKILL.md")
		if _, err := skillsFS.Open(skillMdPath); err != nil {
			return nil // Skip directories without SKILL.md
		}

		info := SkillInfo{
			Name:        name,
			Category:    skillCategories[name],
			Description: skillDescriptions[name],
			IsCastleX:   castleXSkills[name],
			Path:        path,
		}

		if info.Category == "" {
			info.Category = "other"
		}

		skills = append(skills, info)
		return nil
	})

	return skills, err
}

// GetSkill returns a specific skill by name
func GetSkill(name string) (*SkillInfo, error) {
	skills, err := ListSkills()
	if err != nil {
		return nil, err
	}

	for _, s := range skills {
		if s.Name == name {
			return &s, nil
		}
	}

	return nil, nil
}

// SkillExists checks if a skill exists
func SkillExists(name string) bool {
	skill, _ := GetSkill(name)
	return skill != nil
}
