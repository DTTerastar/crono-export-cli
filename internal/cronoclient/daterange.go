// Package cronoclient wraps github.com/jrmycanady/gocronometer with a small
// session helper, typed export methods that return JSON-ready Go values, and
// shared Cobra flag plumbing for date-range selection.
package cronoclient

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

const dateLayout = "2006-01-02"

// DateRange is an inclusive [Start, End] window.  Only the calendar date
// (YYYY-MM-DD) of each endpoint is sent to Cronometer's export endpoints,
// so the time-of-day and zone on these values don't round-trip — but the
// calendar date is resolved in the user's local zone so that --today
// matches the day the user sees in the Cronometer UI.
type DateRange struct {
	Start time.Time
	End   time.Time
}

// AddDateRangeFlags binds --start, --end, --days, --today on cmd.  Each
// subcommand calls this so they all share the same flag vocabulary.
func AddDateRangeFlags(cmd *cobra.Command) {
	cmd.Flags().String("start", "", "start date (YYYY-MM-DD)")
	cmd.Flags().String("end", "", "end date (YYYY-MM-DD), defaults to today")
	cmd.Flags().Int("days", 0, "convenience: last N days ending today")
	cmd.Flags().Bool("today", false, "convenience: today only")
}

// ParseDateRangeFromFlags reads the date-range flags off cmd and resolves
// them into a concrete DateRange.  Default when no flags are passed: the
// last 7 days ending today.  "Today" is the user's local calendar day.
func ParseDateRangeFromFlags(cmd *cobra.Command) (DateRange, error) {
	startStr, _ := cmd.Flags().GetString("start")
	endStr, _ := cmd.Flags().GetString("end")
	days, _ := cmd.Flags().GetInt("days")
	today, _ := cmd.Flags().GetBool("today")
	return resolveDateRange(startStr, endStr, days, today, time.Now())
}

func resolveDateRange(startStr, endStr string, days int, today bool, ref time.Time) (DateRange, error) {
	y, m, d := ref.Date()
	now := time.Date(y, m, d, 0, 0, 0, 0, ref.Location())
	var start, end time.Time

	switch {
	case today:
		start, end = now, now
	case days > 0:
		end = now
		start = now.AddDate(0, 0, -(days - 1))
	case startStr == "" && endStr == "":
		end = now
		start = now.AddDate(0, 0, -6)
	default:
		var err error
		if startStr != "" {
			start, err = time.ParseInLocation(dateLayout, startStr, ref.Location())
			if err != nil {
				return DateRange{}, fmt.Errorf("bad --start: %w", err)
			}
		}
		if endStr != "" {
			end, err = time.ParseInLocation(dateLayout, endStr, ref.Location())
			if err != nil {
				return DateRange{}, fmt.Errorf("bad --end: %w", err)
			}
		} else {
			end = now
		}
		if start.IsZero() {
			start = end
		}
	}

	if end.Before(start) {
		return DateRange{}, fmt.Errorf("--end (%s) is before --start (%s)",
			end.Format(dateLayout), start.Format(dateLayout))
	}
	return DateRange{Start: start, End: end}, nil
}
