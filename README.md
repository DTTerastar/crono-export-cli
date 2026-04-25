# crono-export-cli

Export your personal nutrition, biometric, and food-log data from [Cronometer](https://cronometer.com) as JSON. Built for personal LLM agents and scripts that want structured nutrition data — for example, an LLM-driven bariatric or fitness coach that needs to know how much protein, B12, iron, or calcium you actually got today.

[![Latest Release](https://img.shields.io/github/v/release/quantcli/crono-export-cli)](https://github.com/quantcli/crono-export-cli/releases/latest)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/github/go-mod/go-version/quantcli/crono-export-cli)](go.mod)
![Platforms](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey)

## Features

- **Five export endpoints** — servings (per-food log with full nutrient breakdown), nutrition (daily totals), biometrics (weight, body fat, custom metrics), exercises, and notes
- **JSON on stdout** — pipe straight to `jq`, save to disk, or hand to an LLM tool
- **Date selection** — `--today`, `--days N`, or `--start YYYY-MM-DD --end YYYY-MM-DD` on every subcommand
- **Single static binary** — no Python or Node runtime; drop it in `~/bin/` and go
- **Credentials via env** — `CRONOMETER_USERNAME` / `CRONOMETER_PASSWORD`, no config file needed
- **Built for agents** — designed to be called as a terminal tool by LLMs (Claude, hermes-agent, etc.); run `crono-export prime` for a one-screen orientation (I/O contract, subcommands, jq recipes)

## Quick Start

```sh
# Install with Homebrew
brew tap quantcli/tap
brew install crono-export

# Set credentials and try a query
export CRONOMETER_USERNAME="you@example.com"
export CRONOMETER_PASSWORD="…"
crono-export servings --today
```

## Install

**Homebrew (macOS / Linux):**
```sh
brew tap quantcli/tap
brew install crono-export
```

Or download a pre-built binary from the [releases page](https://github.com/quantcli/crono-export-cli/releases/latest):

**macOS (Apple Silicon):**
```sh
curl -Lo /tmp/crono-export.zip https://github.com/quantcli/crono-export-cli/releases/latest/download/crono-export_darwin_arm64.zip
unzip -jo /tmp/crono-export.zip -d ~/bin && rm /tmp/crono-export.zip
chmod +x ~/bin/crono-export
```

**macOS (Intel):**
```sh
curl -Lo /tmp/crono-export.zip https://github.com/quantcli/crono-export-cli/releases/latest/download/crono-export_darwin_amd64.zip
unzip -jo /tmp/crono-export.zip -d ~/bin && rm /tmp/crono-export.zip
chmod +x ~/bin/crono-export
```

**Linux (amd64):**
```sh
curl -Lo /tmp/crono-export.zip https://github.com/quantcli/crono-export-cli/releases/latest/download/crono-export_linux_amd64.zip
unzip -jo /tmp/crono-export.zip -d ~/bin && rm /tmp/crono-export.zip
chmod +x ~/bin/crono-export
```

**Windows (amd64):**

Download `crono-export_windows_amd64.zip` from the [releases page](https://github.com/quantcli/crono-export-cli/releases/latest), extract it, and add the directory to your PATH.

Make sure `~/bin` is in your `PATH`. If not, add this to your `~/.zshrc` or `~/.bashrc`:
```sh
export PATH="$HOME/bin:$PATH"
```

## Credentials

Set these in your shell, in a `.env` file your runner sources, or in your LLM agent's environment:

```sh
export CRONOMETER_USERNAME="you@example.com"
export CRONOMETER_PASSWORD="your-cronometer-password"
```

The CLI logs in on every invocation; there's no token cache. Cronometer doesn't (yet) offer SSO or API tokens for individuals, so a real password is the only auth option.

## Usage

Every subcommand accepts the same date flags:

| Flag | Meaning |
|---|---|
| `--today` | Just today |
| `--days N` | The last N days, ending today |
| `--start YYYY-MM-DD --end YYYY-MM-DD` | Explicit window (inclusive) |
| *(none)* | Last 7 days, ending today |

### Servings — per-food log

One row per food item logged, with full macro and micronutrient breakdown.

```sh
crono-export servings --today
crono-export servings --days 7
crono-export servings --start 2026-04-01 --end 2026-04-15
```

Sample row (truncated):

```json
{
  "RecordedTime": "2026-04-11T00:00:00Z",
  "Group": "Breakfast",
  "FoodName": "Cheese Cracker",
  "QuantityValue": 20,
  "QuantityUnits": "square",
  "EnergyKcal": 97.8,
  "ProteinG": 1.95,
  "CarbsG": 11.88,
  "FiberG": 0.46,
  "FatG": 4.94,
  "B12Mg": 0.07,
  "CalciumMg": 27.2,
  "IronMg": 0.69
}
```

### Nutrition — daily totals

One row per day, totals across every food logged that day.

```sh
crono-export nutrition --days 30
```

### Biometrics — weight, body fat, custom metrics

```sh
crono-export biometrics --days 30
```

```json
[
  {
    "RecordedTime": "2026-04-10T00:00:00Z",
    "Metric": "Weight",
    "Unit": "lbs",
    "Amount": 237
  }
]
```

### Exercises

```sh
crono-export exercises --days 7
```

### Notes

```sh
crono-export notes --days 30
```

## Output Format

Every subcommand prints pretty-printed JSON to stdout. Errors and progress messages go to stderr, so it's safe to pipe stdout into `jq`, redirect into a file, or feed to an LLM tool without worrying about mixed output.

```sh
crono-export servings --today | jq '[.[] | {food: .FoodName, protein: .ProteinG}]'
```

## About Cronometer

[Cronometer](https://cronometer.com) is a nutrition tracking app with one of the best micronutrient databases of any consumer tool — a major reason it's commonly recommended for bariatric patients, anyone tracking specific vitamin/mineral targets, or athletes managing recovery nutrition.

This CLI is an unofficial tool for exporting your own data. It uses the same web export endpoints the Cronometer SPA uses, via [`jrmycanady/gocronometer`](https://github.com/jrmycanady/gocronometer). It is intended for personal single-user use only — see the upstream library's notes on appropriate use.

## License

MIT — see [LICENSE](LICENSE).

The underlying [`gocronometer`](https://github.com/jrmycanady/gocronometer) library is GPLv2-licensed.
