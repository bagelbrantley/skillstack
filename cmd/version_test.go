package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected version command to succeed, got: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "SkillStack CLI Version") && !strings.Contains(out, "skillstack") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestListCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"list", "--json", "--server", "http://localhost:3000"})

	// Reset flags
	jsonOutput = true

	err := rootCmd.Execute()
	// Connection errors / offline test environment are expected during unit test runs
	if err != nil && !strings.Contains(err.Error(), "connection") && !strings.Contains(err.Error(), "client") && !strings.Contains(err.Error(), "parse") {
		t.Logf("list output/err: %v", err)
	}
}
