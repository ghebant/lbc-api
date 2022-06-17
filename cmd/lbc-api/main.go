package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello, world.")
	http.HandleFunc("/", Health)

	err := http.ListenAndServe(os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	// TODO Remove
	log.Println("/health !")
	w.WriteHeader(http.StatusOK)
}
