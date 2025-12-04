package actions

import (
	"log"
	"strings"
	"time"

	"9fans.net/go/acme"
)

func WatchAcmeLog(wikiroot string) {
	log.Println("watcher starting")

	go retryloop(wikiroot)
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
