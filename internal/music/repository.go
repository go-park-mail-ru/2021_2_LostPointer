package music

import "2021_2_LostPointer/internal/models"

//go:generate moq -out ../mock/music_repo_db_mock.go -internal mock . MusicRepository:MockMusicRepository
type MusicRepository interface {
	GetRandomTracks(amount int, isAuthorized bool) ([]models.Track, error)
	GetRandomAlbums(amount int) ([]models.Album, error)
	GetRandomArtists(amount int) ([]models.Artist, error)
	GetRandomPlaylists(amount int) ([]models.Playlist, error)
}
