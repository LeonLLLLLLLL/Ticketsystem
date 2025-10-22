package main

import (
	"net/http"
	"time"

	"address_module/internal/handlers"
	"address_module/internal/tools"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Configure logging
	//log.SetReportCaller(true)
	//log.SetFormatter(&log.JSONFormatter{})
	//log.SetOutput(os.Stdout)
	tools.Configure()

	// Initialize database connection
	db, err := tools.NewDatabase(15, 3*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize database schema
	err = db.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	err = db.SeedInitialData()
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
	/*
		err = db.InsertTestData()
		if err != nil {
			log.Fatal("Failed to insert test data:", err)
		}

		err = db.InsertUserRolesTestData()
		if err != nil {
			log.Fatal("Faild to instert test data:", err)
		}*/

	/*err = db.RunCRUDTests()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}*/

	defer db.Close()

	// Connect to Postgres Device Management DB
	pgDB, err := tools.NewPostgresDatabase(20, 2*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres (device_management_database): %v", err)
	}
	defer pgDB.Close()

	// Run CRUD test for device management DB
	if err := tools.RunPostgresDeviceCRUDTests(pgDB); err != nil {
		log.Errorf("Postgres device CRUD tests failed: %v", err)
	} else {
		log.Info("Postgres device CRUD tests passed ✅")
	}

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
