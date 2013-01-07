package article

import (
    "testing"
    "time"
)

func Test_Makearticle(t *testing.T) {
    md := MetaData{ "foo.md", "", time.Now(), time.Now(), "", "", false, "", ""}
    if md.FormattedName() != "foo.html" {
        t.Errorf("expected  %s != to actual %s", "foo.html", md.Name)
    }    
}

func Test_UrlForName(t *testing.T) {
    md := MetaData{ "foo.md", "", time.Now(), time.Now(), "", "", false, "", ""}
    s := md.UrlForName("flimmer/blo")
    if s != "file://flimmer/blo/foo.html" {
        t.Errorf("expected  %s != to actual %s", "file://flimmer/blo/foo.html", s)
    }
}

type tuple struct {
    ex string
    err error
    in string
} 

func Test_parseDataUnix(t *testing.T) {
    testdates := []tuple {
        tuple{ "Tue Sep 11 17:34:00 EDT 2012", nil,  "11 Sep 17:34:00 2012"},
        tuple{ "Sat Oct 27 11:39:41 PDT 2012", nil,  "Sat Oct 27 11:39:41 PDT 2012"},
        tuple{ "Wed Jun 15 08:24:39 EDT 2011", nil,  "2011/06/15 08:24:39"},
        tuple{ "Tue Dec 27 17:46:16 EST 2011", nil,  "2011/12/27 17:46:16"},
        tuple{ "Sun Mar 14 08:00:00 EST 2004", nil,  "200403140800"},
        tuple{ "Tue Dec 11 17:34:00 EST 2012", nil,  "11 Dec 17:34:00 2012"},
        tuple{ "Sat Dec  1 17:34:00 EST 2012", nil,  "1 Dec 17:34:00 2012"}     }

    for _, tu := range(testdates) {
        r, err := parseDateUnix(tu.in)
        if tu.err != err {
            t.Errorf("invalid error value in test date %s", tu.in)
        }
        if tu.err == nil && tu.ex != r.Format(time.UnixDate) {
            t.Errorf("bad date: expected %s, received %s", tu.ex, r.Format(time.UnixDate))
        }
    }
}
