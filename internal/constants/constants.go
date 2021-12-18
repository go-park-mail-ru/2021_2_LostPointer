package constants

import "time"

const (
	// Расширения
	ImageExtension                = ".webp"
	VideoExtension                = ".mp4"
	UserAvatarExtension500px      = "_500px.webp"
	UserAvatarExtension150px      = "_150px.webp"
	PlaylistArtworkExtension384px = "_384px.webp"
	PlaylistArtworkExtension100px = "_100px.webp"

	// Atoi
	PasswordRequiredLength = "8"
	MinNicknameLength      = "3"
	MaxNicknameLength      = "15"
	MinPlaylistTitleLength = "3"
	MaxPlaylistTitleLength = "30"

	// Валидация
	PasswordInvalidLengthMessage      = "Password must contain at least " + PasswordRequiredLength + " characters"
	PasswordNoDigitMessage            = "Password must contain at least one digit"
	PasswordNoLetterMessage           = "Password must contain as least one letter"
	NicknameInvalidLengthMessage      = "The length of nickname must be from " + MinNicknameLength + " to " + MaxNicknameLength + " characters"
	NicknameInvalidSyntaxMessage      = "Nickname must contain letters or numbers or '_'"
	PlaylistTitleInvalidLengthMessage = "The length of title must be from " + MinPlaylistTitleLength + " to " + MaxPlaylistTitleLength + " characters"
	EmailInvalidSyntaxMessage         = "Invalid email"

	// Значения по умолчанию
	AvatarDefaultFileName          = "default_avatar"
	PlaylistArtworkDefaultFilename = "default_playlist_artwork"
	PlaylistArtworkDefaultColor    = "#8071c2"

	// Сообщения
	EmailNotUniqueMessage            = "Email is not unique"
	NicknameNotUniqueMessage         = "Nickname is not unique"
	WrongPasswordMessage             = "Old password is wrong"
	OldPasswordFieldIsEmptyMessage   = "Old password field is empty"
	NewPasswordFieldIsEmptyMessage   = "New password field is empty"
	UserIsNotAuthorizedMessage       = "User is not authorized"
	LoggedOutMessage                 = "Logged out"
	SettingsUploadedMessage          = "Settings were uploaded successfully"
	UserCreatedMessage               = "User was created successfully"
	UserAuthorizedMessage            = "User is authorized"
	RequestIDTypeAssertionFailed     = "Type assertion for \"REQUEST_ID\" failed"
	UserIDTypeAssertionFailed        = "Type assertion for \"USER_ID\" failed"
	PlaylistDeletedMessage           = "Playlist deleted"
	TrackAddedToPlaylistMessage      = "Track was successfully added to playlist"
	TrackAlreadyInPlaylistMessage    = "Track was already added to playlist"
	TrackDeletedFromPlaylistMessage  = "Track was successfully deleted from playlist"
	NotPlaylistOwnerMessage          = "You are not owner of this playlist"
	PlaylistNotFoundMessage          = "Playlist doesn't exist"
	PanicRecover                     = "Panic!"
	TrackAlreadyInFavorites          = "Track is already in favorites"
	TrackNotInFavorites              = "Track not in favorites"
	TrackAddedToFavoritesMessage     = "Track was successfully added to favorites"
	TrackDeletedFromFavoritesMessage = "Track was successfully deleted from favorites"
	TrackNotFound                    = "Track not found"

	// Ограничения/лимиты
	ArtistTracksSelectionAmount            = 10
	ArtistAlbumsSelectionAmount            = 8
	HomePageTracksSelectionAmount          = 10
	HomePageAlbumsSelectionAmount          = 4
	HomePageArtistsSelectionAmount         = 4
	SearchTracksAmount                     = 4
	SearchArtistsAmount                    = 4
	SearchAlbumsAmount                     = 3
	HomePageTracksSelectionFavoritesAmount = HomePageTracksSelectionAmount * 2 / 3

	// Прочее
	SaltLength        = 8
	CookieLifetime    = time.Hour * 24 * 30
	CSRFTokenLifetime = 900
)
