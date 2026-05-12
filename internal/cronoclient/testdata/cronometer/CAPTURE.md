# Re-capturing Cronometer wire fixtures

Phase 3 of the clean-room replacement uses captured HTTP behaviour as
its source of truth for the wire layer. Use this playbook when you need
to refresh the captures (e.g., Cronometer rotated the GWT permutation
hash, added an endpoint, or you need a previously-uncaptured failure
mode).

The capture tool itself lives at `tools/wirecapture/`.

## What you need

- A real Cronometer account with a small amount of recent data — at
  minimum, a few logged servings in the last 7 days, so the
  `servings` CSV is non-trivial. 2FA off (the capture tool does not
  prompt for a TOTP code).
- Local `go` toolchain (anything matching the version in `go.mod`).
- A working directory outside the repo for the unredacted dumps. They
  contain raw cookies, session tokens, the password POST body, and your
  Cronometer profile data. They must never be committed.

## Steps

1. Build the capture tool:

   ```bash
   cd tools/wirecapture
   go build -o /tmp/wirecapture .
   ```

2. Pick a local-only output directory:

   ```bash
   mkdir -p ~/cronometer-capture
   chmod 700 ~/cronometer-capture
   ```

3. Run the capture against a recent window. **Inline the credentials
   into the subshell only**; do not `export` them into your interactive
   shell history.

   ```bash
   WIRE_CAPTURE_DIR=$HOME/cronometer-capture \
   CRONOMETER_USERNAME='you@example.com' \
   CRONOMETER_PASSWORD='your-password' \
     /tmp/wirecapture --since 2026-05-04 --until 2026-05-11
   ```

   Expect 14 numbered JSON files in `$WIRE_CAPTURE_DIR`: one per HTTP
   exchange. Each file holds the request method/URL/headers/body and
   the response status/headers/body verbatim.

4. Inspect the unredacted captures locally to understand any new
   behaviour. Do not paste them into PRs, issues, comments, or commit
   messages.

5. Produce the redacted fixtures committed under
   `internal/cronoclient/testdata/cronometer/`. Use the redaction script
   approach in
   [QUA-37 phase 3a PR](https://github.com/quantcli/crono-export-cli/pulls)
   as a reference: strip cookie values, anti-CSRF token, all per-session
   tokens (32-char hex), the password POST body, the username, and any
   real account data from response bodies. Never commit response
   bodies that contain real user data.

## Sensitive content guide

These appear in unredacted captures and **must be stripped before
commit**:

| Where | What | How to redact |
| ----- | ---- | ------------- |
| `Cookie:` request headers | session cookies | replace value with `<REDACTED-COOKIE-VALUES>` |
| `Set-Cookie:` response headers | session cookies (4 in step 1, 2 in step 2, 1 elsewhere) | replace value but keep the cookie name shape |
| step (1) response body | full login HTML incl. `<input name="anticsrf" value="...">` | drop the body; document the input shape in `WIRE_SHAPES.md` |
| step (2) request body | `anticsrf=...&password=...&username=...` | drop the body; document key order in `WIRE_SHAPES.md` |
| step (3) response body | `authenticate` GWT response — full user payload incl. name, weight, height, DOB, timezone, custom-charts JSON, session auth token | drop the body |
| step (4/6/8/10/12) response body | `generateAuthorizationToken` response with the per-export nonce | drop the body |
| step (5/7/9/11/13) response body | CSV with real food/biometric/exercise/note data | drop the body; preserve only the column header row in a `headers/<type>.csv` companion if needed |
| step (5/7/9/11/13) request URL `nonce=` query param | per-export auth token | replace with `<REDACTED>` |
| step (14) request body | `logout` GWT call carrying the session auth token | drop the body |

These are **not sensitive** and may be committed verbatim:

- `X-Gwt-Permutation` header value — public deployment hash served to
  every browser.
- `X-Gwt-Module-Base` header value — public URL.
- `Content-Type`, `Date`, `Strict-Transport-Security`, etc. — server
  policy headers.
- Query string keys (`start`, `end`, `generate`) — public surface.

## How to add a new endpoint

If you need to capture a behaviour the existing capture binary does not
exercise (e.g., a fresh export type, a deliberate error path), edit
`tools/wirecapture/main.go` to add a new step. Each step is a closure
that calls a `gocronometer.Client` method; the recording transport
handles dump generation transparently.
