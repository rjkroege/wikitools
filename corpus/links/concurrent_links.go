package links

import (
	"iter"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)

type command int

const (
	cmdStartRecording command = iota
	cmdAppendStringVectorForwardLinks
	cmdAppendStringVectorBackLinks
	cmdAppendStringVectorOutUrls
	cmdAppendStringVectorDamagedLinks
	cmdBackLinksIterator
	cmdStop
)

type ConcurrentLinks struct {
	cmdCh chan<- commandData
}

type commandData struct {
	cmd  command
	data any
	ret  chan<- any
}

var _ corpus.LinksRecorder = (*ConcurrentLinks)(nil)

func NewConcurrentLinks(mapper wiki.LinkToFile, location string) *ConcurrentLinks {
	cmdCh := make(chan commandData)

	cl := &ConcurrentLinks{
		cmdCh: cmdCh,
	}

	go cl.owner(mapper, location, cmdCh)

	return cl
}

type concurrentLinkRecording struct {
	filepath string
	urls     []corpus.Urllink
	wikis    []corpus.Wikilink
	cmdCh    chan<- commandData
}

var _ corpus.LinkRecording = (*concurrentLinkRecording)(nil)

func (clr *concurrentLinkRecording) RecordUrl(displaytext, url string) {
	clr.urls = append(clr.urls, corpus.MakeUrllink(url, displaytext))
}

func (clr *concurrentLinkRecording) RecordWikilink(displaytext, wikitext string) {
	clr.wikis = append(clr.wikis, corpus.MakeWikilink(wikitext, displaytext))
}

func (clr *concurrentLinkRecording) Commit() {
	ret := make(chan any)
	clr.cmdCh <- commandData{
		cmd: cmdStartRecording,
		data: &linkRecording{
			filepath: clr.filepath,
			urls:     clr.urls,
			wikis:    clr.wikis,
		},
		ret: ret,
	}
	<-ret
}

func (cl *ConcurrentLinks) StartRecordingForFile(filepath string) corpus.LinkRecording {
	return &concurrentLinkRecording{
		filepath: filepath,
		cmdCh:    cl.cmdCh,
	}
}

func (cl *ConcurrentLinks) owner(mapper wiki.LinkToFile, location string, cmdCh <-chan commandData) {
	links := implMakeLinks(mapper, location)

	for cmd := range cmdCh {
		switch cmd.cmd {
		case cmdStartRecording:
			data := cmd.data.(*linkRecording)
			data.links = links
			data.commitOnLinksOwner()
			close(cmd.ret)
		case cmdAppendStringVectorForwardLinks:
			data := cmd.data.(*appendStringVectorData)
			links.AppendStringVectorForwardLinks(func(l corpus.Wikilink) string {
				return data.f(l)
			}, data.articles)
			close(cmd.ret)
		case cmdAppendStringVectorBackLinks:
			data := cmd.data.(*appendStringVectorData)
			links.AppendStringVectorBackLinks(func(l corpus.Wikilink) string {
				return data.f(l)
			}, data.articles)
			close(cmd.ret)
		case cmdAppendStringVectorOutUrls:
			data := cmd.data.(*appendStringVectorData)
			links.AppendStringVectorOutUrls(func(l corpus.Urllink) string {
				return data.f(l)
			}, data.articles)
			close(cmd.ret)
		case cmdAppendStringVectorDamagedLinks:
			data := cmd.data.(*appendStringVectorData)
			links.AppendStringVectorDamagedLinks(func(l corpus.Wikilink) string {
				return data.f(l)
			}, data.articles)
			close(cmd.ret)
		case cmdBackLinksIterator:
			// implement
		case cmdStop:
			close(cmd.ret)
			return
		}
	}
}

type appendStringVectorData struct {
	f        func(any) string
	articles map[string][]string
}

func (cl *ConcurrentLinks) AppendStringVectorForwardLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	ret := make(chan any)
	cl.cmdCh <- commandData{
		cmd: cmdAppendStringVectorForwardLinks,
		data: &appendStringVectorData{
			f:        func(l any) string { return f(l.(corpus.Wikilink)) },
			articles: articles,
		},
		ret: ret,
	}
	<-ret
}

func (cl *ConcurrentLinks) AppendStringVectorBackLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	ret := make(chan any)
	cl.cmdCh <- commandData{
		cmd: cmdAppendStringVectorBackLinks,
		data: &appendStringVectorData{
			f:        func(l any) string { return f(l.(corpus.Wikilink)) },
			articles: articles,
		},
		ret: ret,
	}
	<-ret
}

func (cl *ConcurrentLinks) AppendStringVectorOutUrls(f func(corpus.Urllink) string, articles map[string][]string) {
	ret := make(chan any)
	cl.cmdCh <- commandData{
		cmd: cmdAppendStringVectorOutUrls,
		data: &appendStringVectorData{
			f:        func(l any) string { return f(l.(corpus.Urllink)) },
			articles: articles,
		},
		ret: ret,
	}
	<-ret
}

func (cl *ConcurrentLinks) AppendStringVectorDamagedLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	ret := make(chan any)
	cl.cmdCh <- commandData{
		cmd: cmdAppendStringVectorDamagedLinks,
		data: &appendStringVectorData{
			f:        func(l any) string { return f(l.(corpus.Wikilink)) },
			articles: articles,
		},
		ret: ret,
	}
	<-ret
}

func (cl *ConcurrentLinks) BackLinksIterator() iter.Seq[corpus.LinkTuple] {
	ret := make(chan any)
	cl.cmdCh <- commandData{
		cmd: cmdBackLinksIterator,
		ret: ret,
	}

	return func(yield func(corpus.LinkTuple) bool) {
		for t := range ret {
			if !yield(t.(corpus.LinkTuple)) {
				return
			}
		}
	}
}
