package wiki

import (
	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/testhelpers"
	"testing"
	"time"
)

func Test_ValidBaseName(t *testing.T) {
	testhelpers.AssertString(t, "foo", ValidBaseName([]string{"foo"}))
	testhelpers.AssertString(t, "foo-bar", ValidBaseName([]string{"foo", "bar"}))
	testhelpers.AssertString(t, "fo_o-b-ar-,", ValidBaseName([]string{"fo/o", "b ar", "#"}))
	testhelpers.AssertString(t, "fo-o-bar", ValidBaseName([]string{"fo	o", "bar"}))
	testhelpers.AssertString(t, "one-two-three", ValidBaseName([]string{"one", "two", "three"}))
	testhelpers.AssertString(t, "2012_12_12", ValidBaseName([]string{"2012/12/12"}))
	testhelpers.AssertString(t, "foo-bar", ValidBaseName([]string{"foo,bar"}))
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
	realisticdate, _ := article.ParseDateUnix("1999/03/21 17:00:00")

	nd := &MockSystem{false, time.Time{}}
	testhelpers.AssertString(t, "there.md", UniqueValidName("hello/", "there", ".md", nd))

	nd = &MockSystem{true, realisticdate}
	testhelpers.AssertString(t, "there-19990321-170000.md", UniqueValidName("hello/", "there", ".md", nd))
}
