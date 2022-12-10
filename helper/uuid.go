package helper

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	uuid := uuid.New()
	return uuid.String()
}

// get filename from inserted file path
func GetFilenameWithoutExtension(path string) string {
	split := strings.Split(path, "/")
	filename := split[len(split)-1]
	return filename[:len(filename)-4]
}
