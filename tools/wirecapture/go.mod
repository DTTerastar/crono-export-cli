// Sub-module so the developer-only wirecapture tool can keep importing
// gocronometer (GPL-2.0) without contaminating crono-export-cli's main
// go.mod. The published binary never links this module.
module github.com/quantcli/crono-export-cli/tools/wirecapture

go 1.25.10

require github.com/jrmycanady/gocronometer v1.5.1

require golang.org/x/net v0.46.0 // indirect
