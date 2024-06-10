//go:build !darwin
// +build !darwin

package main

import (
	"github.com/alecthomas/kong"
	"github.com/rjkroege/wikitools/wiki"
)

// startmessageloop just directly executes the _main.
func startmessageloop(ctx *kong.Context, settings *wiki.Settings) {
	_main(ctx, settings)
}
