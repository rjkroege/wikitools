package search

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type nameIndexTestVector struct {
	input_location     string
	input_wikitextfile string
	want               []string
	want_err           error
}

func TestPathsforwikitext(t *testing.T) {
	tests := wikitests

	bp, err := filepath.Abs("testdata")
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
