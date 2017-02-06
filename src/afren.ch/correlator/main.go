package main

import (
	"time"
	"afren.ch/set"
	"afren.ch/db"
	"log"
)

func main() {
	for true {
		incoming := set.PullNew()

		if len(incoming.GetBase()) == 0 {
			// Give crawler time to add more correlations and try again
			time.Sleep(10 * time.Second)
			continue
		}

		if db.AssociationBaseExists(incoming.GetBase()) {
			existing, err := set.PullExisting(incoming.GetBase())
			if err != nil {
				log.Printf("Error pulling existing set: %s", err.Error())
				existing = set.EmptySet()
			}

			incoming = set.Merge(incoming, existing)

			log.Printf("Updated %s", incoming.GetBase())
		} else {
			log.Printf("Created %s", incoming.GetBase())
		}

		incoming.Save()
	}
}
