SHELL			= /bin/sh

prefix			= /usr/local
exec_prefix		= ${prefix}
bindir			= $(DESTDIR)/$(exec_prefix)/bin
mandir			= $(DESTDIR)/$(prefix)/man

.PHONY: clean nuke install all distclean
INSTALL			= /usr/bin/install -c

BINS = wikimake

all: $(BINS)

# TODO(rjkroege): this makefile assumes 64bit intel processors. Fix.
article.6: article.go metadata.go
	6g $^

# According to http://golang.org/doc/go_tutorial.html#tmp_186, 
# we need to be sure to compile buildnote.go first before 
# the compilation of listnotes.go can be successful.
wikimake :  listnotes.go generatemarkup.go article.6
	6g listnotes.go generatemarkup.go
	6l  -o $@ listnotes.6


# Add some unit testing...
# wikitest : wikitesting.go article.6 

# TODO(rjkroege): might want to add some additional tests.
test: wikimake
	rm -f note_list.js
	./wikimake
	diff -q note_list.js note_list.js.baseline

install			: all
	$(INSTALL) -d $(bindir)
	for BIN in $(BINS) ; do \
		$(INSTALL) $$BIN $(bindir)/$$BIN ; \
	done

clean:
	rm -rf *.6 core *.dSYM note_list.js
