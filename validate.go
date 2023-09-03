package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func ChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanedMsg := replaceBadWords(params.Body, badWords)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleanedMsg,
	})
}

func replaceBadWords(msg string, badWords map[string]struct{}) string {

	msgWords := strings.Split(msg, " ")

	for i, word := range msgWords {
		lowWord := strings.ToLower(word)
		if _, ok := badWords[lowWord]; ok {
			msgWords[i] = "****"
		}
	}

	cleanMsg := strings.Join(msgWords, " ")
	return cleanMsg
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
