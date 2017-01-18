package main

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"net/http"
	"strings"
	"fmt"
	"time"
	"encoding/json"
	"os"
)

var correlationsSetChannel chan *map[string]float64

func main() {
	http.HandleFunc("/discover", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving correlations")

		query := r.URL.Query().Get("tags")
		tagSlice := strings.Split(query, ",")

		ignore := r.URL.Query().Get("ignore")
		ignoreSlice := strings.Split(ignore, ",")

		correlationsList := handleQuery(tagSlice, ignoreSlice)

		respJson, err := json.Marshal(correlationsList)

		if err != nil { log.Fatal(err.Error()) }

		w.Write(respJson)
	})

	http.ListenAndServe(":8080", nil)
}

func handleQuery(tagSlice, ignoreSlice []string) []string {
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
		*correlationList = append(*correlationList, Correlation{
			Association: k,
			Ratio: v,
		})
	}

	begin := time.Now()
	sorted := sortCorrelations(correlationList)
	dur := time.Since(begin)
	fmt.Printf("Sort took %s\n", dur.String())

	printed := 0
	hidden := []string{}
	for _, c := range *sorted {
		ignore := ignoreSet[c.Association] || querySet[c.Association]
		if !ignore {
			if printed == 20 { break }
			fmt.Printf("\t%.3f: %s\n", c.Ratio, c.Association)
			hidden = append(hidden, c.Association)
			printed++
		}
	}

	return hidden
}

var connection *sql.DB

type Correlation struct {
	Association string
	Ratio float64
}

type c_set struct {
	Count int            `json:"count"`
	Set   map[string]int `json:"set"`
}

func init() {
	dbUser := os.Getenv("DB_USER")
	if len(dbUser) == 0 { log.Fatal("Unable to get DB_USER from environment") }

	dbPass := os.Getenv("DB_PASS")
	if len(dbPass) == 0 { log.Fatal("Unable to get DB_PASS from environment") }

	dbName := os.Getenv("DB_NAME")
	if len(dbName) == 0 { log.Fatal("Unable to get DB_NAME from environment") }

	var err error
	connection, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName))
	if err != nil { log.Fatalf("Error connecting to database: %s", err.Error()) }
}

func getTagCorrelations(tag string) {
	begin := time.Now()

	row := queryRow("SELECT ASSOCS FROM A WHERE BASE = $1;", tag)

	var setJson string
	row.Scan(&setJson)

	corrs := map[string]float64{}

	if len(setJson) > 0 {
		assocs := c_set{}
		err := json.Unmarshal([]byte(setJson), &assocs)
		if err != nil { log.Fatalf("Failed to unmarshal JSON: %s", err.Error()) }

		for k, v := range assocs.Set {
			corrs[k] = float64(v) / float64(assocs.Count)
		}


		dur := time.Since(begin)
		fmt.Printf("%s query took %s\n", tag, dur.String())
	} else {
		fmt.Printf("No results for %s\n", tag)
	}

	// Add empty correlation set to channel anyway, since we block based on how many goroutines have
	// been launched.
	correlationsSetChannel <- &corrs
}

func queryRow(q string, a ...interface{}) *sql.Row {
	return connection.QueryRow(q, a...)
}

func queryRows(q string, a ...interface{}) *sql.Rows {
	rows, err := connection.Query(q, a...)
	if err != nil { log.Printf("Error querying database: %s", err.Error()) }

	return rows
}
