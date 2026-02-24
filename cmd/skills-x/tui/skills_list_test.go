// Package tui provides terminal interactive UI components
package tui

import (
	"testing"
)

func TestSortSkills(t *testing.T) {
	// Test data: skills in random order
	skills := []SkillItem{
		{FullName: "vercel/react-best-practices"},
		{FullName: "anthropic/frontend-design"},
		{FullName: "superpowers/brainstorming"},
		{FullName: "anthropic/algorithmic-art"},
		{FullName: "vercel/react-native-skills"},
		{FullName: "developer-kit/react-patterns"},
	}

	// Sort
	SortSkills(skills)

	// Expected order
	expected := []string{
		"anthropic/algorithmic-art",
		"anthropic/frontend-design",
		"developer-kit/react-patterns",
		"superpowers/brainstorming",
		"vercel/react-best-practices",
		"vercel/react-native-skills",
	}

	// Verify
	for i, skill := range skills {
		if skill.FullName != expected[i] {
			t.Errorf("Expected position %d to be %s, got %s", i, expected[i], skill.FullName)
		}
	}
}

func TestSortSkillsCaseInsensitive(t *testing.T) {
	// Test case insensitivity
	skills := []SkillItem{
		{FullName: "VERCEL/React-Best-Practices"},
		{FullName: "anthropic/frontend-design"},
		{FullName: "Superpowers/Brainstorming"},
	}

	SortSkills(skills)

	// Should still sort correctly regardless of case
	if len(skills) != 3 {
		t.Errorf("Expected 3 skills, got %d", len(skills))
	}

	// Verify that sorting is consistent
	// anthropic should come before Superpowers, which should come before VERCEL
	if skills[0].FullName != "anthropic/frontend-design" {
		t.Errorf("First skill should be anthropic/frontend-design, got %s", skills[0].FullName)
	}
}

func TestFilterSkills(t *testing.T) {
	skills := []SkillItem{
		{FullName: "anthropic/frontend-design", Description: "Frontend design best practices"},
		{FullName: "anthropic/algorithmic-art", Description: "Creating algorithmic art"},
		{FullName: "vercel/react-best-practices", Description: "React performance optimization"},
	}

	// Filter by source name
	result := FilterSkills(skills, "anthropic")
	if len(result) != 2 {
		t.Errorf("Expected 2 skills matching 'anthropic', got %d", len(result))
	}

	// Filter by skill name
	result = FilterSkills(skills, "react")
	if len(result) != 1 {
		t.Errorf("Expected 1 skill matching 'react', got %d", len(result))
	}

	// Filter by description
	result = FilterSkills(skills, "art")
	if len(result) != 1 {
		t.Errorf("Expected 1 skill matching 'art', got %d", len(result))
	}

	// Empty filter should return all
	result = FilterSkills(skills, "")
	if len(result) != 3 {
		t.Errorf("Expected 3 skills with empty filter, got %d", len(result))
	}
}

func TestGetSkillByFullName(t *testing.T) {
	skills := []SkillItem{
		{FullName: "anthropic/frontend-design", Name: "frontend-design"},
		{FullName: "anthropic/algorithmic-art", Name: "algorithmic-art"},
		{FullName: "vercel/react-best-practices", Name: "react-best-practices"},
	}

	// Found
	skill := GetSkillByFullName(skills, "anthropic/frontend-design")
	if skill == nil {
		t.Error("Expected to find anthropic/frontend-design")
	}
	if skill.Name != "frontend-design" {
		t.Errorf("Expected Name to be frontend-design, got %s", skill.Name)
	}

	// Not found
	skill = GetSkillByFullName(skills, "nonexistent/skill")
	if skill != nil {
		t.Error("Expected nil for nonexistent skill")
	}
}

func TestCheckSkillInstalled(t *testing.T) {
	// Test with non-existent directory
	installed := CheckSkillInstalled("nonexistent-skill", "/tmp/nonexistent-dir-12345")
	if installed {
		t.Error("Expected false for non-existent directory")
	}

	// Test with empty target dir
	installed = CheckSkillInstalled("some-skill", "")
	if installed {
		t.Error("Expected false for empty target dir")
	}
}
