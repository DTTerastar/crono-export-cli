package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const primeText = `crono-export — primer for LLM agents
=====================================

WHAT IT IS
  A CLI that reads your personal Cronometer data (per-food log, daily totals,
  weight/biometrics, exercises, notes) and prints it on stdout.

OUTPUT FORMATS
  Default: narrow, fitdown-style markdown — date-grouped headings, one
  bullet per non-zero field, easy to skim and easy for an LLM to consume
  inline.  Zero-valued nutrients are suppressed in markdown.

  --json  (or --format json)  Pretty-printed JSON ARRAY of full rows.
                              Use this when you want the complete row,
                              when piping to jq, or when round-tripping
                              into other tools.  Nothing is suppressed.

  Errors go to stderr.  You do NOT need '2>&1'.  Exit code is 0 on
  success and non-zero on auth or network failure.  An empty result is
  success — markdown prints "_(no records in window)_", JSON prints '[]'.

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
    Markdown: ## per date, ### per food (group · name · quantity), bullets
    for non-zero nutrients.
    JSON: typed numbers.  Keys (subset):
      RecordedTime, Group, FoodName, QuantityValue, QuantityUnits,
      EnergyKcal, ProteinG, CarbsG, FiberG, FatG, SodiumMg, CalciumMg,
      IronMg, B12Mg, VitaminDUI, Omega3G, ... (60+ nutrients).

  nutrition  — daily totals: one row per day across all foods logged that day.
    Markdown: ## per date, bullets for non-zero columns.
    JSON: string-keyed (raw CSV columns, ALL VALUES ARE STRINGS — cast in jq).
    Keys (subset):
      "Date", "Energy (kcal)", "Protein (g)", "Carbs (g)", "Fat (g)",
      "Fiber (g)", "Sodium (mg)", "Iron (mg)", "Calcium (mg)",
      "B12 (Cobalamin) (µg)", "Cholesterol (mg)", "Completed", ...

  biometrics — weight, body fat, blood pressure, custom metrics.
    Markdown: ## per date, bullet per metric: "- Metric: amount unit".
    JSON keys: RecordedTime, Metric, Unit, Amount.

  exercises  — logged cardio / strength / custom activities.
    Markdown: ## per date, one bullet per session.
    JSON keys: RecordedTime, Exercise, Minutes, CaloriesBurned, Group.

  notes      — user-entered notes per day.  Markdown: ## per date with note
    body.  JSON: string-keyed (raw CSV).

EXAMPLES

  # Today's macros, scannable
  crono-export nutrition --today

  # Today's macros, parsed (numbers via tonumber)
  crono-export nutrition --today --json | jq '.[] | {
    date:    .Date,
    kcal:    (."Energy (kcal)" | tonumber),
    protein: (."Protein (g)"   | tonumber)
  }'

  # 7-day protein total (servings is typed — no tonumber needed)
  crono-export servings --days 7 --json | jq '[.[] | .ProteinG] | add'

  # All foods from today's breakfast
  crono-export servings --today --json | jq '[.[] | select(.Group == "Breakfast") | .FoodName]'

  # Latest weight reading in a 30-day window
  crono-export biometrics --days 30 --json | jq 'map(select(.Metric == "Weight")) | sort_by(.RecordedTime) | last'

GOTCHAS
  - "Today" is your LOCAL calendar day, not UTC.
  - 'nutrition' and 'notes' JSON values are STRINGS (raw CSV) — cast with
    'jq tonumber' when doing math.  'servings', 'biometrics', 'exercises'
    JSON values are already typed numbers.
  - Markdown drops zero-valued nutrients to stay readable.  If you need
    every column (including zeros), use --json.
  - Cronometer logs by calendar day; nothing here is real-time.  The same
    --today call moments apart returns the same data.
`

var primeCmd = &cobra.Command{
	Use:   "prime",
	Short: "Print an LLM-targeted primer (output formats, subcommands, jq recipes)",
	Long: `Print a one-screen primer aimed at LLM agents calling this CLI as a tool.
Covers the output formats (markdown by default, --json for structured),
auth env vars, the subcommands and what their rows look like, the shared
date flags, and a few jq recipes for common questions.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		_, err := fmt.Fprint(cmd.OutOrStdout(), primeText)
		return err
	},
}

func init() {
	rootCmd.AddCommand(primeCmd)
}
