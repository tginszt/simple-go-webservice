package dbhandler

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

//dane połączenia z bazą danych
const (
	host     = "database"
	port     = 5432
	user     = "herbal"
	password = "Huj105"
	dbname   = "postgres"
)

func CreateTable() {

	// tworzymy tabelkę
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS
ziola(
		"id" SERIAL PRIMARY KEY,		
		"nazwa" TEXT,
		"dzialanie" TEXT,
		"wystepowanie" TEXT		
	  );`)

	//łączenie z bazą
	db := connectToDB()
	//dbamy o zamknięcie połączenia
	defer db.Close()

	// przygotowywanie sql'a przez prepare jest bezpieczne, bardzo przydatne, gdy chcemy użyć tej samej kwerendy wiele razy
	log.Println("preparing create table statement")
	statement, err := db.Prepare(createTable)
	checkErr(err)

	defer statement.Close()

	log.Println("executing create table")
	//wywołanie kwerendy
	_, err = statement.Exec()
	checkErr(err)
	log.Println("create table successful")
}

// insercik, wartości z requesta przekazywane w parametrach funkcji
func InsertIntoTable(nazwa string, dzialanie string, wystepowanie string) {

	//łączenie z bazą
	db := connectToDB()
	//dbamy o zamknięcie połączenia
	defer db.Close()

	insertStatementSQL := fmt.Sprintf(`INSERT INTO "ziola"("nazwa", "dzialanie", "wystepowanie") VALUES ($1, $2, $3);`)

	log.Println("prepare insert statement")
	// przygotowywanie sql'a przez prepare jest bezpieczne, bardzo przydatne, gdy chcemy użyć tej samej kwerendy wiele razy
	statement, err := db.Prepare(insertStatementSQL)
	checkErr(err)
	log.Println("successful")

	defer statement.Close()

	// wywołanie kwerendy
	log.Println("executing insert statement")
	_, err = statement.Exec(nazwa, dzialanie, wystepowanie)
	checkErr(err)
	log.Println("insert successful")

}

func PrintFromTable() string {

	//łączenie z bazą
	db := connectToDB()
	//dbamy o zamknięcie połączenia
	defer db.Close()

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

		fullReturn += "Ziolo: " + nazwa + " " + dzialanie + " " + wystepowanie + "\n"
	}
	return fullReturn
}

//Kreator informacji o połączeniu z bazą danych
func connectToDB() *sql.DB {

	// tworzymy połączenie
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// zapisujemy dane połączenia do zmiennej db
	db, _ := sql.Open("postgres", psqlconn)

	//zwracamy dane połączenia
	return db
}

//Sprawdzanie errorów nigdy jeszcze nie było tak szybkie i proste ;)
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
