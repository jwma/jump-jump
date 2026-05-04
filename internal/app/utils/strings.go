package utils

import (
	"math/rand/v2"
	"regexp"
	"strings"
)

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	return string(b)
}

var shortLinkIdMatcher = regexp.MustCompile("[a-zA-Z0-9]+")

func TrimShortLinkId(s string) string {
	return strings.Join(shortLinkIdMatcher.FindAllString(s, -1), "")
}
