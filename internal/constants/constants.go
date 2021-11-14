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
	PasswordValidationNoLetterMessage        = "Password must contain as least one letter"
	InvalidNicknameLengthMessage                   = "The length of nickname must be from " + MinNicknameLength + " to " + MaxNicknameLength + " characters"
	InvalidNicknameMessage                   = "Nickname must contain letters, numbers and '_'"
	InvalidEmailMessage                      = "Invalid email"
	NotUniqueEmailMessage                    = "Email is not unique"
	NotUniqueNicknameMessage                 = "Nickname is not unique"
	WrongPasswordMessage                     = "Old password is wrong"
	OldPasswordFieldIsEmptyMessage           = "Old password field is empty"
	NewPasswordFieldIsEmptyMessage           = "New password field is empty"
	BigArtistPostfix                         = ".webp"
	BigAvatarPostfix                         = "_500px.webp"
	LittleAvatarPostfix                      = "_150px.webp"
	VideoPostfix                             = ".mp4"
	AvatarDefaultFileName                    = "default_avatar"
	UserIsNotAuthorizedMessage               = "User is not authorized"
	LoggedOutMessage                         = "Logged out"
	SettingsUploadedMessage                  = "Settings were uploaded successfully"
	UserCreatedMessage                       = "User was created successfully"
	UserAuthorizedMessage                    = "User is authorized"
	RequestIDTypeAssertionFailed             = "Type assertion for \"REQUEST_ID\" failed"
	UserIDTypeAssertionFailed                = "Type assertion for \"USER_ID\" failed"

	SaltLength                   = 8
	BigAvatarHeight              = 500
	LittleAvatarHeight           = 150
	TracksDefaultAmountForArtist = 10
	AlbumsDefaultAmountForArtist = 8
	TracksCollectionLimit        = 10
	AlbumCollectionLimit         = 4
	ArtistsCollectionLimit       = 4
	CookieLifetime               = time.Hour * 24 * 30
	TracksSearchAmount           = 10
	ArtistsSearchAmount          = 5
	AlbumsSearchAmount           = 5
	CSRFTokenLifetime            = 900
)
