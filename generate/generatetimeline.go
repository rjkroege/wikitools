/*
  Creates the contents of the JavaScript file that lists the entries in the
  timeline.
*/

package generate

import (
	"encoding/json"
	"github.com/rjkroege/wikitools/article"
	"io"
)

func WriteTimeline(fd io.Writer, e []*article.MetaData) {
	io.WriteString(fd, timeline_header)

	enc := json.NewEncoder(fd)
	err := enc.Encode(e)
	if err != nil {
		panic(err)
	}

	io.WriteString(fd, timeline_footer)
}
