package main

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunCallsPostRunOnSuccess(t *testing.T) {
	cmd := &cobra.Command{
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
	cmd.SetArgs([]string{})

	called := false
	code := run(cmd, func() { called = true })
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !called {
		t.Fatal("expected postRun to be called")
	}
}

func TestRunCallsPostRunOnError(t *testing.T) {
	cmd := &cobra.Command{
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("boom")
		},
	}
	cmd.SetArgs([]string{})

	called := false
	code := run(cmd, func() { called = true })
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if !called {
		t.Fatal("expected postRun to be called")
	}
}
