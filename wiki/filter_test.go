package wiki

import (
	"testing"
	"time"

	"github.com/rjkroege/wikitools/testhelpers"
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

func (m MockSystem) Now() time.Time {
	return m.now
}

func Test_UniqueValidName(t *testing.T) {
	realisticdate, _ := ParseDateUnix("1999/03/21 17:00:00")

	nd := &MockSystem{false, time.Time{}}
	testhelpers.AssertString(t, "there.md", UniqueValidName("hello/", "there", ".md", nd))

	nd = &MockSystem{true, realisticdate}
	testhelpers.AssertString(t, "there-19990321-170000.md", UniqueValidName("hello/", "there", ".md", nd))
}
