package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
)

func main() {
	resp, err := http.Get("https://api.instagram.com/v1/users/self/media/recent?access_token=1608698786.e029fea.bd5972b10ca645efb0422d94a7af8801")
	if err != nil { log.Fatal("couldn't do that") }

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}
