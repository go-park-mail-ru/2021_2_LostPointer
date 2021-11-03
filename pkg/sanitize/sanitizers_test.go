package sanitize

import (
	"2021_2_LostPointer/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)
func TestSanitizeUserData(t *testing.T) {
	dangData := models.User{
		Nickname: "<script>nickname</script>",
		Email:    "<script>email</script>",
		Password: "password",
	}
	result := SanitizeUserData(dangData)
	expected := models.User{
		Nickname: "nickname",
		Email:    "email",
		Password: "password",
	}

	assert.Equal(t, expected, result)
}

func TestSanitizeEmail(t *testing.T) {
	dangEmail := "<script>email</script>"
	result := SanitizeEmail(dangEmail)
	expected := "email"

	assert.Equal(t, expected, result)
}

func TestSanitizeNickname(t *testing.T) {
	dangEmail := "<script>nickname</script>"
	result := SanitizeNickname(dangEmail)
	expected := "nickname"

	assert.Equal(t, expected, result)
}