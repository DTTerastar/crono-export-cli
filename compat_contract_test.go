//go:build compat

// Compat-test entry point for crono-export-cli.
//
// This file is only compiled under the `compat` build tag, so it does
// not affect the default `go test ./...` run. CI invokes it as
// `go test -tags=compat ./...` after building the export binary and
// exposing its path through CRONO_EXPORT_BIN.
//
// The actual assertions live in github.com/quantcli/common/compat.
// Drift between this CLI and CONTRACT.md surfaces as a failure here.
package main_test

import (
	"os"
	"testing"

	"github.com/quantcli/common/compat"
	"github.com/quantcli/common/compat/dates"
)

func TestContractDates(t *testing.T) {
	bin := os.Getenv("CRONO_EXPORT_BIN")
	if bin == "" {
		t.Skip("CRONO_EXPORT_BIN not set; skipping compat suite")
	}
	// crono is cobra-based: --since/--until live on each data-producing
	// subcommand, not the root binary. The compat suite dispatches per
	// subcommand under a `subcommand=NAME/...` subtree so any single
	// regression surfaces as a named subtest failure.
	dates.RunContract(t, compat.Runner{
		Binary: bin,
		Subcommands: []string{
			"biometrics",
			"exercises",
			"nutrition",
			"servings",
			"notes",
		},
	})
}
