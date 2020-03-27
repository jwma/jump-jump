package utils

import (
	"math/rand"
	"regexp"
	"strings"
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

var shortLinkIdMatcher = regexp.MustCompile("[a-zA-Z0-9]+")

func TrimShortLinkId(s string) string {
	return strings.Join(shortLinkIdMatcher.FindAllString(s, -1), "")
}
