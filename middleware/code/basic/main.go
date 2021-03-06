package main

import (
	"fmt"
	"log"
	"net/http"
)

// START OMIT
func main() {
	// define simple route
	http.HandleFunc("/hello", helloGET)
	log.Println("App started, listening on port 9090")
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func helloGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
}

// END OMIT
