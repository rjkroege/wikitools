package main

import (
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/rjkroege/wikitools/cmd"
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/corpus/tidy"
	"github.com/rjkroege/wikitools/wiki"
)

var CLI struct {
	// TODO(rjk): Move the default location.
	// TODO(rjk): Move these into settings.
	ConfigFile string `type:"path" help:"Set alternate configuration file" default:"~/.wikinewrc"`
	Debug      bool   `help:"Enable debugging conveniences as needed."`
	Dryrun     bool   `help:"Don't actually modify anything but instead just show what would happen."`

	// This is the action case.
	New struct {
		Tagsandtitle []string `arg:"" name:"tagsandtitle" help:"List of article tags and its title"`
	} `cmd help:"Create new wiki article"`

	// New op for Alfred case. (But note that Alfred requires reprocessing the arguments.)
	Newautocomplete struct {
		Tagsandtitle []string `arg:"" name:"tagsandtitle" help:"List of article tags and its title"`
	} `cmd help:"Create new wiki article"`

	Preview struct {
		// TODO(rjk): Need to make this better.
		Article string `arg:"" name:"article" type:"path" help:"Article to preview"`
	} `cmd help:"Preview a wiki article"`

	Plumb struct {
		Wikitext string `arg:"" name:"wikitext" help:"Wikilink text to open"`
	} `cmd help:"Plumb helper to open another wikitext link"`

	Tidy struct {
		// TODO(rjk): Need to write this.
		All           struct{} `cmd help:"Do all possible tidying."  default:"1"`
		Deepclean     struct{} `cmd help:"Modernize the metadata, fix links, etc."`
		Move          struct{} `cmd help:"Move files into the correct places."`
		Findersync    struct{} `cmd help:"Sync metadata info the Spotlight metadata attributes."`
		Updatetaglist struct{} `cmd help:"Update the cached list of tags for tag autocompletion."`
	} `cmd help:"Clean up wiki aritcles: right structure, corrected metadata, etc."`

	Report struct {
		Articles struct{} `cmd help:"List all articles." default:"1"`
		// TODO(rjk): where it puts it should be configurable right?
		Metadata struct{} `cmd help:"List the metadata versions of each article."`
		Tags     struct{} `cmd help:"List all of the in-use tags."`
		Todos    struct{} `cmd help:"List all outstanding TODO items."`
		Urls     struct{} `cmd help:"Report all the URLs"`
	} `cmd help:"Generate reports on the wiki corpus."`

	Bearimport struct {
		Outputdir string `help:"Output directory for importable files" type:"path" default:"./converted"`

		Filestoprocess []string `arg:"" name:"filestoprocess" help:"List of article tags and its title" type:"path"`
	} `cmd help:"Reprocess articles out of Bear for wiki"`
}

func main() {
	os.Args = cmd.ArgsAlfredPreprocess(os.Args)
	ctx := kong.Parse(&CLI)

	// TODO(rjk): wiki => config
	settings, err := wiki.Read(CLI.ConfigFile)
	if err != nil {
		// TODO(rjk): This is not nice. Set things up sensibly.
		log.Fatal("No configuration file. Fatai:", err)
	}

	// Invokes _main inside on platforms where I assume that functionality (e.g. darwin).
	startmessageloop(ctx, settings)
}

func _main(ctx	*kong.Context,  settings *wiki.Settings) {

	switch ctx.Command() {
	case "new <tagsandtitle>":
		log.Println("should run Wikinew here", CLI.New.Tagsandtitle)
		cmd.Wikinew(settings, CLI.New.Tagsandtitle)
	case "newautocomplete <tagsandtitle>":
		cmd.WikinewAutocomplete(settings, CLI.Newautocomplete.Tagsandtitle)
	case "preview <article>":
		// TODO(rjk): Figure out what this is for.
		cmd.Preview(settings, CLI.Debug)
	case "plumb <wikitext>":
		log.Println("plumb!")
		cmd.PlumberHelper(settings, CLI.Plumb.Wikitext)
	case "bearimport <filestoprocess>":
		log.Println("should run Bearimport here", CLI.Bearimport.Filestoprocess)
		cmd.Bearimport(settings, CLI.Bearimport.Outputdir, CLI.Bearimport.Filestoprocess)

	case "tidy all":
		log.Println("tidy all not implemented")
		// TODO(rjk): Union the other operations.
	case "tidy deepclean":
		// TODO(rjk): Highly likely that this needs some kind of settings.
		tidying, err := tidy.NewMetadataUpdater()
		if err != nil {
			log.Fatal("Can't make a MetadataUpdater( because:", err)
		}
		if err := corpus.Everyfile(settings, tidying); err != nil {
			log.Fatal(err)
		}
		if err := tidying.Summary(); err != nil {
			log.Fatal("report Summary: ", err)
		}
	case "tidy updatetaglist":
		tidying, err := tidy.NewTagsDumper(settings)
		if err != nil {
			log.Fatal("Can't make a TagsDumper( because:", err)
		}
		if err := corpus.Everyfile(settings, tidying); err != nil {
			log.Fatal(err)
		}
		if err := tidying.Summary(); err != nil {
			log.Fatal("report Summary: ", err)
		}
	case "tidy findersync":
		log.Println("tidy findersync not implemented")
		// TODO(rjk): Write me.
	case "tidy move":
		tidying, _ := tidy.NewFilemover(settings, CLI.Dryrun)
		if err := corpus.Everyfile(settings, tidying); err != nil {
			log.Fatal(err)
		}
		if err := tidying.Summary(); err != nil {
			log.Fatal("report Summary: ", err)
		}
	case "report metadata":
		tidying, err := tidy.NewMetadataReporter(settings)
		if err != nil {
			log.Fatal("Can't make a MetadataReporter because:", err)
		}
		if err := corpus.Everyfile(settings, tidying); err != nil {
			log.Fatal(err)
		}
		if err := tidying.Summary(); err != nil {
			log.Fatal("report Summary: ", err)
		}
	case "report tags":
		tidying, err := tidy.NewTagsReporter(settings)
		if err != nil {
			log.Fatal("Can't make a TagsReporter because:", err)
		}
		if err := corpus.Everyfile(settings, tidying); err != nil {
			log.Fatal(err)
		}
		if err := tidying.Summary(); err != nil {
			log.Fatal("report Summary: ", err)
		}
	case "report urls":
		tidying, err := tidy.NewUrlReporter(settings)
		if err != nil {
			log.Fatal("Can't make a NewUrlReporter because:", err)
		}
		if err := corpus.Everyfile(settings, tidying); err != nil {
			log.Fatal(err)
		}
		if err := tidying.Summary(); err != nil {
			log.Fatal("report Summary: ", err)
		}
	case "report todos":
		log.Println("report todos not implemented")
	case "report articles":
		tid := corpus.NewListAllWikiFilesTidying()
		if err := corpus.Everyfile(settings, tid); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Missing command: ", ctx.Command())
	}

}