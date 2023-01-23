package util

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// A RandomInt returns a random integer in the range [min, max].
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// A RandomString returns a random string of the given length.
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

// A RandomOwner returns a random owner name.
func RandomOwner() string {
	return RandomString(6)
}

// A RandomMoney returns a random amount of money.
func RandomMoney() int64 {
	return RandomInt(0, 1000000)
}

// A RandomCurrency returns a random currency.
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP","KES"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
