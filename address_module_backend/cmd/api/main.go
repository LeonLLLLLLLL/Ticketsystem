package main

import (
	"net/http"
	"os"
	"time"

	"address_module/internal/handlers"
	"address_module/internal/tools"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Configure logging
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Initialize database connection
	db, err := tools.NewDatabase(5, 3*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database schema
	err = db.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	/*err = db.InsertTestData()
	if err != nil {
		log.Fatal("Failed to insert test data:", err)
	}*/

	// Create a new router
	r := chi.NewRouter()

	// Apply CORS Middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"}, // Match frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            false, // Enable for debugging if needed
	})
	r.Use(corsMiddleware.Handler)

	// Register API Routes
	handlers.Handler(r)

	// Log server start
	log.Info("Starting GO API backend service on port 8000...")

	// Start HTTP Server with improved error handling
	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: r,
	}

	log.Info("Server is ready to handle requests")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
