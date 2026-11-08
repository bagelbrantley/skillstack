package cmd

import (
	"fmt"

	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/bagelbrantley/skillstack/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the SkillStack CLI version and build metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		info := version.GetInfo()

		if IsJSONMode() {
			return style.PrintJSON(out, info)
		}

		style.PrintHeader(out, "SkillStack CLI Version")
		fmt.Fprintln(out, info.String())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
