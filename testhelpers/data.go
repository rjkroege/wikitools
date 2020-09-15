package testhelpers

const Test_header_1 = `title: What I want
date: 2012/03/19 06:51:15

I need to figure out what I want. 
`

const Test_header_1_dash = `---
title: What I want
date: 2012/03/19 06:51:15
---

I need to figure out what I want. 
`

const Test_header_2 = `title: What I want
date: 2012/03/19 06:51:15
tags: @journal

I need to figure out what I want. 
`
const Test_header_3 = `date: 2012/03/19 06:51:15
title: What I want
tags: @journal

I need to figure out what I want. 
`
const Test_header_4 = `I need
to figure out what I want. And code it.
`

const Test_header_5 = `Date: 2012/03/19 06:51:15
Title: What I want
tags: @journal
`

const Test_header_6 = `plastic: yes
Date: 2012/03/19 06:51:15
Tag: empty
Title: What I want
tags: @journal

I need to figure out what to code
`

const Test_header_6_dash = `---
plastic: yes
Date: 2012/03/19 06:51:15
Tag: empty
Title: What I want
tags: @journal
---

I need to figure out what to code
`

const Test_header_7 = `plastic: yes
Date: 2012/03/19 06:51:15
Tag: empty
Title: What I want
tags: @journal   @fiddle

I need to figure out what to code
`
const Test_header_8 = `plastic: yes
Date: 2012/03/19 06:51:15
Tag: empty
Title: What I want
tags: @journal  @hello     @bye

I need to figure out what to code
`

const Test_header_9 = `title: Business Korea
date: 2012/03/19 06:51:15
tags: @book
bib-bibkey: kenna97
bib-author: Peggy Kenna and Sondra Lacy
bib-title: Business Korea
bib-publisher: Passport Books
bib-year:  1997

Business book.
`

const Test_header_9_dash = `---
title: Business Korea
date: 2012/03/19 06:51:15
tags: @book
bib-bibkey: kenna97
bib-author: Peggy Kenna and Sondra Lacy
bib-title: Business Korea
bib-publisher: Passport Books
bib-year:  1997
---

Business book.
`

const Test_header_10_dash = `---
title: Business Korea
date: 2012/03/19 06:51:15
tags: @book #business #korea
bib-bibkey: kenna97
bib-author: Peggy Kenna and Sondra Lacy
bib-title: Business Korea
bib-publisher: Passport Books
bib-year:  1997
---

Business book.
`
const Test_header_10 = `
title: Business Korea
date: 2012/03/19 06:51:15
tags: @book #business #korea
bib-bibkey: kenna97
bib-author: Peggy Kenna and Sondra Lacy
bib-title: Business Korea
bib-publisher: Passport Books
bib-year:  1997

Business book.
`