package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	mux.HandleFunc("POST /user", createUser)
	mux.HandleFunc("GET /user/{id}", getUser)
	mux.HandleFunc("DELETE /user/{id}", deleteUser)
	fmt.Println("server listening to 8080")
	http.ListenAndServe(":8080", mux)
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

func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cacheMutex.RLock()
	user, ok := userCache[id]
	cacheMutex.RUnlock()
	if !ok {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := userCache[id]; !ok {
		http.Error(w, "User not found", http.StatusBadRequest)
	}
	return

	cacheMutex.Lock()
	delete(userCache, id)
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
