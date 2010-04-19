SHELL			= /bin/sh

prefix			= /usr/local
exec_prefix		= ${prefix}
bindir			= $(DESTDIR)/$(exec_prefix)/bin
mandir			= $(DESTDIR)/$(prefix)/man

.PHONY: clean nuke install all distclean
INSTALL			= /usr/bin/install -c


SRCS = generatemarkup.go metadata.go
BINS = extractmeta entrylist

all: $(BINS)

entrylist		:  listnotes.go $(SRCS)
	6g -I . listnotes.go $(SRCS)
	6l  -o entrylist listnotes.6

extractmeta		:  extractone.go $(SRCS)
	6g -I . extractone.go $(SRCS)
	6l  -o extractmeta extractone.6

# TODO(rjkroege): might want to add some additional tests.
test: entrylist
	rm -f note_list.js
	./entrylist
	diff -q note_list.js note_list.js.baseline

install			: all
	$(INSTALL) -d $(bindir)
	for BIN in $(BINS) ; do \
		$(INSTALL) $$BIN $(bindir)/$$BIN ; \
	done

clean:
	rm -rf *.6 core *.dSYM note_list.js
