package util

import (
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

func GenerateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

func GetFilenameWithoutExtension(path string) string {
	split := strings.Split(path, "/")
	filename := split[len(split)-1]
	return filename[:len(filename)-4]
}
