---
title: PromptBacklinksRewrite
date: Wed 24 Dec 2025, 04:44:46 MST
tags: #agentic #planning #software #wikitools
---

You are in the root of a golang project. Below is a markdown list of
tasks. Do each task not marked as done one at a time. Once `go test
./...` passes for each task, git commit the change with a description
of the change to the current branch. Stop in the task list if you
cannot get the tests to pass. Continue through the tasks until
are completed. Update this file (prompt.md) to record the completion
of each task.


- [x] create two functions from _urlReportGen: _urlReportGenHtml and
_urlReportGenMarkdown. Make the _urlReportGenHtml only do the html ==
true case while _urlReportGenMarkdown does the Markdown case. Update
the callsites to _urlReportGen appropriately to use the new pair of
functions. Remove _urlReportGen.
- [x] make a new type struct LinkTuple containing a string and corpus.LinkMap[corpus.Wikilink]. put
this defn in links.go
- [x] create a new member function  BackLinksIterator in link_helper.go on Links that returns a iter.Seq[T] where T can be a LinkTuple
and the successive elements are LinkTuple objects corresponding to the k,v contents of the
BackLinks member in the Links argument. be sure to write a test for BackLinksIterator and place
this test in link_helpers_test.go
- [x] add BackLinksIterator to the LinksRecorder interface
- [x] rewrite linkUpdate by iterating over the iter.Seq[T] returned by BackLinksIterator
- [x] remove GetBackLinks from the code
- [ ] create a new file concurrent_links.go that contains a complete channel based proxy
implementation of LinkRecording with a Links instance in an owning go routine and a ConcurrentLinks struct
that implements LinkRecording. ConcurrentLinks's BackLinksIterator should use the channel iterator wrapping
pattern.
- [ ] create the ConcurrentLinks implementation and go routine in wikiserve/main.go. create the Links
implementation in wikitools/main.go. Place the LinkRecorder instance in Settings. Use the LinkRecorder
from Settings across all of the non-test code.
