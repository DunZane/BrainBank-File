package utils

import (
	"testing"
)

var generator = IdGenerator{}

func TestGenerateUUID(t *testing.T) {
	id := generator.GenerateUUID()

	t.Logf("GenerateUUID:%s", id)
}

func TestGenerateFileID(t *testing.T) {
	fileId := generator.GenerateFileID()

	t.Logf("GenerateFileID:%s", fileId)
}

func TestGenerateUserID(t *testing.T) {
	userId := generator.GenerateUserID()

	t.Logf("GenerateUserID:%s", userId)
}
