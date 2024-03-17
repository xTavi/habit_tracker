package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":8080"

func main() {
	http.HandleFunc("/habits", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		fmt.Fprintf(w, "Hello World")
	})
	//TODO How can I add another route like GET /habits and POST /habits
	

	err := http.ListenAndServe(portNumber, nil)
	if err != nil {
		log.Fatal(err)
	}
}
