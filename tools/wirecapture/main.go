// Command wirecapture drives the gocronometer client through the five
// export operations consumed by crono-export-cli, recording every HTTP
// exchange to disk for later use as Phase 3 fixtures.
//
// Authored without consulting gocronometer source. Inputs to authorship:
//
//   - the public surface of github.com/jrmycanady/gocronometer (exported
//     names, signatures, and field types as visible via go doc), and
//   - crono-export-cli's own MIT-licensed call sites
//     (internal/cronoclient/client.go).
//
// Use locally only:
//
//	CRONOMETER_USERNAME=... CRONOMETER_PASSWORD=... \
//	WIRE_CAPTURE_DIR=$HOME/cronometer-capture \
//	go run ./tools/wirecapture
//
// Output files contain raw cookies, full responses, and may contain the
// password POST body. Do not commit them. Run the redaction step
// separately to produce sharable fixtures.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jrmycanady/gocronometer"
)

func main() {
	var (
		dir     = flag.String("dir", os.Getenv("WIRE_CAPTURE_DIR"), "directory to write captured exchanges to")
		sinceS  = flag.String("since", "", "start date YYYY-MM-DD (default: 7 days ago, local)")
		untilS  = flag.String("until", "", "end date YYYY-MM-DD (default: today, local)")
		timeout = flag.Duration("timeout", 60*time.Second, "per-request timeout")
	)
	flag.Parse()

	if *dir == "" {
		fail("--dir or WIRE_CAPTURE_DIR is required")
	}
	user := os.Getenv("CRONOMETER_USERNAME")
	pass := os.Getenv("CRONOMETER_PASSWORD")
	if user == "" || pass == "" {
		fail("CRONOMETER_USERNAME and CRONOMETER_PASSWORD must be set")
	}
	if err := os.MkdirAll(*dir, 0o700); err != nil {
		fail("mkdir %s: %v", *dir, err)
	}
	since, err := parseDate(*sinceS, time.Now().AddDate(0, 0, -7))
	if err != nil {
		fail("--since: %v", err)
	}
	until, err := parseDate(*untilS, time.Now())
	if err != nil {
		fail("--until: %v", err)
	}
	if since.After(until) {
		fail("--since after --until")
	}

	client := gocronometer.NewClient(nil)
	if client.HTTPClient == nil {
		client.HTTPClient = &http.Client{}
	}
	rec := newRecordingTransport(*dir, client.HTTPClient.Transport)
	client.HTTPClient.Transport = rec
	if *timeout > 0 {
		client.HTTPClient.Timeout = *timeout
	}

	ctx := context.Background()
	logf("login as %s …", user)
	if err := client.Login(ctx, user, pass); err != nil {
		fail("login: %v", err)
	}

	type step struct {
		label string
		run   func(context.Context) error
	}
	steps := []step{
		{"servings", func(ctx context.Context) error {
			_, err := client.ExportServingsParsedWithLocation(ctx, since, until, time.Local)
			return err
		}},
		{"exercises", func(ctx context.Context) error {
			_, err := client.ExportExercisesParsedWithLocation(ctx, since, until, time.Local)
			return err
		}},
		{"biometrics", func(ctx context.Context) error {
			_, err := client.ExportBiometricRecordsParsedWithLocation(ctx, since, until, time.Local)
			return err
		}},
		{"daily_nutrition", func(ctx context.Context) error {
			_, err := client.ExportDailyNutrition(ctx, since, until)
			return err
		}},
		{"notes", func(ctx context.Context) error {
			_, err := client.ExportNotes(ctx, since, until)
			return err
		}},
	}

	for _, s := range steps {
		logf("%s [%s..%s] …", s.label, since.Format("2006-01-02"), until.Format("2006-01-02"))
		if err := s.run(ctx); err != nil {
			logf("  step %s failed: %v (continuing — partial capture preserved)", s.label, err)
		}
	}

	logf("logout …")
	if err := client.Logout(ctx); err != nil {
		logf("  logout failed: %v (continuing)", err)
	}

	logf("done; %d exchanges written to %s", rec.seq.Load(), *dir)
}

func parseDate(s string, fallback time.Time) (time.Time, error) {
	if s == "" {
		y, m, d := fallback.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.Local), nil
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func logf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "[wirecapture] "+format+"\n", a...)
}

func fail(format string, a ...any) {
	fmt.Fprintf(os.Stderr, "wirecapture: "+format+"\n", a...)
	os.Exit(1)
}
