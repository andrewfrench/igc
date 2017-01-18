package main

import (
	"fmt"
	"os"
	"time"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
	"runtime"
)

const (
	RecentMedia = iota
	Likes
)

type RequestBundle struct {
	RequestType int
	RequestUrl string
}

type ResponseBundle struct {
	ResponseType int
	ResponseData interface{}
}

type User struct {
	Name string
	Id string
}

type Media struct {
	Id string
	Url string
}

var usersChannel chan *User
var mediaChannel chan *Media
var correlationsChannel chan *[]string

var accessToken string

func main() {
	log.Println("Crawling Instagram...")

	accessToken = os.Getenv("ACCESS_TOKEN")
	if len(accessToken) == 0 { log.Fatal("Unable to get ACCESS_TOKEN from environment") }

	userSeed := os.Getenv("USER_SEED")
	if len(userSeed) == 0 { log.Fatal("Unable to get USER_SEED from environment") }

	correlationWorkersString := os.Getenv("CORRELATION_WORKERS")
	correlationWorkers, err := strconv.Atoi(correlationWorkersString)
	if err != nil || correlationWorkers < 1 { correlationWorkers = 1 } else {
		log.Printf("Using %d correlation workers on %d CPUs\n", correlationWorkers, runtime.NumCPU())
	}

	// Launch routine to store users
	usersChannel = make(chan *User, 1024)
	go storeUsers()

	// Launch routine to store media
	mediaChannel = make(chan *Media, 1024)
	go storeMedia()

	// Launch routine to calculate and store correlations
	correlationsChannel = make(chan *[]string, 16384)
	for i := 0; i < correlationWorkers; i++ {
		go storeCorrelations()
	}

	// Launch routine that will report crawl state
	go reportCrawlState()

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	tick := time.Tick(8 * time.Second)
	for range tick {
		requestBundle := getRequestBundle()

		resp, err := client.Get(requestBundle.RequestUrl)
		if err != nil { log.Fatalf("Error on request: %s\n", err.Error()) }
		if resp.StatusCode != 200 { continue }

		responseStruct := new(interface{})
		json.NewDecoder(resp.Body).Decode(responseStruct)

		switch requestBundle.RequestType {
		case RecentMedia:
			go handleRecentMedia(responseStruct)

		case Likes:
			go handleLikes(responseStruct)
		}
	}
}

func reportCrawlState() {
	tick := time.Tick(60 * time.Second)
	for range tick {
		log.Printf("%d users,\t %d media,\t %d correlation rows,\t %d goroutines\n", numUnvisitedUsers(), numUnvisitedMedia(), corrTableTotalCount(), runtime.NumGoroutine())
	}
}

func storeUsers() {
	for user := range usersChannel {
		storeUserId(user.Id)
	}
}

func storeMedia() {
	for media := range mediaChannel {
		storeMediaId(media.Id)
	}
}

func storeCorrelations() {
	for csp := range correlationsChannel {
		for i, a := range *csp {
			for _, b := range (*csp)[i:] {
				if a == b { continue }
				insertCorrelation(a, b)
				insertCorrelation(b, a)
			}
		}
	}
}

func getRequestBundle() *RequestBundle {
	// Get a new user to scrape recent posts from
	userId, empty := getUser()

	// Define new request bundle
	requestBundle := RequestBundle{}

	if empty {
		// No more unvisited users
		// Get media ID to request user like list
		mediaId, empty := getUnvisitedMediaId()
		if empty { log.Fatal("Exhausted user and media queue!") }

		requestBundle.RequestType = Likes
		requestBundle.RequestUrl = fmt.Sprintf("https://api.instagram.com/v1/media/%s/likes?access_token=%s", mediaId, accessToken)
	} else {
		requestBundle.RequestType = RecentMedia
		requestBundle.RequestUrl = fmt.Sprintf("https://api.instagram.com/v1/users/%s/media/recent?access_token=%s", userId, accessToken)
	}

	return &requestBundle
}

func handleRecentMedia(data interface{}) {
	mediaList := (*(data.(*interface{}))).(map[string]interface{})["data"].([]interface{})

	for _, media := range mediaList {
		m := media.(map[string]interface{})
		tags := m["tags"].([]interface{})
		id := m["id"].(string)

		// Send media to be added to database
		mediaChannel <- &Media{
			Id: id,
		}

		// Do a type assertion on each element of tags
		tagsAsStrings := []string{}
		for _, t := range tags { tagsAsStrings = append(tagsAsStrings, t.(string)) }

		if len(tagsAsStrings) > 1 { correlationsChannel <- &tagsAsStrings }
	}
}

func handleLikes(data interface{}) {
	likesList := (*(data.(*interface{}))).(map[string]interface{})["data"].([]interface{})
	for _, user := range likesList {
		userElem := user.(map[string]interface{})
		u := User{
			Name: userElem["username"].(string),
			Id: userElem["id"].(string),
		}

		usersChannel <- &u
	}
}
