package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bagelbrantley/skillstack/pkg/client"
	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var submitFlag bool

var runCmd = &cobra.Command{
	Use:   "run [lesson_uuid]",
	Short: "Run local lesson tests and verification checks",
	Long:  "Executes local tests for the current guided project lesson. Pass -s / --submit to run checks and submit in one step.",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		server := GetServerURL()
		apiClient := client.NewClient(server)

		lessonID := "proj-lesson-1-1"
		if len(args) > 0 {
			lessonID = args[0]
		}

		// Detect working directory (check if inside engine-kv or current dir)
		workDir := "."
		if _, err := os.Stat("engine-kv"); err == nil {
			workDir = "engine-kv"
		}

		testPkg := "./internal/framer/..."
		if lessonID == "proj-lesson-1-2" || lessonID == "1.2" || lessonID == "2" {
			testPkg = "./internal/storage/..."
		} else if lessonID == "proj-lesson-2-1" || lessonID == "2.1" || lessonID == "3" {
			testPkg = "./internal/index/..."
		} else if lessonID == "proj-lesson-1-1" || lessonID == "1.1" || lessonID == "1" {
			testPkg = "./internal/framer/..."
		}

		// Execute local test command
		execCmd := exec.Command("go", "test", "-v", testPkg)
		execCmd.Dir = workDir
		outputBytes, err := execCmd.CombinedOutput()
		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = 1
			}
		}

		outputStr := string(outputBytes)

		if !IsJSONMode() {
			style.PrintHeader(out, fmt.Sprintf("Running Checkpoint Verification [%s]", lessonID))
			fmt.Fprintln(out, style.AccentStyle.Render("$ go test -v "+testPkg))
			fmt.Fprintln(out, style.MutedStyle.Render(strings.Repeat("-", 60)))
			fmt.Fprintln(out, outputStr)
			fmt.Fprintln(out, style.MutedStyle.Render(strings.Repeat("-", 60)))
		}

		if exitCode == 0 {
			if !IsJSONMode() {
				style.PrintSuccess(out, "[SkillStack CLI] Checkpoint verification: SUCCESS")
			}
		} else {
			if !IsJSONMode() {
				style.PrintError(out, fmt.Sprintf("[SkillStack CLI] Tests failed with exit code %d", exitCode))
			}
			if submitFlag {
				return fmt.Errorf("local checks failed (exit code %d). Cannot submit until tests pass", exitCode)
			}
		}

		if submitFlag {
			if !IsJSONMode() {
				fmt.Fprintln(out)
				style.PrintHeader(out, "Submitting Checkpoint to SkillStack Server")
			}

			subRes, err := apiClient.SubmitCheckpoint(client.SubmissionPayload{
				ProjectID:        "project-kv-store",
				LessonID:         lessonID,
				CheckpointNumber: "1.1",
				ExitCode:         exitCode,
				Output:           outputStr,
			})
			if err != nil {
				return fmt.Errorf("submission error: %w", err)
			}

			if IsJSONMode() {
				return style.PrintJSON(out, subRes)
			}

			style.PrintSuccess(out, subRes.Message)
			style.PrintKeyValue(out, "Account", subRes.User)
			style.PrintKeyValue(out, "XP Awarded", fmt.Sprintf("+%d XP", subRes.XPAwarded))
			style.PrintKeyValue(out, "Verified At", subRes.Timestamp)
			return nil
		}

		if IsJSONMode() {
			return style.PrintJSON(out, map[string]interface{}{
				"lesson_id": lessonID,
				"exit_code": exitCode,
				"output":    outputStr,
				"passed":    exitCode == 0,
			})
		}

		if exitCode == 0 {
			fmt.Fprintln(out)
			style.PrintKeyValue(out, "Ready to Submit", "Run `skillstack submit` or `skillstack run -s`")
		}
		return nil
	},
}

func init() {
	runCmd.Flags().BoolVarP(&submitFlag, "submit", "s", false, "Run checks and submit checkpoint in one command")
	rootCmd.AddCommand(runCmd)
}

// Suppress unused import warning if filepath is needed
var _ = filepath.Join
