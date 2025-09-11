package main

import (
    "encoding/json"
	"flag"
    "fmt"
    "log"
    "net/http"
)

// ---------- tiny helpers that pretend to do the real work --------------------

// cmd helpers
func wikinew() any                  { return map[string]string{"cmd": "Wikinew"} }
func wikinewAutocomplete() any      { return map[string]string{"cmd": "WikinewAutocomplete"} }
func preview() any                  { return map[string]string{"cmd": "Preview"} }
func plumberHelper() any            { return map[string]string{"cmd": "PlumberHelper"} }
func bearimport() any                { return map[string]string{"cmd": "Bearimport"} }

// tidy helpers
func newMetadataUpdater() any       { return map[string]string{"tidy": "NewMetadataUpdater"} }
func newTagsDumper() any            { return map[string]string{"tidy": "NewTagsDumper"} }
func newBacklinkwriter() any        { return map[string]string{"tidy": "NewBacklinkwriter"} }
func newFilemover() any              { return map[string]string{"tidy": "NewFilemover"} }
func newMetadataReporter() any      { return map[string]string{"tidy": "NewMetadataReporter"} }
func newTagsReporter() any          { return map[string]string{"tidy": "NewTagsReporter"} }
func newUrlReporter() any           { return map[string]string{"tidy": "NewUrlReporter"} }

// corpus helpers (invoked by tidy actions)
func everyfile() any               { return map[string]string{"corpus": "Everyfile"} }
func summary() any                 { return map[string]string{"tidying": "Summary"} }

// listAllWikiFilesTidying is a constructor, so we expose it too
func listAllWikiFilesTidying() any { return map[string]string{"corpus": "NewListAllWikiFilesTidying"} }

// ---------- routing table -----------------------------------------------------

type route struct {
    pattern string
    handler http.HandlerFunc
}

var routes = []route{
    // cmd namespace
    {"/cmd/wikinew", wrap(wikinew)},
    {"/cmd/wikinew-autocomplete", wrap(wikinewAutocomplete)},
    {"/cmd/preview", wrap(preview)},
    {"/cmd/plumber-helper", wrap(plumberHelper)},
    {"/cmd/bearimport", wrap(bearimport)},

    // tidy namespace
    {"/tidy/new-metadata-updater", wrap(newMetadataUpdater)},
    {"/tidy/new-tags-dumper", wrap(newTagsDumper)},
    {"/tidy/new-backlinkwriter", wrap(newBacklinkwriter)},
    {"/tidy/new-filemover", wrap(newFilemover)},
    {"/tidy/new-metadata-reporter", wrap(newMetadataReporter)},
    {"/tidy/new-tags-reporter", wrap(newTagsReporter)},
    {"/tidy/new-url-reporter", wrap(newUrlReporter)},

    // corpus helpers used by tidy
    {"/corpus/everyfile", wrap(everyfile)},
    {"/tidying/summary", wrap(summary)},
    {"/corpus/new-list-all-wiki-files-tidying", wrap(listAllWikiFilesTidying)},
}

// ---------- generic JSON responder ------------------------------------------

func wrap(f func() any) http.HandlerFunc {
    return func(w http.ResponseWriter, _ *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(f()); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}

// ---------- entry point -------------------------------------------------------

func main() {
    port := flag.String("port", "8080", "HTTP server port")
    flag.Parse()
    addr := ":" + *port

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }
        fmt.Fprintln(w, "<h1>Hello, world</h1>")
    })

    // register REST endpoints
    for _, rt := range routes {
        http.HandleFunc(rt.pattern, rt.handler)
    }

    log.Printf("Listening on %s …", addr)
    log.Fatal(http.ListenAndServe(addr, nil))
}
