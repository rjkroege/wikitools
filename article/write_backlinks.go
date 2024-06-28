package article

import (
	"encoding/json"
	"bytes"
	"fmt"

	"golang.org/x/sys/unix"
	"github.com/rjkroege/wikitools/corpus"
)

// remember xattr to dump the extended attributes
const backlinkkey = "org.liqui.wikitoolsback"

func WriteBacklinks(fname string, backmap map[corpus.Wikilink]corpus.Empty) error {
	links := make([]corpus.Wikilink,0,  len(backmap))
	for k, _ := range backmap {
		links = append(links, k)
	}
	buffy := new(bytes.Buffer)
	encoder := json.NewEncoder(buffy)
	if err := encoder.Encode(links); err != nil {
		return fmt.Errorf("Can't WriteBacklinks to %q because %v", fname, err)
	}

	// Write it to the file
	if err := unix.Setxattr(fname, backlinkkey,  buffy.Bytes(), 0); err != nil {
		return fmt.Errorf("Can't WriteBacklinks to %q because %v", fname, err)
	}
	return nil
}

