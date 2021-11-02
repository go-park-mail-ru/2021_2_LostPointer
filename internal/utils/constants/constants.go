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
	InvalidNicknameMessage                   = "The length of nickname must be from " + MinNicknameLength + " to " + MaxNicknameLength + " characters"
	InvalidEmailMessage                      = "Invalid email"
	EmailRegexPattern                        = `^[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+$`
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
	NoArtists                                = "No artists"
	InvalidParameter                         = "Invalid parameter"
	DatabaseNotResponding                    = "Database not responding"
	NoPlaylists                              = "No playlists"
	WrongCredentials                         = "Wrong credentials"

	SaltLength                   = 8
	BigAvatarHeight              = 500
	LittleAvatarHeight           = 150
	TracksDefaultAmountForArtist = 10
	AlbumsDefaultAmountForArtist = 8
	TracksCollectionLimit        = 10
	AlbumCollectionLimit         = 4
	PlaylistsCollectionLimit     = 4
	ArtistsCollectionLimit       = 4
	SiteID                       = 0
	CookieLifetime               = time.Hour * 24 * 30
)
