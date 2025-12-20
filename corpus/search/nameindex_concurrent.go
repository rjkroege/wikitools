package search

// Largely AI generated.

import (
	"github.com/rjkroege/wikitools/wiki"
)

type 	pathReply struct {
		path string
		err  error
	}

type pathReq struct {
		location, lsd, wikitext string
		reply                     chan pathReply
	}

type allpathsReq struct {
		location, lsd, wikitext string
		reply                  chan allpathsReply
	}

type	allpathsReply struct {
		paths []string
		err   error
	}


type	wikitextReq struct {
		frompath, topath string
		reply            chan wikitextReply
	}

type wikitextReply struct {
		text string
		err  error
	}


// channelIndexer wraps a *wikilinkIndexerimpl that lives in its own
// goroutine.  Every exported method forwards the request over a channel
// and waits for the reply.
type channelIndexer struct {
	// ---- channels used to talk to the owning goroutine -----------------
	pathChan     chan pathReq
	allpathsChan chan allpathsReq
	wikitextChan chan wikitextReq

	// ---- lifecycle ------------------------------------------------------
	quit chan struct{} // closed to stop the goroutine
	done chan struct{} // closed when the goroutine has returned
}

// NewChannelIndexer builds a channelIndexer around a *wikilinkIndexerimpl
// and starts the goroutine that owns the wrapped implementation.
func NewChannelIndexer(wikiroot string) *channelIndexer {
	ci := &channelIndexer{
		pathChan:     make(chan pathReq),
		allpathsChan: make(chan allpathsReq),
		wikitextChan: make(chan wikitextReq),
		quit:         make(chan struct{}),
		done:         make(chan struct{}),
	}

	go ci.run(wikiroot)
	return ci
}

// run is the goroutine that owns the *wikilinkIndexerimpl.
func (ci *channelIndexer) run(wikiroot string) {
	defer close(ci.done)
	impl := implMakeWikilinkNameIndex(wikiroot)

	for {
		select {
		case <-ci.quit:
			return

		case req := <-ci.pathChan:
			p, err := impl.Path(req.location, req.lsd, req.wikitext)
			req.reply <- pathReply{path: p, err: err}

		case req := <-ci.allpathsChan:
			ps, err := impl.Allpaths(req.location, req.lsd, req.wikitext)
			req.reply <- allpathsReply{paths: ps, err: err}

		case req := <-ci.wikitextChan:
			t, err := impl.Wikitext(req.frompath, req.topath)
			req.reply <- wikitextReply{text: t, err: err}
		}
	}
}

// Close shuts the indexer down and waits for its goroutine to finish.
func (ci *channelIndexer) Close() {
	close(ci.quit)
	<-ci.done
}

// ---- LinkToFile implementation -----------------------------------------

var _ wiki.LinkToFile = (*channelIndexer)(nil)

func (ci *channelIndexer) Path(location, lsd, wikitext string) (string, error) {
	reply := make(chan pathReply, 1)
	ci.pathChan <- pathReq{location: location, lsd: lsd, wikitext: wikitext, reply: reply}
	r := <-reply
	return r.path, r.err
}

func (ci *channelIndexer) Allpaths(location, lsd, wikitext string) ([]string, error) {
	reply := make(chan allpathsReply, 1)
	ci.allpathsChan <- allpathsReq{location: location, lsd: lsd, wikitext: wikitext, reply: reply}
	r := <-reply
	return r.paths, r.err
}

func (ci *channelIndexer) Wikitext(frompath, topath string) (string, error) {
	reply := make(chan wikitextReply, 1)
	ci.wikitextChan <- wikitextReq{frompath: frompath, topath: topath, reply: reply}
	r := <-reply
	return r.text, r.err
}
