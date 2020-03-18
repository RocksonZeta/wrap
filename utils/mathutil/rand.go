package mathutil

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var alpnum = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var alpnumLow = []byte("abcdefghijklmnopqrstuvwxyz0123456789")
var num = []byte("0123456789")

func RandomStr(n int, ignoreCase bool) string {
	var letters = alpnum
	if ignoreCase {
		letters = alpnumLow
	}
	return RandomIn(n, letters)
}
func RandomIn(n int, from []byte) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = from[rand.Intn(len(from))]
	}
	return string(b)
}

func RandomStr32() string {
	return RandomStr(32, true)
}

func RandomInt(n int) string {
	return RandomIn(n, num)
}
