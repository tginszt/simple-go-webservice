package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/dbhandler"
	"net/http"
)

// struktura będąca odpowiednikiem do struktury bazy danych
type ziolaStruct struct {
	Nazwa        string
	Dzialanie    string
	Wystepowanie string
}

/*
TODO:
1.
- (^>^) GET zwraca liste rzeczy z DB
- (^>^) Zrob nowy package, z funkcjami dotyczącymi DB- np. db.go. Zaimportuj go i używaj.
- (^>^) Implementacja endpoitów osobnych funkcjach, nie w mainie
- (^>^) Dockerfile. Musi być możliwość zrobienia `docker run .` w folderze webservice, i serwis ma działać.
2.
- (^>^) Zmiana bazy danych na postgresql (Dockerfile przestanie dzialac)
- (^>^) Docker-compose uzywajacy istniejacego Dockerfile (ma znow dzialac)

*/
func main() {

	// tworzymy strukturę bazy (tabelki)
	dbhandler.CreateTable()

	// poniższe funkcje stawiają nam serwer
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)
	log.Println("Handlery post i get działają na porcie 1234")
	log.Fatal(http.ListenAndServe(":1234", nil))
}

// serwer post
func postHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("Post request incoming")

	// czytanie jsona do struktury ziolaS
	decoder := json.NewDecoder(req.Body)
	var ziolaS ziolaStruct
	err := decoder.Decode(&ziolaS)
	checkErr(err)

	log.Println("Json successfully decoded")

	// INSERT RECORDS
	dbhandler.InsertIntoTable(ziolaS.Nazwa, ziolaS.Dzialanie, ziolaS.Wystepowanie)
	log.Println("Uwaga, dodano -> Nazwa: ", ziolaS.Nazwa, " Dzialanie: ", ziolaS.Dzialanie, " Wystepowanie: ", ziolaS.Wystepowanie)
}

// serwer get
func getHandler(w http.ResponseWriter, req *http.Request) {

	// ładnie wyświetla select * from database
	fmt.Fprintf(w, dbhandler.PrintFromTable())
}

//Sprawdzanie errorów nigdy jeszcze nie było tak szybkie i proste ;)
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
