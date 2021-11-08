package music

import (
	"2021_2_LostPointer/internal/microservices/music/proto"
)

type Storage interface {
	RandomTracks(int64, bool) (*proto.Tracks, error)
	RandomAlbums(int64) (*proto.Albums, error)
	RandomArtists(int64) (*proto.Artists, error)
}
