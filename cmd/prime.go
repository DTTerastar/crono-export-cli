package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const primeText = `crono-export — primer for LLM agents
=====================================

WHAT IT IS
  A CLI that reads your personal Cronometer data (per-food log, daily totals,
  weight/biometrics, exercises, notes) and prints it as JSON on stdout.

I/O CONTRACT
  - Output: pretty-printed JSON ARRAY on stdout. THIS IS THE ONLY MODE.
    There is no --json / --format flag; JSON is always the format.
  - Errors: human-readable text on stderr. You do NOT need '2>&1'.
  - Empty result '[]' = success with zero rows in the window, not an error.
  - Exit code: 0 success, non-zero on auth or network failure.
  - Filter with jq. Don't pipe to python — output is already JSON.

AUTH
  Set both env vars before invoking. No config file or token cache; the CLI
  logs in on every run.
    CRONOMETER_USERNAME   your Cronometer email
    CRONOMETER_PASSWORD   your Cronometer password

DATE FLAGS  (every export subcommand accepts these)
  --today                              just today (LOCAL calendar date)
  --days N                             last N days, ending today
  --start YYYY-MM-DD --end YYYY-MM-DD  explicit inclusive window
  (no flag)                            last 7 days, ending today

SUBCOMMANDS

  servings   — per-food log: one row per food eaten, full nutrient breakdown.
    Typed numbers. Keys (subset):
      RecordedTime, Group, FoodName, QuantityValue, QuantityUnits,
      EnergyKcal, ProteinG, CarbsG, FiberG, FatG, SodiumMg, CalciumMg,
      IronMg, B12Mg, VitaminDUI, Omega3G, Omega6G, ... (60+ nutrients).

  nutrition  — daily totals: one row per day across all foods logged that day.
    String-keyed (raw CSV columns, ALL VALUES ARE STRINGS — cast in jq).
    Keys (subset):
      "Date", "Energy (kcal)", "Protein (g)", "Carbs (g)", "Fat (g)",
      "Fiber (g)", "Sodium (mg)", "Iron (mg)", "Calcium (mg)",
      "B12 (Cobalamin) (µg)", "Cholesterol (mg)", "Completed", ...

  biometrics — weight, body fat, blood pressure, custom metrics.
    Typed. Keys: RecordedTime, Metric, Unit, Amount.

  exercises  — logged cardio / strength / custom activities.
    Typed. Keys: RecordedTime, Exercise, Minutes, CaloriesBurned, Group.

  notes      — user-entered notes per day. String-keyed (raw CSV).

EXAMPLES

  # Today's macros, as numbers
  crono-export nutrition --today | jq '.[] | {
    date:    .Date,
    kcal:    (."Energy (kcal)" | tonumber),
    protein: (."Protein (g)"   | tonumber),
    carbs:   (."Carbs (g)"     | tonumber),
    fat:     (."Fat (g)"       | tonumber)
  }'

  # 7-day protein total (servings is typed — no tonumber needed)
  crono-export servings --days 7 | jq '[.[] | .ProteinG] | add'

  # All foods from today's breakfast
  crono-export servings --today | jq '[.[] | select(.Group == "Breakfast") | .FoodName]'

  # Latest weight reading in a 30-day window
  crono-export biometrics --days 30 | jq 'map(select(.Metric == "Weight")) | sort_by(.RecordedTime) | last'

GOTCHAS
  - "Today" is your LOCAL calendar day, not UTC.
  - 'nutrition' and 'notes' values are STRINGS (raw CSV) — cast with
    'jq tonumber' when doing math. 'servings', 'biometrics', 'exercises'
    are already typed numbers.
  - Cronometer logs by calendar day; nothing here is real-time. The same
    --today call moments apart returns the same data.
`

var primeCmd = &cobra.Command{
	Use:   "prime",
	Short: "Print an LLM-targeted primer (I/O contract, subcommands, jq recipes)",
	Long: `Print a one-screen primer aimed at LLM agents calling this CLI as a tool.
Covers the output contract (JSON-on-stdout, no --json flag), auth env vars,
the subcommands and what their rows look like, the shared date flags, and a
few jq recipes for common questions.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		_, err := fmt.Fprint(cmd.OutOrStdout(), primeText)
		return err
	},
}

func init() {
	rootCmd.AddCommand(primeCmd)
}
