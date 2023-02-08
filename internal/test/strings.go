package test

import (
	"math/rand"
	"strings"
)

const (
	charSet = "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
)

func RandomString(random *rand.Rand, length int) (result string) {
	buffer := strings.Builder{}
	for i := 0; i < length; i++ {
		buffer.WriteByte(charSet[random.Intn(len(charSet))])
	}
	return buffer.String()
}
