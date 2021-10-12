package music

import "2021_2_LostPointer/pkg/models"

//go:generate moq -out ./mock/music_repo_db_mock.go -pkg mock . MusicRepositoryInterface:MockMusicRepositoryInterface
type MusicRepositoryInterface interface {
	IsGenreExist(genres []string) (bool, error)

	CreateTracksRequestWithParameters(gettingWith uint8, parameters interface{}, distinctOn uint8) string
	CreateAlbumsDefaultRequest(amount int) string
	CreateArtistsDefaultRequest(amount int) string
	CreatePlaylistsDefaultRequest(amount int) string

	GetTracks(request string) ([]models.Track, error)
	GetAlbums(request string) ([]models.Album, error)
	GetArtists(request string) ([]models.Artist, error)
	GetPlaylists(request string) ([]models.Playlist, error)
}