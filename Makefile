SHELL			= /bin/sh

prefix			= /usr/local
exec_prefix		= ${prefix}
bindir			= $(DESTDIR)/$(exec_prefix)/bin
mandir			= $(DESTDIR)/$(prefix)/man

.PHONY: clean nuke install all distclean
INSTALL			= /usr/bin/install -c


SRCS			= listnotes.go generatemarkup.go metadata.go
BINS = entrylist

entrylist		:  $(SRCS)
	6g -I . $(SRCS)
	6l  -o entrylist listnotes.6


all: $(BINS)

install			: all
	$(INSTALL) -d $(bindir)
	for BIN in $(BINS) ; do \
		$(INSTALL) $$BIN $(bindir)/$$BIN ; \
	done


clean:
	rm -rf *.6 core *.dSYM
