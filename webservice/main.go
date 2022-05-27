package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/dbhandler"
	"net/http"
)

// struct corresponding to the structure of the database
type herbStruct struct {
	Name         string
	Effects      string
	Occurrence   string
}

func main() {

	// creating db structure (tables)
	dbhandler.CreateTable()

	// the following functions are setting up get/post server
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/get", getHandler)
	log.Println("Post and get handlers are working on port 1234")
	log.Fatal(http.ListenAndServe(":1234", nil))
}

// server post
func postHandler(w http.ResponseWriter, req *http.Request) {

	log.Println("Post request incoming")

	// reading json to the herbStruct
	decoder := json.NewDecoder(req.Body)
	var herbs herbStruct
	err := decoder.Decode(&herbs)
	checkErr(err)

	log.Println("Json successfully decoded")

	// INSERT RECORDS
	dbhandler.InsertIntoTable(herbs.Name, herbs.Effects, herbs.Occurrence)
	log.Println("Added -> Name: ", herbs.Name, " Effects: ", herbs.Effects, " Occurrence: ", herbs.Occurrence)
}

// server get
func getHandler(w http.ResponseWriter, req *http.Request) {

	// printing select * from database
	fmt.Fprintf(w, dbhandler.PrintFromTable())
}

// Simple and quick error handling
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
