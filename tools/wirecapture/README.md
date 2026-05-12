# wirecapture — developer-only Cronometer capture tool

`wirecapture` drives the existing `gocronometer`-backed export flow
through a recording `http.RoundTripper` and writes one JSON file per
HTTP exchange to disk. It is the source of truth for the redacted wire
fixtures used by the clean-room client tests in Phase 3b of [QUA-12](https://github.com/quantcli/common).

## Why a separate Go module

This module imports `github.com/jrmycanady/gocronometer` (GPL-2.0). If
it lived in the main `crono-export-cli` module, the GPL transitive
dependency would re-contaminate the published binary's license claim.
Separate `go.mod` keeps `tools/wirecapture/` out of the production
module's `go.sum` and out of `go mod why`.

After Phase 3b drops `gocronometer` from the main module, this tool
remains buildable independently and continues to work for re-capturing
fixtures.

## Authorship

`wirecapture` was authored without consulting `gocronometer`'s source.
Inputs to authorship were:

- the public surface of `github.com/jrmycanady/gocronometer` as visible
  via `go doc` (exported names, signatures, struct field names and
  types), and
- `crono-export-cli`'s own MIT-licensed call sites
  (`internal/cronoclient/client.go`, `cmd/auth.go`).

The recording transport implementation in `recorder.go` is fresh-authored
under MIT and contains no fragments of `gocronometer` source.

## Usage

See the playbook in
[`internal/cronoclient/testdata/cronometer/CAPTURE.md`](../../internal/cronoclient/testdata/cronometer/CAPTURE.md).

Output bytes contain raw session cookies, the anti-CSRF token, the
password POST body, all per-session auth tokens, and the captured user's
profile data. Never commit them. The `.gitignore` rules in this repo
already block `cronometer-capture/` and `tools/wirecapture/captures/`.
