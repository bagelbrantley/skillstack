package cmd

import (
	"fmt"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"env", "doctor"},
	Short:   "Check local workstation environment and software dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		server := GetServerURL()

		goVer := getCommandOutput("go", "version")
		gitVer := getCommandOutput("git", "--version")

		netStatus := "PASS"
		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get(server + "/api/cli/projects")
		if err != nil || resp.StatusCode != http.StatusOK {
			netStatus = "FAIL (Could not reach " + server + ")"
		} else {
			resp.Body.Close()
		}

		if IsJSONMode() {
			return style.PrintJSON(out, map[string]interface{}{
				"os":            runtime.GOOS,
				"arch":          runtime.GOARCH,
				"go_version":    goVer,
				"git_version":   gitVer,
				"server_url":    server,
				"server_status": netStatus,
			})
		}

		style.PrintHeader(out, "SkillStack Workstation Diagnostics")
		style.PrintKeyValue(out, "Operating System", fmt.Sprintf("%s / %s", runtime.GOOS, runtime.GOARCH))
		style.PrintKeyValue(out, "Go Compiler", goVer)
		style.PrintKeyValue(out, "Git CLI", gitVer)
		style.PrintKeyValue(out, "SkillStack Server", fmt.Sprintf("%s [%s]", server, netStatus))
		fmt.Fprintln(out)

		if strings.HasPrefix(netStatus, "PASS") && goVer != "Not installed" {
			style.PrintSuccess(out, "Workstation environment ready for Guided Projects.")
		} else {
			style.PrintWarn(out, "Environment check completed with warnings. Verify Go installation and server connectivity.")
		}
		return nil
	},
}

func getCommandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	bytes, err := cmd.Output()
	if err != nil {
		return "Not installed"
	}
	return strings.TrimSpace(string(bytes))
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
