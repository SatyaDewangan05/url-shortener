package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"url-shortener/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URLHandler struct {
	pool *pgxpool.Pool
}

func NewURLHandler(pool *pgxpool.Pool) *URLHandler {
	return &URLHandler{pool: pool}
}

type CreateURLRequest struct {
	URL string `json:"url"`
}

type URL struct {
	ID          int    `json:"id"`
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
}

// Package-level variables
var urls = make(map[string]URL)

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to URL Shortenerr"})
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server is running!")
}

// Create URL Handler
func (h *URLHandler) CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	var request CreateURLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	url := URL{
		ID:          len(urls) + 1,
		ShortCode:   utils.GenerateShortCode(6),
		OriginalURL: request.URL,
	}

	_, resErr := h.pool.Exec(
		r.Context(),
		"INSERT INTO urls (short_code, original_url, access_count) VALUES ($1, $2, 0)", url.ShortCode, url.OriginalURL,
	)

	if resErr != nil {
		fmt.Println("Database error: ", resErr)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	urls[url.ShortCode] = url

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(url)
}

// Get Original URL Handler
func (h *URLHandler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("short_code") == "" {
		http.Error(w, "Missing short_code parameter", http.StatusBadRequest)
		return
	}

	var url URL

	row := h.pool.QueryRow(
		r.Context(),
		"SELECT id, short_code, original_url FROM urls WHERE short_code = $1", r.URL.Query().Get("short_code"),
	)

	resErr := row.Scan(&url.ID, &url.ShortCode, &url.OriginalURL)

	if resErr != nil {
		fmt.Println("Database error: ", resErr)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// shortCode := r.URL.Query().Get("short_code")
	// urlLokal, exists := urls[shortCode]

	// if !exists {
	// 	http.Error(w, "Short URL not found", http.StatusNotFound)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(url.OriginalURL)
}

// Update URL Handler
func (h *URLHandler) UpdateURLHandler(w http.ResponseWriter, r *http.Request) {

	shortCode := r.URL.Query().Get("short_code")
	if shortCode == "" {
		http.Error(w, "Missing short_code parameter", http.StatusBadRequest)
		return
	}

	_, exists := urls[shortCode]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	var newURL CreateURLRequest
	err := json.NewDecoder(r.Body).Decode(&newURL)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update the original URL for the given short code
	url := urls[shortCode]
	url.OriginalURL = newURL.URL
	urls[shortCode] = url
}

// Delete URL Handler
func (h *URLHandler) DeleteURLHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Query().Get("short_code")
	if shortCode == "" {
		http.Error(w, "Missing short_code parameter", http.StatusBadRequest)
		return
	}

	_, exists := urls[shortCode]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	delete(urls, shortCode)
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
