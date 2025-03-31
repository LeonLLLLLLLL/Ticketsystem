package handlers

import (
	"address_module/internal/model"
	"address_module/internal/tools"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	log "github.com/sirupsen/logrus"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if req.Email == "" || req.Password == "" || req.Username == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password: ", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	db, err := tools.NewDatabase(5, 3*time.Second)
	if err != nil {
		log.Error("DB connection error: ", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Create user
	user := model.User{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: string(hashed),
		CreatedAt:      time.Now(),
	}

	userID, err := db.InsertUser(user)
	if err != nil {
		log.Error("Failed to insert user during registration: ", err)
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}
	log.Infof("✅ Registered new user %s (ID: %d)", user.Username, userID)

	// Optionally assign default role
	defaultRoleName := "user"
	role, err := db.GetRoleByName(defaultRoleName)
	if err == nil {
		err = db.InsertUserRole(model.UserRole{
			UserID: userID,
			RoleID: role.ID,
		})
		if err != nil {
			log.Warnf("Could not assign default role %s to new user: %v", defaultRoleName, err)
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registration successful",
	})
}
