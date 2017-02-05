package request

import (
	"afren.ch/db"
	"afren.ch/crawler/request/users"
	"afren.ch/crawler/request/media"
	"log"
	"time"
)

type Request interface {
	Url() string
	DoRequest() error
}

func GetRequest() Request {
	var req Request

	if db.NewUsersExist() {
		req = media.Create()
	} else if db.NewMediaExists() {
		log.Print("No new media exist, checking for new users")
		req = users.Create()
	} else {
		log.Fatal("Exhausted new user and media queues")
	}

	return req
}

func RequestLoop() {
	minInterval := 8 * time.Second

	for true {
		begin := time.Now()

		req := GetRequest()
		req.DoRequest()

		diff := time.Since(begin)

		if diff < minInterval {
			remaining := minInterval - diff

			log.Printf("Sleeping for %s", remaining.String())
			time.Sleep(remaining)
		}
	}
}
