package cmd

import (
	"fmt"

	"github.com/bagelbrantley/skillstack/pkg/client"
	"github.com/bagelbrantley/skillstack/pkg/style"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:     "show [project_id or lesson_id]",
	Aliases: []string{"view", "instructions"},
	Short:   "Show details and instructions for a project or lesson",
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		server := GetServerURL()
		apiClient := client.NewClient(server)

		targetID := "project-kv-store"
		if len(args) > 0 {
			targetID = args[0]
		}

		proj, err := apiClient.GetProject("project-kv-store")
		if err != nil {
			return err
		}

		// Check if target matches a specific lesson ID
		var matchedLesson *client.Lesson
		var matchedChapter *client.Chapter
		for _, ch := range proj.Chapters {
			for _, l := range ch.Lessons {
				if l.ID == targetID || l.Number == targetID || l.CheckpointTitle == targetID {
					lessonCopy := l
					chapterCopy := ch
					matchedLesson = &lessonCopy
					matchedChapter = &chapterCopy
					break
				}
			}
		}

		if matchedLesson != nil {
			if IsJSONMode() {
				return style.PrintJSON(out, matchedLesson)
			}
			style.PrintHeader(out, fmt.Sprintf("Lesson %s: %s", matchedLesson.Number, matchedLesson.Title))
			style.PrintKeyValue(out, "Chapter", fmt.Sprintf("CH%d: %s", matchedChapter.Number, matchedChapter.Title))
			style.PrintKeyValue(out, "Checkpoint", matchedLesson.CheckpointTitle)
			style.PrintKeyValue(out, "CLI Command", matchedLesson.CLICommand)
			fmt.Fprintln(out)
			fmt.Fprintln(out, style.AccentStyle.Render("Instructions:"))
			fmt.Fprintln(out, matchedLesson.Instructions)
			fmt.Fprintln(out)
			fmt.Fprintln(out, style.MutedStyle.Render("Expected stdout:"))
			fmt.Fprintln(out, matchedLesson.ExpectedOutput)
			return nil
		}

		if IsJSONMode() {
			return style.PrintJSON(out, proj)
		}

		style.PrintHeader(out, proj.Title)
		style.PrintKeyValue(out, "Project ID", proj.ID)
		style.PrintKeyValue(out, "Level", proj.Level)
		style.PrintKeyValue(out, "Language", proj.Language)
		style.PrintKeyValue(out, "Est. Duration", proj.Duration)
		fmt.Fprintln(out)
		fmt.Fprintln(out, style.MutedStyle.Render(proj.Summary))
		fmt.Fprintln(out)

		style.PrintHeader(out, "Chapters & Checkpoints")
		for _, ch := range proj.Chapters {
			fmt.Fprintln(out, style.AccentStyle.Render(fmt.Sprintf("Chapter %d: %s", ch.Number, ch.Title)))
			for _, l := range ch.Lessons {
				fmt.Fprintln(out, fmt.Sprintf("  • Checkpoint %s [%s]: %s", l.Number, l.ID, l.Title))
				fmt.Fprintln(out, style.MutedStyle.Render(fmt.Sprintf("    Command: %s", l.CLICommand)))
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
