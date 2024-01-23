package random

import (
	"math/rand"
)

func GenerateRandomString(n int) string {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	res := make([]byte, n)
	for i := range res {
		res[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(res)
}
