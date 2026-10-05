package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"url-shortener/internal/models"
	"url-shortener/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type URLHandler struct {
	pool *pgxpool.Pool
}

func NewURLHandler(pool *pgxpool.Pool) *URLHandler {
	return &URLHandler{pool: pool}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome to URL Shortenerr"})
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server is running!")
}

// Create URL Handler
func (h *URLHandler) CreateURLHandler(w http.ResponseWriter, r *http.Request) {
	var request models.CreateURLRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newShortCode := utils.GenerateShortCode(6)
	_, resErr := h.pool.Exec(
		r.Context(),
		"INSERT INTO urls (short_code, original_url) VALUES ($1, $2)", newShortCode, request.URL,
	)
	if resErr != nil {
		fmt.Println("Database error: ", resErr)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"short_code": newShortCode,
	})
}

// Get Original URL Handler
func (h *URLHandler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {

	if r.URL.Query().Get("short_code") == "" {
		http.Error(w, "Missing short_code parameter", http.StatusBadRequest)
		return
	}

	var url models.URL
	resErr := h.pool.QueryRow(
		r.Context(),
		"SELECT id, short_code, original_url, access_count + 1, created_at, updated_at FROM urls WHERE short_code = $1", r.URL.Query().Get("short_code"),
	).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.AccessCount, &url.CreatedAt, &url.UpdatedAt)
	// Short Code not found in the database (pgx.ErrNoRows)
	if resErr != nil && resErr.Error() == "no rows in result set" {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	if resErr != nil {
		fmt.Println("Database error: ", resErr)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Increament the access count for the retrieved URL
	_, updateErr := h.pool.Exec(
		r.Context(),
		"UPDATE urls SET access_count = access_count + 1 WHERE short_code = $1", r.URL.Query().Get("short_code"),
	)
	if updateErr != nil {
		fmt.Println("Database error: ", updateErr)
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

	var newURL models.CreateURLRequest
	err := json.NewDecoder(r.Body).Decode(&newURL)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	result, errPool := h.pool.Exec(
		r.Context(),
		"UPDATE urls SET original_url = $1, access_count = access_count + 1, updated_at = NOW() WHERE short_code = $2;", newURL.URL, r.URL.Query().Get("short_code"))

	if errPool != nil {
		fmt.Println("Database error: ", errPool)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "URL updated successfully",
		"new_result": result.RowsAffected(),
	})
}

// Delete URL Handler
func (h *URLHandler) DeleteURLHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Query().Get("short_code")
	if shortCode == "" {
		http.Error(w, "Missing short_code parameter", http.StatusBadRequest)
		return
	}

	result, err := h.pool.Exec(
		r.Context(),
		"DELETE FROM urls WHERE short_code = $1", shortCode,
	)
	if result.RowsAffected() == 0 {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Println("Database error: ", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
