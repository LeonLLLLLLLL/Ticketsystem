package middleware

import (
	"address_module/internal/tools"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// RequirePermission enforces that the logged-in user has a specific permission
func RequirePermission(requiredPermission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uid := r.Context().Value(UserIDKey)
			if uid == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userID := uid.(int64)

			db, err := tools.NewDatabase(5, 3*time.Second)
			if err != nil {
				log.Error("Failed to connect to DB in RequirePermission: ", err)
				http.Error(w, "Server error", http.StatusInternalServerError)
				return
			}
			defer db.Close()

			perms, err := db.GetUserPermissions(userID)
			if err != nil {
				log.Warnf("Failed to get permissions for user %d: %v", userID, err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			for _, p := range perms {
				if p.Name == requiredPermission {
					// User has the required permission
					next.ServeHTTP(w, r)
					return
				}
			}

			log.Warnf("User %d does not have required permission: %s", userID, requiredPermission)
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
