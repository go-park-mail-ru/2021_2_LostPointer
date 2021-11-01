package sanitize

import (
	"2021_2_LostPointer/internal/models"
	"github.com/kennygrant/sanitize"
)

func SanitizeUserData(userData models.User) models.User {
	var sanitizedData models.User

	sanitizedData.Nickname = sanitize.HTML(userData.Nickname)
	sanitizedData.Email = sanitize.HTML(userData.Email)
	sanitizedData.Password = userData.Password

	return sanitizedData
}

func SanitizeEmail(email string) string {
	return sanitize.HTML(email)
}

func SanitizeNickname(nickname string) string {
	return sanitize.HTML(nickname)
}
