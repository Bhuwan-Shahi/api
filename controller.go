package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// sever home rout
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>This is the landing page</h1>"))
}

func getAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("All the journals")
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(journal)
}
