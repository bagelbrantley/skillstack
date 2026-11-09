package cmd

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"env", "doctor"},
	Short:   "Check local workspace environment and software dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		server := GetServerURL()

		goVer := GetSystemCommandOutput("go", "version")
		gitVer := GetSystemCommandOutput("git", "--version")

		netStatus := "PASS"
		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get(server + "/api/cli/projects")
		if err != nil || resp.StatusCode != http.StatusOK {
			netStatus = "FAIL (Could not reach " + server + ")"
		} else {
			resp.Body.Close()
		}

		if IsJSONMode() {
			return style.PrintJSON(os.Stdout, map[string]interface{}{
				"os":            runtime.GOOS,
				"arch":          runtime.GOARCH,
				"go_version":    goVer,
				"git_version":   gitVer,
				"server_url":    server,
				"server_status": netStatus,
			})
		}

		style.PrintHeader(os.Stdout, "SkillStack Workspace Diagnostics")
		style.PrintKeyValue(os.Stdout, "Operating System", fmt.Sprintf("%s / %s", runtime.GOOS, runtime.GOARCH))
		style.PrintKeyValue(os.Stdout, "Go Compiler", goVer)
		style.PrintKeyValue(os.Stdout, "Git CLI", gitVer)
		style.PrintKeyValue(os.Stdout, "SkillStack Server", fmt.Sprintf("%s [%s]", server, netStatus))
		fmt.Fprintln(os.Stdout)

		if strings.HasPrefix(netStatus, "PASS") && goVer != "Not installed" {
			style.PrintSuccess(os.Stdout, "Workspace environment ready for Guided Projects.")
		} else {
			style.PrintWarn(os.Stdout, "Environment check completed with warnings. Verify Go installation and server connectivity.")
		}
		_ = os.Stdout.Sync()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
