// Package cmd wires the Cobra command tree for crono-export.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "crono-export",
	Short: "Export Cronometer nutrition, biometrics, and food log data as JSON",
	Long: `crono-export reads your personal Cronometer data via the same export
endpoints the web app uses and prints it as JSON on stdout.

Credentials must be supplied via environment variables:
  CRONOMETER_USERNAME  your Cronometer email
  CRONOMETER_PASSWORD  your Cronometer password

Designed for use by personal LLM agents and scripts that want structured
nutrition data — for example, an LLM-driven bariatric or fitness coach.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute runs the root command.  Called by main.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// emitJSON pretty-prints v as JSON to stdout.  Used by every subcommand.
func emitJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
