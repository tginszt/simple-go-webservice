package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "net/http"
    "os"
)

// struktura będąca odpowiednikiem do struktury bazy danych
type ziolaStruct struct {
    Nazwa string
    Dzialanie string
    Wystepowanie string
}
/*
TODO:
1.
- (^>^) GET zwraca liste rzeczy z DB
- Zrob nowy package, z funkcjami dotyczącymi DB- np. db.go. Zaimportuj go i używaj.
- Implementacja endpoitów osobnych funkcjach, nie w mainie
- Dockerfile. Musi być możliwość zrobienia `docker run .` w folderze webservice, i serwis ma działać.
2.
- Zmiana bazy danych na postgresql (Dockerfile przestanie dzialac)
- Docker-compose uzywajacy istniejacego Dockerfile (ma znow dzialac)


*/
func main() {
    // czyścimy baze danych żeby było wszystko widać
    os.Remove("sqlite-database.db")

    // tworzymy baze danych
    log.Println("Creating sqlite-database.db...")
    file, err := os.Create("sqlite-database.db") // Create SQLite file
        checkErr(err)
    file.Close()
    log.Println("sqlite-database.db created")

    // otwarcie i zdeferowane zamknięcie pliku z baną dazych
    sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
    defer sqliteDatabase.Close()

    // tworzymy strukturę bazy (tabelki)
    createTable(sqliteDatabase)

    // post
    postHandler := func(w http.ResponseWriter, req *http.Request) {

        // czytanie jsona do struktury ziolaS
        decoder := json.NewDecoder(req.Body)
        var ziolaS ziolaStruct
        err := decoder.Decode(&ziolaS)
            checkErr(err)

        // INSERT RECORDS
        log.Println("Uwaga, dodano -> Nazwa: ", ziolaS.Nazwa," Dzialanie: ", ziolaS.Dzialanie," Wystepowanie: ", ziolaS.Wystepowanie)
        insertIntoTable(sqliteDatabase, ziolaS.Nazwa, ziolaS.Dzialanie, ziolaS.Wystepowanie)
    }

    // get
    getHandler := func(w http.ResponseWriter, req *http.Request) {

        // ładnie wyświetla select * from database
        fmt.Fprintf(w,printFromTable(sqliteDatabase))
    }

    // poniższe funkcje stawiają nam serwer
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

    //wywołanie kwerendy
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

    // wywołanie kwerendy
    _, err = statement.Exec(nazwa, dzialanie, wystepowanie)
        checkErr(err)

}

func printFromTable(db *sql.DB) string {
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

//Sprawdzanie errorów nigdy jeszcze nie było tak szybkie i proste ;)
func checkErr(err error){
    if err != nil {
        log.Fatalln(err.Error())
    }
}
