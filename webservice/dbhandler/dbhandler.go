package dbhandler

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func CreateTable(db *sql.DB) {
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

	//wywołanie kwerendy
	_, err = statement.Exec()
	checkErr(err)
}

// insercik, wartości z requesta przekazywane w parametrach funkcji
func InsertIntoTable(db *sql.DB, nazwa string, dzialanie string, wystepowanie string) {
	insertStatementSQL := fmt.Sprintf(`INSERT INTO ziola(nazwa, dzialanie, wystepowanie) VALUES (?, ?, ?)`)

	// przygotowywanie sql'a przez prepare jest bezpieczne, bardzo przydatne, gdy chcemy użyć tej samej kwerendy wiele razy
	statement, err := db.Prepare(insertStatementSQL)
	checkErr(err)

	defer statement.Close()

	// wywołanie kwerendy
	_, err = statement.Exec(nazwa, dzialanie, wystepowanie)
	checkErr(err)

}

func PrintFromTable(db *sql.DB) string {
	//pobieramy po rzędzie dane bazych
	row, err := db.Query(fmt.Sprintf("SELECT * FROM ziola"))
	checkErr(err)
	defer row.Close()

	// do tego stringa zwracamy wszystkie dane ładnie ułożone
	var fullReturn string

	// czytamy po rzędzie z bazy, wrzucamy wartości do fullReturn
	for row.Next() {
		var id int
		var nazwa string
		var dzialanie string
		var wystepowanie string
		row.Scan(&id, &nazwa, &dzialanie, &wystepowanie)

		fullReturn += "Ziolo: " + nazwa + " " + dzialanie + " " + wystepowanie +"\n"
	}
	return fullReturn
}

func StartDB(){
	// czyścimy baze danych żeby było wszystko widać
	os.Remove("sqlite-database.db")

	// tworzymy baze danych
	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	checkErr(err)
	file.Close()
	log.Println("sqlite-database.db created")
}

//Sprawdzanie errorów nigdy jeszcze nie było tak szybkie i proste ;)
func checkErr(err error){
	if err != nil {
		log.Fatalln(err.Error())
	}
}