package random

import (
	"math/rand"
	"time"
)

func GenerateRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())

	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	res := make([]byte, n)
	for i := range res {
		res[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(res)
}
