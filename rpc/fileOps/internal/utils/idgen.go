package utils

import (
	"github.com/google/uuid"
	"strings"
)

type IdGenerator struct {
}

// GenerateUUID generates a new UUID
func (ig *IdGenerator) GenerateUUID() string {
	return uuid.NewString()
}

// GenerateFileID generates a unique file ID
func (ig *IdGenerator) GenerateFileID() string {
	u := uuid.New()
	return "file-" + strings.ReplaceAll(u.String(), "-", "")[:30]
}

// GenerateUserID generates a unique user ID
func (ig *IdGenerator) GenerateUserID() string {
	return "user-" + uuid.NewString()
}
