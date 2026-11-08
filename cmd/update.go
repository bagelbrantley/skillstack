package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the SkillStack CLI to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		// Get current running executable path
		currentExe, _ := os.Executable()
		if currentExe != "" {
			if sym, err := filepath.EvalSymlinks(currentExe); err == nil {
				currentExe = sym
			}
		}

		if !IsJSONMode() {
			style.PrintHeader(out, "Updating SkillStack CLI")
			fmt.Fprintln(out, "Fetching latest release directly from GitHub...")
		}

		// Use GOPROXY=direct to bypass proxy.golang.org caching
		c := exec.Command("go", "install", "github.com/bagelbrantley/skillstack@latest")
		c.Env = append(os.Environ(), "GOPROXY=direct")
		if !IsJSONMode() {
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.ErrOrStderr()
		}

		err := c.Run()
		if err != nil {
			if IsJSONMode() {
				return style.PrintJSON(out, map[string]interface{}{"success": false, "error": err.Error()})
			}
			style.PrintError(out, fmt.Sprintf("Failed to update CLI via go install: %v", err))
			return err
		}

		// Determine where Go installed the new binary
		goBinOut, _ := exec.Command("go", "env", "GOBIN").Output()
		goBin := strings.TrimSpace(string(goBinOut))
		if goBin == "" {
			goPathOut, _ := exec.Command("go", "env", "GOPATH").Output()
			goPath := strings.TrimSpace(string(goPathOut))
			if goPath == "" {
				home, _ := os.UserHomeDir()
				goPath = filepath.Join(home, "go")
			}
			goBin = filepath.Join(goPath, "bin")
		}

		installedBinary := filepath.Join(goBin, "skillstack")

		// If current running executable is outside GOBIN (e.g. /usr/local/bin or local ./bin), update it as well
		if currentExe != "" && currentExe != installedBinary {
			if _, err := os.Stat(installedBinary); err == nil {
				_ = copyFile(installedBinary, currentExe)
			}
		}

		if IsJSONMode() {
			return style.PrintJSON(out, map[string]interface{}{
				"success": true,
				"message": "SkillStack CLI updated successfully",
				"binary":  installedBinary,
			})
		}

		fmt.Fprintln(out)
		style.PrintSuccess(out, "SkillStack CLI updated successfully!")
		style.PrintKeyValue(out, "Installed Binary", installedBinary)
		if currentExe != "" && currentExe != installedBinary {
			style.PrintKeyValue(out, "Active Executable", currentExe)
		}
		return nil
	},
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, info.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
