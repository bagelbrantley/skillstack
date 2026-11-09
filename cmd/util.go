package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FindGoBinary dynamically locates the Go compiler executable
func FindGoBinary() string {
	if p, err := exec.LookPath("go"); err == nil {
		return p
	}
	home, _ := os.UserHomeDir()
	candidates := []string{
		"/usr/local/go/bin/go",
		"/root/.local/opt/go/bin/go",
		filepath.Join(home, "go", "bin", "go"),
		filepath.Join(home, ".local", "opt", "go", "bin", "go"),
		"/usr/bin/go",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return "go"
}

// GetSystemCommandOutput runs a command safely and returns its trimmed output
func GetSystemCommandOutput(name string, args ...string) string {
	if name == "go" {
		name = FindGoBinary()
	} else if p, err := exec.LookPath(name); err == nil {
		name = p
	}
	cmd := exec.Command(name, args...)
	bytes, err := cmd.Output()
	if err != nil {
		return "Not installed"
	}
	return strings.TrimSpace(string(bytes))
}
