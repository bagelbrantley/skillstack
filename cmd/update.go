package cmd

import (
	"fmt"
	"os/exec"

	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the SkillStack CLI to the latest version via go install",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		if IsJSONMode() {
			c := exec.Command("go", "install", "github.com/bagelbrantley/skillstack@latest")
			err := c.Run()
			if err != nil {
				return style.PrintJSON(out, map[string]interface{}{"success": false, "error": err.Error()})
			}
			return style.PrintJSON(out, map[string]interface{}{"success": true, "message": "SkillStack CLI updated successfully"})
		}

		style.PrintHeader(out, "Updating SkillStack CLI")
		fmt.Fprintln(out, "Running: go install github.com/bagelbrantley/skillstack@latest ...")

		c := exec.Command("go", "install", "github.com/bagelbrantley/skillstack@latest")
		c.Stdout = cmd.OutOrStdout()
		c.Stderr = cmd.ErrOrStderr()
		err := c.Run()
		if err != nil {
			style.PrintError(out, fmt.Sprintf("Failed to update CLI: %v", err))
			return err
		}

		fmt.Fprintln(out)
		style.PrintSuccess(out, "SkillStack CLI updated successfully to latest version!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
