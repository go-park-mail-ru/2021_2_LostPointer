package music

import (
	"2021_2_LostPointer/internal/microservices/music/proto"
)

type Storage interface {
	RandomTracks(int64, bool) (*proto.Tracks, error)
	RandomAlbums(int64) (*proto.Albums, error)
	RandomArtists(int64) (*proto.Artists, error)
	GetArtistInfo(int64) (*proto.Artist, error)
	GetArtistTracks(int64, bool, int64) ([]*proto.Track, error)
	GetArtistAlbums(int64, int64) ([]*proto.Album, error)
	IncrementListenCount(int64) error
	AlbumData(int64) (*proto.AlbumPageResponse, error)
	AlbumTracks(int64, bool) ([]*proto.AlbumTrack, error)
	FindTracksByFullWord(string, bool) ([]*proto.Track, error)
	FindTracksByPartial(string, bool) ([]*proto.Track, error)
	FindArtists(string) ([]*proto.Artist, error)
	FindAlbums(string) ([]*proto.Album, error)
	GetUserPlaylists(int64) ([]*proto.PlaylistData, error)
	IsPlaylistOwner(int64, int64) (bool, error)
	GetPlaylistTracks(int64) ([]*proto.Track, error)
	GetPlaylistInfo(int64) (*proto.PlaylistData, error)
	DoesPlaylistExist(int64) (bool, error)
}
