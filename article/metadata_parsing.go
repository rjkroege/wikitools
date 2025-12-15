package article

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/rjkroege/wikitools/wiki"
)

var metadataMatcher = regexp.MustCompile("^([-A-Za-z]*):[ \t]*(.*)$")

func trim(line string) string {
	if len(line) > 0 {
		return line[0 : len(line)-1]
	}
	return line
}

func (md *MetaData) rootThroughFileForMetadataImpl(rd *bufio.Reader) error {
	lc := 0
	md.mdtype = MdInvalid

	keys := make(map[string]string)
	for lc < 5 || md.mdtype != MdInvalid {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("rootThroughFileForMetadataImpl can't ReadString: %v", err)
		}
		line = trim(line)

		switch {
		case lc == 0 && line == "---":
			// We're one of the modern metadata formats MdIaWriter, MdModern
			// MdModern is reserved for the situation where the title and tags have been
			// modernized.
			md.mdtype = MdUnterminatedIaWriterOrModern
		case lc == 0 && line != "---":
			// We don't know yet what kind of metadata is present. But assume that
			// the first line is the title if we don't have metadata.
			// We might replace this below.
			md.Title = line
			md.mdtype = MdUnterminatedLegacy
			handleLine(line, md, keys)
		case lc > 0 && md.mdtype == MdUnterminatedIaWriterOrModern && line == "---":
			md.mdtype = MdUnterminatedIaWriterOrModernBlank
		case lc > 0 && md.mdtype == MdUnterminatedLegacy && line == "":
			md.mdtype = MdLegacy
			processKeys(keys, md)
			return nil
		case lc > 0 && md.mdtype == MdUnterminatedIaWriterOrModernBlank && line == "":
			md.mdtype = MdModern
			processKeys(keys, md)
			return nil
		default:
			handleLine(line, md, keys)
		}
		lc++
	}

	md.mdtype = MdInvalid
	return nil
}

func handleLine(line string, md *MetaData, keys map[string]string) {
	m1 := metadataMatcher.FindStringSubmatch(line)
	if len(m1) > 0 {
		k := strings.ToLower(m1[1])
		v := strings.TrimSpace(m1[2])

		if _, ok := keys[k]; ok {
			// Having duplicate keys is an error.
			md.mdtype = MdInvalid
		}
		keys[k] = v
	} else {
		md.mdtype = MdInvalid
	}
}

func processKeys(kvpairs map[string]string, md *MetaData) {
	// Valid metadata must include a title and date.
	hastitle := false
	hasdate := false
	for k, v := range kvpairs {
		switch k {
		case "title":
			md.Title = v
			hastitle = true
			delete(kvpairs, k)
		case "date":
			delete(kvpairs, k)
			date, de := wiki.ParseDateUnix(strings.TrimSpace(v))
			if de == nil {
				hasdate = true
				md.DateFromMetadata = date
			}
		case "tags":
			// carve out the correct values here?
			delete(kvpairs, k)
			processTags(v, md)
		}
	}
	md.extraKeys = kvpairs
	if !hastitle || !hasdate {
		md.mdtype = MdInvalid
	}
}

func processTags(tagstring string, md *MetaData) {
	moderntag := false
	legacytag := false
	badtag := false

	tags := make([]string, 0)

	// A tag must start with # or @ and be at least 1 character long.
	for _, u := range strings.Fields(tagstring) {
		switch {
		case len(u) > 1 && u[0] == '#':
			tags = append(tags, u[1:])
			moderntag = true
		case len(u) > 1 && u[0] == '@':
			tags = append(tags, u[1:])
			legacytag = true
		default:
			badtag = true
		}
	}
	md.Tags = tags

	switch {
	case md.mdtype == MdModern &&
		moderntag == false &&
		legacytag == false &&
		badtag == false:
		md.mdtype = MdModern
	case md.mdtype == MdModern &&
		moderntag == true &&
		legacytag == false &&
		badtag == false:
		md.mdtype = MdModern
	case md.mdtype == MdModern &&
		legacytag == true &&
		badtag == false:
		md.mdtype = MdIaWriter
	case md.mdtype == MdModern &&
		badtag == true:
		md.mdtype = MdInvalid
	case md.mdtype == MdLegacy &&
		badtag == true:
		md.mdtype = MdInvalid
	}
}

// RootThroughFileForMetadata opens a specified file and attempts to
// extract metadata. There are two possibilities for metadata. Without
// either, dates fallback to the modification date of the file and the
// the first line as the fallback.
//
// 1. The date is in a metadata segment at the top of the file as
// defined for MetaMarkdown. This format consists of key: value with
// a following blank line.
//
// 2. The date is contained in a comment as a sequence of numbers.
// To keep this from being too inefficient, it must be found in the top
// 5 lines.
//
// 3. An iAWriter metadata block: metadata is some number of key: value within ---.
//
// TODO(rjk): Consider having some kind of error response? There
// could be I/O errors.
// TODO(rjk): Have it say what kind of metadata the file has
func (md *MetaData) RootThroughFileForMetadata(reader io.Reader) {
	rd := bufio.NewReader(reader)
	md.rootThroughFileForMetadataImpl(rd)
}
