package wiki

import (
	"testing"
)

func TestSplitActualName(t *testing.T) {

	for _, tv := range []struct {
		input string
		n     string
		u     string
		e     string
	}{
		{
			input: "foo",
			n:     "foo",
		},
		{
			input: "foo.md",
			n:     "foo",
			e:     ".md",
		},
		{
			input: "Wikitidy-20221220-055304.md",
			n:     "Wikitidy",
			u:     "-20221220-055304",
			e:     ".md",
		},
		{
			input: "All-About-Wikitidy-20221220-055304.md",
			n:     "All-About-Wikitidy",
			u:     "-20221220-055304",
			e:     ".md",
		},
		{
			input: "All-About-Wikitidy-2022.md",
			n:     "All-About-Wikitidy-2022",
			u:     "",
			e:     ".md",
		},
	} {
		n, u, e := SplitActualName(tv.input)
		if got, want := n, tv.n; got != want {
			t.Errorf("%s: got %#v but want %#v", "name", got, want)
		}
		if got, want := u, tv.u; got != want {
			t.Errorf("%s: got %#v but want %#v", "name", got, want)
		}
		if got, want := e, tv.e; got != want {
			t.Errorf("%s: got %#v but want %#v", "name", got, want)
		}
	}
}
