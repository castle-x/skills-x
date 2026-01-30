// Package gitutil provides git operations for cloning repositories
package gitutil

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	// CloneTimeout is the maximum time for a git clone operation
	CloneTimeout = 60 * time.Second
	// SparseCloneTimeout is the maximum time for a sparse clone operation
	SparseCloneTimeout = 30 * time.Second
	// TempDirPrefix is the prefix for temporary directories
	TempDirPrefix = "skills-"
	// MaxRetries is the maximum number of retry attempts for network operations
	MaxRetries = 3
	// RetryDelay is the delay between retry attempts
	RetryDelay = 2 * time.Second
)

// CloneError represents a git clone error
type CloneError struct {
	URL         string
	Message     string
	IsTimeout   bool
	IsAuthError bool
}

func (e *CloneError) Error() string {
	return e.Message
}

// CloneResult contains the result of a clone operation
type CloneResult struct {
	TempDir string // Path to the cloned repository
	Repo    string // Original repo identifier
}

// CloneRepo clones a git repository to a temporary directory
// Uses shallow clone (--depth 1) for efficiency
// Supports caching: if repo exists, reuses it directly (no network request)
// Includes retry logic for transient network failures
func CloneRepo(gitURL string, repoName string) (*CloneResult, error) {
	return CloneRepoWithRefresh(gitURL, repoName, false)
}

// CloneRepoWithRefresh clones a git repository with optional cache refresh
// If refresh is true, it will fetch the latest changes even if cache exists
func CloneRepoWithRefresh(gitURL string, repoName string, refresh bool) (*CloneResult, error) {
	// Create a deterministic temp directory based on repo name
	// This allows caching/reuse across invocations
	tempDir := getTempDir(repoName)

	// If directory already exists and has content
	if dirExists(tempDir) && hasGitContent(tempDir) {
		if refresh {
			// Force refresh: update the cached repo
			if err := updateShallowRepo(tempDir); err != nil {
				// If update fails, remove and re-clone
				os.RemoveAll(tempDir)
			} else {
				return &CloneResult{
					TempDir: tempDir,
					Repo:    repoName,
				}, nil
			}
		} else {
			// Use cache directly - no network request
			return &CloneResult{
				TempDir: tempDir,
				Repo:    repoName,
			}, nil
		}
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(tempDir), 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Clone with retry logic for transient network failures
	var lastErr error
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		result, err := cloneWithTimeout(gitURL, tempDir)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Don't retry on auth errors - they won't succeed
		if cloneErr, ok := err.(*CloneError); ok && cloneErr.IsAuthError {
			return nil, err
		}

		// Clean up failed attempt
		os.RemoveAll(tempDir)

		// Wait before retrying (except on last attempt)
		if attempt < MaxRetries {
			time.Sleep(RetryDelay)
		}
	}

	return nil, lastErr
}

// cloneWithTimeout performs a single clone attempt with timeout
func cloneWithTimeout(gitURL string, tempDir string) (*CloneResult, error) {
	cmd := exec.Command("git", "clone", "--depth", "1", gitURL, tempDir)

	// Set up timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, parseGitError(err, gitURL)
		}
	case <-time.After(CloneTimeout):
		// Kill the process on timeout
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return nil, &CloneError{
			URL:       gitURL,
			Message:   fmt.Sprintf("Clone timed out after %v. Check your network or credentials.", CloneTimeout),
			IsTimeout: true,
		}
	}

	return &CloneResult{
		TempDir: tempDir,
		Repo:    filepath.Base(tempDir),
	}, nil
}

// SparseCloneRepo clones only specific directories from a git repository
// This is much faster for large repositories when you only need specific paths
func SparseCloneRepo(gitURL string, repoName string, sparsePaths []string) (*CloneResult, error) {
	// Create a deterministic temp directory based on repo name + sparse paths
	tempDir := getTempDirSparse(repoName, sparsePaths)

	// If directory already exists and has content, reuse it
	if dirExists(tempDir) && hasGitContent(tempDir) {
		// For sparse checkout, we don't pull - just reuse
		return &CloneResult{
			TempDir: tempDir,
			Repo:    repoName,
		}, nil
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(tempDir), 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Step 1: Initialize empty repo with sparse checkout
	if err := runWithTimeout(exec.Command("git", "init", tempDir), SparseCloneTimeout); err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to init repo: %w", err)
	}

	// Step 2: Add remote
	if err := runWithTimeout(exec.Command("git", "-C", tempDir, "remote", "add", "origin", gitURL), SparseCloneTimeout); err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to add remote: %w", err)
	}

	// Step 3: Enable sparse checkout
	if err := runWithTimeout(exec.Command("git", "-C", tempDir, "config", "core.sparseCheckout", "true"), SparseCloneTimeout); err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to enable sparse checkout: %w", err)
	}

	// Step 4: Configure sparse checkout paths
	sparseCheckoutFile := filepath.Join(tempDir, ".git", "info", "sparse-checkout")
	if err := os.MkdirAll(filepath.Dir(sparseCheckoutFile), 0755); err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to create sparse-checkout dir: %w", err)
	}

	// Write sparse paths to config
	sparseContent := strings.Join(sparsePaths, "\n") + "\n"
	if err := os.WriteFile(sparseCheckoutFile, []byte(sparseContent), 0644); err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to write sparse-checkout config: %w", err)
	}

	// Step 5: Fetch only the needed data (depth 1)
	fetchCmd := exec.Command("git", "-C", tempDir, "fetch", "--depth", "1", "origin", "HEAD")
	if err := runWithTimeout(fetchCmd, CloneTimeout); err != nil {
		os.RemoveAll(tempDir)
		return nil, parseGitError(err, gitURL)
	}

	// Step 6: Checkout
	checkoutCmd := exec.Command("git", "-C", tempDir, "checkout", "FETCH_HEAD")
	if err := runWithTimeout(checkoutCmd, SparseCloneTimeout); err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to checkout: %w", err)
	}

	return &CloneResult{
		TempDir: tempDir,
		Repo:    repoName,
	}, nil
}

// runWithTimeout runs a command with a timeout
func runWithTimeout(cmd *exec.Cmd, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return &CloneError{
			Message:   fmt.Sprintf("command timed out after %v", timeout),
			IsTimeout: true,
		}
	}
}

// CleanupTempDir removes a temporary directory
// Only removes directories within the system temp directory for safety
func CleanupTempDir(dir string) error {
	tmpDir := os.TempDir()

	// Normalize paths for comparison
	normalizedDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("invalid directory path: %w", err)
	}

	normalizedTmpDir, err := filepath.Abs(tmpDir)
	if err != nil {
		return fmt.Errorf("invalid temp directory path: %w", err)
	}

	// Safety check: only allow deletion within temp directory
	if !strings.HasPrefix(normalizedDir, normalizedTmpDir) {
		return fmt.Errorf("attempted to clean up directory outside of temp directory: %s", dir)
	}

	return os.RemoveAll(dir)
}

// CleanupAllSkillsDirs removes all skills-* temp directories
func CleanupAllSkillsDirs() error {
	tmpDir := os.TempDir()
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), TempDirPrefix) {
			path := filepath.Join(tmpDir, entry.Name())
			os.RemoveAll(path)
		}
	}
	return nil
}

// getTempDir returns a deterministic temp directory path for a repo
func getTempDir(repoName string) string {
	// Create a short hash of the repo name for uniqueness
	hash := sha256.Sum256([]byte(repoName))
	shortHash := fmt.Sprintf("%x", hash[:4])

	// Replace special characters
	safeName := strings.ReplaceAll(repoName, "/", "-")
	safeName = strings.ReplaceAll(safeName, ".", "-")

	return filepath.Join(os.TempDir(), TempDirPrefix+safeName+"-"+shortHash)
}

// getTempDirSparse returns a deterministic temp directory path for a sparse checkout
func getTempDirSparse(repoName string, sparsePaths []string) string {
	// Include sparse paths in the hash for uniqueness
	hashInput := repoName + ":" + strings.Join(sparsePaths, ",")
	hash := sha256.Sum256([]byte(hashInput))
	shortHash := fmt.Sprintf("%x", hash[:4])

	// Replace special characters
	safeName := strings.ReplaceAll(repoName, "/", "-")
	safeName = strings.ReplaceAll(safeName, ".", "-")

	return filepath.Join(os.TempDir(), TempDirPrefix+safeName+"-sparse-"+shortHash)
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// hasGitContent checks if directory has .git folder
func hasGitContent(path string) bool {
	gitDir := filepath.Join(path, ".git")
	return dirExists(gitDir)
}

// pullRepo pulls latest changes in a repo (for non-shallow clones)
func pullRepo(dir string) error {
	cmd := exec.Command("git", "-C", dir, "pull", "--ff-only")
	return cmd.Run()
}

// updateShallowRepo updates a shallow clone by fetching latest and resetting
// This is the correct way to update a --depth 1 clone
func updateShallowRepo(dir string) error {
	// Fetch latest with depth 1
	fetchCmd := exec.Command("git", "-C", dir, "fetch", "--depth", "1", "origin")
	if err := runWithTimeout(fetchCmd, CloneTimeout); err != nil {
		return err
	}

	// Get the default branch name
	branchCmd := exec.Command("git", "-C", dir, "symbolic-ref", "refs/remotes/origin/HEAD", "--short")
	output, err := branchCmd.Output()
	if err != nil {
		// Fallback: try to reset to origin/main or origin/master
		resetCmd := exec.Command("git", "-C", dir, "reset", "--hard", "origin/HEAD")
		return runWithTimeout(resetCmd, SparseCloneTimeout)
	}

	// Reset to the fetched head
	branch := strings.TrimSpace(string(output))
	resetCmd := exec.Command("git", "-C", dir, "reset", "--hard", branch)
	return runWithTimeout(resetCmd, SparseCloneTimeout)
}

// parseGitError converts git errors to CloneError
func parseGitError(err error, url string) *CloneError {
	msg := err.Error()

	isAuth := strings.Contains(msg, "Authentication failed") ||
		strings.Contains(msg, "could not read Username") ||
		strings.Contains(msg, "Permission denied") ||
		strings.Contains(msg, "Repository not found")

	if isAuth {
		return &CloneError{
			URL:         url,
			Message:     fmt.Sprintf("Authentication failed for %s. Check your credentials or repository access.", url),
			IsAuthError: true,
		}
	}

	return &CloneError{
		URL:     url,
		Message: fmt.Sprintf("Failed to clone %s: %v", url, err),
	}
}

// GetCachedDir returns the cached directory path if it exists
func GetCachedDir(repoName string) (string, bool) {
	tempDir := getTempDir(repoName)
	if dirExists(tempDir) && hasGitContent(tempDir) {
		return tempDir, true
	}
	return "", false
}

// GetCachedDirSparse returns the cached sparse directory path if it exists
func GetCachedDirSparse(repoName string, sparsePaths []string) (string, bool) {
	tempDir := getTempDirSparse(repoName, sparsePaths)
	if dirExists(tempDir) && hasGitContent(tempDir) {
		return tempDir, true
	}
	return "", false
}
