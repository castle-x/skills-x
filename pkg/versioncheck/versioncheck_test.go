package versioncheck

import "testing"

func TestNormalizeVersion(t *testing.T) {
	cases := map[string]string{
		"v0.2.10":       "0.2.10",
		"0.2.10":        "0.2.10",
		"0.2.10-dirty":  "0.2.10",
		"0.2.10+build1": "0.2.10",
		"dev":           "dev",
	}

	for input, expected := range cases {
		if got := NormalizeVersion(input); got != expected {
			t.Fatalf("NormalizeVersion(%q) = %q, want %q", input, got, expected)
		}
	}
}

func TestLatestFromNpmJSON(t *testing.T) {
	data := []byte(`{"dist-tags":{"latest":"0.2.10"}}`)
	latest, err := LatestFromNpmJSON(data)
	if err != nil {
		t.Fatalf("LatestFromNpmJSON returned error: %v", err)
	}
	if latest != "0.2.10" {
		t.Fatalf("LatestFromNpmJSON = %q, want %q", latest, "0.2.10")
	}
}

func TestShouldNotify(t *testing.T) {
	if !ShouldNotify("0.2.9", "0.2.10") {
		t.Fatal("expected update prompt for older version")
	}
	if ShouldNotify("v0.2.10", "0.2.10") {
		t.Fatal("did not expect prompt for same version")
	}
	if ShouldNotify("0.2.10-dirty", "0.2.10") {
		t.Fatal("did not expect prompt for dirty but same version")
	}
	if ShouldNotify("dev", "0.2.10") {
		t.Fatal("did not expect prompt for dev version")
	}
}
