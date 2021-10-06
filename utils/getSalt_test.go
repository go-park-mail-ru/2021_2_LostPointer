package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRandomString(t *testing.T) {
	length := 10
	str := GetRandomString(length)
	assert.Equal(t, len(str), length)
}

func TestRandInt(t *testing.T) {
	min := -10
	max := 10
	res := RandInt(min, max)
	assert.True(t, res > min && res < max)
}
