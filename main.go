package main

import (
	"log"

	"github.com/alecthomas/kong"
)

var CLI struct {
	ConfigFile string `type:"path" help:"Set alternate configuration file" default:"~/.config/wiki/wiki.json"`

	// Might have a special mode for invoking from Alfred right?
	New struct {
		Config string `arg name:"config" help:"Defined configuration for instance"`
		Name   string `arg name:"name" help:"Name of instance"`
	} `cmd help:"Create new wiki article"`

	Test struct {
	} `cmd help:"Stub action, does nothing."`

}

func main() {
	ctx := kong.Parse(&CLI)

	// TODO(rjk): Read the configuration.

	switch ctx.Command() {
	case "test":
		log.Println("Got a test command")
	default:
		panic(ctx.Command())
	}
}
