package search

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type nameIndexTestVector struct {
	input_location     string
	input_wikitextfile string
	want               []string
	want_err           error
}

func TestPathsforwikitext(t *testing.T) {
	tests := wikitests

	bp, err := filepath.Abs("../testdata")
	if err != nil {
		t.Fatalf("test can't run: %v", err)
	}

	spix := MakeWikilinkNameIndex(bp)

	for _, tc := range tests {
		got, err := spix.pathsforwikitext(filepath.Join(bp, tc.input_location), tc.input_wikitextfile)

		want := []string{}
		for _, s := range tc.want {
			want = append(want, filepath.Join(bp, s))
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("path mismatch (-want +got):\n%s", diff)
		}
		if diff := cmp.Diff(tc.want_err, err); diff != "" {
			t.Errorf("error mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestSplitPathParts(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: "/foo",
			want:  []string{"/foo"},
		},
		{
			input: "/foo//bar",
			want:  []string{"/foo", "bar"},
		},
		{
			input: "foo//bar",
			want:  []string{"foo", "bar"},
		},
		{
			input: "bar",
			want:  []string{"bar"},
		},
		{
			input: "bar/goodness/wiki",
			want:  []string{"bar", "goodness", "wiki"},
		},
	}

	for _, tc := range tests {
		got := splitPathParts(tc.input)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("dump mismatch (-want +got):\n%s", diff)
		}

	}
}

type WikitextTestCase struct {
	frompath string
	topath string
	want string
	wanterr error
	nexterr error
}

func TestWikitext(t *testing.T) {
	bp, err := filepath.Abs("../testdata")
	if err != nil {
		t.Fatalf("test can't run: %v", err)
	}
	wikiroot := filepath.Join(bp, "wiki")

	spix := MakeWikilinkNameIndex(bp)

	tests := []WikitextTestCase{
		{
			frompath: "wiki/unsorted/Saturday.md",
			topath: "wiki/unsorted/EveningJournal.md",
			want: "EveningJournal.md",
			wanterr: nil,
		},
		{
			frompath: "wiki/unsorted/Saturday.md",
			topath: "wiki/2023/02-Feb/28/Saturday.md",
			want: "28/Saturday.md",
			wanterr: nil,
		},
		{
			frompath: "wiki/2023/04-Apr/20/Thursday.md",
			topath: "wiki/2023/04-Apr/24/Monday.md",
			want: "24/Monday.md",
			wanterr: nil,
		},
		{
			frompath: "wiki/2023/06-Jun/5/Monday.md",
			topath: "wiki/2023/06-Jun/Monday.md",
			want: "06-Jun/Monday.md",
			wanterr: nil,
		},
		{
			frompath: "wiki/2023/06-Jun/Monday.md",
			topath: "wiki/2023/06-Jun/5/Monday.md",
			want: "5/Monday.md",
			wanterr: nil,
		},
		{
			frompath: "wiki/unsorted/EveningJournal.md",
			topath: "wiki/2023/Decisions.md",
			want: "2023/Decisions.md",
			wanterr: nil,
		},
		{
			frompath: "wiki/2023/06-Jun/Monday.md",
			topath: "wiki/2023/12-Dec/1/PlottingTools.md",
			want: "12-Dec/1/PlottingTools.md",
			wanterr: nil,
		},
	}

	for _, tc := range tests {
		got, err := spix.Wikitext(
			filepath.Join(bp, tc.frompath), 
			filepath.Join(bp, tc.topath))

		want := tc.want
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Wikitext mismatch (-want +got):\n%s", diff)
			continue
		}
		if diff := cmp.Diff(tc.wanterr, err, cmpopts.EquateErrors()); diff != "" {
			t.Errorf("Wikitext error mismatch (-want +got):\n%s", diff)
			continue
		}

		// Verify invertability
		igot, ierr := spix.Path(
			wikiroot,
			filepath.Dir(filepath.Join(bp, tc.frompath)),
			got)
		
		want = filepath.Join(bp, tc.topath)
		if diff := cmp.Diff(want, igot); diff != "" {
			t.Logf("want: %q igot: %q, (wikitext) got: %q, frompath: %q topath: %q", want, igot, got, tc.frompath, tc.topath)
			t.Errorf("Path mismatch (-want +got):\n%s", diff)
			continue
		}
		if diff := cmp.Diff(tc.nexterr, ierr, cmpopts.EquateErrors()); diff != "" {
			t.Errorf("Path error mismatch (-want +got):\n%s", diff)
			continue
		}
	}
}
