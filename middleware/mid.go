package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mineracail/guardApi/models"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Try to authenticate as Staff first
	var staffUser models.Staff
	if err := db.Where("email = ? AND password = ?", req.Email, req.Password).First(&staffUser).Error; err == nil {
		// Generate JWT token for Staff
		sendTokenResponse(w, staffUser.Email, staffUser.ID.String())
		return
	}

	// If not found as Staff, try to authenticate as Parent
	var parentUser models.Parent
	if err := db.Where("email = ? AND password = ?", req.Email, req.Password).First(&parentUser).Error; err == nil {
		// Generate JWT token for Parent
		sendTokenResponse(w, parentUser.Email, parentUser.ID.String())
		return
	}

	// If  found as Parent or log the user found
	// Log that neither Staff nor Parent was found
	log.Printf("Failed login attempt with email: %s", req.Email)
	// If neither Staff nor Parent were found, return unauthorized
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

// Helper function to send the token response
func sendTokenResponse(w http.ResponseWriter, email string, userID string) {
	token, err := GenerateToken(email, userID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Send the token as a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set status code to 200 OK
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
