package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHash(t *testing.T) {
	str := "abc"
	expectedHash := "a9993e364706816aba3e25717850c26c9cd0d89d"
	assert.Equal(t, GetHash(str), expectedHash)
}
