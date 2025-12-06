package actions

import (
	"log"
	"strings"
	"time"

	"9fans.net/go/acme"
)

type AcmeWindow struct {
	ID   int
	Name string
}

type AcmeWatcher struct {
	add      chan AcmeWindow
	remove   chan int
	snapshot chan chan []AcmeWindow
	winds map[int]AcmeWindow
}

func NewAcmeWatcher(wikiroot string) *AcmeWatcher {
	log.Println("watcher starting")
	wm := &AcmeWatcher{
		add:      make(chan AcmeWindow),
		remove:   make(chan int),
		snapshot: make(chan chan []AcmeWindow),
	}
	go wm.run()
	go retryloop(wikiroot)
	return wm
}

func (wm *AcmeWatcher) run() {
	am := make(map[int]AcmeWindow)

	for {
		select {
		case win := <-wm.add:
			am[win.ID] = win
		case id := <-wm.remove:
			delete(am, id)
		case resp := <-wm.snapshot:
			snapshot := make([]AcmeWindow, 0, len(am))
			for _, win := range am {
				snapshot = append(snapshot, win)
			}
			resp <- snapshot
		}
	}
}

func (wm *AcmeWatcher) Add(win AcmeWindow) {
	wm.add <- win
}

func (wm *AcmeWatcher) Remove(id int) {
	wm.remove <- id
}

func (wm *AcmeWatcher) Snapshot() []AcmeWindow {
	resp := make(chan []AcmeWindow)
	wm.snapshot <- resp
	return <-resp
}

func retryloop(wikiroot string) {
	backoff := 1 * time.Second
	maxBackoff := 60 * time.Second

	for {
		if err := readwindows(wikiroot); err != nil {
			log.Printf("readwindows failed, retrying in %v", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		if err := watchacmelog(wikiroot); err != nil {
			log.Printf("watchacmelog failed, retrying in %v", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		backoff = 1 * time.Second
	}
}

func readwindows(wikiroot string) error  {
              wins, err := acme.Windows()
               if err != nil {
                    log.Printf("can't get acme windows; probably acme is not running: %v", err)
			return err
               }

               // TODO(rjk): acme.Windows might not correctly handle Name instances
               for _, w := range wins {
			if strings.HasPrefix(w.Name, wikiroot) {
				log.Println("w.Name", w.Name)
				// TODO(rjk) do tag update for each of these things
			}
               }

	return nil
}

func watchacmelog(wikiroot string) error {
	r, err := acme.Log()
	if err != nil {
		log.Printf("can't open acme; probably it's not running: %v", err)
		return err
	}
	defer r.Close()

		for {
			ev, err := r.Read()
			if err != nil {
				log.Printf("can't read events from acme; probably it's not running: %v", err)
			return err
				}

			if strings.HasPrefix(ev.Name, wikiroot) {
				log.Println("watching... got an event", ev)
			}
		}
}

