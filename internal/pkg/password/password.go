package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultCost is the default bcrypt cost
	DefaultCost = bcrypt.DefaultCost
)

var (
	// ErrMismatchedHashAndPassword is returned when password verification fails
	ErrMismatchedHashAndPassword = errors.New("mismatched hash and password")
)

// Hasher handles password hashing operations
type Hasher struct {
	cost int
}

// NewHasher creates a new password hasher
func NewHasher() *Hasher {
	return &Hasher{
		cost: DefaultCost,
	}
}

// NewHasherWithCost creates a new password hasher with custom cost
func NewHasherWithCost(cost int) *Hasher {
	return &Hasher{
		cost: cost,
	}
}

// Hash hashes a plain text password
func (h *Hasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// Verify verifies a plain text password against a hash
func (h *Hasher) Verify(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrMismatchedHashAndPassword
		}
		return fmt.Errorf("password verification failed: %w", err)
	}
	return nil
}

// IsValidPassword checks if a password meets basic requirements
func IsValidPassword(password string) bool {
	// At least 8 characters
	if len(password) < 8 {
		return false
	}
	// Maximum 72 characters (bcrypt limitation)
	if len(password) > 72 {
		return false
	}
	return true
}
