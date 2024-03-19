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

type HttpResponseError struct {
	ErrorMessage string `json:"errorMessage"`
}

func main() {

	http.HandleFunc("POST /habits", handleHabitCreation)
	http.HandleFunc("POST /habits/{id}/track", handleHabitRecordCreation)
	http.HandleFunc("GET /habits", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("%s %s", r.Method, r.URL)
		fmt.Fprintf(w, "Hit the habits endpoint")
	})

	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleHabitCreation(w http.ResponseWriter, r *http.Request) {
	printRequestInfo(*r)
	w.Header().Set("Content-Type", "application/json")
	var habit Habit
	err := json.NewDecoder(r.Body).Decode(&habit)

	if err != nil {
		log.Default().Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if habit.Name == "" {
		message, err := json.Marshal(HttpResponseError{ErrorMessage: "Name is required"})

		if err != nil {
			log.Default().Printf("Error marshalling response: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, string(message[:]), http.StatusBadRequest)
		return
	}

	habit.ID = guid.New().String()

	log.Default().Print(habit)
	fmt.Println(habit.ID)
	fmt.Println(habit.Name)

	if data, err := json.Marshal(habit); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
		r.Body.Close()
	}

}

func handleHabitRecordCreation(w http.ResponseWriter, r *http.Request) {
	printRequestInfo(*r)
	w.Header().Set("Content-Type", "application/json")
	incomingId := r.PathValue("id")
	if !guid.IsGuid(incomingId) {
		message, err := json.Marshal(HttpResponseError{ErrorMessage: "Invalid habitId"})

		if err != nil {
			log.Default().Printf("Error marshalling response: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, string(message[:]), http.StatusBadRequest)
		return
	}

	var habitRecord HabitRecord
	err := json.NewDecoder(r.Body).Decode(&habitRecord)

	if err != nil {
		log.Default().Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	habitRecord.ID = guid.New().String()
	habitRecord.HabitId = incomingId
	habitRecord.Date = time.Now()

	log.Default().Print("Created habitRecord: ", habitRecord)
	fmt.Println(habitRecord.ID)
	fmt.Println(habitRecord.HabitId)
	fmt.Println(habitRecord.Date)
	fmt.Println(habitRecord.Note)

	createdHabitRecord, err := json.Marshal(habitRecord)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(createdHabitRecord))
}

func printRequestInfo(r http.Request) {
	fmt.Printf("%s %s \n", r.Method, r.URL)
}
