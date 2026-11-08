package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bagelbrantley/skillstack/pkg/client"
	"github.com/bagelbrantley/skillstack/pkg/config"
	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var codeFlag string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate CLI with your SkillStack account",
	Long:  "Generates and prompts for a single-use login code from SkillStack to store credentials locally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		server := GetServerURL()
		apiClient := client.NewClient(server)

		if !IsJSONMode() {
			style.PrintBanner(out, "Welcome to the SkillStack CLI!", "Local Workstation Guided Project Integration")
			fmt.Fprintln(out)
			fmt.Fprintln(out, "Please navigate to:")
			fmt.Fprintln(out, style.AccentStyle.Render(server+"/cli/login"))
			fmt.Fprintln(out)
		}

		code := strings.TrimSpace(codeFlag)
		if code == "" {
			if IsJSONMode() {
				return fmt.Errorf("in --json mode, please pass --code flag")
			}
			fmt.Fprint(out, "Paste your login code: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read login code: %w", err)
			}
			code = strings.TrimSpace(input)
		}

		if code == "" {
			return fmt.Errorf("login code cannot be empty")
		}

		res, err := apiClient.VerifyCode(code)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		creds := config.Credentials{
			Token:     res.Token,
			UID:       res.User.UID,
			Email:     res.User.Email,
			ServerURL: server,
			SavedAt:   time.Now().Format(time.RFC3339),
		}

		if err := config.SaveCredentials(creds); err != nil {
			return fmt.Errorf("failed to store credentials: %w", err)
		}

		if IsJSONMode() {
			return style.PrintJSON(out, map[string]interface{}{
				"authenticated": true,
				"user":          res.User,
			})
		}

		fmt.Fprintln(out)
		style.PrintSuccess(out, fmt.Sprintf("Successfully logged in as %s (%s)", res.User.Email, res.User.UID))
		style.PrintKeyValue(out, "Config File", "~/.config/skillstack/credentials.json")
		style.PrintKeyValue(out, "Permissions", "0600 (Owner Read/Write Only)")
		return nil
	},
}

func init() {
	loginCmd.Flags().StringVarP(&codeFlag, "code", "c", "", "Login code from SkillStack web app")
	rootCmd.AddCommand(loginCmd)
}
