package initcmd

import "testing"

func TestNewCommandDoesNotExposeIncludeXFlag(t *testing.T) {
	cmd := NewCommand()
	if cmd.Flags().Lookup("include-x") != nil {
		t.Fatalf("include-x flag should be removed")
	}
}
