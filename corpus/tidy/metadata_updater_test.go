package tidy

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/rjkroege/wikitools/wiki"
)

func TestUpdateMetadata(t *testing.T) {
	tmpd, err := ioutil.TempDir("", "testupdatemetadata")
	if err != nil {
		t.Fatal("no tempdir", err)
	}
	// defer os.RemoveAll(tmpd)
	// I'll want to go read them for sure.
	t.Log(tmpd)

	settings := &wiki.Settings{
		Wikidir: tmpd,
	}

	abc, err := makeMetadataUpdaterImpl(settings)
	if err != nil {
		t.Fatal("no makeMetadataUpdaterImpl:", err)
	}

	for i, tc := range []struct {
		inputfile    string
		fname        string
		errordetails string
		expected     string
		skipped      bool
	}{
		{
			inputfile:    "Test_header_1",
			fname:        "test_header_1.md",
			errordetails: "",
			expected:     "---\ntitle: What I want\ndate: Mon 19 Mar 2012, 06:51:15 EDT\n---\n\nI need to figure out what I want. \n",
		},
		{
			inputfile:    "Test_header_3",
			fname:        "test_header_3.md",
			errordetails: "",
			expected:     "---\ntitle: What I want\ndate: Mon 19 Mar 2012, 06:51:15 EDT\ntags: #journal\n---\n\nI need to figure out what I want. \n",
		},
		{
			inputfile:    "Test_header_6",
			fname:        "test_header_6.md",
			errordetails: "",
			expected:     "---\ntitle: What I want\ndate: Mon 19 Mar 2012, 06:51:15 EDT\ntags: #journal\nplastic: yes\ntag: empty\n---\n\nI need to figure out what to code\n",
		},
		{
			inputfile:    "Test_header_9",
			fname:        "test_header_9.md",
			errordetails: "",
			expected:     "---\ntitle: Business Korea\ndate: Mon 19 Mar 2012, 06:51:15 EDT\ntags: #book\nbib-author: Peggy Kenna and Sondra Lacy\nbib-bibkey: kenna97\nbib-publisher: Passport Books\nbib-title: Business Korea\nbib-year: 1997\n---\n\nBusiness book.\n",
		},
		{
			inputfile:    "Test_header_6_dash",
			fname:        "test_header_6_dash.md",
			errordetails: "",
			expected:     "",
			skipped:      true,
		},
		{
			inputfile:    "Test_header_10",
			fname:        "test_header_10.md",
			errordetails: "",
			expected:     "---\ntitle: Business Korea\ndate: Mon 19 Mar 2012, 06:51:15 EDT\ntags: #book #business #korea\nbib-author: Peggy Kenna and Sondra Lacy\nbib-bibkey: kenna97\nbib-publisher: Passport Books\nbib-title: Business Korea\nbib-year: 1997\n---\n\nBusiness book.\n",
		},
	} {
		// setup test
		path := filepath.Join(tmpd, tc.fname)
		fd, err := os.Create(path)
		if err != nil {
			t.Fatal("can't make", path, err)
		}
		inpath := filepath.Join("../../testdata", tc.inputfile)
		rd, err := os.Open(inpath)
		if err != nil {
			t.Fatalf("can't open %q: %v", inpath, err)
		}
		if _, err := io.Copy(fd, rd); err != nil {
			t.Fatal("can't write input write tmp file", err)
		}
		rd.Close()
		fd.Close()

		npath, err := abc.updateMetadata(path)
		switch {
		case err == nil && tc.errordetails != "":
			t.Errorf("[%d] unexpected pass", i)
			continue
		case err != nil && tc.errordetails == "":
			t.Errorf("[%d] unexpected error %v", i, err)
			continue
		case err != nil && err.Error() == tc.errordetails:
			// expected error
			continue
		case err != nil && err.Error() != tc.errordetails:
			t.Errorf("[%d] expected error got %v want %v", i, err, tc.errordetails)
		}
		// err == nil && tc.errordetails == "" means that we validate results.

		if tc.skipped {
			if npath == "" {
				// No output should be generated.
				continue
			} else if npath != "" {
				t.Errorf("[%d] expected to do nothing but got %v", i, npath)
			}
		}

		// validate that the generated is correct
		fd, err = os.Open(npath)
		if err != nil {
			t.Errorf("[%d] can't open ouput %s: %v", i, npath, err)
		}
		nval, err := ioutil.ReadAll(fd)
		if err != nil {
			t.Errorf("[%d] can't read %s: %v", i, npath, err)
		}
		if got, want := string(nval), tc.expected; got != want {
			t.Errorf("[%d] update failed got %#v, want %#v", i, got, want)
		}
	}
}
