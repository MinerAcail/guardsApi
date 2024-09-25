package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JwtSecret is the secret key used for signing JWT tokens.
// Replace this with your own secret key and keep it secure.
var JwtSecret = []byte("SecretInHere")

// TokenStruct defines the structure of the JWT token claims.
type TokenStruct struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	jwt.StandardClaims
}

// Define context keys as custom types to avoid conflicts.
type contextKey string

const (
	IDContextKey       contextKey = "ID"
	UserTypeContextKey contextKey = "userType"
)

// Middleware function for handling authentication and setting context values.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			userID, userType, _, err := ParseToken(tokenStr)
			if err == nil {
				ctx = context.WithValue(ctx, IDContextKey, userID)
				ctx = context.WithValue(ctx, UserTypeContextKey, userType)
			}
		}

		// Pass the request to the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ParseToken parses the provided JWT token and returns the claims if valid.
func ParseToken(tokenStr string) (string, string, map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})
	if err != nil {
		return "", "", nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ID, ok := claims["id"].(string)
		if !ok {
			return "", "", nil, fmt.Errorf("ID not found in the token claims or has an unexpected type")
		}

		Type, ok := claims["type"].(string)
		if !ok {
			return "", "", nil, fmt.Errorf("type not found in the token claims or has an unexpected type")
		}

		// Extract optional claims, if they exist.
		optionalClaims := make(map[string]interface{})
		for key, value := range claims {
			if key != "id" && key != "type" {
				optionalClaims[key] = value
			}
		}

		return ID, Type, optionalClaims, nil
	} else {
		return "", "", nil, fmt.Errorf("invalid token")
	}
}

// ExtractCtxInfoForAllAccess extracts and validates information from the request context.
func ExtractCtxInfoForAllAccess(ctx context.Context) error {
	_, ok := ctx.Value(IDContextKey).(string)
	if !ok {
		return fmt.Errorf("token expired or you don't have access to the request context, try logging in again")
	}

	userType, ok := ctx.Value(UserTypeContextKey).(string)
	if !ok {
		return fmt.Errorf("userType not found in request context")
	}

	userType = strings.ToLower(userType) // Convert userType to lowercase

	// Check if the userType is not one of the allowed values.
	if userType != "staff" && userType != "student" {
		return fmt.Errorf("%s is not allowed", userType)
	}

	return nil
}

// ValidateTokens validates the provided JWT token and returns the claims if the token is valid.
func ValidateTokens(signedToken string) (*TokenStruct, error) {
	// Parse the JWT token with claims and validation function.
	token, err := jwt.ParseWithClaims(
		signedToken,
		&TokenStruct{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method used in the token.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	// Verify if the token is valid and not expired.
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract and type-assert the claims.
	claims, ok := token.Claims.(*TokenStruct)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify if the token has expired.
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, fmt.Errorf("token has expired")
	}

	// Token is valid, return the claims.
	return claims, nil
}

// GenerateToken creates a new JWT token with the specified type, ID, and optional claims.
func GenerateToken(userType string, ID string, optionalClaims ...jwt.MapClaims) (string, error) {
	// Create a new JWT token with the signing method HS256 (HMAC SHA-256)
	token := jwt.New(jwt.SigningMethodHS256)

	// Convert the token's claims to a map
	claims := token.Claims.(jwt.MapClaims)

	// Set the "type" and "ID" claims in the token
	claims["type"] = userType
	claims["id"] = ID

	// Set the "exp" (expiration) claim to the current time + 24 hours
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Set optional claims if provided
	for _, optionalClaim := range optionalClaims {
		for key, value := range optionalClaim {
			claims[key] = value
		}
	}

	// Sign the token using a secret key
	tokenString, err := token.SignedString(JwtSecret)
	if err != nil {
		return "", err
	}

	// Return the signed token as a string
	return tokenString, nil
}

/*
// SetTokenForContext generates a new token and sets it in the context’s cookie.
func SetTokenForContext(ctx context.Context, userType, id string) error {
	// Generate a new token
	token, err := GenerateToken(userType, id)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Retrieve the CookieAccess object from the context
	ca, ok := ctx.Value(CookieAccessKeyCtx).(*CookieAccess)
	if !ok {
		return fmt.Errorf("failed to retrieve CookieAccess from context")
	}

	// Set the generated token as a cookie using the CookieAccess object
	ca.SetToken(token)
	return nil
}
*/
// GetIDFromContext retrieves the ID from the context.
func GetIDFromContext(ctx context.Context) (string, error) {
	id, ok := ctx.Value(IDContextKey).(string)
	if !ok {
		return "", fmt.Errorf("ID not found in context")
	}
	return id, nil
}
