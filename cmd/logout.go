package cmd

import (
	"fmt"

	"github.com/bagelbrantley/skillstack/pkg/config"
	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Sign out and clear local SkillStack credentials",
	Long:  "Removes local stored credentials file (~/.config/skillstack/credentials.json).",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		creds, err := config.LoadCredentials()
		if err != nil {
			return fmt.Errorf("error checking local credentials: %w", err)
		}

		if creds == nil {
			if IsJSONMode() {
				return style.PrintJSON(out, map[string]interface{}{"status": "already_logged_out"})
			}
			style.PrintWarn(out, "No active SkillStack credentials found.")
			return nil
		}

		if err := config.ClearCredentials(); err != nil {
			return fmt.Errorf("failed to clear credentials: %w", err)
		}

		if IsJSONMode() {
			return style.PrintJSON(out, map[string]interface{}{"status": "success", "cleared_user": creds.Email})
		}

		style.PrintSuccess(out, fmt.Sprintf("Successfully signed out (%s). Saved credentials removed.", creds.Email))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
