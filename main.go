package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
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

var habitsDatabase = make([]Habit, 0)
var habitRecordDatabase = make([]HabitRecord, 0)

func main() {
	http.HandleFunc("POST /habits", handleHabitCreation)
	http.HandleFunc("POST /habits/{id}/track", handleHabitRecordCreation)
	http.HandleFunc("GET /habits", handleGetHabits)

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
		err2 := insertNewHabitIntoDb(habit)
		if err2 != nil {
			fmt.Println(err2)
			fmt.Println("Failed to create habit")
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(data)
		r.Body.Close()

	}

}

func handleGetHabits(w http.ResponseWriter, r *http.Request) {
	printRequestInfo(*r)
	w.Header().Set("Content-Type", "application/json")
	habits, err := json.Marshal(getAllHabits())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(habits)
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

	isHabitExistent := slices.ContainsFunc(habitsDatabase, func(habit Habit) bool {
		return habit.ID == incomingId
	})

	if !isHabitExistent {
		message, err := json.Marshal(HttpResponseError{ErrorMessage: "Habit not found"})

		if err != nil {
			log.Default().Printf("Error marshalling response: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, string(message[:]), http.StatusBadRequest)
		return
	}

	todayDate := time.Now()
	fmt.Println(todayDate, todayDate.Day())
	isHabitRecordedTodayAlready := slices.ContainsFunc(habitRecordDatabase, func(habitRecord HabitRecord) bool {
		return habitRecord.HabitId == incomingId && habitRecord.Date.Day() == todayDate.Day() && habitRecord.Date.Month() == todayDate.Month() && habitRecord.Date.Year() == todayDate.Year()
	})

	if isHabitRecordedTodayAlready {
		message, err := json.Marshal(HttpResponseError{ErrorMessage: "Habit already recorded today"})

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

	habitRecordDatabase = append(habitRecordDatabase, habitRecord)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(createdHabitRecord))
}

func printRequestInfo(r http.Request) {
	fmt.Printf("%s %s \n", r.Method, r.URL)
}
