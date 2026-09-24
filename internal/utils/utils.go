package utils

import "math/rand"

// Helper Functions

func GenerateShortCode(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	randBytes := make([]byte, length)

	for i := range randBytes {
		n := rand.Intn(len(chars))
		randBytes[i] = chars[n]
	}
	return string(randBytes)
}
