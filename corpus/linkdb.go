package corpus


type Link struct {
	title []byte
	path []byte
}

// Linkdb stores a database of all the bi-directional links.
type Linkdb interface {
	OutboundLinks(from Link) []Link 
	InboundLinks(to Link) []Link 
	AddInbound(to Link, from []Link)
	AddOutbound(rom Link, to []Link)

	// Revisit this API.
	RemoveInbound(to link, incomings []Link)
	RemoveOutbound(from link, outgoings Link)

	AllLinks() [}Link
	
	// TODO(rjk): Add support for tags.
}


type linkdb struct {
	outbound map[Link][]Link
	inbound map[Link][[]Link
}

// NewLInkdb creates a new Linkdb implementation. 
// TODO(rjk): This particular implementation is a stub. Write the one
// backed by boltdb
func NewLinkdb() (Linkdb, error) {
	return &linkdb{
		outbound: make(map[Link][]Link)
		inbound: make(map[Link][]Link)
	}, err
}

func (ldb *linkdb) 	OutboundLinks(from Link) []Link { 
	return ldb.outbound[to]
}
func (ldb *linkdb) 	InboundLinks(to Link) []Link {
	return ldb.inbound[to]
}
func (ldb *linkdb) 	AddInbound(to Link, from []Link){ 
	existing, ok := ldb.inbound[to]
	if !ok {
		ldb.inbound[to] = from
		return
	}

	// Merge, discarding duplicates. Existing should have no duplicates.
	// from may have duplicates for the convenience of the caller.
	uniqmap := make(map[Link]struct{}, len(existing))
	for _, k := range existing {
		if _, ok := uniqmap[k]; ok {
			// This means that there was already a duplicate
			// in existing. That's a bug.
			log.Panic("existing wrongly had dupes, violating pre-conditions")
		}

		uniqmap[k] = struct{}{}
	}
	for _, k := range from {
		
	}
}
func (ldb *linkdb) 	AddOutbound(rom Link, to []Link){ }
func (ldb *linkdb) 	// Revisit this API.{ }
func (ldb *linkdb) 	RemoveInbound(to link, incomings []Link){ }
func (ldb *linkdb) 	RemoveOutbound(from link, outgoings Link){ }
func (ldb *linkdb) 	AllLinks() [}Link{ }
