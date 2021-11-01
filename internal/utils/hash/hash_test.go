package hash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRandomString(t *testing.T) {
	const length = 10

	randomStr := GetRandomString(length)

	assert.True(t, len(randomStr) > 0)
}

func TestRandInt(t *testing.T) {
	const (
		min = 3
		max = 2017
	)

	randInt := RandInt(min, max)

	assert.True(t, randInt > min)
}

func TestGetHash(t *testing.T) {
	const str = "varfwq"

	hash := GetHash(str)

	assert.True(t, len(hash) > 0)
	assert.NotEqual(t, hash, str)
}
