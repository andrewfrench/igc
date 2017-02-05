package media

import (
	"afren.ch/db"
	"log"
	"fmt"
	"net/http"
	"io/ioutil"
	"afren.ch/crawler/response/media"
	"errors"
	"afren.ch/env"
)

type UsersMediaRequest struct {
	id string
}

func Create() *UsersMediaRequest {
	r := new(UsersMediaRequest)

	r.id = db.GetUserId()

	return r
}

func (r *UsersMediaRequest) Url() string {
	accessToken := env.RequiredString("ACCESS_TOKEN")
	return fmt.Sprintf("https://api.instagram.com/v1/users/%s/media/recent?access_token=%s", r.id, accessToken)
}

func (r *UsersMediaRequest) DoRequest() error {
	log.Printf("Requesting user media: %s", r.id)

	resp, err := http.Get(r.Url())
	if err != nil { return err }
	if resp.StatusCode != 200 { return errors.New(fmt.Sprintf("Request failed with error: %d", resp.StatusCode)) }

	body, err := ioutil.ReadAll(resp.Body)
	media.Handle(body)

	return err
}