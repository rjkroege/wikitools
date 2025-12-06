package actions

import (
	"log"
	"strings"
	"time"

	"9fans.net/go/acme"
	"github.com/rjkroege/wikitools/corpus"
)

type AcmeWatcher struct {
	add      chan corpus.AcmeWindow
	remove   chan int
	snapshot chan chan []corpus.AcmeWindow
	winds map[int]corpus.AcmeWindow
}

func NewAcmeWatcher(wikiroot string) *AcmeWatcher {
	log.Println("watcher starting")
	wm := &AcmeWatcher{
		add:      make(chan corpus.AcmeWindow),
		remove:   make(chan int),
		snapshot: make(chan chan []corpus.AcmeWindow),
	}
	go wm.run()
	go retryloop(wm, wikiroot)
	return wm
}

func (wm *AcmeWatcher) run() {
	am := make(map[int]corpus.AcmeWindow)

	for {
		select {
		case win := <-wm.add:
			am[win.ID] = win
		case id := <-wm.remove:
			delete(am, id)
		case resp := <-wm.snapshot:
			snapshot := make([]corpus.AcmeWindow, 0, len(am))
			for _, win := range am {
				snapshot = append(snapshot, win)
			}
			resp <- snapshot
		}
	}
}

func (wm *AcmeWatcher) Add(win corpus.AcmeWindow) {
	wm.add <- win
}

func (wm *AcmeWatcher) Remove(id int) {
	wm.remove <- id
}

func (wm *AcmeWatcher) Snapshot() []corpus.AcmeWindow {
	resp := make(chan []corpus.AcmeWindow)
	wm.snapshot <- resp
	return <-resp
}

func retryloop(wm *AcmeWatcher, wikiroot string) {
	backoff := 1 * time.Second
	maxBackoff := 60 * time.Second

	for {
		if err := readwindows(wm, wikiroot); err != nil {
			log.Printf("readwindows failed, retrying in %v", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		if err := watchacmelog(wm, wikiroot); err != nil {
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

func readwindows(wm *AcmeWatcher, wikiroot string) error  {
              wins, err := acme.Windows()
               if err != nil {
                    log.Printf("can't get acme windows; probably acme is not running: %v", err)
			return err
               }

               // TODO(rjk): acme.Windows might not correctly handle Name instances
               for _, w := range wins {
			if strings.HasPrefix(w.Name, wikiroot) {
				wm.Add(corpus.AcmeWindow{
					ID: w.ID,
					Name: w.Name,
				})
			}
               }

	return nil
}

func watchacmelog(wm *AcmeWatcher, wikiroot string) error {
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


			if ev.Op == "new" {
				wm.Add(corpus.AcmeWindow{
					ID: ev.ID,
					Name: ev.Name,
				})
			}

			if ev.Op == "del" {
				wm.Remove(ev.ID)
			}


			}
		}
}

