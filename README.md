# Overview
Tools to help manage a personal wiki of notes stored in `git`.

# Theory of Operation
I have several years of accumulated notes writtin in Markdown format.
I want to make them viewable and searchable across:

* iOS devices
* Macs with and w/o iCloud
* Linux

I have tried multiple schemes to keep these in sync and have concluded
what I should have known in the first place: just use `git`. So tools
here provide some structuring of a `git` repo and `git` (or an iOS
specific Git client) does the work of keeping the repositories in
sync.

The repo's `templates` directory holds starting point templates for notes.

 `unsorted` holds notes without any structuring. The single directory
has gotten unwieldy after several years. More about that below.

The  `wikinew` command makes a new article based on a specified
template and opens it via the Plan9 `plumber` in the configured editor.
`wikinew` sets the paths of new notes to be in  the `unsorted` directory but
no file is created until one saves the note from the editor.

`wikitidy` arranges the notes in the `unsorted` directory into a directory
hierarchy of *year*/*month*/*day* based on the date found in the
article metadata.

The `wikiedit` command searches all the available articles and opens
one for editing via `plumb`. *Future: it should pull and index before
searching*

The `wikiread` command searches all the available articles and
opens one for reading. *Future: it should search the git repository
and use that view (which can be decorated in some fashion) or
fall back to displaying the article with Marked*

The `wikipp` is a Markdown processor that can convert wiki Markdown
articles into HTML. It's currently setup to support converting wiki
articles for the Marked2 Markdown previewer.

# Issues
Stuff to worry about:

* how to handle pictures? Given an article `foo.md`, it can have a
companion directory `foo` that contains pictures, diagrams, 
etc  

* supporting this at first will be a non-goal

* I will later expand `wikitidy` to rewrite the metadata to work with the
Git server, iAWriter, etc.


## iOS
Figuring out the iOS workflow remains open. I currently have settled
on the following:

* There will be a shortcut equivalent to `wikinew` implemented in Scriptable
or Siri shortcuts
* Workingcopy is the Git client. It has a feature where it can make a Git
repository appear in the iOS Files app so can be used with most other
applications
* iAWriter to actually edit 

* there needs to be an equivalent to `wikinew` for  iOS

# Code Structure
`wikitools` was my first attempt at writing Go code when the Go language
was new. It's highly non-idiomatic. The code is not structured well. Aspirational
structure that I'm considering.

`article`
: article parsing functionality including metadata extraction and Markdown processing

`corpus`
: interacting with the central database of articles, tidying ops across all articles but not
the actual code to mutate the article. That would go in `article`. In particular, the KV store
implementation lives here.

`wiki`
: command line parsers, configuration file parsing, etc. This is the new name for
the `wiki` directory. Note that template expansion is an article concept. 

`generate`
: code to make reports, calendars, linkmaps, etc. In particular, the code to make
link maps or the like lives here.


