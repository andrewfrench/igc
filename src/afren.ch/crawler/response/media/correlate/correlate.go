package correlate

import (
	"afren.ch/db"
)

func Correlate(tags []string) int {
	if len(tags) < 2 {
		return 0
	}

	var count int
	for i, base := range tags {
		for _, assoc := range tags[i:] {
			if base == assoc {
				continue
			}

			db.StoreCorrelationPair(base, assoc)
			db.StoreCorrelationPair(assoc, base)
			count += 2
		}
	}

	return count
}
