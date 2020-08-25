package article

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rjkroege/wikitools/config"
)

type mockFileInfo struct {
	name string
}

func (m *mockFileInfo) Name() string       { return m.name }
func (m *mockFileInfo) Size() int64        { return int64(0) }
func (m *mockFileInfo) Mode() os.FileMode  { return os.FileMode(0) }
func (m *mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool        { return false }
func (m *mockFileInfo) Sys() interface{}   { return nil }

func TestSkipper(t *testing.T) {
	for i, tc := range []struct {
		base     string
		abs      string
		expected bool
	}{
		{
			"Summer-2020.md",
			filepath.Join(config.Basepath, "2020/7/19/Summer-2020.md"),
			false,
		},
		{
			"README.md",
			filepath.Join(config.Basepath, "README.md"),
			true,
		},
		{
			"foo.md",
			filepath.Join(config.Basepath, "templates/foo.md"),
			true,
		},
		{
			"foo.md",
			filepath.Join(config.Basepath, "templates/subdir/foo.md"),
			true,
		},
	} {
		if got, want := skipper(tc.abs, &mockFileInfo{tc.base}), tc.expected; got != want {
			t.Errorf("[%d] skipper %s got %v want %v", i, tc.abs, got, want)
		}

	}
}

func TestUpdateMetadata(t *testing.T) {
	tmpd, err := ioutil.TempDir("", "testupdatemetadata")
	if err != nil {
		t.Fatal("no tempdir", err)
	}
	// defer os.RemoveAll(tmpd)
	// I'll want to go read them for sure.
	log.Println(tmpd)

	abc, err := makeMetadataUpdaterImpl()
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
			inputfile:    test_header_1,
			fname:        "test_header_1.md",
			errordetails: "",
			expected:     "---\ntitle: What I want\ndate: Mon 19 Mar 2012, 06:51:15 EDT\n---\n\nI need to figure out what I want. \n",
		},
		{
			inputfile:    test_header_3,
			fname:        "test_header_3.md",
			errordetails: "",
			expected:     "---\ntitle: What I want\ndate: Mon 19 Mar 2012, 06:51:15 EDT\ntags: #journal\n---\n\nI need to figure out what I want. \n",
		},
		{
			inputfile:    test_header_6,
			fname:        "test_header_6.md",
			errordetails: "",
			expected:     "---\ntitle: What I want\ndate: Mon 19 Mar 2012, 06:51:15 EDT\ntags: #journal\nplastic: yes\ntag: empty\n---\n\nI need to figure out what to code\n",
		},
		{
			inputfile:    test_header_9,
			fname:        "test_header_9.md",
			errordetails: "",
			expected:     "---\ntitle: Business Korea\ndate: Mon 19 Mar 2012, 06:51:15 EDT\ntags: #book\nbib-author: Peggy Kenna and Sondra Lacy\nbib-bibkey: kenna97\nbib-publisher: Passport Books\nbib-title: Business Korea\nbib-year: 1997\n---\n\nBusiness book.\n",
		},
		{
			inputfile:    test_header_6_dash,
			fname:        "test_header_6_dash.md",
			errordetails: "",
			expected:     "",
			skipped:      true,
		},
		{
			inputfile:    test_header_10,
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
		if length, err := fd.WriteString(tc.inputfile); err != nil || length != len(tc.inputfile) {
			t.Fatal("can't write input", err)
		}
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
