package handlers

import (
	"address_module/internal/model"
	"address_module/internal/tools"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

var jwtSecret = []byte("super_secret_change_me")

// LoginHandler supports login with either username or email + password
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := tools.NewDatabase(5, 3*time.Second)
	if err != nil {
		log.Error("Database connection failed: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var user *model.User

	if strings.Contains(creds.Identifier, "@") {
		user, err = db.GetUserByEmail(creds.Identifier)
	} else {
		user, err = db.GetUserByUsername(creds.Identifier)
	}

	if err != nil || user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// TODO: Replace this with bcrypt.CompareHashAndPassword()
	if user.HashedPassword != creds.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		log.Error("JWT generation failed: ", err)
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	resp := model.LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// JWT generator
func generateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
