package main

import (
	"fmt"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}

func (db *Database) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := db.GetUsers()
	if err != nil {
		fmt.Printf("Error fetching users: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}