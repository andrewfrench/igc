package users

import (
	"testing"
	"io/ioutil"
)

func TestUsers(t *testing.T) {
	// Setup
	validData, _ := ioutil.ReadFile("users_response_valid.json")
	if len(validData) == 0 { t.Fatal("Unable to load valid user response sample") }

	invalidData, _ := ioutil.ReadFile("users_response_invalid.json")
	if len(invalidData) == 0 { t.Fatal("Unable to load invalid user response sample") }

	emptyData := []byte{}

	// Test cases
	t.Run("Empty", func(t *testing.T) {
		createUsersResponseStruct(emptyData)
	})

	t.Run("RealData", func(t *testing.T) {
		ur := createUsersResponseStruct(validData)
		if ur.Meta.Code != 200 { t.Fail() }
	})

	t.Run("FakeData", func(t *testing.T) {
		createUsersResponseStruct(invalidData)
	})
}
