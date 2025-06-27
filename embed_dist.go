//go:build dist

package main

import "embed"

//go:embed all:dist
var distFS embed.FS
