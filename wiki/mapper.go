 package wiki

// LinkToFile is implemented by objects that can return a unique or all file paths corresponding
// to a given wikilink.
type LinkToFile interface {
	// Returns a single unique path corresponding to the wikitext found in
	// file lsd limiting the search for target paths to files in location or
	// error if this is impossible.
	Path(location, lsd, wikitext string) (string, error)

	// Returns all (absolute) paths in the wiki that would match wikitext.
	// TODO(rjk): Why is lsd here?
	Allpaths(location, lsd, wikitext string) ([]string, error)

	// Wikitext returns a wikitext such that clicking on it in file frompath
	// will open file topath or an error if it was impossible to do so. In
	// particular: Path(wikiroot, frompath, Wikitext(frompath, topath)) ==
	// topath
	Wikitext(frompath, topath string) (string, error)
}
