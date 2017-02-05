package users

import (
	"fmt"
	"os"
	"log"
	"afren.ch/db"
	"net/http"
	"io/ioutil"
	"afren.ch/crawler/response/users"
	"errors"
	"afren.ch/env"
)

type LikingUsersRequest struct {
	id string
}

func Create() *LikingUsersRequest {
	r := new(LikingUsersRequest)

	r.id = db.GetMediaId()

	return r
}

func (r *LikingUsersRequest) Url() string {
	accessToken := env.RequiredString("ACCESS_TOKEN")
	return fmt.Sprintf("https://api.instagram.com/v1/media/%s/likes?access_token=%s", r.id, accessToken)
}

func (r *LikingUsersRequest) DoRequest() error {
	log.Printf("Requesting media likes: %s", r.id)

	resp, err := http.Get(r.Url())
	if err != nil { return err }
	if resp.StatusCode != 200 { return errors.New(fmt.Sprintf("Request failed with error: %d", resp.StatusCode)) }

	body, err := ioutil.ReadAll(resp.Body)
	users.Handle(body)

	return err
}
