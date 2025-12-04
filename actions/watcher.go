package actions

import (
	"log"

	"9fans.net/go/acme"
)

func WatchAcmeLog() {
	log.Println("watcher starting")

	go retryloop()
}

func retryloop() {
	readwindows()
	watchacmelog()
}

func readwindows() error  {
              wins, err := acme.Windows()
               if err != nil {
                    log.Printf("can't get acme windows; probably acme is not running: %v", err)
			return err
               }

               // TODO(rjk): acme.Windows might not correctly handle Name instances
               for _, w := range wins {
			log.Println("w.Name", w.Name)
			// TODO(rjk) do tag update for each of these things
               }

	return nil
}

func watchacmelog() error {
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


			log.Println("watching... got an event", ev)
		}
}
