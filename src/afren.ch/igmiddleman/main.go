package main

import "time"

type c_set struct {
	Count int		`json:"count"`
	Set   map[string]int	`json:"set"`
}

func main() {
	for true {
		base, new := rawSet()

		if len(base) == 0 {
			// Give crawler time to add more correlations and try again
			time.Sleep(10 * time.Second)
			continue
		}

		var old *c_set
		if setExists(base) {
			old = cookedSet(base)
		} else {
			old = &c_set{}
			old.Set = map[string]int{}
		}

		old.Count += new.Count
		for a, c := range (*new).Set {
			(*old).Set[a] += c
		}

		insertSet(base, old)
	}
}
