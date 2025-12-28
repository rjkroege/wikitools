package links

import (
	"iter"
	"maps"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)

type ownedWork struct {
	ret      chan<- struct{}
	op 	func(l *Links)
}

type ConcurrentLinks struct {
	taskCh	chan ownedWork
	stopCh               chan struct{}
}

func NewConcurrentLinks(mapper wiki.LinkToFile, location string) *ConcurrentLinks {
	cl := &ConcurrentLinks{
		taskCh:	make(chan ownedWork),
		stopCh:               make(chan struct{}),
	}
	go cl.owner(mapper, location)
	return cl
}
// Must clone this in wikiserve/main

func (cl *ConcurrentLinks) owner(mapper wiki.LinkToFile, location string) {
	links := implMakeLinks(mapper, location)

	for {
		select {
		case task := <-cl.taskCh:
			task.op(links)
			close(task.ret)
		case <-cl.stopCh:
			return
		}
	}
}

func (cl *ConcurrentLinks) AppendStringVectorForwardLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	ret := make(chan struct{})
	cl.taskCh <- ownedWork{
		ret: ret,
		op:        func(l *Links)  { 
			l.AppendStringVectorForwardLinks(f, articles)
		},

	}		
	<-ret
}

func (cl *ConcurrentLinks) AppendStringVectorBackLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	ret := make(chan struct{})
	cl.taskCh <- ownedWork{
		ret: ret,
		op:        func(l *Links)  { 
			l.AppendStringVectorBackLinks(f, articles)
		},

	}		
	<-ret
}

func (cl *ConcurrentLinks) AppendStringVectorOutUrls(f func(corpus.Urllink) string, articles map[string][]string) {
	ret := make(chan struct{})
	cl.taskCh <- ownedWork{
		ret: ret,
		op:        func(l *Links)  { 
			l.AppendStringVectorOutUrls(f,  articles)
		},

	}		
	<-ret
}

func (cl *ConcurrentLinks) AppendStringVectorDamagedLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	ret := make(chan struct{})
	cl.taskCh <- ownedWork{
		ret: ret,
		op:        func(l *Links)  { 
			l.AppendStringVectorDamagedLinks(f,  articles)
		},

	}		
	<-ret
}

func (cl *ConcurrentLinks) BackLinksIterator() iter.Seq[corpus.LinkTuple] {
	retarg := make(chan map[string]corpus.LinkMap[corpus.Wikilink])
	defer close(retarg)
	cl.taskCh <- ownedWork{
		op:        func(l *Links)  {
			backlinks := maps.Clone(l.BackLinks)
			for k, v := range backlinks {
				backlinks[k] = maps.Clone(v)
			}
			retarg <- backlinks
		},
	}
	return backLinksIteratorImpl(<-retarg)
}

func (cl *ConcurrentLinks) Stop() {
	cl.stopCh <- struct{}{}
}

func (cl *ConcurrentLinks) Commit(lri corpus.LinkRecording) {
	ret := make(chan struct{})
	cl.taskCh <- ownedWork{
		ret: ret,
		op:        func(l *Links)  { 
			l.Commit(lri)
		},

	}		
	<-ret
}
