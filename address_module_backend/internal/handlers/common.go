package handlers

import (
	"address_module/internal/tools"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// getDBInstance returns a live database connection or handles the error and response
func getDBInstance(w http.ResponseWriter) (*tools.MySQLDB, bool) {
	db, err := tools.NewDatabase(5, 3*time.Second)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, false
	}
	return db, true
}

// getPostgresDBInstance returns a live PostgreSQL database connection or handles the error and response
func getPostgresDBInstance(w http.ResponseWriter) (*tools.PostgresDB, bool) {
	db, err := tools.NewPostgresDatabase(5, 3*time.Second)
	if err != nil {
		log.Errorf("Failed to connect to PostgreSQL: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, false
	}
	return db, true
}
