package db

import (
	"log"
	"afren.ch/env"
)

func init() {
	mediaSeed := env.OptionalString("MEDIA_SEED", "")

	if len(mediaSeed) == 0 { return }
	if !MediaExists(mediaSeed) {
		log.Printf("Storing media seed: %s", mediaSeed)
		StoreMediaId(mediaSeed)
	} else {
		log.Printf("Media seed already exists: %s", mediaSeed)
	}
}

func GetMediaId() string {
	row := queryRow(queries.GetNewMediaId)

	var id string
	row.Scan(&id)

	return id
}

func StoreMediaId(id string) {
	execute(queries.StoreMediaId, id)
}

func NewMediaCount() int {
	row := queryRow(queries.NewMediaCount)

	var count int
	row.Scan(&count)

	return count
}

func NewMediaExists() bool {
	return NewMediaCount() > 0
}

func MediaExists(id string) bool {
	row := queryRow(queries.MediaExists, id)

	var count int
	row.Scan(&count)

	return count > 0
}
