package constants

import "time"

const (
	PasswordRequiredLength                   = "8"
	MinNicknameLength                        = "3"
	MaxNicknameLength                        = "15"
	PasswordValidationInvalidLengthMessage   = "Password must contain at least " + PasswordRequiredLength + " characters"
	PasswordValidationNoDigitMessage         = "Password must contain at least one digit"
	PasswordValidationNoUppercaseMessage     = "Password must contain at least one uppercase letter"
	PasswordValidationNoLowerCaseMessage     = "Password must contain at least one lowercase letter"
	PasswordValidationNoSpecialSymbolMessage = "Password must contain as least one special character"
	NickNameValidationInvalidLengthMessage   = "The length of nickname must be from " + MinNicknameLength + " to " + MaxNicknameLength + " characters"
	InvalidEmailMessage                      = "Invalid email"
	EmailRegexPattern                        = `[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`
	NotUniqueEmailMessage                    = "Email is not unique"
	NotUniqueNicknameMessage                 = "Nickname is not unique"
	WrongPasswordMessage                     = "Old password is wrong"
	OldPasswordFieldIsEmptyMessage           = "Old password field is empty"
	NewPasswordFieldIsEmptyMessage           = "New password field is empty"
	BigAvatarPostfix                         = "_500px.webp"
	LittleAvatarPostfix                      = "150px.webp"
	AvatarDefaultFileName                    = "default_avatar"
	UserIsNotAuthorizedMessage               = "User is not authorized"
	LoggedOutMessage                         = "Logged out"
	SettingsUploadedMessage                  = "Settings were uploaded successfully"
	UserCreatedMessage                       = "User was created successfully"
	UserAuthorizedMessage                    = "User is authorized"

	SaltLength         = 8
	BigAvatarHeight    = 500
	LittleAvatarHeight = 150

	NoArtists                = "No artists"
	InvalidParameter         = "Invalid parameter"
	DatabaseNotResponding    = "Database not responding"
	TracksDefaultAmountForArtist     = 20
	AlbumsDefaultAmountForArtist     = 8
	NoMusic                  = "No music today"
	TracksCollectionLimit    = 10
	AlbumCollectionLimit     = 4
	PlaylistsCollectionLimit = 4
	ArtistsCollectionLimit   = 4
	NoPlaylists              = "No playlists"

	SiteID                = 0

	CookieLifetime = time.Hour * 24 * 30

	WrongCredentials  = "Wrong credentials"
	NotUniqueEmail    = "Email is not unique"
	NotUniqueNickname = "Nickname is not unique"
)
