// Package errmsg provides unified command-line error formatting
package errmsg

import (
	"fmt"

	"github.com/castle-x/skills-x/cmd/skills-x/i18n"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// Error is a custom error type for formatted error messages
type Error struct {
	// Title (displayed in red)
	Title string
	// Detail (displayed in red)
	Detail string
	// Conditions (displayed with yellow header)
	Conditions []string
	// Solutions (displayed with yellow header)
	Solutions []string
	// Documentation URL
	DocURL string
	// Original error
	Cause error
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Title
}

// Unwrap returns the original error
func (e *Error) Unwrap() error {
	return e.Cause
}

// Print outputs the formatted error message to terminal
func (e *Error) Print() {
	fmt.Println()

	// Red: Error title
	fmt.Printf("%s%s%s: %s%s\n", colorBold, colorRed, i18n.T("err_title"), e.Title, colorReset)

	// Red: Error detail
	if e.Detail != "" {
		fmt.Printf("%s  %s%s\n", colorRed, e.Detail, colorReset)
	}
	fmt.Println()

	// Yellow: Conditions
	if len(e.Conditions) > 0 {
		fmt.Printf("%s%s:%s\n", colorYellow, i18n.T("err_conditions"), colorReset)
		for _, c := range e.Conditions {
			fmt.Printf("  %s• %s%s\n", colorGray, c, colorReset)
		}
		fmt.Println()
	}

	// Yellow: Solutions
	if len(e.Solutions) > 0 {
		fmt.Printf("%s%s:%s\n", colorYellow, i18n.T("err_solutions"), colorReset)
		for i, s := range e.Solutions {
			fmt.Printf("  %s%d. %s%s\n", colorCyan, i+1, s, colorReset)
		}
		fmt.Println()
	}

	// Yellow: Documentation URL
	if e.DocURL != "" {
		fmt.Printf("%s%s: %s%s%s\n", colorYellow, i18n.T("err_doc"), colorCyan, e.DocURL, colorReset)
		fmt.Println()
	}
}

// IsCustomError checks if the error is a custom error
func IsCustomError(err error) bool {
	_, ok := err.(*Error)
	return ok
}

// PrintError prints the error, formatted if custom error, plain otherwise
func PrintError(err error) {
	if e, ok := err.(*Error); ok {
		e.Print()
	} else {
		fmt.Printf("\n%s%s%s: %s%s\n\n", colorBold, colorRed, i18n.T("err_title"), err.Error(), colorReset)
	}
}

// ============================================================================
// Predefined Error Constructors
// ============================================================================

// SkillNotFound returns an error when skill is not found
func SkillNotFound(name string) *Error {
	return &Error{
		Title:  i18n.T("err_skill_not_found"),
		Detail: i18n.Tf("err_skill_not_found_detail", name),
		Conditions: []string{
			i18n.T("err_skill_not_found_cond1"),
			i18n.T("err_skill_not_found_cond2"),
		},
		Solutions: []string{
			i18n.T("err_skill_not_found_sol1"),
			i18n.T("err_skill_not_found_sol2"),
		},
		DocURL: "https://github.com/castle-x/skills-x",
	}
}

// MissingArgument returns an error when argument is missing
func MissingArgument(argName string) *Error {
	return &Error{
		Title: i18n.Tf("err_missing_argument", argName),
	}
}

// TargetDirCreateError returns an error when cannot create target directory
func TargetDirCreateError(path string) *Error {
	return &Error{
		Title:  i18n.T("err_target_dir_create"),
		Detail: i18n.Tf("err_target_dir_create_detail", path),
		Solutions: []string{
			i18n.T("err_target_dir_create_sol1"),
			i18n.T("err_target_dir_create_sol2"),
		},
	}
}

// CopyFailed returns an error when copy operation fails
func CopyFailed(skillName string) *Error {
	return &Error{
		Title:  i18n.T("err_copy_failed"),
		Detail: i18n.Tf("err_copy_failed_detail", skillName),
		Solutions: []string{
			i18n.T("err_copy_failed_sol1"),
			i18n.T("err_copy_failed_sol2"),
		},
	}
}

// NoSkillsEmbedded returns an error when no skills data is embedded
func NoSkillsEmbedded() *Error {
	return &Error{
		Title:  i18n.T("err_no_skills_embedded"),
		Detail: i18n.T("err_no_skills_embedded_detail"),
		Solutions: []string{
			i18n.T("err_no_skills_embedded_sol1"),
		},
		DocURL: "https://github.com/castle-x/skills-x",
	}
}
