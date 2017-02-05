package db

import (
	"log"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"afren.ch/env"
)

type Queries struct {
	GetNewUserId         string `json:"get_new_user_id"`
	StoreUserId          string `json:"store_user_id"`
	NewUserCount         string `json:"new_user_count"`
	UserExists           string `json:"user_exists"`

	GetNewMediaId        string `json:"get_new_media_id"`
	StoreMediaId         string `json:"store_media_id"`
	NewMediaCount        string `json:"new_media_count"`
	MediaExists          string `json:"media_exists"`

	StoreCorrelationPair string `json:"store_correlation_pair"`

	GetAssociationCount  string `json:"get_association_count"`
	InsertAssociationSet string `json:"insert_association_set"`
	UpdateAssociationSet string `json:"update_association_set"`
	QueryIncomingSet     string `json:"query_incoming_set"`
	QueryAssociationSet  string `json:"query_association_set"`
}

var connection *sql.DB
var queries *Queries

func init() {
	dbUser := env.RequiredString("DB_USER")
	dbPass := env.RequiredString("DB_PASS")
	dbName := env.RequiredString("DB_NAME")

	data, err := ioutil.ReadFile("queries.json")
	queries = new(Queries)
	json.Unmarshal(data, queries)

	connection, err = sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			dbUser,
			dbPass,
			dbName))

	execute("set statement_timeout to 30000;");

	if err != nil { log.Fatalf("Error connecting to database: %s", err.Error()) }
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
