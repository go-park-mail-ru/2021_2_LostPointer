package models

import "2021_2_LostPointer/internal/microservices/music/proto"

type SearchResult struct {
	Tracks  []Track  `json:"tracks,omitempty"`
	Albums  []Album  `json:"albums,omitempty"`
	Artists []Artist `json:"artists,omitempty"`
}

func (s *SearchResult) BindProto(search *proto.FindResponse) {
	tracks := make([]Track, 0)
	for _, protoTrack := range search.Tracks {
		track := Track{}
		track.BindProto(protoTrack)
		tracks = append(tracks, track)
	}
	albums := make([]Album, 0)
	for _, protoAlbum := range search.Albums {
		album := Album{}
		album.BindProto(protoAlbum)
		albums = append(albums, album)
	}
	artists := make([]Artist, 0)
	for _, currentArtist := range search.Artists {
		artist := Artist{}
		artist.BindProto(currentArtist)
		artists = append(artists, artist)
	}

	bindedSearchResult := SearchResult{
		Tracks:  tracks,
		Artists: artists,
		Albums:  albums,
	}

	*s = bindedSearchResult
}
