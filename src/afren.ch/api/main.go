package main

import (
	"log"
	"net/http"
	"strings"
	"fmt"
	"time"
	"encoding/json"
	"afren.ch/cache"
	"afren.ch/correlator/set"
	"afren.ch/env"
)

var correlationsSetChannel chan *map[string]float64

type tagPair struct {
	Hashtag   string  `json:"hashtag"`
	Relevance float64 `json:"relevance"`
}

type responseStruct struct {
	ResultCount int       `json:"count"`
	TagSet      []tagPair `json:"set"`
	Time        string    `json:"time"`
}

func main() {
	log.Print("Starting server")

	http.HandleFunc("/discover", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("tags")
		tagSlice := strings.Split(query, ",")

		ignore := r.URL.Query().Get("ignore")
		ignoreSlice := strings.Split(ignore, ",")

		beginTime := time.Now()
		correlationsList, totalCount := handleQuery(tagSlice, ignoreSlice)

		timeDiff := time.Since(beginTime)
		rs := responseStruct{
			ResultCount: totalCount,
			TagSet: correlationsList,
			Time: timeDiff.String(),
		}

		respJson, err := json.Marshal(rs)
		if err != nil { log.Fatal(err.Error()) }

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(respJson)

		log.Printf("Inputs: %s", tagSlice)
		if len(ignoreSlice) > 0 { log.Printf("Ignore: %s", ignoreSlice) }
		log.Printf("Query took: %s", timeDiff.String())
	})

	port := env.RequiredString("PORT");
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func handleQuery(tagSlice, ignoreSlice []string) ([]tagPair, int) {
	correlationsSetChannel = make(chan *map[string]float64)
	querySet := map[string]bool{}
	ignoreSet := map[string]bool{}
	correlations := map[string]float64{}

	for _, word := range tagSlice {
		querySet[word] = true
	}

	for _, word := range ignoreSlice {
		ignoreSet[word] = true
	}

	for k := range querySet {
		go getTagCorrelations(k)
	}

	for range querySet {
		cs := <- correlationsSetChannel

		for a, c := range *cs {
			correlations[a] += c
		}
	}

	correlationList := new([]Correlation)
	for k, v := range correlations {
		ignore := ignoreSet[k] || querySet[k]
		if !ignore {
			*correlationList = append(*correlationList, Correlation{
				Association: k,
				Ratio: v,
			})
		}
	}

	sorted := sortCorrelations(correlationList)

	var maxRelevance float64

	if len(*sorted) > 0 {
		maxRelevance = (*sorted)[0].Ratio
	}

	returned := 0
	returnSet := []tagPair{}
	for _, c := range *sorted {
		if returned == 20 { break }
		returnSet = append(returnSet, tagPair{
			Hashtag: c.Association,
			Relevance: c.Ratio / maxRelevance,
		})
		returned++
	}

	return returnSet, len(*sorted)
}

type Correlation struct {
	Association string
	Ratio float64
}

func getTagCorrelations(base string) {
	var assocs *set.CorrelationSet

	// Check the local cache to see if we can avoid a database query
	if cache.Contains(base) {
		assocs = cache.Get(base)
	} else {
		var err error
		assocs, err = set.PullExisting(base)
		if err != nil { log.Printf("Failed to pull from database: %s", base) }

		// Add this to cache while we have it
		cache.Insert(assocs)
	}

	corrs := map[string]float64{}

	for k, v := range assocs.AssocMap {
		corrs[k] = float64(v) / float64(assocs.Count)
	}

	correlationsSetChannel <- &corrs
}
