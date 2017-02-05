package users

import (
	"encoding/json"
	"afren.ch/db"
	"log"
)

type UsersResponse struct {
	Meta Meta
	Data []Data
}

type Meta struct {
	Code int `json:"code"`
}

type Data struct {
	Id string
	Username string
}

func Handle(data []byte) error {
	ur := createUsersResponseStruct(data)

	saveUsers(ur.Data)

	log.Printf("Handled %d users", len(ur.Data))

	return nil
}

func createUsersResponseStruct(data []byte) *UsersResponse {
	ur := new(UsersResponse)

	json.Unmarshal(data, ur)

	return ur
}

func saveUsers(users []Data) {
	for _, user := range users {
		db.StoreUser(user.Id, user.Username)
	}
}
