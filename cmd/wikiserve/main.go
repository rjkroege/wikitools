package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/corpus/tidy"
	"github.com/rjkroege/wikitools/wiki"
	"github.com/rjkroege/wikitools/actions"
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
func everyfile() any  { return map[string]string{"corpus": "Everyfile"} }
func summary() any    { return map[string]string{"tidying": "Summary"} }
func tagsReport() any { return map[string]string{"tidy": "NewTagsReporter"} }

// listAllWikiFilesTidying is a constructor, so we expose it too
func listAllWikiFilesTidying() any { return map[string]string{"corpus": "NewListAllWikiFilesTidying"} }

// ---------- routing table -----------------------------------------------------

type TidyingPassFactory func(*wiki.Settings, *http.Request) (corpus.Tidying, error)

type route struct {
	pattern string
	handler TidyingPassFactory
}

var routes = []route{
	// cmd namespace
	// 	{"/cmd/wikinew", wikinew},
	// 	{"/cmd/wikinew-autocomplete", wikinewAutocomplete},
	// 	{"/cmd/preview", preview},
	// 	{"/cmd/plumber-helper", plumberHelper},
	// 	{"/cmd/bearimport", bearimport},
	//
	// 	// tidy namespace
	// 	{"/tidy/new-metadata-updater", newMetadataUpdater},
	// 	{"/tidy/new-tags-dumper", newTagsDumper},
	// 	{"/tidy/new-backlinkwriter", newBacklinkwriter},
	// 	{"/tidy/new-filemover", newFilemover},
	// 	{"/tidy/new-metadata-reporter", newMetadataReporter},
	// 	{"/tidy/new-url-reporter", newUrlReporter},
	//
	// 	// corpus helpers used by tidy
	// 	{"/corpus/everyfile", everyfile},
	// 	{"/tidying/summary", summary},

	{"/corpus/list", tidy.NewListAllWikiFilesTidying},
	{"/corpus/list/{tags}", tidy.NewListAllWikiFilesTidying},
	{"/corpus/tags", tidy.NewTagsReporter},
	{"/corpus/urls", tidy.NewUrlReporter},
	{"/corpus/meta", tidy.NewMetadataReporter},
	{"/tidy/backlinks", tidy.NewBacklinkwriter},
	{"/tidy/tags", tidy.NewTagsDumper},
	{"/tidy/meta", tidy.NewMetadataUpdater},
	{"/tidy/move", tidy.NewFilemover},
}

func wrap(settings *wiki.Settings, f func() any) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("settings %v", settings)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(f()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func figureoutoutputformat(r *http.Request) int {
	// Set by parameter.
	queryParams := r.URL.Query()
	if _f, ok := queryParams["_f"]; ok {
		switch _f[0] {
		case "json":
			return wiki.OutputJSON
		case "html":
			return wiki.OutputHTML
		default:
			return wiki.OutputCLI
		}
	}

	// Don't know yet. Guess. If it's a browser, return HTML
	if _h, ok := r.Header["User-Agent"]; ok {
		ua := _h[0]
		log.Printf("ua: %q", ua)
		switch {
		case strings.HasPrefix(ua, "curl"):
			return wiki.OutputCLI
		default:
			return wiki.OutputHTML
		}
	}
	return wiki.OutputCLI
}


func figureoutdryrun(r *http.Request) bool {
	return r.URL.Query().Has("_dry")
}

// tidywrap returns an http.HandlerFunc corresponding to the specified tidying
// structured pass.
func tidywrap(settings *wiki.Settings, f TidyingPassFactory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Every request gets a private settings.
		reqsettings := *settings
		reqsettings.OutputType = figureoutoutputformat(r)
		reqsettings.Dryrun = figureoutdryrun(r)
		log.Printf("settings %v", reqsettings)
		

		tidying, err := f(&reqsettings, r)
		if err != nil {
			log.Printf("Can't make a tidying object for this request because: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := corpus.Everyfile(settings, tidying); err != nil {
			log.Printf("Can't Everyfile for this request because: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch reqsettings.OutputType {
		case wiki.OutputCLI:
			w.Header().Set("Content-Type", "text/plain")
			log.Printf("tidywrap running cli output OutputCLI")
			if err := tidying.SummaryWrite(w); err != nil {
				log.Printf("tidywrap %s %v", "OutputCLI", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case wiki.OutputJSON:
			log.Printf("tidywrap running cli output OutputJSON")
			w.Header().Set("Content-Type", "application/json")
			encoder := json.NewEncoder(w)
			if err := tidying.SummaryEncode(encoder); err != nil {
				log.Printf("tidywrap %s %v", "OutputJSON", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case wiki.OutputHTML:
			log.Printf("tidywrap running cli output OutputHTML")
			w.Header().Set("Content-Type", "text/html")
			if err := tidying.SummaryWrite(w); err != nil {
				log.Printf("tidywrap %s %v", "OutputHTML", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

// TODO(rjk): Move this into a separate file.
const homepage = `<html>
<head>
<title>Wiki</title>
<style>
/* Pretty button-style link */
a.pretty-button {
  display: inline-block;
  margin: 4px 2px;
  padding: 10px 22px;
  font: 600 14px/1.2 system-ui, sans-serif;
  color: #fff;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 6px;
  text-decoration: none;
  box-shadow: 0 4px 12px rgba(102, 126, 234, .35);
  transition: all .25s ease;
}

/* Hover highlight */
a.pretty-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, .55);
  background: linear-gradient(135deg, #7c8ff0 0%, #865fca 100%);
}
</style>
</head>
<body>
<h1>Wiki</h1>
<h2>Reports</h2>
<ul>
<li><a href="/corpus/list">List articles</a></li>
<li><a href="/corpus/tags">List tags</a></li>
<li><a href="/corpus/urls">List urls</a></li>
<li><a href="/corpus/meta">Metadata report</a></li>
</ul>
<h2>Tidying Passes</h2>
<em>Warning: these operations do not yet update the indexes</em>
<ul>
<li><a class="pretty-button" href="/tidy/backlinks?_dry=1">Preview backlinks update</a> <a  class="pretty-button" href="/tidy/backlinks">Do it!</a></li>
<li><a class="pretty-button" href="/tidy/tags?_dry=1">Preview tag writing</a> <a  class="pretty-button" href="/tidy/tags">Do it!</a></li>
<li><a class="pretty-button" href="/tidy/meta?_dry=1">Preview meta writing</a> <a  class="pretty-button" href="/tidy/meta">Do it!</a></li>
<li><a class="pretty-button" href="/tidy/move?_dry=1">Preview file relocations</a> <a  class="pretty-button" href="/tidy/move">Do it!</a></li>
</ul>
<h2>Article of the Day</h2>
not yet implemented
</body>
</html>
`

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

	actions.WatchAcmeLog(settings.Wikidir)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, homepage)
	})

	// TODO(rjk): In the future
	// register REST endpoints
	for _, rt := range routes {
		http.HandleFunc(rt.pattern, tidywrap(settings, rt.handler))
	}

	log.Printf("Listening on %s …", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
