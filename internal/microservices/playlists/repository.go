package playlists

import "2021_2_LostPointer/internal/microservices/playlists/proto"

type Storage interface {
	CreatePlaylist(int64, string) (*proto.CreatePlaylistResponse, error)
	UpdatePlaylist(int64, string) error
	DeletePlaylist(int64) error
	UserPlaylists(int64) ([]*proto.Playlist, error)
	AddTrack(int64, int64) error
	IsAdded(int64, int64) (bool, error)
	IsOwner(int64, int64) (bool, error)
}
