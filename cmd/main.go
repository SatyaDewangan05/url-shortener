package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"url-shortener/internal/database"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

var mux = http.NewServeMux()

func main() {

	// Connect to DB
	const databaseURL = "postgres://postgres:postgres@localhost:5432/url_shortener"

	pool, err := database.ConnectToDatabase(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// Check DB connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database Connected Successfully")
	
	urlRepository := repository.NewURLRepository(pool)
	urlService := service.NewURLService(urlRepository)
	urlHandler := handler.NewURLHandler(urlService)

	// Welcome
	mux.HandleFunc("GET /", handler.Welcome)

	// Health
	mux.HandleFunc("GET /health", handler.HandleHealth)

	// Create URL
	mux.HandleFunc("POST /api/v1/url", urlHandler.CreateURLHandler)
	
	// Retrieve URL
	mux.HandleFunc("GET /api/v1/url", urlHandler.GetOriginalURL)

	// Update URL
	mux.HandleFunc("PUT /api/v1/url", urlHandler.UpdateURLHandler)
	// Delete URL
	mux.HandleFunc("DELETE /api/v1/url", urlHandler.DeleteURLHandler)

	// Stats on URL usage

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Server failed: ", err)
	} else {
		fmt.Println("Server is running on https://localhost:8080")
	}

	pool.Close()
}





