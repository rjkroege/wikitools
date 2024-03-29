package article

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/rjkroege/wikitools/wiki"
)

var metadataMatcher = regexp.MustCompile("^([-A-Za-z]*):[ \t]*(.*)$")
var commentDataMatcher = regexp.MustCompile("<!-- *([0-9]*) *-->")

func trim(line string) string {
	if len(line) > 0 {
		return line[0 : len(line)-1]
	}
	return line
}

func (md *MetaData) rootThroughFileForMetadataImpl(rd *bufio.Reader) error {
	lc := 0
	md.mdtype = MdInvalid

	var date time.Time
	var de error

	for lc < 5 || md.mdtype != MdInvalid {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("rootThroughFileForMetadataImpl can't ReadString: %v", err)
		}
		if err != nil || md.mdtype != MdInvalid && line == "\n" {
			md.DateFromMetadata = date
			return nil
		}
		line = trim(line)

		if lc == 0 && line == "---" {
			// We're one of the modern metadata formats MdIaWriter, MdModern
			// MdModern is reserved for the situation where the title and tags have been
			// modernized.
			md.mdtype = MdIaWriter
		} else if lc == 0 {
			// We don't know yet what kind of metadata is present. But assume that
			// the first line is the title if we don't have metadata.
			md.Title = line
		}

		// fmt.Print("running regexp matcher...\n")
		m1 := metadataMatcher.FindStringSubmatch(line)
		m2 := commentDataMatcher.FindStringSubmatch(line)
		if len(m1) > 0 {
			s := strings.ToLower(m1[1])
			if s == "title" {
				md.Title = m1[2]
			} else if s == "date" {
				date, de = wiki.ParseDateUnix(strings.TrimSpace(m1[2]))
			} else if s == "tags" {
				for _, u := range strings.Split(strings.TrimSpace(m1[2]), " ") {
					if u != "" && len(u) > 1 && (u[0] == '#' || u[0] == '@') {
						md.Tags = append(md.Tags, u[1:])
					}
				}
			} else {
				md.extraKeys[s] = strings.TrimSpace(m1[2])
			}
			// If we have some combination of structured data but not MdIaWriter
			// metadata, then, we're MdLegacy
			if md.mdtype == MdInvalid {
				md.mdtype = MdLegacy
			}
		} else if len(m2) > 0 {
			// fmt.Print("matched for  <" + m2[1] + ">\n");
			date, de = wiki.ParseDateUnix(m2[1])
		}

		// I have no test that actually enforces that this is valid.
		// push to a helper
		if de != nil || date.IsZero() {
			//fmt.Print("date is zero, trying whole resultLine: <" + resultLine + ">\n");
			date, de = wiki.ParseDateUnix(md.Title)
		}
		lc++
	}
	md.DateFromMetadata = date
	return nil
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
