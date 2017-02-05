package media

import (
	"encoding/json"
	"afren.ch/db"
	"log"
	"afren.ch/crawler/response/media/correlate"
)

type MediaResponse struct {
	Meta Meta `json:"meta"`
	Media []Media `json:"data"`
}

type Meta struct {
	Code int `json:"code"`
}

type Media struct {
	Caption string `json:"caption"`
	Tags []string `json:"tags"`
	Location string `json:"location"`
	Likes Likes `json:"likes"`
	Id string `json:"id"`
}

type Likes struct {
	Count int `json:"count"`
}

func Handle(data []byte) error {
	mr := createMediaResponseStruct(data)

	saveMedia(mr.Media)

	var correlationCount int
	for _, m := range mr.Media {
		correlationCount += correlate.Correlate(m.Tags)
	}

	log.Printf("Handled %d media", len(mr.Media))
	log.Printf("Handled %d correlations", correlationCount)

	return nil
}

func createMediaResponseStruct(data []byte) *MediaResponse {
	mr := new(MediaResponse)

	json.Unmarshal(data, mr)

	return mr
}

func saveMedia(media []Media) {
	for _, m := range media {
		db.StoreMediaId(m.Id)
	}
}