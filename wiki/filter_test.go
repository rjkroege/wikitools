package wiki

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidBaseName(t *testing.T) {
	for _, tv := range []struct {
		want string
		args []string
	}{
		{
			"foo",
			[]string{"foo"},
		},
		{
			"foo-bar",
			[]string{"foo", "bar"},
		},
		{
			"fo_o-b-ar-_",
			[]string{"fo/o", "b ar", "#"},
		},
		{
			"fo-o-bar",
			[]string{"fo	o", "bar"},
		},
		{
			"one-two-three",
			[]string{"one", "two", "three"},
		},
		{
			"2012_12_12",
			[]string{"2012/12/12"},
		},
		{
			"foo_bar",
			[]string{"foo,bar"},
		},
		{
			"foo_s-bar",
			[]string{"foo's", "bar"},
		},
		{
			"こんにちは-bar-হয়-ন_-",
			[]string{"こんにちは", "bar", "হয় না।\n"},
		},
	} {
		if got, want := ValidBaseName(tv.args), tv.want; got != want {
			t.Errorf("got %#v but want %#v", got, want)
		}
	}
}

type MockSystem struct {
	exists bool
	now    time.Time
}

func (m MockSystem) Exists(path string) bool {
	return m.exists
}

var mocktime time.Time

func mocknow() time.Time {
	return mocktime
}

func TestUniquingExtension(t *testing.T) {
	s := &Settings{
		Wikidir: t.TempDir(),
	}

	realisticdate, _ := ParseDateUnix("1999/03/21 17:00:00")

	mocktime = time.Time{}
	nowfunc = mocknow

	// TODO(rjk): fix this up

	if got, want := s.UniquingExtension("", "there"), ""; got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}

	mocktime = realisticdate
	if err := os.WriteFile(filepath.Join(s.Wikidir, "there.md"), []byte("testfile"), 0644); err != nil {
		t.Errorf("can't write temp file: %v", err)
	}

	if got, want := s.UniquingExtension("", "there"), "-19990321-170000"; got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
