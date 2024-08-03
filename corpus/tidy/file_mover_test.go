package tidy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)

// Cover the following additional cases
//
// file with only unix date set in the filesystem
// summary reports are generated correctly

func writeFile(t *testing.T, settings *wiki.Settings, path, contents string) {
	t.Helper()

	// Make the directory hierarchy if I need to
	abspath := filepath.Join(settings.Wikidir, path)
	absdir := filepath.Dir(abspath)
	if err := os.MkdirAll(absdir, 0755); err != nil {
		t.Errorf("can't mkdir: %s: %v", absdir, err)
	}

	if err := os.WriteFile(abspath, []byte(contents), 0644); err != nil {
		t.Errorf("can't write %s: %v", path, err)
	}
}

func readFile(t *testing.T, settings *wiki.Settings, oldpath, newpath, want string) {
	t.Helper()

	if oldpath != newpath {
		if _, err := os.ReadFile(filepath.Join(settings.Wikidir, oldpath)); err == nil {
			t.Errorf("oldpath %s still exists", oldpath)
		}
	}

	gotbytes, err := os.ReadFile(filepath.Join(settings.Wikidir, newpath))
	if err != nil {
		t.Errorf("nothing at newpath %s: %v", newpath, err)
	}

	if string(gotbytes) != want {
		// TODO(rjk): I should probably diff them as I expand test coverage.
		t.Errorf("wrong contents of newpath %s", newpath)
	}
}

func TestEachFile(t *testing.T) {
	s := &wiki.Settings{
		Wikidir: t.TempDir(),
	}

	writeFile(t, s, wrongplacefilename, wrongplacecontents)
	writeFile(t, s, rightplacefilename, rightplacecontents)
	writeFile(t, s, unnecessaryuniquingfilename, unnecessaryuniquingcontents)

	fm, err := NewFilemover(s, false)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if err := corpus.Everyfile(s, fm); err != nil {
		t.Errorf("Everyfile didn't succeed: %v", err)
	}

	readFile(t, s, wrongplacefilename, "2020/11-Nov/06/Session.md", wrongplacecontents)
	readFile(t, s, rightplacefilename, rightplacefilename, rightplacecontents)
	readFile(t, s, unnecessaryuniquingfilename, "2022/11-Nov/29/Inversion-of-Control.md", unnecessaryuniquingcontents)

}

// Some data
const wrongplacecontents = `---
title: Session
date: Fri  6 Nov 2020, 15:04:52 EST
tags: #entry #therapy
---

`
const wrongplacefilename = "foo.md"

const rightplacefilename = "2022/12-Dec/09/Wikitidy.md"

const rightplacecontents = `---
title: Wikitidy
date: Fri  9 Dec 2022, 06:44:23 EST
tags: #entry #planning #wikitools
---


`

const unnecessaryuniquingfilename = "2022/11-Nov/29/Inversion-of-Control-20221220-055304.md"
const unnecessaryuniquingcontents = `---
title: Inversion of Control
date: Tue 29 Nov 2022, 07:04:34 EST
tags: #entry #edwood #nurmi #planning
---


`
