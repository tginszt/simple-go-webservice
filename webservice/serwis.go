package main

import (
    "encoding/json"
    "fmt"
    "database/sql"
    "os"
    "log"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
)

type ziolaStruct struct {
    Nazwa string
    Dzialanie string
    Wystepowanie string
}

func main() {
    os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records. SQLite is a file based database.
    log.Println("Creating sqlite-database.db...")
    file, err := os.Create("sqlite-database.db") // Create SQLite file
        checkErr(err)

    file.Close()
    log.Println("sqlite-database.db created")

    sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
    defer sqliteDatabase.Close() // Defer Closing the database

    createTable(sqliteDatabase) // Create Database Table

    // post
    postHandler := func(w http.ResponseWriter, req *http.Request) {

        decoder := json.NewDecoder(req.Body)
        var ziolaS ziolaStruct
        err := decoder.Decode(&ziolaS)
            checkErr(err)

        // INSERT RECORDS
        log.Println("Uwaga, dodano -> Nazwa: ", ziolaS.Nazwa," Dzialanie: ", ziolaS.Dzialanie," Wystepowanie: ", ziolaS.Wystepowanie)
        insertIntoTable(sqliteDatabase, ziolaS.Nazwa, ziolaS.Dzialanie, ziolaS.Wystepowanie)

        return 1
    }
    // get
    getHandler := func(w http.ResponseWriter, req *http.Request) {
        // ładnie wyświetla select * from database
        printFromTable(sqliteDatabase)
    }

    http.HandleFunc("/post", postHandler)
    http.HandleFunc("/get", getHandler)
    log.Fatal(http.ListenAndServe(":1234", nil))

}

func createTable(db *sql.DB) {
    // tworzymy tabelkę
    createTableSQL := fmt.Sprintf(`CREATE TABLE ziola(
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"nazwa" TEXT,
		"dzialanie" TEXT,
		"wystepowanie" TEXT		
	  );`)

    // przygotowywanie sql'a przez prepare jest bezpieczne, bardzo przydatne, gdy chcemy użyć tej samej kwerendy wiele razy
    statement, err := db.Prepare(createTableSQL)
        checkErr(err)

    defer statement.Close()

    _, err = statement.Exec()
        checkErr(err)
}

// insercik, wartości z requesta przekazywane w parametrach funkcji
func insertIntoTable(db *sql.DB, nazwa string, dzialanie string, wystepowanie string) {
    insertStatementSQL := fmt.Sprintf(`INSERT INTO ziola(nazwa, dzialanie, wystepowanie) VALUES (?, ?, ?)`)

    // przygotowywanie sql'a przez prepare jest bezpieczne, bardzo przydatne, gdy chcemy użyć tej samej kwerendy wiele razy
    statement, err := db.Prepare(insertStatementSQL)
        checkErr(err)

    defer statement.Close()

    _, err = statement.Exec(nazwa, dzialanie, wystepowanie)
        checkErr(err)

}

func printFromTable(db *sql.DB) {
    row, err := db.Query(fmt.Sprintf("SELECT * FROM ziola"))
        checkErr(err)
    defer row.Close()

    for row.Next() {
        var id int
        var nazwa string
        var dzialanie string
        var wystepowanie string
        row.Scan(&id, &nazwa, &dzialanie, &wystepowanie)
        log.Println("Ziolo: ", nazwa, " ", dzialanie, " ", wystepowanie)
    }
}

func checkErr(err error){
    if err != nil {
        log.Fatalln(err.Error())
    }
}
