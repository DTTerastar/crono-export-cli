package cmd

import (
	"github.com/spf13/cobra"

	"github.com/quantcli/crono-export-cli/internal/cronoclient"
)

var notesCmd = &cobra.Command{
	Use:   "notes",
	Short: "Export user-entered notes",
	RunE: func(cmd *cobra.Command, _ []string) error {
		rng, err := cronoclient.ParseDateRangeFromFlags(cmd)
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		c, err := cronoclient.NewLoggedIn(ctx)
		if err != nil {
			return err
		}
		defer c.Logout()
		rows, err := c.Notes(ctx, rng)
		if err != nil {
			return err
		}
		return emitJSON(rows)
	},
}

func init() {
	cronoclient.AddDateRangeFlags(notesCmd)
	rootCmd.AddCommand(notesCmd)
}
