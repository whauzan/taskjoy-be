package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents the JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// TokenManager handles JWT token operations
type TokenManager struct {
	secretKey     []byte
	expiryHours   int
	signingMethod jwt.SigningMethod
}

// NewTokenManager creates a new TokenManager
func NewTokenManager(secretKey string, expiryHours int) *TokenManager {
	return &TokenManager{
		secretKey:     []byte(secretKey),
		expiryHours:   expiryHours,
		signingMethod: jwt.SigningMethodHS256,
	}
}

// TokenResponse contains the generated token and its expiration time
type TokenResponse struct {
	Token     string
	ExpiresAt time.Time
}

// GenerateToken generates a new JWT token for the given user
func (tm *TokenManager) GenerateToken(userID uuid.UUID, email string) (*TokenResponse, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(tm.expiryHours) * time.Hour)

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "todo-api",
		},
	}

	token := jwt.NewWithClaims(tm.signingMethod, claims)
	signedToken, err := token.SignedString(tm.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &TokenResponse{
		Token:     signedToken,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateToken validates a JWT token and returns the claims
func (tm *TokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if token.Method.Alg() != tm.signingMethod.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tm.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// RefreshToken generates a new token with extended expiry
func (tm *TokenManager) RefreshToken(tokenString string) (*TokenResponse, error) {
	claims, err := tm.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Generate a new token with the same user info
	return tm.GenerateToken(claims.UserID, claims.Email)
}
