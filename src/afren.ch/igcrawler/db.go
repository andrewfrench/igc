package main

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"os"
	"fmt"
)

var connection *sql.DB

type Correlation struct {
	Association string
	Ratio float64
}

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

func storeMediaId(id string) {
	if !mediaExists(id) { execute("INSERT INTO MEDIA (ID, ADDED, VISITED) VALUES ($1, NOW(), FALSE);", id) }
}

func mediaExists(id string) bool {
	row := queryRow("SELECT COUNT(*) FROM MEDIA WHERE ID = $1;")

	var count int
	row.Scan(&count)

	return count > 1
}

func getUnvisitedMediaId() (string, bool) {
	row := queryRow("UPDATE MEDIA SET VISITED = TRUE WHERE ID = (SELECT ID FROM MEDIA WHERE VISITED = FALSE LIMIT 1) RETURNING ID;")

	var id string
	row.Scan(&id)

	return id, id == ""
}

func userExists(id string) bool {
	row := queryRow("SELECT COUNT(*) FROM USERS WHERE ID = $1;", id)

	var count int
	row.Scan(&count)

	return count > 0
}

func storeUserId(id string) {
	if !userExists(id) { execute("INSERT INTO USERS (ID, ADDED, VISITED) VALUES ($1, NOW(), FALSE);", id) }
}

func getUser() (string, bool) {
	row := queryRow("UPDATE USERS SET VISITED = TRUE WHERE ID = (SELECT ID FROM USERS WHERE VISITED = FALSE LIMIT 1) RETURNING ID;")

	var id string
	row.Scan(&id)

	return id, id == ""
}

func numUnvisitedUsers() int {
	row := queryRow("SELECT COUNT(*) FROM USERS WHERE VISITED = FALSE;")

	var count int
	row.Scan(&count)

	return count
}

func numUnvisitedMedia() int {
	row := queryRow("SELECT COUNT(*) FROM MEDIA WHERE VISITED = FALSE;")

	var count int
	row.Scan(&count)

	return count
}

func getTagCorrelations(tag string) *[]Correlation {
	tagId := getTagId(tag)

	rows := queryRows("SELECT HASHTAG.NAME, (CAST (CORRELATIONS.COUNT AS FLOAT)) / (CAST (HASHTAGS.COUNT AS FLOAT)) AS RATIO FROM HASHTAGS, CORRELATIONS WHERE CORRELATIONS.ASSOCIATION_ID = HASHTAGS.ID WHERE CORRELATIONS.BASE_ID = $1 AND CORRELATIONS.COUNT > 1 ORDER BY RATIO DESC, CORRELATIONS.COUNT DESC LIMIT 20;", tagId)

	correlations := []Correlation{}

	for rows.Next() {
		correlation := Correlation{}

		rows.Scan(&correlation.Association, &correlation.Ratio)

		correlations = append(correlations, correlation)
	}

	return &correlations
}

func numTags() int {
	row := queryRow("SELECT COUNT(*) FROM HASHTAGS;")

	var count int
	row.Scan(&count)

	return count
}

func numCorrelations() int {
	row := queryRow("SELECT COUNT(*) FROM CORRELATIONS;")

	var count int
	row.Scan(&count)

	return count
}

func tagExists(tag string) bool {
	row := queryRow("SELECT COUNT(*) FROM HASHTAGS WHERE NAME = $1;", tag)

	var count int
	row.Scan(&count)

	return count > 0
}

func getTagId(tag string) int {
	row := queryRow("SELECT ID FROM HASHTAGS WHERE NAME = $1;", tag)

	var id int
	row.Scan(&id)

	return id
}

func addTag(tag string) {
	execute("INSERT INTO HASHTAGS (NAME, COUNT, ADDED) VALUES ($1, $2, NOW());", tag, 1)
}

func incrementTag(tag string) {
	tagId := getTagId(tag)

	execute("UPDATE HASHTAGS SET COUNT = COUNT + 1 WHERE ID = $1;", tagId)
}

func associationExists(base, assoc string) bool {
	baseId := getTagId(base)
	assocId := getTagId(assoc)

	row := queryRow("SELECT COUNT(*) FROM CORRELATIONS WHERE BASE_ID = $1 AND ASSOCIATION_ID = $2;", baseId, assocId)

	var count int
	row.Scan(&count)

	return count > 0
}

func incrementAssociation(base, assoc string) {
	baseId := getTagId(base)
	assocId := getTagId(assoc)

	execute("UPDATE CORRELATIONS SET COUNT = COUNT + 1, UPDATED = NOW() WHERE BASE_ID = $1 AND ASSOCIATION_ID = $2;", baseId, assocId)
}

func addAssociation(base, assoc string) {
	baseId := getTagId(base)
	assocId := getTagId(assoc)

	execute("INSERT INTO CORRELATIONS (BASE_ID, ASSOCIATION_ID, COUNT, UPDATED) VALUES ($1, $2, $3, NOW());", baseId, assocId, 1)
}

func queryRow(q string, a ...interface{}) *sql.Row {
	return connection.QueryRow(q, a...)
}

func queryRows(q string, a ...interface{}) *sql.Rows {
	rows, err := connection.Query(q, a...)
	if err != nil { log.Printf("Error querying database: %s", err.Error()) }

	return rows
}

func execute(q string, a ...interface{}) *sql.Result {
	res, err := connection.Exec(q, a...)
	if err != nil { log.Printf("Could not execute SQL: %s\n", err.Error()) }

	return &res
}

func insertCorrelation(base, assoc string) {
	execute("INSERT INTO H (BASE, ASSOC, AT) VALUES ($1, $2, NOW());", base, assoc)
}

func corrTableTotalCount() int {
	row := queryRow("SELECT COUNT(*) FROM H;")

	var count int
	row.Scan(&count)

	return count
}
