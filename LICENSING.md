# Licensing history

`crono-export-cli` is MIT-licensed end to end as of the release that
ships [QUA-37](https://github.com/quantcli/common). This document
records how the project got there.

## Current state

- `crono-export-cli`'s own code: **MIT** (see [LICENSE](LICENSE)).
- The Cronometer HTTP client in `internal/cronoapi/` is fresh-authored
  in-tree, also MIT.
- `go mod why github.com/jrmycanady/gocronometer` returns "not needed
  by main module" — the GPL-2.0 client is no longer a transitive
  dependency.
- Remaining direct/indirect dependencies are tracked in
  [`go.mod`](go.mod); all use permissive (MIT / BSD / Apache-2.0)
  licences.

## History

Releases **prior to the MIT-clean cut** depended on
[`github.com/jrmycanady/gocronometer`](https://github.com/jrmycanady/gocronometer)
(GPL-2.0). Because Go programs link statically, those binaries linked
GPL-2.0 code and inherited the GPL-2.0 source-availability obligations
for the resulting combined work. Users on those builds who never
redistributed binaries were unaffected in practice; redistributors
needed to honour GPL-2.0.

The replacement `internal/cronoapi` package was authored under the
clean-room rules captured in [QUA-12](https://github.com/quantcli/common)
plan v4:

- The Phase 1 protocol spec in [`docs/cronometer-protocol.md`](docs/cronometer-protocol.md)
  was derived only from `crono-export-cli`'s own MIT-licensed call
  sites — no `gocronometer` source was consulted.
- Phase 3a recorded the live HTTP behaviour of `cronometer.com` via a
  fresh-authored `http.RoundTripper` and committed redacted captures
  to [`internal/cronoclient/testdata/cronometer/`](internal/cronoclient/testdata/cronometer/).
  Real session cookies, anti-CSRF tokens, per-export nonces, and
  account data were stripped before commit; only metadata (URL,
  status, headers) survived.
- Phase 3b authored `internal/cronoapi` against the wire-shape document
  produced in Phase 3a and the public API surface (`go doc`) used by
  our own MIT call sites. `gocronometer` source remained out of scope
  for the duration.

The `tools/wirecapture/` sub-module pins `gocronometer` for capture-only
use; it is **not** imported by the production binary and stays out of
the main module's `go.mod`. It is retained so future wire-shape
recaptures (e.g. after Cronometer rotates the GWT permutation hash or
adds new nutrient columns) can be performed without re-introducing the
GPL dependency to production.

## Reproducing the audit

```sh
# Should show the gocronometer dep is no longer needed by the main
# module:
go mod why github.com/jrmycanady/gocronometer

# Should not list any GPL-licensed dependency. license-detection
# tools that read Go module metadata (e.g., go-licenses) will return
# only MIT / BSD / Apache-2.0 packages.
go list -m all
```

## ToS

Cronometer's Terms of Service §10(b) governs automated access. The
board has accepted that risk on record as part of
[QUA-12](https://github.com/quantcli/common) plan v4; this CLI's
licensing posture is independent of that policy question and is
resolved here at the source-code layer only.
