package crypto

import gonanoid "github.com/matoous/go-nanoid/v2"

var (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	size     = 14
)

func GenerateID() (string, error) {
	return gonanoid.Generate(alphabet, size)
}
