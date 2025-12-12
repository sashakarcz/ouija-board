package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// App holds application dependencies
type App struct {
	config  *Config
	storage Storage
	ollama  *OllamaClient
}

// AskRequest represents the incoming question request
type AskRequest struct {
	Question string `json:"question"`
}

// AskResponse represents the answer response
type AskResponse struct {
	Answer string `json:"answer"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// indexHandler serves the main HTML page
func (app *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// askHandler handles question submissions
func (app *App) askHandler(w http.ResponseWriter, r *http.Request) {
	// Validate content type
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		respondWithError(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	// Parse request
	var req AskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate question length
	if len(req.Question) > 1000 {
		respondWithError(w, "Question too long (max 1000 characters)", http.StatusBadRequest)
		return
	}

	// Validate question is not empty after trimming
	if strings.TrimSpace(req.Question) == "" {
		respondWithError(w, "Question cannot be empty", http.StatusBadRequest)
		return
	}

	// Generate answer using Ollama
	answer, err := app.ollama.GenerateAnswer(r.Context(), req.Question)
	if err != nil {
		log.Printf("Error generating answer: %v", err)
		respondWithError(w, "Failed to generate answer", http.StatusInternalServerError)
		return
	}

	// Store Q&A pair
	pair := QAPair{
		Question: req.Question,
		Answer:   answer,
	}

	if err := app.storage.Add(pair); err != nil {
		log.Printf("Error storing Q&A pair: %v", err)
		// Don't fail the request if storage fails, just log it
	}

	// Respond with answer
	respondWithJSON(w, AskResponse{Answer: answer}, http.StatusOK)
}

// historyHandler returns all Q&A history
func (app *App) historyHandler(w http.ResponseWriter, r *http.Request) {
	pairs, err := app.storage.GetAll()
	if err != nil {
		log.Printf("Error retrieving history: %v", err)
		respondWithError(w, "Failed to retrieve history", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, pairs, http.StatusOK)
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, payload interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	respondWithJSON(w, ErrorResponse{Error: message}, statusCode)
}
