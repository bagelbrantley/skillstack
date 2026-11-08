package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bagelbrantley/skillstack/pkg/client"
	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var submitFlag bool

const (
	DefaultProjectUUID    = "7a2b9c3d-4e5f-41a2-9b8c-1d2e3f4a5b6c"
	DefaultCheckpoint1UUID = "c8d9e0f1-2a3b-4c5d-9e8f-7a6b5c4d3e2f"
	DefaultCheckpoint2UUID = "f1a2b3c4-5d6e-4f7a-8b9c-0d1e2f3a4b5c"
	DefaultCheckpoint3UUID = "a9b8c7d6-5e4f-3a2b-1c0d-9e8f7a6b5c4d"
)

var runCmd = &cobra.Command{
	Use:   "run [lesson_uuid]",
	Short: "Run local lesson tests and verification checks",
	Long:  "Executes local tests for the current guided project lesson. Pass -s / --submit to run checks and submit in one step.",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		server := GetServerURL()
		apiClient := client.NewClient(server)

		lessonID := DefaultCheckpoint1UUID
		if len(args) > 0 {
			lessonID = args[0]
		}

		workDir := "."
		if _, err := os.Stat("engine-kv"); err == nil {
			workDir = "engine-kv"
		}

		testPkg := "./internal/framer/..."
		checkpointNum := "1.1"

		if lessonID == DefaultCheckpoint2UUID || lessonID == "proj-lesson-1-2" || lessonID == "1.2" || lessonID == "2" {
			lessonID = DefaultCheckpoint2UUID
			testPkg = "./internal/storage/..."
			checkpointNum = "1.2"
		} else if lessonID == DefaultCheckpoint3UUID || lessonID == "proj-lesson-2-1" || lessonID == "2.1" || lessonID == "3" {
			lessonID = DefaultCheckpoint3UUID
			testPkg = "./internal/index/..."
			checkpointNum = "2.1"
		} else {
			lessonID = DefaultCheckpoint1UUID
			testPkg = "./internal/framer/..."
			checkpointNum = "1.1"
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
				ProjectID:        DefaultProjectUUID,
				LessonID:         lessonID,
				CheckpointNumber: checkpointNum,
				ExitCode:         exitCode,
				Output:           outputStr,
			})
			if err != nil {
				return fmt.Errorf("submission error: %w", err)
			}

			if IsJSONMode() {
				return style.PrintJSON(out, subRes)
			}

			style.PrintSuccess(out, fmt.Sprintf("Checkpoint %s successfully submitted!", checkpointNum))
			style.PrintKeyValue(out, "Project UUID", subRes.ProjectID)
			style.PrintKeyValue(out, "Lesson UUID", subRes.LessonID)
			style.PrintKeyValue(out, "XP Awarded", fmt.Sprintf("+%d XP", subRes.XPAwarded))
			style.PrintKeyValue(out, "Timestamp", subRes.Timestamp)
		}

		return nil
	},
}

func init() {
	runCmd.Flags().BoolVarP(&submitFlag, "submit", "s", false, "Submit verified progress to SkillStack after tests pass")
	rootCmd.AddCommand(runCmd)
}
