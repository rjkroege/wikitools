package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"9fans.net/go/acme"
	"github.com/rjkroege/gozen"
	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/corpus/search"
	"github.com/rjkroege/wikitools/wiki"
)

func PlumberHelper(settings *wiki.Settings, lsd, wikitext string) {
	log.Println("PlumberHelper", wikitext)

	// Remember that on darwin that we are running in a secondary Go routine.
	// TODO(rjk): Consider renaming this later.
	mapper := search.MakeWikilinkNameIndex(settings.Wikidir)

	if filepath.Ext(wikitext) == "" {
		wikitext = wikitext + ".md"
	}

	fp, err := mapper.Path(settings.Wikidir, lsd, wikitext)
	if err != nil {
		log.Printf("indexer.Path errored on %q: %v", wikitext, err)
		woe := fmt.Sprintf("trying to follow [[%s]] failed: %v\n", wikitext, err)
		writewikierror(settings, []byte(woe))
		return
	}

	backlinks := makebacklinkstring(fp)

	log.Println(fp)
	gozen.Editinacme(fp, gozen.Addtotag(backlinks), gozen.Blinktag(""))

}

func makebacklinkstring(fp string) string {
	backlinks, err := article.ReadBacklinks(fp)
	if err != nil {
		log.Printf("can't read backlinks from %q, maybe none exist: %v", fp, err)
		return ""
	}
	// Add the backlinks to Edwood.
	var b strings.Builder
	for k := range backlinks {
		b.WriteString(k.Markdown())
		b.WriteRune(' ')
	}
	return b.String()
}

const logfile = "+Wikierror"

// writewikierror appends text to the special file for wiki errors.
func writewikierror(settings *wiki.Settings, text []byte) error {
	// shove stuff into Edwood/Acme error area
	fn := filepath.Join(settings.Wikidir, logfile)

	// Two choices: we already have the Window open.
	wins, err := acme.Windows()
	if err != nil {
		return fmt.Errorf("plumbhelper acme.Windows")
	}

	win := (*acme.Win)(nil)
	for _, wi := range wins {
		if wi.Name == fn {
			win, err = acme.Open(wi.ID, nil)
			if err != nil {
				return fmt.Errorf("plumbhelper acme.Open")
			}
			break
		}
	}

	if win == nil {
		var err error
		win, err = acme.New()
		if err != nil {
			return fmt.Errorf("writewikierror acme.New: %v", err)
		}

		if err := win.Ctl("nomark"); err != nil {
			return fmt.Errorf("writewikierror win.Ctl nomark: %v", err)
		}

		if err := win.Name(fn); err != nil {
			return fmt.Errorf("writewikierror win.Name: %v", err)
		}

		if err = win.Ctl("mark"); err != nil {
			return fmt.Errorf("writewikierror %q: %v", "mark", err)
		}

		if err = win.Ctl("clean"); err != nil {
			return fmt.Errorf("writewikierror %q: %v", "clean", err)
		}
	}

	// TODO(rjk): This code could be more ergonomic.
	if _, err := win.Write("body", text); err != nil {
		return fmt.Errorf("writewikierror %q: %v", "write", err)
	}

	return nil
}
