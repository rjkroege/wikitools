package wiki

import (
	"log"
	"time"
)

const (
	lsdate         = "_2 Jan 15:04:05 2006"
	lstring        = "20060102150405"
	slashdate      = "2006/01/02 15:04:05"
	sstring        = "200601021504"
	unixlike       = "Mon _2 Jan 2006 15:04:05"
	almostunixlike = "Monday _2 Jan 2006 15:04:05"
	short          = "Monday, Jan _2, 2006"
	df             = "Mon _2 Jan 2006, 15:04:05"
	iso            = "2006-01-02 15:04:05.999999999 -0700"
	unixlikezoned  = "Mon _2 Jan 2006, 15:04:05 -0700"

// odf = "Mon Jan _2 15:04:05 MST 2006"
)

var easternTimeZone *time.Location

// ParseDateUnix attempts to extract a valid time from the provided string in
// a number of formats.
func ParseDateUnix(ds string) (t time.Time, err error) {
	timeformats := []string{
		lsdate,
		lstring,
		slashdate,
		sstring,
		unixlike,
		almostunixlike,
		unixlikezoned,
		short,
		df,
		iso,
	}

	timeswithzones := []string{
		time.UnixDate,
	}

	if easternTimeZone == nil {
		easternTimeZone, err = time.LoadLocation("America/Toronto")
		if err != nil {
			log.Fatal("no eastern time zone?")
		}
	}

	for _, fs := range timeformats {
		// Have an explicit timezone
		t, err = time.Parse(fs+" MST", ds)
		if err == nil {
			return
		}

		// Try to parse without a time zone.
		t, err = time.ParseInLocation(fs, ds, easternTimeZone)
		if err == nil {
			return
		}
	}

	for _, fs := range timeswithzones {
		// These always have an explicit timezone.
		t, err = time.Parse(fs, ds)
		if err == nil {
			return
		}
	}

	// log.Print("Invalid time string ", ds, "\n")
	return
}
