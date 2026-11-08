package cmd

import (
	"fmt"

	"github.com/bagelbrantley/skillstack/pkg/client"
	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"projects"},
	Short:   "List all available Guided Projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		server := GetServerURL()
		apiClient := client.NewClient(server)

		projects, err := apiClient.GetProjects()
		if err != nil {
			return fmt.Errorf("could not fetch projects list: %w", err)
		}

		if IsJSONMode() {
			return style.PrintJSON(out, projects)
		}

		style.PrintHeader(out, "Available Guided Projects")
		if len(projects) == 0 {
			fmt.Fprintln(out, "No projects available.")
			return nil
		}

		headers := []string{"PROJECT ID / SLUG", "LEVEL", "LANGUAGE", "TITLE"}
		var rows [][]string
		for _, p := range projects {
			rows = append(rows, []string{p.ID, p.Level, p.Language, p.Title})
		}

		fmt.Fprintln(out, style.FormatTable(headers, rows))
		fmt.Fprintln(out, style.MutedStyle.Render("Use `skillstack show <id>` to view project details or `skillstack init <id>` to set up locally."))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
