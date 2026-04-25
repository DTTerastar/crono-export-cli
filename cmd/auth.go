package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print one-line auth readiness state and exit 0 if usable",
	Long: `Print a one-line summary of whether the CLI has the credentials it needs
to talk to Cronometer. Exit code 0 if usable, 1 if missing.

This is a local check — no network call. It only verifies the
CRONOMETER_USERNAME and CRONOMETER_PASSWORD environment variables are set.
The CLI logs in fresh on every invocation, so a successful status here is
necessary but not sufficient: a wrong password will still fail at run time.

Per the quantcli shared contract:
https://github.com/quantcli/common/blob/main/CONTRACT.md#5-auth`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		user := os.Getenv("CRONOMETER_USERNAME")
		pass := os.Getenv("CRONOMETER_PASSWORD")
		switch {
		case user == "" && pass == "":
			return fmt.Errorf("missing CRONOMETER_USERNAME and CRONOMETER_PASSWORD")
		case user == "":
			return fmt.Errorf("missing CRONOMETER_USERNAME")
		case pass == "":
			return fmt.Errorf("missing CRONOMETER_PASSWORD")
		}
		fmt.Fprintf(cmd.OutOrStdout(), "credentials present for %s (env-var auth, no token cache)\n", user)
		return nil
	},
}

func init() {
	authCmd.AddCommand(authStatusCmd)
	rootCmd.AddCommand(authCmd)
}
