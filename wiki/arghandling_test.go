package wiki

import "testing"

func TestSplit_Empty(t *testing.T) {
    ar, tg := Split([]string { })

    if len(ar) != 0  && len(tg) != 0 {
        t.Error("empty args not handled correctly")
    }
}

func expect(t *testing.T, expected []string , actual []string) {
    for i, _ := range(actual) {
        if expected[i] != actual[i] {
            t.Errorf("expected  %s != to actual %s", expected[i], actual[i]);
        }
    }
}

func TestSplit_Basic(t *testing.T) {
    ar, tg := Split([]string { "@one", "@two", "three", "four"})

    expect(t, []string{ "@one", "@two" }, tg)
    expect(t, []string{ "three", "four" }, ar)
}

func TestSplit_Unordered(t *testing.T) {
    ar, tg := Split([]string { "@one", "three", "four", "@two"})

    expect(t, []string{ "@one", "@two" }, tg)
    expect(t, []string{ "three", "four" }, ar)
}

func TestPicktemplate_firstarg(t *testing.T) {
    ar, tg := Split([]string { "@flong", "journal", "@fling" })
    tm, ar, tg := Picktemplate(ar, tg)
    if tm != journaltmpl {
        t.Error("didn't pick correct template, instead chose: " + tm)
    }
    expect(t, []string{"@flong", "@fling", "@journal"}, tg)
    if len(ar) != 0 {
        t.Error("should not have any args")
    }
}

func TestPicktemplate_tagpriority(t *testing.T) {
    ar, tg := Split([]string { "@flong", "journal", "@book" })
    tm, ar, tg := Picktemplate(ar, tg)
    if tm != booktmpl {
        t.Error("didn't pick correct template")
    }
    expect(t, []string{"@flong", "@book"}, tg)
    expect(t, []string{"journal"}, ar)
}
