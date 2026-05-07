package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, status int, payload any) error {
	js, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func RespondWithError(w http.ResponseWriter, status int, msg string) error {
	if status > 499 {
		log.Println("Responding with 5XX error: ", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	return RespondWithJSON(w, status, errorResponse{
		Error: msg,
	})
}
