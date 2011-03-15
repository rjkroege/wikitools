# This Makefile is very primitive and could be replaced

include $(GOROOT)/src/Make.inc

TARG = wikimake

# Add the rest of the files here.
# GOFILES = article.go generatemarkup.go listnotes.go metadata.go
GOFILES = listnotes.go generatemarkup.go

# hypothesis: 2011/2/28
# extractone.go is unnecessary.
# assumption: let's make it this way.

# GC+=-I .

CLEANFILES+=mkdtest note_list.js

include $(GOROOT)/src/Make.cmd

# need to set some kind of include path
# $GOROOT/src/pkg/github.com/knieriem/markdown


# Is the test necessary?
# Yes. but doe we have to code like this?
# TODO(rjkroege): fix this test up
# mkdtest :  install testdriver.go
#	6g -I./_obj -I../sre2/_obj testdriver.go
#	6l -o $@ testdriver.6

# Runs the test too.
all: article wikimake targettest
	
article: article/article.go article/metadata.go
	make -C article install

targettest: wikimake
	rm -f note_list.js
	cd testdata ; ../wikimake
	diff -q testdata/note_list.js testdata/note_list.js.baseline

install			: all
	$(INSTALL) -d $(bindir)
	for BIN in $(BINS) ; do \
		$(INSTALL) $$BIN $(bindir)/$$BIN ; \
	done

clean:
	rm -rf *.6 core *.dSYM note_list.js
