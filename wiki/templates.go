package wiki

const (
journaltmpl = 
`title: {{.Title}}
date: {{.DetailedDate}}
tags: {{.Tagstring}}

Yo dawg! Write stuff here.
`

booktmpl =
`title: {{.Title}}
date: {{.DetailedDate}}
tags: {{.Tagstring}}
bib-address: Publisher's address (usually just the city, but can be the full address for lesser-known publishers)
bib-author: The name(s) of the author(s) (in the case of more than one author, separated by and)
bib-booktitle: The title of the book, if only part of it is being cited
bib-chapter: The chapter number
bib-edition: The edition of a book, long form (such as "first" or "second")
bib-editor: The name(s) of the editor(s)
bib-eprint: A specification of an electronic publication, often a preprint or a technical report
bib-howpublished: How it was published, if the publishing method is nonstandard
bib-institution: The institution that was involved in the publishing, but not necessarily the publisher
bib-journal: The journal or magazine the work was published in
bib-month: The month of publication (or, if unpublished, the month of creation)
bib-number: The "(issue) number" of a journal, magazine, or tech-report, if applicable. (Most publications have a "volume", but no "number" field.)
bib-organization: The conference sponsor
bib-pages: Page numbers, separated either by commas or double-hyphens.
bib-publisher: The publisher's name
bib-school: The school where the thesis was written
bib-series: The series of books the book was published in (e.g. "The Hardy Boys" or "Lecture Notes in Computer Science")
bib-title: {{.Title}}
bib-type: The type of tech-report, for example, "Research Note"
bib-url: The WWW address
bib-volume: The volume of a journal or multi-volume book
bib-year: The year of publication (or, if unpublished, the year of creation)
bib-isbn: The isbn number
bib-dewey: The categorization number the book has.

Yo dawg! Put the bookreview here. Delete the undesired tags.
`

)
