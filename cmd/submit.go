package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bagelbrantley/skillstack/pkg/client"
	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var submitCmd = &cobra.Command{
	Use:   "submit [lesson_uuid]",
	Short: "Submit verified lesson progress to SkillStack",
	Long:  "Executes local tests and sends authenticated completion proof to SkillStack.",
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

		if exitCode != 0 {
			style.PrintError(out, fmt.Sprintf("Local checks failed with exit code %d. Submission aborted.", exitCode))
			fmt.Fprintln(out, outputStr)
			return fmt.Errorf("checkpoint tests did not pass")
		}

		subRes, err := apiClient.SubmitCheckpoint(client.SubmissionPayload{
			ProjectID:        DefaultProjectUUID,
			LessonID:         lessonID,
			CheckpointNumber: checkpointNum,
			ExitCode:         exitCode,
			Output:           outputStr,
		})
		if err != nil {
			return fmt.Errorf("submission failed: %w", err)
		}

		if IsJSONMode() {
			return style.PrintJSON(out, subRes)
		}

		style.PrintBanner(out, "Checkpoint Submitted & Verified!", fmt.Sprintf("+%d XP Awarded", subRes.XPAwarded))
		style.PrintSuccess(out, subRes.Message)
		style.PrintKeyValue(out, "Account Email", subRes.User)
		style.PrintKeyValue(out, "Project UUID", subRes.ProjectID)
		style.PrintKeyValue(out, "Lesson UUID", subRes.LessonID)
		style.PrintKeyValue(out, "Timestamp", subRes.Timestamp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
