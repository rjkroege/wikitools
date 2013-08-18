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

func Test_ExtractBibTeXEntryType(t *testing.T) {
	te, e := ExtractBibTeXEntryType([]string{ "@book" })
	if te != "book" && e != nil {
		t.Error("book is not the default entry type")
	}
	
	te, e = ExtractBibTeXEntryType([]string{ "@bibtex-article", "@book" })
	if te != "article" && e != nil {
		t.Error("entry type is 'article' but have selected: ", te)
	}

	te, e = ExtractBibTeXEntryType([]string{ "@book", "@bibtex-article" })
	if te != "article" && e != nil {
		t.Error("article specified, not selected")
	}
	
	te, e = ExtractBibTeXEntryType([]string{ "@journal", "@bibtex-article" })
	if e == nil || e.Error() != "No book tag present."  {
		t.Error("bad error for missing @book")
	}

	te, e = ExtractBibTeXEntryType([]string{ "@book", "@bibtex-article",  "@bibtex-phd" })
	if e == nil || (e.Error() != "More than one supplementary @bibtex-(.*) tag." && te != "phd" ) {
		t.Error("bad error for missing @book")
	}
	
}