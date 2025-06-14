package util

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestPasswordHash(t *testing.T) {
	inputPass := "123456"
	expectectedHash := "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"
	actualHash := HashPassword(inputPass)
	assert.Equal(t, expectectedHash, actualHash)
}

func TestPasswordHashEmpty(t *testing.T) {
	inputPass := ""
	expectectedHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	actualHash := HashPassword(inputPass)
	assert.Equal(t, expectectedHash, actualHash)
}
