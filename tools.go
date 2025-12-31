//go:build tools
// +build tools

// This file ensures that development dependencies are tracked in go.mod
// even though they're not imported in the main codebase.
// See: https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	_ "github.com/gravityblast/fresh"
)

