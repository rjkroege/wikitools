//go:build darwin
// +build darwin

package search

import (
	"fmt"
	"log"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/dispatch"
//	"github.com/progrium/macdriver/objc"
)

// spotlightWikilinkIndexer hides all of the darwin-specific code needed
// to support indexing using spotlight.
// Changes to this structure need to be synchronized correctly.
type spotlightWikilinkIndexer struct {
}

// TODO(rjk): put the waiter structure inside the timer function.
// There are subtle rules about queues (and custom queues) that I will need to learn.
func (_ *spotlightWikilinkIndexer) Allpaths(wikitext string)  ([]string, error) {
	log.Println("starting a timer off runloop")
	// Post a message to the runloop
	waiterchan := make(chan int)

	// Post task to mainqueue:
	q := dispatch.MainQueue()
	q.DispatchAsync(func() {
		timer(waiterchan, 1.5)
	})

	<-waiterchan
	return nil, fmt.Errorf("spotlightWikilinkIndexer not implemented")
}

func MakeWikilinkNameIndex() corpus.WikilinkNameIndex {
	log.Println("hi from darwin")

	spidx := &spotlightWikilinkIndexer{}
	return spidx
}


// timer can be run from an arbitrary Go routine to register a foundation
// timer that will signal provided channel fired after time seconds.
// TODO(rjk): This pattern might be amenable to something more
// more sophisticated.
func timer(fired chan int, duration float64) {
		// Start a timer. This method seems to work. The first method that I
		// tried did not function successfully. But this example works.
		foundation.Timer_ScheduledTimerWithTimeIntervalRepeatsBlock(foundation.TimeInterval(duration), false, func(timer foundation.Timer) {
			log.Println("timer fired!, spiffy!", timer)
			fired <- 1
		})	
}

