//go:build !dist

package main

import "embed"

// Empty embed.FS for development mode
var distFS embed.FS
