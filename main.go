package main

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/rjkroege/wikitools/cmd"
	"github.com/rjkroege/wikitools/wiki"
	"github.com/rjkroege/wikitools/corpus"
)

var CLI struct {
	// TODO(rjk): Move the default location.
	ConfigFile string `type:"path" help:"Set alternate configuration file" default:"~/.wikinewrc"`
	Debug      bool   `help:"Enable debugging conveniences as needed."`

	// Might have a different command for Alfred vs Non-alfted case?
	New struct {
		Tagsandtitle []string `arg:"" name:"tagsandtitle" help:"List of article tags and its title"`
	} `cmd help:"Create new wiki article"`

	Preview struct {
		Article string `arg:"" name:"article" type:"path" help:"Article to preview"`
	} `cmd help:"Preview a wiki article"`

	Tidy struct {
		Dryrun    bool `help:"Don't actually move the files, just show what would happen"`
		Deepclean bool `help:"Rewrite the metadata, move files into improved directories"`
		Report    bool `help:"Generate the metadata status report."`
	} `cmd help:"Clean up wiki aritcles: right structure, corrected metadata, etc."`

	List struct {
	} `cmd help:"List all articles using the corpus facilities"`

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
		log.Fatal("No configuration file. Fatai:", err)
	}

	switch ctx.Command() {
	case "new <tagsandtitle>":
		log.Println("should run Wikinew here", CLI.New.Tagsandtitle)
		cmd.Wikinew(settings, CLI.New.Tagsandtitle)
	case "preview <article>":
		// TODO(rjk): Figure out what this is for.
		cmd.Preview(settings, CLI.Debug)
	case "bearimport <filestoprocess>":
		log.Println("should run Bearimport here", CLI.Bearimport.Filestoprocess)
		cmd.Bearimport(settings, CLI.Bearimport.Outputdir, CLI.Bearimport.Filestoprocess)
	case "tidy":
		log.Println("should run Tidy here")
		cmd.Tidy(settings, CLI.Tidy.Dryrun, CLI.Tidy.Deepclean, CLI.Tidy.Report)
	case "test":
		log.Println("Got a test command")
	case "list":
		log.Println("Got a test command")
		tid := corpus.NewListAllWikiFilesTidying()
		if err := corpus.Everyfile(settings, tid); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Missing command: ", ctx.Command())
	}
}
