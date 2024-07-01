package article

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/rjkroege/wikitools/corpus"
	"golang.org/x/sys/unix"
)

// remember xattr to dump the extended attributes
const backlinkkey = "org.liqui.wikitoolsback"

func WriteBacklinks(fname string, backmap map[corpus.Wikilink]corpus.Empty) error {
	links := make([]corpus.Wikilink, 0, len(backmap))
	for k := range backmap {
		links = append(links, k)
	}
	buffy := new(bytes.Buffer)
	encoder := json.NewEncoder(buffy)
	if err := encoder.Encode(links); err != nil {
		return fmt.Errorf("Can't WriteBacklinks to %q because %w", fname, err)
	}

	// Write it to the file
	if err := unix.Setxattr(fname, backlinkkey, buffy.Bytes(), 0); err != nil {
		return fmt.Errorf("Can't WriteBacklinks to %q because %w", fname, err)
	}
	return nil
}

// It's arguable that a list is all that's necessary? No. I need the original map back
// to update the links.
func ReadBacklinks(fname string) (map[corpus.Wikilink]corpus.Empty, error) {
	by := make([]byte, 1<<16)
	sz, err := unix.Getxattr(fname, backlinkkey, by)
	if err != nil {
		return nil, fmt.Errorf("Can't ReadBacklinks to %q because %w", fname, err)
	}
	// I would be surprised if this didn't happen. But still.
	by = by[0:sz]
	buffy := bytes.NewBuffer(by)

	links := make([]corpus.Wikilink, 0)
	decoder := json.NewDecoder(buffy)
	if err := decoder.Decode(&links); err != nil {
		return nil, fmt.Errorf("Can't ReadBacklinks to %q because %w", fname, err)
	}

	backmap := make(map[corpus.Wikilink]corpus.Empty)
	for _, v := range links {
		backmap[v] = corpus.Empty{}
	}

	return backmap, nil
}
