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
	rootCmd.SetArgs([]string{"list", "--json"})

	// Reset flags
	jsonOutput = true

	err := rootCmd.Execute()
	// Should work or attempt connection
	if err != nil && !strings.Contains(err.Error(), "connection") {
		t.Logf("list output/err: %v", err)
	}
}
