// Package web embeds the frontend assets into the Go binary so the whole
// application ships as a single executable.
package web

import "embed"

//go:embed static
var FS embed.FS
