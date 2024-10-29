package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type User struct {
	Name string `json :"name"`
}

var userCache = make(map[int]User)

// to prevent race condition
var cacheMutex sync.RWMutex

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /user", createUser)

	mux.HandleFunc("/", handlerRoot)

	fmt.Println("server listening to 8080")
	http.ListenAndServe(":8080", mux)
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Name == "" {
		http.Error(w, "Name is error!!", http.StatusBadRequest)
		return
	}
	cacheMutex.Lock()
	userCache[len(userCache)+1] = user
	cacheMutex.Unlock()
	w.WriteHeader(http.StatusNoContent)
}
