package versioncheck

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// NormalizeVersion trims common prefixes/suffixes for comparison.
func NormalizeVersion(version string) string {
	v := strings.TrimSpace(version)
	if v == "" {
		return ""
	}
	if strings.HasPrefix(v, "v") || strings.HasPrefix(v, "V") {
		v = v[1:]
	}
	if idx := strings.IndexAny(v, "-+"); idx >= 0 {
		v = v[:idx]
	}
	return v
}

// ShouldNotify decides whether to show an update prompt.
func ShouldNotify(currentVersion, latestVersion string) bool {
	if latestVersion == "" {
		return false
	}
	current := NormalizeVersion(currentVersion)
	latest := NormalizeVersion(latestVersion)
	if current == "" || current == "dev" || current == "unknown" {
		return false
	}
	return current != latest
}

// LatestFromNpmJSON parses npm registry JSON and returns dist-tags.latest.
func LatestFromNpmJSON(data []byte) (string, error) {
	var payload struct {
		DistTags map[string]string `json:"dist-tags"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return "", err
	}
	latest := ""
	if payload.DistTags != nil {
		latest = payload.DistTags["latest"]
	}
	if latest == "" {
		return "", errors.New("latest version not found")
	}
	return latest, nil
}

// FetchLatestVersion queries npm registry for the latest version.
func FetchLatestVersion(ctx context.Context, packageName string) (string, error) {
	if packageName == "" {
		return "", errors.New("package name is empty")
	}
	url := fmt.Sprintf("https://registry.npmjs.org/%s", packageName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return LatestFromNpmJSON(data)
}
