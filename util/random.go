package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var result strings.Builder

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(len(alphabet))]
		result.WriteByte(c)
	}

	return result.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomBalance() int64 {
	return rand.Int63n(1000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD}
	n := rand.Intn(len(currencies))

	return currencies[n]
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}
