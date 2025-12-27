package links

import (
	"iter"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)

// commitCmd is the data sent on the commit channel.
type commitCmd struct {
	filepath string
	urls     []corpus.Urllink
	wikis    []corpus.Wikilink
	ret      chan<- struct{}
}

// appendStringVectorCmd is the data sent on the appendStringVector channels.
type appendStringVectorCmd struct {
	f        func(any) string
	articles map[string][]string
	ret      chan<- struct{}
}

// backLinksIteratorCmd is the data sent on the backLinksIterator channel.
type backLinksIteratorCmd struct {
	ret chan<- iter.Seq[corpus.LinkTuple]
}

type ConcurrentLinks struct {
	commitCh             chan<- commitCmd
	appendFwdLinksCh     chan<- appendStringVectorCmd
	appendBackLinksCh    chan<- appendStringVectorCmd
	appendOutUrlsCh      chan<- appendStringVectorCmd
	appendDamagedLinksCh chan<- appendStringVectorCmd
	backLinksIteratorCh  chan<- backLinksIteratorCmd
	stopCh               chan<- struct{}
}

var _ corpus.LinksRecorder = (*ConcurrentLinks)(nil)

func NewConcurrentLinks(mapper wiki.LinkToFile, location string) *ConcurrentLinks {
	commitCh := make(chan commitCmd)
	appendFwdLinksCh := make(chan appendStringVectorCmd)
	appendBackLinksCh := make(chan appendStringVectorCmd)
	appendOutUrlsCh := make(chan appendStringVectorCmd)
	appendDamagedLinksCh := make(chan appendStringVectorCmd)
	backLinksIteratorCh := make(chan backLinksIteratorCmd)
	stopCh := make(chan struct{})

	cl := &ConcurrentLinks{
		commitCh:             commitCh,
		appendFwdLinksCh:     appendFwdLinksCh,
		appendBackLinksCh:    appendBackLinksCh,
		appendOutUrlsCh:      appendOutUrlsCh,
		appendDamagedLinksCh: appendDamagedLinksCh,
		backLinksIteratorCh:  backLinksIteratorCh,
		stopCh:               stopCh,
	}

	go cl.owner(mapper, location,
		commitCh,
		appendFwdLinksCh,
		appendBackLinksCh,
		appendOutUrlsCh,
		appendDamagedLinksCh,
		backLinksIteratorCh,
		stopCh)

	return cl
}

type concurrentLinkRecording struct {
	filepath string
	urls     []corpus.Urllink
	wikis    []corpus.Wikilink
	commitCh chan<- commitCmd
}

var _ corpus.LinkRecording = (*concurrentLinkRecording)(nil)

func (clr *concurrentLinkRecording) RecordUrl(displaytext, url string) {
	clr.urls = append(clr.urls, corpus.MakeUrllink(url, displaytext))
}

func (clr *concurrentLinkRecording) RecordWikilink(displaytext, wikitext string) {
	clr.wikis = append(clr.wikis, corpus.MakeWikilink(wikitext, displaytext))
}

func (clr *concurrentLinkRecording) Commit() {
	ret := make(chan struct{})
	clr.commitCh <- commitCmd{
		filepath: clr.filepath,
		urls:     clr.urls,
		wikis:    clr.wikis,
		ret:      ret,
	}
	<-ret
}

func (cl *ConcurrentLinks) StartRecordingForFile(filepath string) corpus.LinkRecording {
	return &concurrentLinkRecording{
		filepath: filepath,
		commitCh: cl.commitCh,
	}
}

func (cl *ConcurrentLinks) owner(
	mapper wiki.LinkToFile,
	location string,
	commitCh <-chan commitCmd,
	appendFwdLinksCh <-chan appendStringVectorCmd,
	appendBackLinksCh <-chan appendStringVectorCmd,
	appendOutUrlsCh <-chan appendStringVectorCmd,
	appendDamagedLinksCh <-chan appendStringVectorCmd,
	backLinksIteratorCh <-chan backLinksIteratorCmd,
	stopCh <-chan struct{},
) {
	links := implMakeLinks(mapper, location)

	for {
		select {
		case cmd := <-commitCh:
			// This is a little weird because I'm mixing two kinds of objects.
			lr := &linkRecording{
				filepath: cmd.filepath,
				links:    links,
				urls:     cmd.urls,
				wikis:    cmd.wikis,
			}
			lr.commitOnLinksOwner()
			close(cmd.ret)

		case cmd := <-appendFwdLinksCh:
			links.AppendStringVectorForwardLinks(func(l corpus.Wikilink) string {
				return cmd.f(l)
			}, cmd.articles)
			close(cmd.ret)

		case cmd := <-appendBackLinksCh:
			links.AppendStringVectorBackLinks(func(l corpus.Wikilink) string {
				return cmd.f(l)
			}, cmd.articles)
			close(cmd.ret)

		case cmd := <-appendOutUrlsCh:
			links.AppendStringVectorOutUrls(func(l corpus.Urllink) string {
				return cmd.f(l)
			}, cmd.articles)
			close(cmd.ret)

		case cmd := <-appendDamagedLinksCh:
			links.AppendStringVectorDamagedLinks(func(l corpus.Wikilink) string {
				return cmd.f(l)
			}, cmd.articles)
			close(cmd.ret)

		case cmd := <-backLinksIteratorCh:
			cmd.ret <- links.BackLinksIterator()

		case <-stopCh:
			return
		}
	}
}

func (cl *ConcurrentLinks) AppendStringVectorForwardLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	ret := make(chan struct{})
	cl.appendFwdLinksCh <- appendStringVectorCmd{
		f:        func(l any) string { return f(l.(corpus.Wikilink)) },
		articles: articles,
		ret:      ret,
	}
	<-ret
}

func (cl *ConcurrentLinks) AppendStringVectorBackLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	ret := make(chan struct{})
	cl.appendBackLinksCh <- appendStringVectorCmd{
		f:        func(l any) string { return f(l.(corpus.Wikilink)) },
		articles: articles,
		ret:      ret,
	}
	<-ret
}

func (cl *ConcurrentLinks) AppendStringVectorOutUrls(f func(corpus.Urllink) string, articles map[string][]string) {
	ret := make(chan struct{})
	cl.appendOutUrlsCh <- appendStringVectorCmd{
		f:        func(l any) string { return f(l.(corpus.Urllink)) },
		articles: articles,
		ret:      ret,
	}
	<-ret
}

func (cl *ConcurrentLinks) AppendStringVectorDamagedLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	ret := make(chan struct{})
	cl.appendDamagedLinksCh <- appendStringVectorCmd{
		f:        func(l any) string { return f(l.(corpus.Wikilink)) },
		articles: articles,
		ret:      ret,
	}
	<-ret
}

func (cl *ConcurrentLinks) BackLinksIterator() iter.Seq[corpus.LinkTuple] {
	ret := make(chan iter.Seq[corpus.LinkTuple])
	cl.backLinksIteratorCh <- backLinksIteratorCmd{
		ret: ret,
	}
	return <-ret
}

func (cl *ConcurrentLinks) Stop() {
	cl.stopCh <- struct{}{}
}
