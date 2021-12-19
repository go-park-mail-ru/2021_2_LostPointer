package music

import (
	"2021_2_LostPointer/internal/microservices/music/proto"
	"2021_2_LostPointer/internal/models"
)

//go:generate moq -out ./mock/music_repo_mock.go -pkg mock . Storage:MockStorage
type Storage interface {
	RandomTracks(int64, int64, bool) ([]*proto.Track, error)
	RandomAlbums(int64) (*proto.Albums, error)
	RandomArtists(int64) (*proto.Artists, error)
	ArtistInfo(int64) (*proto.Artist, error)
	ArtistTracks(int64, int64, bool, int64) ([]*proto.Track, error)
	ArtistAlbums(int64, int64) ([]*proto.Album, error)
	IncrementListenCount(int64) error
	AlbumData(int64) (*proto.AlbumPageResponse, error)
	AlbumTracks(int64, int64, bool) ([]*proto.AlbumTrack, error)
	FindTracksByFullWord(string, int64, bool) ([]*proto.Track, error)
	FindTracksByPartial(string, int64, bool) ([]*proto.Track, error)
	FindArtists(string) ([]*proto.Artist, error)
	FindAlbums(string) ([]*proto.Album, error)
	UserPlaylists(int64) ([]*proto.PlaylistData, error)
	IsPlaylistOwner(int64, int64) (bool, error)
	IsPlaylistPublic(int64) (bool, error)
	PlaylistTracks(int64, int64) ([]*proto.Track, error)
	PlaylistInfo(int64) (*proto.PlaylistData, error)
	DoesPlaylistExist(int64) (bool, error)
	AddTrackToFavorites(int64, int64) error
	DeleteTrackFromFavorites(int64, int64) error
	GetFavorites(int64) ([]*proto.Track, error)
	IsTrackInFavorites(int64, int64) (bool, error)
	GetCompilation(int64) (*models.Selection, error)
	GetTracksCompilation(int64, []string) ([]string, error)
	GetFavoriteTracksID(int64) ([]string, error)
	StoreCompilation(int64, *models.Selection) error
	GetTracksByID([]string, int64, bool) ([]*proto.Track, error)
}
