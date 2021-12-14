package main

import (
    "database/sql"
    "dbhandler.com/dbhandler"
    "encoding/json"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "net/http"
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
- (^>^) Zrob nowy package, z funkcjami dotyczącymi DB- np. db.go. Zaimportuj go i używaj.
- Implementacja endpoitów osobnych funkcjach, nie w mainie
- Dockerfile. Musi być możliwość zrobienia `docker run .` w folderze webservice, i serwis ma działać.
2.
- Zmiana bazy danych na postgresql (Dockerfile przestanie dzialac)
- Docker-compose uzywajacy istniejacego Dockerfile (ma znow dzialac)


*/
func main(){

    dbhandler.StartDB()

    // otwarcie i zdeferowane zamknięcie pliku z baną dazych
    sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
    defer sqliteDatabase.Close()

    // tworzymy strukturę bazy (tabelki)
    dbhandler.CreateTable(sqliteDatabase)

    // post
    postHandler := func(w http.ResponseWriter, req *http.Request) {

        // czytanie jsona do struktury ziolaS
        decoder := json.NewDecoder(req.Body)
        var ziolaS ziolaStruct
        err := decoder.Decode(&ziolaS)
            checkErr(err)

        // INSERT RECORDS
        log.Println("Uwaga, dodano -> Nazwa: ", ziolaS.Nazwa," Dzialanie: ", ziolaS.Dzialanie," Wystepowanie: ", ziolaS.Wystepowanie)
        dbhandler.InsertIntoTable(sqliteDatabase, ziolaS.Nazwa, ziolaS.Dzialanie, ziolaS.Wystepowanie)
    }

    // get
    getHandler := func(w http.ResponseWriter, req *http.Request) {

        // ładnie wyświetla select * from database
        fmt.Fprintf(w, dbhandler.PrintFromTable(sqliteDatabase))
    }

    // poniższe funkcje stawiają nam serwer
    http.HandleFunc("/post", postHandler)
    http.HandleFunc("/get", getHandler)
    log.Fatal(http.ListenAndServe(":1234", nil))
}

//Sprawdzanie errorów nigdy jeszcze nie było tak szybkie i proste ;)
func checkErr(err error){
    if err != nil {
        log.Fatalln(err.Error())
    }
}
