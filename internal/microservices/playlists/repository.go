package playlists

import "2021_2_LostPointer/internal/microservices/playlists/proto"

type Storage interface {
	CreatePlaylist(int64, string, string, string) (*proto.CreatePlaylistResponse, error)
	UpdatePlaylist(int64, string, string, string) error
	GetOldArtwork(int64) (string, error)
	DeletePlaylist(int64) error
	AddTrack(int64, int64) error
	DeleteTrack(int64, int64) error
	IsAdded(int64, int64) (bool, error)
	IsOwner(int64, int64) (bool, error)
	DoesPlaylistExist(int64) (bool, error)
	UpdatePlaylistTitle(int64, string) error
	UpdatePlaylistArtwork(int64, string, string) error
	DeletePlaylistArtwork(int64) error
}
