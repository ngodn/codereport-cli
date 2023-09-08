package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ngodn/codereport-cli/cmd/summarize/commits"
	"github.com/spf13/cobra"
)

var (
	summarizeDateFilterStart string
	summarizeDateFilterEnd   string
	summarizeOutputJSON      bool
)

func init() {
	summarizeCommitsCmd.Flags().StringVarP(&summarizeDateFilterStart, "start", "s", "", "specify a start date to filter by. Can be of format YYYY-MM-DD, or a SQLite \"date modifier,\" relative to 'now'")
	summarizeCommitsCmd.Flags().StringVarP(&summarizeDateFilterEnd, "end", "e", "", "specify an end date to filter by. Can be of format YYYY-MM-DD, or a SQLite \"date modifier,\" relative to 'now'")
	summarizeCommitsCmd.Flags().BoolVar(&summarizeOutputJSON, "json", false, "output as JSON")
}

var summarizeCommitsCmd = &cobra.Command{
	Use:   "commits [file pattern]",
	Short: "Print a summary of commit activity",
	Long: `Prints a summary of commit activity in the default repository (either the current directory or supplied by --repo).
Specify a file pattern as an argument to filter for commits that only modified a certain file or directory.
The path is used in a SQL LIKE clause, so use '%' as a wildcard.
Read more here: https://sqlite.org/lang_expr.html#the_like_glob_regexp_and_match_operators
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var pathPattern string
		if len(args) > 0 {
			pathPattern = args[0]
		}

		var ui *commits.TermUI
		var err error
		if ui, err = commits.NewTermUI(pathPattern, summarizeDateFilterStart, summarizeDateFilterEnd); err != nil {
			handleExitError(err)
		}
		defer func() {
			if err := ui.Close(); err != nil {
				handleExitError(err)
			}
		}()

		if summarizeOutputJSON {
			fmt.Println(ui.PrintJSON())
			return
		}

		// check if output is a terminal (https://rosettacode.org/wiki/Check_output_device_is_a_terminal#Go)
		if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
			if _, err := tea.NewProgram(ui).Run(); err != nil {
				handleExitError(err)
			}
		} else {
			fmt.Print(ui.PrintNoTTY())
		}
	},
}
