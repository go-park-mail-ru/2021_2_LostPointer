package validation

import (
	"regexp"
	"strconv"

	"github.com/asaskevich/govalidator"

	"2021_2_LostPointer/internal/constants"
)

func ValidatePassword(password string) (bool, string, error) {
	patterns := map[string]string{
		`^.{` + constants.PasswordRequiredLength + `,}$`: constants.PasswordInvalidLengthMessage,
		`[0-9]`:    constants.PasswordNoDigitMessage,
		`[a-zA-Z]`: constants.PasswordNoLetterMessage,
	}

	for pattern, errorMessage := range patterns {
		isValid, err := regexp.MatchString(pattern, password)
		if err != nil {
			return false, "", err
		}
		if !isValid {
			return false, errorMessage, err
		}
	}

	return true, "", nil
}

func ValidateRegisterCredentials(email string, password string, nickname string) (bool, string, error) {
	isEmailValid := govalidator.IsEmail(email)
	if !isEmailValid {
		return false, constants.EmailInvalidSyntaxMessage, nil
	}

	isNicknameValid, err := regexp.MatchString(`^([\wА-Яа-я]+)$`, nickname)
	if err != nil {
		return false, "", err
	}
	if !isNicknameValid {
		return false, constants.NicknameInvalidSyntaxMessage, nil
	}
	minLength, _ := strconv.Atoi(constants.MinNicknameLength)
	maxLength, _ := strconv.Atoi(constants.MaxNicknameLength)
	if len(nickname) < minLength || len(nickname) > maxLength {
		return false, constants.NicknameInvalidLengthMessage, nil
	}

	passwordValid, message, err := ValidatePassword(password)
	if err != nil {
		return false, "", err
	}
	if !passwordValid {
		return false, message, nil
	}

	return true, "", nil
}

func ValidatePlaylistTitle(title string) (bool, string, error) {
	minLength, _ := strconv.Atoi(constants.MinPlaylistTitleLength)
	maxLength, _ := strconv.Atoi(constants.MaxPlaylistTitleLength)

	if len(title) < minLength || len(title) > maxLength {
		return false, constants.PlaylistTitleInvalidLengthMessage, nil
	}

	return true, "", nil
}
