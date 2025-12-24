# wikitools

`wikitools` is a Go project designed to manage a personal wiki of notes, primarily stored in Git. Its core purpose is to facilitate the viewing, searching, and organization of Markdown-formatted notes across various platforms including iOS, macOS, and Linux, leveraging Git for synchronization.

## Key Features:
- **Article Creation (`wikinew`):** Generate new articles from predefined templates.
- **Note Organization (`wikitidy`):** Automatically arrange unsorted notes into a structured `year/month/day` directory hierarchy based on article metadata.
- **Article Management (`wikiedit`, `wikiread`):** Search and open articles for editing or reading.
- **Markdown Processing (`wikipp`):** Convert wiki Markdown articles into HTML, useful for previewers.

The project aims to provide a robust set of tools for personal knowledge management, emphasizing a Git-centric workflow for version control and cross-device syncing. The code is structured into modules handling article parsing, corpus interaction, command-line utilities, and report generation.