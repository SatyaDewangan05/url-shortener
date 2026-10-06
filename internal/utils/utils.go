package utils

import "math/rand"

// Helper Functions
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode(length int) string {

	randBytes := make([]byte, length)

	for i := range randBytes {
		randBytes[i] = chars[rand.Intn(len(chars))]
	}
	return string(randBytes)
}
