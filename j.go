package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Journal struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func main() {

	mux := http.NewServeMux()
	// mux.HandleFunc("POST /{id}", createJournal)
	mux.HandleFunc("POST /create", createJournal)
	mux.HandleFunc("GET /v1", getJOurnal)
	mux.HandleFunc("GET /v2", helper)

	fmt.Println("Hello World")

	http.ListenAndServe(":8080", mux)
}

func helper(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You app is up and listening to the :8080  port")
}

func createJournal(w http.ResponseWriter, r *http.Request) {
	//Checking which method is sent
	if r.Method != "POST" {
		http.Error(w, "INvalid request", http.StatusBadRequest)
		return
	}

	//Json decoding part
	var journal Journal
	err := json.NewDecoder(r.Body).Decode(&journal)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	file, err := os.OpenFile("journal.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "faild to write to file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var journals []Journal

	journalDataFile, err := os.ReadFile("journal.json")
	if err != nil && err != os.ErrNotExist {
		http.Error(w, "Faild to read for file", http.StatusInternalServerError)
		return
	}

	if len(journalDataFile) > 0 {
		err = json.Unmarshal(journalDataFile, &journals)
		if err != nil {
			http.Error(w, "faild to unmarshal json", http.StatusInternalServerError)
			return
		}
	}

	journals = append(journals, journal)

	journalDataFile, err = json.MarshalIndent(journals, "", "\t")
	if err != nil {
		http.Error(w, "failws to marshal JSON ", http.StatusInternalServerError)
		return

	}

	err = file.Truncate(0)
	if err != nil {
		http.Error(w, "Faild to truncate file", http.StatusInternalServerError)
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		http.Error(w, "Faild to seek file", http.StatusInternalServerError)
		return
	}

	_, err = file.Write(journalDataFile)
	if err != nil {
		http.Error(w, "Faild to write to file", http.StatusInternalServerError)
		return
	}

}

func getJOurnal(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	journalData, err := os.ReadFile("journal.json")
	if err != nil {
		if os.IsNotExist(err) {
			journal := Journal{
				Title:       "Example Journal",
				Description: "This is an example journal",
			}
			json.NewEncoder(w).Encode(journal)
			return
		}
		http.Error(w, "Faild to read for file", http.StatusInternalServerError)
		return
	}
	var journals []Journal

	err = json.Unmarshal(journalData, &journals)
	if err != nil {
		http.Error(w, "fialde to unmarshal json", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(journals)
}
