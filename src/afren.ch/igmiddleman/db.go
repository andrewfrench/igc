package main

import (
	"os"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"encoding/json"
)

var connection *sql.DB

func init() {
	dbUser := os.Getenv("DB_USER")
	if len(dbUser) == 0 { log.Fatal("Unable to get DB_USER from environment") }

	dbPass := os.Getenv("DB_PASS")
	if len(dbPass) == 0 { log.Fatal("Unable to get DB_PASS from environment") }

	dbName := os.Getenv("DB_NAME")
	if len(dbName) == 0 { log.Fatal("Unable to get DB_NAME from environment") }

	var err error
	connection, err = sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			dbUser,
			dbPass,
			dbName))

	if err != nil { log.Fatalf("Error connecting to database: %s", err.Error()) }
}

func rawSet() (string, *c_set) {
	rows, err := connection.Query("DELETE FROM H WHERE BASE = (SELECT BASE FROM H GROUP BY BASE, ID LIMIT 1) RETURNING BASE, ASSOC;")
	if err != nil { log.Fatalf("Failed to query asociations: %s", err.Error()) }

	var base string
	var count int
	set := map[string]int{}

	for rows.Next() {
		var assoc string
		rows.Scan(&base, &assoc)

		set[assoc]++
		count++
	}

	return base, &c_set{Count: count, Set: set}
}

func cookedSet(base string) *c_set {
	row := connection.QueryRow("SELECT ASSOCS FROM A WHERE BASE = $1;", base)

	var setJson string
	row.Scan(&setJson)

	set := c_set{}
	set.Set = map[string]int{}

	if len(setJson) > 0 {
		err := json.Unmarshal([]byte(setJson), &set)
		if err != nil { log.Fatalf("Unable to unmarshal JSON: %s", err.Error()) }
	}

	return &set
}

func setExists(base string) bool {
	row := connection.QueryRow("SELECT COUNT(*) FROM A WHERE BASE = $1;", base)

	var count int
	row.Scan(&count)

	return count > 0
}

func insertSet(base string, set *c_set) {
	setJson, err := json.Marshal(*set)

	if err != nil { log.Fatalf("Unable to marshal set: %s", err.Error()) }

	if setExists(base) {
		_, err = connection.Exec("UPDATE A SET ASSOCS = $2, UPDATED = NOW() WHERE BASE = $1;", base, setJson)
		if err != nil { log.Fatalf("Failed to update assocs: %s", err.Error()) }

		fmt.Printf("Update %s\n", base)
	} else {
		_, err = connection.Exec("INSERT INTO A (BASE, ASSOCS, UPDATED) VALUES ($1, $2, NOW());", base, setJson)
		if err != nil { log.Fatalf("Failed to insert into assocs: %s", err.Error()) }

		fmt.Printf("Create %s\n", base)
	}
}