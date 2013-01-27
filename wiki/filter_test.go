package wiki

import (
    "github.com/rjkroege/wikitools/testhelpers"
    "testing"
)

func Test_Validname(t *testing.T) {
    testhelpers.AssertString(t, "foo", Validname([]string{"foo"}))
    testhelpers.AssertString(t, "foo-bar", Validname([]string{"foo", "bar"}))
    testhelpers.AssertString(t, "fo,o-b-ar-,", Validname([]string{"fo/o", "b ar", "#"}))
    testhelpers.AssertString(t, "fo-o-bar", Validname([]string{"fo	o", "bar"}))
    testhelpers.AssertString(t, "one-two-three", Validname([]string{"one", "two", "three"}))
    testhelpers.AssertString(t, "2012,12,12", Validname([]string{"2012/12/12"}))
}
