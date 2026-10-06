package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"url-shortener/internal/models"
	"url-shortener/internal/service"

	"github.com/jackc/pgx/v5"
)

// type URLHandler struct {
// 	pool *pgxpool.Pool
// }

// func NewURLHandler(pool *pgxpool.Pool) *URLHandler {
// 	return &URLHandler{pool: pool}
// }

type URLHandler struct {
	service *service.URLService
}

func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "Welcome to URL Shortener",
		},
	)
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(
		map[string]string{
			"status": "healthy",
		},
	)
}

// Create URL Handler
func (h *URLHandler) CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	var request models.URLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if request.URL == "" {
		http.Error(
			w,
			"URL cannot be empty",
			http.StatusBadRequest,
		)
		return
	}

	shortCode, err := h.service.CreateURL(
		r.Context(),
		request.URL,
	)
	if err != nil {
		// fmt.Println("Database error: ", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"short_code": shortCode,
	})
}

// Get Original URL Handler
func (h *URLHandler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {

	shortCode := r.URL.Query().Get("short_code")

	if shortCode == "" {
		http.Error(w, "Missing short_code parameter", http.StatusBadRequest)
		return
	}

	url, err := h.service.GetURL(r.Context(), shortCode)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	if err != nil {
		fmt.Println("Database error: ", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url)
}

// Update URL Handler
func (h *URLHandler) UpdateURLHandler(w http.ResponseWriter, r *http.Request) {

	shortCode := r.URL.Query().Get("short_code")
	if shortCode == "" {
		http.Error(w, "Missing short_code parameter", http.StatusBadRequest)
		return
	}

	var newURL models.URLRequest
	err := json.NewDecoder(r.Body).Decode(&newURL)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateURL(
		r.Context(),
		shortCode,
		newURL.URL,
	)

	if err != nil {
		// fmt.Println("Database error: ", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "URL updated successfully",
	})
}

// Delete URL Handler
func (h *URLHandler) DeleteURLHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Query().Get("short_code")
	if shortCode == "" {
		http.Error(w, "Missing short_code parameter", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteURL(r.Context(), shortCode)

	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	if err != nil {
		// fmt.Println("Database error: ", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
