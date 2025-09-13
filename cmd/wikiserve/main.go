package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/rjkroege/wikitools/wiki"
)

// ---------- tiny helpers that pretend to do the real work --------------------

// cmd helpers
func wikinew() any             { return map[string]string{"cmd": "Wikinew"} }
func wikinewAutocomplete() any { return map[string]string{"cmd": "WikinewAutocomplete"} }
func preview() any             { return map[string]string{"cmd": "Preview"} }
func plumberHelper() any       { return map[string]string{"cmd": "PlumberHelper"} }
func bearimport() any          { return map[string]string{"cmd": "Bearimport"} }

// tidy helpers
func newMetadataUpdater() any  { return map[string]string{"tidy": "NewMetadataUpdater"} }
func newTagsDumper() any       { return map[string]string{"tidy": "NewTagsDumper"} }
func newBacklinkwriter() any   { return map[string]string{"tidy": "NewBacklinkwriter"} }
func newFilemover() any        { return map[string]string{"tidy": "NewFilemover"} }
func newMetadataReporter() any { return map[string]string{"tidy": "NewMetadataReporter"} }
func newUrlReporter() any      { return map[string]string{"tidy": "NewUrlReporter"} }

// corpus helpers (invoked by tidy actions)
func everyfile() any { return map[string]string{"corpus": "Everyfile"} }
func summary() any   { return map[string]string{"tidying": "Summary"} }
func tagsReport() any { return map[string]string{"tidy": "NewTagsReporter"} }

// listAllWikiFilesTidying is a constructor, so we expose it too
func listAllWikiFilesTidying() any { return map[string]string{"corpus": "NewListAllWikiFilesTidying"} }


// ---------- routing table -----------------------------------------------------

type route struct {
	pattern string
	handler func() any
}

var routes = []route{
	// cmd namespace
	{"/cmd/wikinew", wikinew},
	{"/cmd/wikinew-autocomplete", wikinewAutocomplete},
	{"/cmd/preview", preview},
	{"/cmd/plumber-helper", plumberHelper},
	{"/cmd/bearimport", bearimport},

	// tidy namespace
	{"/tidy/new-metadata-updater", newMetadataUpdater},
	{"/tidy/new-tags-dumper", newTagsDumper},
	{"/tidy/new-backlinkwriter", newBacklinkwriter},
	{"/tidy/new-filemover", newFilemover},
	{"/tidy/new-metadata-reporter", newMetadataReporter},
	{"/tidy/new-tags-reporter", tagsReport},
	{"/tidy/new-url-reporter", newUrlReporter},

	// corpus helpers used by tidy
	{"/corpus/everyfile", everyfile},
	{"/tidying/summary", summary},
	{"/corpus/new-list-all-wiki-files-tidying", listAllWikiFilesTidying},
}

// ---------- generic JSON responder ------------------------------------------

func wrap(settings *wiki.Settings,  f func() any) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("settings %v", settings)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(f()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// ---------- entry point -------------------------------------------------------

func main() {
	port := flag.String("port", "8080", "HTTP server port")
	configfile := flag.String("path", "~/.wikinewrc", "Set alternate configuration file")
	flag.Parse()
	addr := ":" + *port

	// TODO(rjk): wiki => config
	settings, err := wiki.Read(*configfile)
	if err != nil {
		// TODO(rjk): This is not nice. Set things up sensibly and
		// proceed with reasonable defaults.
		log.Fatal("No configuration file. Fatal:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintln(w, "<h1>Wiki</h1>")
	})

	// register REST endpoints
	for _, rt := range routes {
		http.HandleFunc(rt.pattern, wrap(settings, rt.handler))
	}

	log.Printf("Listening on %s …", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
