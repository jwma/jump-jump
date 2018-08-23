package utils

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyz")
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return string(b)
}
