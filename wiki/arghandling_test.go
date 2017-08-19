package wiki

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSplit_Empty(t *testing.T) {
	ar, tg := Split([]string{})

	if len(ar) != 0 && len(tg) != 0 {
		t.Error("empty args not handled correctly")
	}
}

func expect(t *testing.T, expected []string, actual []string) {
	for i, _ := range actual {
		if expected[i] != actual[i] {
			t.Errorf("expected  %s != to actual %s", expected[i], actual[i])
		}
	}
}

func TestSplit_Basic(t *testing.T) {
	ar, tg := Split([]string{"@one", "@two", "three", "four"})

	expect(t, []string{"@one", "@two"}, tg)
	expect(t, []string{"three", "four"}, ar)
}

func TestSplit_Unordered(t *testing.T) {
	ar, tg := Split([]string{"@one", "three", "four", "@two"})

	expect(t, []string{"@one", "@two"}, tg)
	expect(t, []string{"three", "four"}, ar)
}

func TestPicktemplate_firstarg(t *testing.T) {
	tmpls := NewTemplatePalette()

	journaltimepicker = func() bool { return true }
	defer func() { journaltimepicker = BeforeNoon }()

	ar, tg := Split([]string{"@flong", "journal", "@fling"})
	tm, ar, tg := tmpls.Picktemplate(ar, tg)
	if tm != tmpls["journalam"] {
		t.Errorf("didn't pick correct template, instead chose: %v", tm)
	}
	expect(t, []string{"@flong", "@fling", "@journal"}, tg)
	if len(ar) != 0 {
		t.Error("should not have any args")
	}

	ar, tg = Split([]string{"@flong", "@code"})
	tm, ar, tg = tmpls.Picktemplate(ar, tg)
	if tm != tmpls["code"] {
		t.Error("didn't pick correct template, instead chose: %v", tm)
	}
	expect(t, []string{"@flong", "@code"}, tg)
	if len(ar) != 0 {
		t.Error("should not have any args")
	}

	journaltimepicker = func() bool { return false }
	ar, tg = Split([]string{"@flong", "@journal"})
	tm, ar, tg = tmpls.Picktemplate(ar, tg)
	if tm != tmpls["journalpm"] {
		t.Error("didn't pick correct template, instead chose: %v", tm)
	}
	expect(t, []string{"@flong", "@journal"}, tg)
	if len(ar) != 0 {
		t.Error("should not have any args")
	}

}

func TestPicktemplate_tagpriority(t *testing.T) {
	tmpls := NewTemplatePalette()
	ar, tg := Split([]string{"@flong", "journal", "@book"})
	tm, ar, tg := tmpls.Picktemplate(ar, tg)
	if tm != tmpls["book"] {
		t.Error("didn't pick correct template")
	}
	expect(t, []string{"@flong", "@book"}, tg)
	expect(t, []string{"journal"}, ar)
}

const (
	codetemplate = `
hello codetemplate
`

	totallynew = `
hello totallynew
`
)

func TestAddDynamcTemplates(t *testing.T) {
	tmpls := NewTemplatePalette()

	// Create some test data.
	dir, err := ioutil.TempDir("", "configuration")
	if err != nil {
		t.Fatalf("Couldn't make tempdir %v", err)
	}
	defer os.RemoveAll(dir)

	codetemplatepath := filepath.Join(dir, "codetemplate")
	fd, err := os.Create(codetemplatepath)
	if err != nil {
		t.Fatalf("Couldn't open the codetemplate %v", err)
	}
	if _, err := io.WriteString(fd, codetemplate); err != nil {
		fd.Close()
		t.Fatalf("Couldn't write the tmpfile %v", err)
	}
	fd.Close()

	totallynewpath := filepath.Join(dir, "totallynew")
	fd, err = os.Create(totallynewpath)
	if err != nil {
		t.Fatalf("Couldn't open the totallynew %v", err)
	}
	if _, err := io.WriteString(fd, totallynew); err != nil {
		fd.Close()
		t.Fatalf("Couldn't write the tmpfile %v", err)
	}
	fd.Close()

	missingfilepath := filepath.Join(dir, "missing_file")

	config := map[string]string{
		"code":       codetemplatepath,
		"book":       codetemplatepath,
		"totallynew": totallynewpath,
		"journalam":  missingfilepath,
	}

	tmpls.AddDynamcTemplates(config)

	if got, want := tmpls["code"].Template, basetmpl; got != want {
		t.Errorf("got %v but wanted %v", got, want)
	}
	if got, want := tmpls["code"].Custombody, codetemplate; got != want {
		t.Errorf("got %v but wanted %v", got, want)
	}

	if got, want := tmpls["book"].Template, booktmpl; got != want {
		t.Errorf("got %v but wanted %v", got, want)
	}
	if got, want := tmpls["book"].Custombody, codetemplate; got != want {
		t.Errorf("got %v but wanted %v", got, want)
	}

	if _, ok := tmpls["totallynew"]; !ok {
		t.Fatalf("no added key %s", "totallynew")
	}
	if got, want := tmpls["totallynew"].Template, basetmpl; got != want {
		t.Errorf("got %v but wanted %v", got, want)
	}
	if got, want := tmpls["totallynew"].Custombody, totallynew; got != want {
		t.Errorf("got %v but wanted %v", got, want)
	}

	if got, want := tmpls["journalam"].Template, basetmpl; got != want {
		t.Errorf("got %v but wanted %v", got, want)
	}
	if got, want := tmpls["journalam"].Custombody, "File "+missingfilepath+" for key journalam had error: open "+missingfilepath+": no such file or directory"; got != want {
		t.Errorf("got %#v but wanted %v", got, want)
	}
}
