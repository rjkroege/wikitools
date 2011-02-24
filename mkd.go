/*
	A wrapper for the discount library for munging markup.
*/

package mkd

/*
	Scheme: let's get this to compile. Then, let's add the
	rest of the features.
*/


// #include <stdio.h>
// #include <mkdio.h>
import "C"

import (
		"unsafe";
		"io";
)

// TODO(rjkroege): I don't think this is necessary because
// it comes out of the include file.
// const (
// 	NOLINKS	uint32 = 0x00000001	/* don't do link processing, block <a> tags  */
// 	NOIMAGE	uint32 = 0x00000002	/* don't do image processing, block <img> */
// 	NOPANTS	uint32 = 0x00000004	/* don't run smartypants() */
// 	NOHTML	uint32 = 0x00000008	/* don't allow raw html through AT ALL */
// 	STRICT	uint32 = 0x00000010	/* disable SUPERSCRIPT, RELAXED_EMPHASIS */
// 	TAGTEXT	uint32 = 0x00000020	/* process text inside an html tag; no
// 	 <em>, no <bold>, no html or [] expansion */
// 	NO_EXT	uint32 = 0x00000040	/* don't allow pseudo-protocols */
// 	CDATA		uint32 = 0x00000080	/* generate code for xml ![CDATA[...]] */
// 	NOSUPERSCRIPT uint32 = 0x00000100	/* no A^B */
// 	NORELAXED			uint32 = 0x00000200	/* emphasis happens /everywhere/ */
// 	NOTABLES			uint32 = 0x00000400	/* disallow tables */
// 	NOSTRIKETHROUGH uint32 = 0x00000800	/* forbid ~~strikethrough~~ */
// 	TOC				    uint32 = 0x00001000	/* do table-of-contents processing */
// 	1_COMPAT	    uint32 = 0x00002000	/* compatability with MarkdownTest_1.0 */
// 	AUTOLINK	    uint32 = 0x00004000	/* make http://foo.com link even without <>s */
// 	SAFELINK	    uint32 = 0x00008000	/* paranoid check for link protocol */
// 	NOHEADER	    uint32 = 0x00010000	/* don't process header blocks */
// 	TABSTOP		    uint32 = 0x00020000	/* expand tabs to 4 spaces */
// 	NODIVQUOTE	  uint32 = 0x00040000	/* forbid >%class% blocks */
// 	NOALPHALIST	  uint32 = 0x00080000	/* forbid alphabetic lists */
// 	NODLIST	      uint32 = 0x00100000	/* forbid definition lists */
// 	EMBED	uint32 = MKD_NOLINKS|MKD_NOIMAGE|MKD_TAGTEXT
// )


// TODOD(rjkroege): figure out the way of this...
// Can I do this?  Do I want to?
type Flags uint32;

type Markdown struct {
	// Level of indirection may not be necessary.
	m *C.MMIOT;
}

// Reads a markdown input file and returns a Markdown containing the preprocessed document. (which
// is then fed to markdown() for final formatting.)
//func In(r *io.Reader, Flags uint32) *Markdown {
//	m := make(Markdown)
//  would need to get a filedescriptor for the reader... might not always work...
//	m.m = C.mkd_in(x, flags)
//	return m;
//}

// Processes the given string into a markdown blob suitable for feeding to markdown()
func String(s string, flags uint32) *Markdown {
 	m := make(Markdown)
	p := C.CString(s)
 	// name := C.CString(s);
 	// defer C.free(unsafe.Pointer(name))

 	m.m = C.mkd_string(p, C.uint(flags))
	C.free(unsafe.Pointer(p))
 	return m;
}


func (mkd *Markdown) Basename(f string) {
 	fmt.Println("Hi there" + f);
}

