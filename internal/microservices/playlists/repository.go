package playlists

import "2021_2_LostPointer/internal/microservices/playlists/proto"

//go:generate moq -out ./mock/playlists_repo_mock.go -pkg mock . PlaylistsStorage:MockPlaylistsStorage
type PlaylistsStorage interface {
	CreatePlaylist(int64, string, string, string, bool) (*proto.CreatePlaylistResponse, error)
	GetOldPlaylistSettings(int64) (string, error)
	DeletePlaylist(int64) error
	AddTrack(int64, int64) error
	DeleteTrack(int64, int64) error
	IsAdded(int64, int64) (bool, error)
	IsOwner(int64, int64) (bool, error)
	DoesPlaylistExist(int64) (bool, error)
	UpdatePlaylistTitle(int64, string) error
	UpdatePlaylistArtwork(int64, string, string) error
	DeletePlaylistArtwork(int64) error
	UpdatePlaylistAccess(int64, bool) error
}
