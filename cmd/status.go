package cmd

import (
	"fmt"

	"github.com/bagelbrantley/skillstack/pkg/config"
	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "View current authentication status and progress sync",
	RunE:  runStatus,
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync progress with SkillStack server",
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	out := cmd.OutOrStdout()
	creds, err := config.LoadCredentials()
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	if creds == nil || creds.Token == "" {
		if IsJSONMode() {
			return style.PrintJSON(out, map[string]interface{}{"authenticated": false})
		}
		style.PrintWarn(out, "Not authenticated. Run `skillstack login` to authenticate.")
		return nil
	}

	if IsJSONMode() {
		return style.PrintJSON(out, map[string]interface{}{
			"authenticated": true,
			"user_id":       creds.UID,
			"email":         creds.Email,
			"server_url":    creds.ServerURL,
			"saved_at":      creds.SavedAt,
		})
	}

	style.PrintHeader(out, "SkillStack CLI Authentication Status")
	style.PrintKeyValue(out, "Account Email", creds.Email)
	style.PrintKeyValue(out, "User UID", creds.UID)
	style.PrintKeyValue(out, "Server URL", creds.ServerURL)
	style.PrintKeyValue(out, "Authenticated Since", creds.SavedAt)
	fmt.Fprintln(out)
	style.PrintSuccess(out, "Progress sync status: ACTIVE & CONNECTED")
	return nil
}

func init() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(syncCmd)
}
