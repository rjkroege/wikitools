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

func Test_Intersectsorted_empty(t *testing.T) {
	m := intersectsorted([]string{}, []string{})
	testhelpers.AssertStringArray(t, []string{}, m)
}

func Test_VerifyRequiredFields(t *testing.T) {
	// Lots of cases

	err := VerifyRequiredFields("article", []string{"bibkey", "author", "title", "journal", "year"})
	if err != nil {
		t.Error("error should be nil for exact match")
	}

	err = VerifyRequiredFields("article", []string{"bibkey", "author", "title", "journal", "year", "extra-stuff"})
	if err != nil {
		t.Error("error should be nil for additional entry")
	}

	err = VerifyRequiredFields("article", []string{})
	if err == nil {
		t.Error("missing fields should have a non-nil error")
	} else {
		testhelpers.AssertString(t, "Missing required fields: author bibkey journal title year for entry type article", err.Error())
	}
}

func Test_VerifyRequiredFields_bookeditor(t *testing.T) {
	err := VerifyRequiredFields("book", []string{"bibkey", "author", "title", "publisher", "year"})
	if err != nil {
		t.Error("wrongly claim that valid author book keys are invalid: " + err.Error())		
	}	

	err = VerifyRequiredFields("book", []string{"bibkey", "editor", "title", "publisher", "year"})
	if err != nil {
		t.Error("wrongly claim that valid editor book keys are invalid: " + err.Error())	
	}	

	err = VerifyRequiredFields("book", []string{"bibkey", "flimflam", "title", "publisher", "year"})
	if err == nil {
		t.Error("book needs to have editor or author")
	}
}

func Test_VerifyRequiredFields_inbook_editor(t *testing.T) {
	err := VerifyRequiredFields("inbook", []string{"bibkey", "author", "title", "chapter", "publisher", "year"})
	if err != nil {
		t.Error("wrongly claim that valid author/chapter inbook keys are invalid: " + err.Error())		
	}	

	err = VerifyRequiredFields("inbook", []string{"bibkey", "editor", "title", "chapter", "publisher", "year"})
	if err != nil {
		t.Error("wrongly claim that valid editor/chapter inbook keys are invalid: " + err.Error())	
	}

	err = VerifyRequiredFields("inbook", []string{"bibkey", "author", "title", "pages", "publisher", "year"})
	if err != nil {
		t.Error("wrongly claim that valid author/pages inbook keys are invalid: " + err.Error())		
	}	

	err = VerifyRequiredFields("inbook", []string{"bibkey", "editor", "title", "pages", "publisher", "year"})
	if err != nil {
		t.Error("wrongly claim that valid editor/pages inbook keys are invalid: " + err.Error())	
	}

	err = VerifyRequiredFields("inbook", []string{"bibkey", "editor", "title", "publisher", "year"})
	if err == nil {
		t.Error("for editor inbook, wrongly claim missing both pages and chapter is correct")
	} else {
		testhelpers.AssertString(t, "Missing required fields: chapter pages for entry type inbook", err.Error())
	}

	err = VerifyRequiredFields("inbook", []string{"bibkey", "author", "title", "publisher", "year"})
	if err == nil {
		t.Error("for author inbook, wrongly claim missing both pages and chapter is correct")
	} else {
		testhelpers.AssertString(t, "Missing required fields: chapter pages for entry type inbook", err.Error())
	}
}

const (
output1 =
`@book(jones2013,
	editor = "Peyton Jones",
	publisher = "Penguin",
	title = "Collected Angst",
	year = "2013",
)
`
)

func Test_ExploreTemplating(t *testing.T) {
	s, e := CreateBibTexEntry([]string{"@book"}, map[string]string{"bib-bibkey": "jones2013", "bib-editor": "Peyton Jones", "bib-title": "Collected Angst", "bib-publisher": "Penguin", "bib-year": "2013"})
	if e != nil {
		t.Error("CreateBibTexEntry wrongly failed with: " + e.Error())
	}
	testhelpers.AssertString(t, output1, s)

	s, e = CreateBibTexEntry([]string{"@book",  "@bibtex-article"}, map[string]string{"bib-bibkey": "jones2013", "bib-editor": "Peyton Jones", "bib-title": "Collected Angst", "bib-publisher": "Penguin", "bib-year": "2013"})
	if e == nil {
		t.Error("CreateBibTexEntry wrongly succeeded")
	} else {
		testhelpers.AssertString(t, "", s)
		testhelpers.AssertString(t, "Missing required fields: author journal for entry type article", e.Error())
	}
}
