//go:build darwin
// +build darwin

package search

import (
	"fmt"
	"log"

	"github.com/progrium/macdriver/dispatch"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
	"github.com/rjkroege/wikitools/corpus"
)

// spotlightWikilinkIndexer hides all of the darwin-specific code needed
// to support indexing using spotlight.
// Changes to this structure need to be synchronized correctly.
type spotlightWikilinkIndexer struct {
}

// There are subtle rules about queues (and custom queues) that I will need to learn.
func (_ *spotlightWikilinkIndexer) Allpaths(wikitext string) ([]string, error) {
	log.Println("Allpaths: the goroutine")
	waiterchan := make(chan foundation.MetadataQuery)
	qs := fmt.Sprintf("kMDItemFSName == '%s'", wikitext)
	predicate := foundation.Predicate_PredicateFromMetadataQueryString(qs)

	// Post task to runloop.
	q := dispatch.MainQueue()
	q.DispatchAsync(func() {
		log.Println("Allpaths: on runloop")
		// Create new query.
		// TODO(rjk): Could I create this outside?
		// TODO(rjk): Can I run the query from an arbitrary thread?
		// TODO(rjk): Can I get the notification back on a different thread?
		query := foundation.NewMetadataQuery().Init()
		// query persists beyond a single event cycle and therefore (I believe) needs
		// to be retained.
		objc.Retain(&query)

		// SetPredicate has the @property(copy) so predicate is copied here.
		query.SetPredicate(predicate)

		nc := foundation.NotificationCenter_DefaultCenter()
		var token objc.IObject
		token = nc.AddObserverForNameObjectQueueUsingBlock(
			foundation.MetadataQueryDidFinishGatheringNotification,
			nil,
			foundation.OperationQueue_CurrentQueue(),
			func(notification foundation.Notification) {
				// This runs on the runloop thread.
				log.Println("Allpaths sez finished gathering on runloop!")
				nc.RemoveObserver(token)
				query.StopQuery()

				// I am passing responsibility for cleanup to a different thread along
				// with the object.
				waiterchan <- query
			},
		)
		// TODO(rjk): Perhaps I need to check some kind of return code.
		// There can only be one of these in flight at a time.
		query.StartQuery()
	})

	// TODO(rjk): I might need to return some kind of status. Things can go wrong.
	// TODO(rjk): generate actual output here. The "best" path finding is probably common code.
	afterQueryDone(<-waiterchan)

	return nil, fmt.Errorf("spotlightWikilinkIndexer not implemented")
}

func MakeWikilinkNameIndex() corpus.WikilinkNameIndex {
	// The Apple docs imply (very strongly) that there can only be a single
	// query running at a time.
	spidx := &spotlightWikilinkIndexer{}
	return spidx
}

// afterQueryDone is code to run on the goroutine to process the results from the query.
// TODO(rjk): I am assuming that nothing here (i.e. methods on query) need to run on
// the runloop thread.
func afterQueryDone(query foundation.MetadataQuery) {
	defer query.Release()
	rc := query.ResultCount()
	for i := 0; uint(i) < rc; i++ {
		// I dislike this syntax but I disassembled the output and it's
		// effectively a nop. I believe that I can just fold this together.
		mdi := &foundation.MetadataItem{query.ResultAtIndex(uint(i))}
		values := mdi.ValuesForAttributes(mdi.Attributes())

		for k, v := range values {
			// Can refer to the embedded class data like this. I think that this
			// pattern is generally right.
			if v.IsKindOfClass(foundation.StringClass.Class) {
				// This works. Is there a nicer way? I'm creating an additional
				// pointer sized object to note the type of something? Couldn't
				// we have generic APIs?
				s := foundation.String{v}
				log.Println(k, s.Description())
			} else {
				log.Println(k, "not a string")
			}
		}

		// The keys are just strings. However, they do not appear to be defined in progrium.
		// I can just define the path and pass it in. See MDItem.h for the actual string values
		// of the keys.
		//
		// The construction of the foundation.String is unnecessary because I'm
		// going to turn around and just pass it to ToGoString.
		s := mdi.ValueForAttribute("kMDItemPath")

		// This technique for converting the string works and is code-minimizing.
		// I think that this is as close to "toll-free" bridging as possible. The
		// Go documentation says that the underlying implementation will
		// duplicate the string. Conclusion: this is the right way to implement
		// getting a string value from an NSString instance.
		log.Println("path?", objc.ToGoString(s.Ptr()))
	}
}
