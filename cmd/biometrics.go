package cmd

import (
	"github.com/spf13/cobra"

	"github.com/DTTerastar/crono-export-cli/internal/cronoclient"
)

var biometricsCmd = &cobra.Command{
	Use:   "biometrics",
	Short: "Export biometric records (weight, body fat, blood pressure, custom metrics)",
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
		recs, err := c.Biometrics(ctx, rng)
		if err != nil {
			return err
		}
		return emitJSON(recs)
	},
}

func init() {
	cronoclient.AddDateRangeFlags(biometricsCmd)
	rootCmd.AddCommand(biometricsCmd)
}
