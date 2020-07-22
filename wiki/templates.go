package wiki

// Templates for new journal entries.
const (
	basetmpl = `---
title: {{.Title}}
date: {{.DetailedDate}}
tags: {{.Tagstring}}
--

{{.Dynamicstring}}
`

	journalamtmpl = `# Today I am grateful for
1. 
2. 
3. 

# What would make today great?
1. 
2.  
3.  

# Daily affirmations. I am
1. 
2.  
3. 
`

	journalpmtmpl = `# Three Amazing Things Today
1. 
2. 
3. 

# How Could I Have Made Today Better
1. 
2.  
3.  

`
	// For notes and such.
	entrytmpl = `Yo dawg! Write stuff here.
`

	codetmpl = `# Summary of change
*what's happening here anyway?*
*can I divide this?*

# Dillgence
## Landscape
*what process are we in?*
*what thread*
*callstack*
*plumbing map*

## Testing
*how to tell if I'm done*
*how to tell that it's right*
`

	booktmpl = `---
title: {{.Title}}
date: {{.DetailedDate}}
tags: {{.Tagstring}}
bib-bibkey: The reference key for this book. Required.
bib-author: The name(s) of the author(s) (in the case of more than one author, separated by and)
bib-title: {{.Title}}
bib-publisher: The publisher's name
bib-year: The year of publication (or, if unpublished, the year of creation)
bib-address: (Optional) Publisher's address (usually just the city, but can be the full address for lesser-known publishers)
bib-number: (Optional) The "(issue) number" of a journal, magazine, or tech-report, if applicable. (Most publications have a "volume", but no "number" field.)
bib-volume: (Optional) The volume of a journal or multi-volume book
bib-series: (Optional) The series of books the book was published in (e.g. "The Hardy Boys" or "Lecture Notes in Computer Science")
bib-edition: (Optional) The edition of a book, long form (such as "first" or "second")
bib-month:  (Optional) The month the book was issued.
bib-isbn: (Optional, non-standard) The ISBN for this book.
bib-url: (Optional, non-standard) A URL for an e-book.
---

Yo dawg! Put the bookreview here. Delete the undesired tags. Remove the blank line.
`

	articletmpl = `---
title: {{.Title}}
date: {{.DetailedDate}}
tags: {{.Tagstring}}
bib-bibkey: The reference key for this book. Required.
bib-author: The name(s) of the author(s) (in the case of more than one author, separated by and)
bib-title: {{.Title}}
bib-journal: The journal or magazine the work was published in
bib-year: The year of publication (or, if unpublished, the year of creation)
bib-volume: (Optional) The volume of a journal or multi-volume book
bib-number: (Optional) The "(issue) number" of a journal, magazine, or tech-report, if applicable. (Most publications have a "volume", but no "number" field.)
bib-pages: (Optional) Page numbers, separated either by commas or double-hyphens.
bib-month:  (Optional) The month the book was issued.
bib-url: (Optional, non-standard) A URL for an e-journal.
---

Yo dawg! Put the bookreview here. Delete the undesired tags. Remove the blank line.
`
)

// TODO(rjkroege): consider adding the remaining types.
/*
bib-booktitle: The title of the book, if only part of it is being cited
bib-chapter: The chapter number
bib-editor: The name(s) of the editor(s)
bib-eprint: A specification of an electronic publication, often a preprint or a technical report
bib-howpublished: How it was published, if the publishing method is nonstandard
bib-institution: The institution that was involved in the publishing, but not necessarily the publisher
bib-month: The month of publication (or, if unpublished, the month of creation)
bib-organization: The conference sponsor
bib-school: The school where the thesis was written
bib-type: The type of tech-report, for example, "Research Note"
bib-url: The WWW address
bib-isbn: The isbn number
bib-dewey: The categorization number the book has.
*/
