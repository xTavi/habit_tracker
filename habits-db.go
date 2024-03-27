package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initDb() *sql.DB {
	const file string = "activities.db"
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func insertNewHabitIntoDb(habit Habit) error {
	db := initDb()
	defer db.Close()
	log.Default().Print("Creating habit called ", habit.Name, "and ID ", habit.ID)
	_, err := db.Exec("INSERT INTO habits (id, name) VALUES (?,?)", habit.ID, habit.Name)
	if err != nil {
		log.Fatal(err)
		return errors.New("Failed to create habit")
	}
	return nil
}

func getAllHabits() []Habit {
	db := initDb()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM habits")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var habits []Habit
	for rows.Next() {
		var habit Habit
		err := rows.Scan(&habit.ID, &habit.Name)
		if err != nil {
			log.Fatal(err)
		}
		habits = append(habits, habit)
	}

	return habits
}
