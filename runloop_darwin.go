//go:build darwin
// +build darwin

package main

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/progrium/macdriver/macos"
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/rjkroege/wikitools/wiki"
)

// startmessageloop starts the mac application runloop and runs the core main
// in a go routine.
func startmessageloop(ctx	*kong.Context,  settings *wiki.Settings) {
	// runs macOS application event loop with a callback on success
	macos.RunApp(func(app appkit.Application, delegate *appkit.ApplicationDelegate) {
		// TODO(rjk): Assorted setup here to figure out.
		go func() {
			_main(ctx, settings)
			endmessageloop()
		}()
	})
	// This function never exits.
}

// endmessageloop ends the runloop and exits the application. main code
// after `macos.Run` is not executed according to both the Apple
// documentation and empirical test.
func endmessageloop() {
	log.Println("requesting to end message loop!")
	// Get the application. It ought to already exist. Don't
	// invoke this if the application doesn't exist.
	app := appkit.Application_SharedApplication()

	// This should shutdown the runloop.
	app.Terminate(nil)
}
