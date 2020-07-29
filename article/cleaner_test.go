package article

import (
	"testing"
	"path/filepath"
	"os"
	"time"

	"github.com/rjkroege/wikitools/config"
	
)

type mockFileInfo struct {
	name string
}

func (m *mockFileInfo) Name() string     {	return m.name }
func (m *mockFileInfo) Size() int64     { 	return int64(0) }
func (m *mockFileInfo) Mode() os.FileMode     {	return os.FileMode(0) }
func (m *mockFileInfo) ModTime() time.Time     { 	return time.Time{} }
func (m *mockFileInfo) IsDir() bool {	return false }
func (m *mockFileInfo) Sys() interface{} {	return nil }


func TestSkipper(t *testing.T) {
	for i, tc  := range []struct{
		base string
		abs string
		expected bool
	} {
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
