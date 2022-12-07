package main

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/rjkroege/wikitools/cmd"
	"github.com/rjkroege/wikitools/wiki"
)

var CLI struct {
	// TODO(rjk): Move the default location.
	ConfigFile string `type:"path" help:"Set alternate configuration file" default:"~/.wikinewrc"`

	// Might have a special mode for invoking from Alfred right?
	New struct {
		Tagsandtitle []string `arg:"" name:"tagsandtitle" help:"List of article tags and its title"`
	} `cmd help:"Create new wiki article"`

	Bearimport struct {
		Outputdir string `help:"Output directory for importable files" type:"path" default:"./converted"`

		Filestoprocess []string `arg:"" name:"filestoprocess" help:"List of article tags and its title" type:"path"`
	} `cmd help:"Reprocess articles out of Bear for wiki"`
	

	Test struct {
	} `cmd help:"Stub action, does nothing."`
}

func main() {
	ctx := kong.Parse(&CLI)

	// TODO(rjk): wiki => config
	settings, err := wiki.Read(CLI.ConfigFile)
	if err != nil {
		// TODO(rjk): This is not nice. Set things up sensibly.
		log.Panic("No configuration file. Fatai:", err)
	}

	switch ctx.Command() {
	case "new <tagsandtitle>":
		log.Println("should run Wikinew here", CLI.New.Tagsandtitle)
		cmd.Wikinew(settings, CLI.New.Tagsandtitle)
	case "bearimport <filestoprocess>":
		log.Println("should run Bearimport here", CLI.New.Tagsandtitle)
		cmd.Bearimport(settings, CLI.Bearimport.Outputdir, CLI.Bearimport.Filestoprocess)
	case "test":
		log.Println("Got a test command")
	default:
		log.Panic("Missing command: ", ctx.Command())
	}
}
