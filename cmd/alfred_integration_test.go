package cmd

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestArgsAlfredPreprocess(t *testing.T) {
	in := []string{"my", "wiki", "article is good"}
	if got, want := ArgsAlfredPreprocess(in), in; !cmp.Equal(got, want) {
		t.Errorf("no env: %v", cmp.Diff(got, want))
	}

	// Now the forked-from-Alfred case.
	if err := os.Setenv("alfred_workflow_uid", "1"); err != nil {
		t.Errorf("can't set environment: %v", err)
	}
	defer os.Unsetenv("alfred_workflow_uid")

	for _, tv := range []struct {
		in   []string
		want []string
	}{
		{
			in:   []string{"wikitools"},
			want: []string{"wikitools"},
		},
		{
			in:   []string{"wikitools", "this is a good article"},
			want: []string{"wikitools", "newautocomplete", "this", "is", "a", "good", "article"},
		},
		{
			in:   []string{"wikitools", "in"},
			want: []string{"wikitools", "newautocomplete", "in"},
		},
		{
			in:   []string{"wikitools", "#actionphase this is a good article"},
			want: []string{"wikitools", "new", "this", "is", "a", "good", "article"},
		},
	} {
		if got, want := ArgsAlfredPreprocess(tv.in), tv.want; !cmp.Equal(got, want) {
			t.Errorf("no env: %v", cmp.Diff(got, want))
		}
	}
}
