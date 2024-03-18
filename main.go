package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/beevik/guid"
)

const portNumber = ":8081"

type Habit struct {
	ID   string
	Name string `json:"name"`
}

func CreateHabit(name string) Habit {
	log.Default().Print("Creating habit called ", name)
	return Habit{
		ID:   guid.New().String(),
		Name: name,
	}
}

type HabitRecord struct {
	ID      string    `json:"id"`
	HabitId string    `json:"habitId"`
	Date    time.Time `json:"date"`
	Note    string    `json:"note"`
}

func main() {

	http.HandleFunc("POST /habits", func(w http.ResponseWriter, r *http.Request) {

		// respond to the client with the error message and a 400 status code.
		var habit Habit
		err := json.NewDecoder(r.Body).Decode(&habit)

		if err != nil {
			log.Default().Printf("Error decoding request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if habit.Name == "" {
			http.Error(w, json.Marshal(name: "Name is required"), http.StatusBadRequest)
			return
		}

		habit.ID = guid.New().String()

		fmt.Printf("%s %s \n", r.Method, r.URL)
		log.Default().Print(habit)
		fmt.Println(habit.ID)
		fmt.Println(habit.Name)
		fmt.Fprintf(w, "Hit the habits endpoint")
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
