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

var (
	codeFlag  string
	forceFlag bool
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate CLI with your SkillStack account",
	Long:  "Generates and prompts for a single-use login code from SkillStack to store credentials locally.",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		server := GetServerURL()
		apiClient := client.NewClient(server)

		// Check if already authenticated locally
		if !forceFlag {
			if existing, err := config.LoadCredentials(); err == nil && existing != nil && existing.Token != "" {
				if IsJSONMode() {
					return style.PrintJSON(out, map[string]interface{}{
						"already_logged_in": true,
						"user":              existing.Email,
						"message":           "Already logged in. Run `skillstack logout` to sign out first or pass --force.",
					})
				}
				style.PrintWarn(out, fmt.Sprintf("You are already logged in as %s.", existing.Email))
				style.PrintKeyValue(out, "Config File", "~/.config/skillstack/credentials.json")
				fmt.Fprintln(out)
				fmt.Fprintln(out, style.MutedStyle.Render("Run `skillstack logout` to sign out first, or pass `--force` to re-authenticate."))
				return nil
			}
		}

		code := strings.TrimSpace(codeFlag)
		if code == "" {
			if !IsJSONMode() {
				style.PrintBanner(os.Stdout, "Welcome to the SkillStack CLI!", "Local Workspace Guided Project Integration")
				fmt.Fprintln(os.Stdout)
				fmt.Fprintln(os.Stdout, "Enter your 8-character code from the SkillStack web app (/cli/login or Account Settings):")
				fmt.Fprint(os.Stdout, "Paste your login code: ")
				_ = os.Stdout.Sync()
			}
			if IsJSONMode() {
				return fmt.Errorf("in --json mode, please pass --code flag (-c XXXX-XXXX)")
			}
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil && len(input) == 0 {
				return fmt.Errorf("failed to read login code from terminal stdin: %w", err)
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
	loginCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Force re-authentication even if already logged in")
	rootCmd.AddCommand(loginCmd)
}
