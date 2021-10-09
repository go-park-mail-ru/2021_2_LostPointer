package utils

import (
	"2021_2_LostPointer/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePasswordValid(t *testing.T) {
	password := "AlexRulitTankom2005!"
	isValidPassword, _, err := ValidatePassword(password)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with password %s\n",
			err, password)
	}
	assert.True(t, isValidPassword)
}

func TestValidatePasswordInvalid(t *testing.T) {
	password := "AlexRulitTankom"
	isValidPassword, _, err := ValidatePassword(password)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with password %s\n",
			err, password)
	}
	assert.False(t, isValidPassword)
}

func TestValidateSignUpValid(t *testing.T) {
	user := &models.User{
		ID: 1,
		Email: "alexeikasenke@gmail.com",
		Password: "AlexRulitTankom2005!",
		Nickname: "Kasenka",
		Salt: GetRandomString(SaltLength),
	}
	isValid, _, err := ValidateSignUp(user)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, user)
	}
	assert.True(t, isValid)
}

func TestValidateSignUpInvalidEmail(t *testing.T) {
	user := &models.User{
		ID: 1,
		Email: "alexeikasenkegmail.com",
		Password: "AlexRulitTankom2005!",
		Nickname: "Kasenka",
		Salt: GetRandomString(SaltLength),
	}
	isValid, _, err := ValidateSignUp(user)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, user)
	}
	assert.False(t, isValid)
}

func TestValidateSignUpInvalidPassword(t *testing.T) {
	user := &models.User{
		ID: 1,
		Email: "alexeikasenke@gmail.com",
		Password: "AlexRulitTankom2005",
		Nickname: "Kasenka",
		Salt: GetRandomString(SaltLength),
	}
	isValid, _, err := ValidateSignUp(user)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, user)
	}
	assert.False(t, isValid)
}

func TestValidateSignUpInvalidName(t *testing.T) {
	user := &models.User{
		ID: 1,
		Email: "alexeikasenkegmail.com",
		Password: "AlexRulitTankom2005!",
		Nickname: "Kasenka_1",
		Salt: GetRandomString(SaltLength),
	}
	isValid, _, err := ValidateSignUp(user)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, user)
	}
	assert.False(t, isValid)
}
