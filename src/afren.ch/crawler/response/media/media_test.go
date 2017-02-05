package media

import (
	"testing"
	"io/ioutil"
)

func TestMedia(t *testing.T) {
	// Setup
	validData, _ := ioutil.ReadFile("media_response_valid.json")
	if len(validData) == 0 { t.Fatal("Unable to load valid media response sample") }

	invalidData, _ := ioutil.ReadFile("media_response_invalid.json")
	if len(invalidData) == 0 { t.Fatal("Unable to load invalid media response sample") }

	emptyData := []byte{}

	// Test cases
	t.Run("Empty", func(t *testing.T) {
		createMediaResponseStruct(emptyData)
	})

	t.Run("ValidData", func(t *testing.T) {
		mr := createMediaResponseStruct(validData)
		if mr.Meta.Code != 200 { t.Fail() }
	})

	t.Run("InvalidData", func(t *testing.T) {
		createMediaResponseStruct(invalidData)
	})
}
