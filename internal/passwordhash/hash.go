package passwordhash

import (
	"crypto/sha256"
	"encoding/hex"
)

// Calculate the password hash
// using Hex(SHA256(input)) transformation
func HashPassword(input string) string {
	sha256Hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sha256Hash[:])
}
