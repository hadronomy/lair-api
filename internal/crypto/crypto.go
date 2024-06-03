package crypto

import gonanoid "github.com/matoous/go-nanoid/v2"

var (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	size     = 16
)

// With the help of the go-nanoid library, we can generate a unique ID for each user.
//
// With this alphabet and a size of 16:
//
// ~4 million years or 30T IDs needed, 
// in order to have a 1% probability of at least one collision.
func GenerateID() (string, error) {
	return gonanoid.Generate(alphabet, size)
}
