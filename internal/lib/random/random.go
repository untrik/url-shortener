package random

import "math/rand/v2"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewRandomString(aliasLength int) string {
	aliasRune := make([]rune, 0, aliasLength)
	for i := 0; i < aliasLength; i++ {
		aliasRune = append(aliasRune, rune((letters[rand.IntN(len(letters))])))
	}
	return string(aliasRune)
}
