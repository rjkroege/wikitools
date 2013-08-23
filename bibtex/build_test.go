package bibtex

import (
    "testing"
    "github.com/rjkroege/wikitools/testhelpers"
)


func Test_FilterExtrakeys_Empty(t *testing.T) {
	m, v := FilterExtrakeys(map[string]string{})
	testhelpers.AssertStringArray(t, []string{}, v)
	testhelpers.AssertStringMap(t, map[string]string{}, m)
}

func Test_FilterExtrakeys_Removing(t *testing.T) {
	m, v := FilterExtrakeys(map[string]string{"foo": "hello"})
	testhelpers.AssertStringArray(t, []string{},v)
	testhelpers.AssertStringMap(t, map[string]string{}, m)
}

func Test_FilterExtrakeys_Keeping(t *testing.T) {
	m, v := FilterExtrakeys(map[string]string{"foo": "hello", "bib-bar": "bye"})
	 testhelpers.AssertStringArray(t, v, []string{"bar"})
	 testhelpers.AssertStringMap(t, map[string]string{"bar": "bye"}, m)
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

//
//func Test_VerifyRequiredFields(t *testing.T) {
//	
//}
