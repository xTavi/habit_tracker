package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

type Habit struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type HabitRecord struct {
	ID      string    `json:"id"`
	HabitId string    `json:"habitId"`
	Date    time.Time `json:"date"`
	Note    string    `json:"note"`
}

func main() {

	http.HandleFunc("POST /habits", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s", r.Method, r.URL)

	})

	http.HandleFunc("GET /habits", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("%s %s", r.Method, r.URL)
		fmt.Fprintf(w, "Hit the habits endpoint")
	})

	//TODO How can I add another route like GET /habits and POST /habits

	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		log.Fatal(err)
	}
}
