package main

import (
	"incident-report-app/database"
	"incident-report-app/handlers"
	authmw "incident-report-app/middleware" // alias to avoid clash with chi middleware
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	database.Init()

	port := getEnv("PORT", "8080")
	frontendOrigin := getEnv("FRONTEND_ORIGIN", "http://localhost:5173")

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{frontendOrigin},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	r.Post("/api/auth/register", handlers.Register)
	r.Post("/api/auth/login", handlers.Login)

	// Protected — valid JWT required for all incident routes
	r.Group(func(r chi.Router) {
		r.Use(authmw.RequireAuth)
		r.Route("/api/incidents", func(r chi.Router) {
			r.Get("/", handlers.GetIncidents)
			r.Post("/", handlers.CreateIncident)
			r.Put("/{id}", handlers.UpdateIncident)
			r.Delete("/{id}", handlers.DeleteIncident)
		})
	})

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Server running on http://localhost:%s (CORS: %s)\n", port, frontendOrigin)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}