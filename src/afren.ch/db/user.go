package db

import (
	"log"
	"afren.ch/env"
)

func init() {
	userSeed := env.OptionalString("USER_SEED", "")

	if len(userSeed) == 0 { return }

	if !UserExists(userSeed) {
		log.Printf("Storing user seed: %s", userSeed)
		StoreUser(userSeed, "")
	} else {
		log.Printf("User seed already exists: %s", userSeed)
	}
}

func GetUserId() string {
	row := queryRow(queries.GetNewUserId)

	var id string
	row.Scan(&id)

	return id
}

func StoreUser(id, username string) {
	execute(queries.StoreUserId, id, username)
}

func NewUsersCount() int {
	row := queryRow(queries.NewUserCount)

	var count int
	row.Scan(&count)

	return count
}

func NewUsersExist() bool {
	return NewUsersCount() > 0
}

func UserExists(id string) bool {
	row := queryRow(queries.UserExists, id)

	var count int
	row.Scan(&count)

	return count > 0
}
