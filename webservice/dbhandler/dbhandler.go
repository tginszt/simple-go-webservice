package dbhandler

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
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

	// tworzymy połączenie
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// połączenie i zdeferowanie zamknięcia pliku z bazą danych
	db, _ := sql.Open("postgres", psqlconn)
	defer db.Close()

	// przygotowywanie sql'a przez prepare jest bezpieczne, bardzo przydatne, gdy chcemy użyć tej samej kwerendy wiele razy
	statement, err := db.Prepare(createTable)
	checkErr(err)

	defer statement.Close()

	//wywołanie kwerendy
	_, err = statement.Exec()
	checkErr(err)
}

// insercik, wartości z requesta przekazywane w parametrach funkcji
func InsertIntoTable(nazwa string, dzialanie string, wystepowanie string) {

	// tworzymy połączenie
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// połączenie i zdeferowanie zamknięcia pliku z bazą danych
	log.Println("connecting to database.db...")
	db, _ := sql.Open("postgres", psqlconn)
	defer db.Close()
	log.Println("Database connected!")

	insertStatementSQL := fmt.Sprintf(`INSERT INTO "ziola"("nazwa", "dzialanie", "wystepowanie") VALUES ($1, $2, $3);`)

	log.Println("prepare insert statement")
	log.Println(insertStatementSQL)
	// przygotowywanie sql'a przez prepare jest bezpieczne, bardzo przydatne, gdy chcemy użyć tej samej kwerendy wiele razy
	statement, err := db.Prepare(insertStatementSQL)
	checkErr(err)
	log.Println("succesful")

	defer statement.Close()

	// wywołanie kwerendy
	log.Println("execute")
	_, err = statement.Exec(nazwa, dzialanie, wystepowanie)
	checkErr(err)

}

func PrintFromTable() string {

	// tworzymy połączenie
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// połączenie i zdeferowanie zamknięcia pliku z bazą danych
	log.Println("connecting to database.db...")
	db, _ := sql.Open("postgres", psqlconn)
	defer db.Close()
	log.Println("Database connected!")

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

func StartDB() {
	// czyścimy baze danych żeby było wszystko widać
	os.Remove("database.db")

	// tworzymy baze danych
	log.Println("Creating database.db...")
	file, err := os.Create("database.db")
	checkErr(err)

	// kończenie operacji
	file.Close()
	log.Println("database.db created")
}

//Sprawdzanie errorów nigdy jeszcze nie było tak szybkie i proste ;)
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
