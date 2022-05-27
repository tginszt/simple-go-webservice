package dbhandler

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

// db connection credentials
const (
	host     = "database"
	port     = 5432
	user     = "herbal"
	password = "webservice"
	dbname   = "postgres"
)

func CreateTable() {

	// creating a table
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS
herbs(
		"id" SERIAL PRIMARY KEY,		
		"name" TEXT,
		"effects" TEXT,
		"occurrence" TEXT
	  );`)

	// connecting to db
	db := connectToDB()
	// closing the connection
	defer db.Close()

	// preparing sql statement
	log.Println("preparing create table statement")
	statement, err := db.Prepare(createTable)
	checkErr(err)

	defer statement.Close()

	log.Println("executing create table")
	// executing the statement
	_, err = statement.Exec()
	checkErr(err)
	log.Println("create table successful")
}

// inserting, values from request are passed as function parameters
func InsertIntoTable(name string, effects string, occurrence string) {

	// connecting to db
	db := connectToDB()
	// closing the connection
	defer db.Close()

	insertStatementSQL := fmt.Sprintf(`INSERT INTO "herbs"("name", "effects", "occurrence") VALUES ($1, $2, $3);`)

	log.Println("prepare insert statement")
	// preparing sql statement
	statement, err := db.Prepare(insertStatementSQL)
	checkErr(err)
	log.Println("successful")

	defer statement.Close()

	// executing the statement
	log.Println("executing insert statement")
	_, err = statement.Exec(name, effects, occurrence)
	checkErr(err)
	log.Println("insert successful")

}

func PrintFromTable() string {

	//connecting to db
	db := connectToDB()
	// closing the connection
	defer db.Close()

	// collecting data from db row by row
	row, err := db.Query(fmt.Sprintf("SELECT * FROM herbs"))
	checkErr(err)
	defer row.Close()

	// this string is used to return the data
	var fullReturn string

	// collecting data from db row by row into fullReturn
	for row.Next() {
		var id int
		var name string
		var effects string
		var occurrence string
		row.Scan(&id, &name, &effects, &occurrence)

		fullReturn += "herb: " + name + " " + effects + " " + occurrence + "\n"
	}
	return fullReturn
}

// db connection creator
func connectToDB() *sql.DB {

	// creating the connection
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// saving connection data to db
	db, _ := sql.Open("postgres", psqlconn)

	// returning the data
	return db
}

// Simple and quick error handling
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
