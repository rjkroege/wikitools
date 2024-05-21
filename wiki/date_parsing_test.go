package wiki

import (
	"testing"
	"time"
)

type pdSR struct {
	ex  string
	err error
	in  string
}

func Test_ParseDateUnix(t *testing.T) {
	testdates := []pdSR{
		{"Mon Mar 11 00:00:00 EDT 2013", nil, "Monday, Mar 11, 2013"},
		{"Tue Sep 11 17:34:00 EDT 2012", nil, "11 Sep 17:34:00 2012"},
		{"Sat Oct 27 11:39:41 PDT 2012", nil, "Sat Oct 27 11:39:41 PDT 2012"},
		{"Wed Jun 15 08:24:39 EDT 2011", nil, "2011/06/15 08:24:39"},
		{"Tue Dec 27 17:46:16 EST 2011", nil, "2011/12/27 17:46:16"},
		{"Sun Mar 14 08:00:00 EST 2004", nil, "200403140800"},
		{"Tue Dec 11 17:34:00 EST 2012", nil, "11 Dec 17:34:00 2012"},
		{"Fri Jun 14 07:25:48 EDT 2013", nil, "Fri 14 Jun 2013, 07:25:48 EDT"},
		{"Sat Dec  1 17:34:00 EST 2012", nil, "1 Dec 17:34:00 2012"},
		{"Tue Sep  5 11:14:03 PDT 2006", nil, "Tue Sep  5 11:14:03 PDT 2006"},
		{"Tue Feb  5 08:52:22 -0700 2019", nil, "2019-02-05 08:52:22.000000000 -0700"},
		{"Fri Sep 13 07:19:07 -0600 2019", nil, "Fri 13 Sep 2019, 07:19:07 -0600"},
		{"Wed Aug 10 13:46:00 PDT 2016", nil, "20160810 1346 PDT"},
	}

	for _, tu := range testdates {
		r, err := ParseDateUnix(tu.in)
		if tu.err != err {
			t.Errorf("invalid error value in test date %s", tu.in)
		}
		if tu.err == nil && tu.ex != r.Format(time.UnixDate) {
			t.Errorf("bad date: expected %s, received %s", tu.ex, r.Format(time.UnixDate))
		}
	}
}
