package bibtex

import (
    "testing"
    "github.com/rjkroege/wikitools/testhelpers"
)


func Test_FilterExtrakeys_Empty(t *testing.T) {
	 testhelpers.AssertStringArray(t, []string{}, FilterExtrakeys([]string{}))
}

func Test_FilterExtrakeys_Removing(t *testing.T) {
	 testhelpers.AssertStringArray(t, []string{}, FilterExtrakeys([]string{"foo"}))
}

func Test_FilterExtrakeys_Keeping(t *testing.T) {
	 testhelpers.AssertStringArray(t, []string{"bar"}, FilterExtrakeys([]string{"foo", "bib-bar"}))
}