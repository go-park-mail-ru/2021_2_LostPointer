package usecase

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/music/mock"
	"2021_2_LostPointer/pkg/music/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMusicUseCase_GetTracksForCollection(t *testing.T) {
	track := models.Track{
		Id:          1,
		Title:       "awa",
		Artist:      "awa",
		Album:       "awa",
		Explicit:    true,
		Genre:       "awa",
		Number:      1,
		ListenCount: 1,
		Duration:    1,
		Cover:       "awa",
	}

	tests := []struct {
		name          string
		dbMock        *mock.MockMusicRepositoryIFace
		amount        int
		expected      []models.Track
		expectedError bool
	}{
		{
			name: "get 4 tacks",
			dbMock: &mock.MockMusicRepositoryIFace{
				CreateTracksRequestWithParametersFunc: func(gettingWith uint8, parameters []string, distinctOn uint8) (string, error) {
					switch gettingWith {
					case repository.GettingWithGenres:
						request := `SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE g.name IN ('Pop', 'Rock')`
						return request, nil
					case repository.GettingWithID:
						request := `SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE tracks.id IN (1, 2, 3, 4)`
						return request, nil
					default:
						return "", nil
					}
				},
				GetTracksFunc: func(request string) ([]models.Track, error) {
					switch request {
					case `SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE g.name IN ('Pop', 'Rock')`:
						var tracks []models.Track
						for i := 0; i < 50; i++ {
							tracks = append(tracks, track)
							track.Id += 1
						}
						return tracks, nil
					case `SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE tracks.id IN (1, 2, 3, 4)`:
						var tracks []models.Track
						track.Id = 1
						for i := 0; i < 4; i++ {
							tracks = append(tracks, track)
						}
						return tracks, nil
					default:
						return nil, nil
					}

				},
			},
			amount:        4,
			expected:      []models.Track{track, track, track, track},
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			usecase := NewMusicUseCase(test.dbMock)
			result, err := usecase.GetTracksForCollection(test.amount, []string{"Pop", "Rock"})
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestMusicUseCase_GetAlbumsForCollection(t *testing.T) {
	album := models.Album{
		Id:             1,
		Title:          "awa",
		Year:           1,
		Artist:         "awa",
		ArtWork:        "awa",
		TracksCount:    1,
		TracksDuration: 1,
	}

	tests := []struct {
		name          string
		dbMock        *mock.MockMusicRepositoryIFace
		amount        int
		expected      []models.Album
		expectedError bool
	}{
		{
			name: "get albums",
			dbMock: &mock.MockMusicRepositoryIFace{
				CreateAlbumsDefaultRequestFunc: func(amount int) string {
					return `SELECT a.id, a.title, a.year, art.name,
									a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a
									LEFT JOIN artists art ON art.id = a.artist
									JOIN tracks t on t.album = a.id
									WHERE art.name = 'Земфира'
									GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count
									LIMIT 4`
				},
				GetAlbumsFunc: func(request string) ([]models.Album, error) {
					var albums []models.Album
					for i := 0; i < 4; i++ {
						albums = append(albums, album)
					}
					return albums, nil
				},
			},
			amount:        4,
			expected:      []models.Album{album, album, album, album},
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			usecase := NewMusicUseCase(test.dbMock)
			result, err := usecase.GetAlbumsForCollection(test.amount)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestMusicUseCase_GetArtistsForCollection(t *testing.T) {
	artist := models.Artist{
		Id:     1,
		Name:   "awa",
		Avatar: "awa",
	}

	tests := []struct {
		name          string
		dbMock        *mock.MockMusicRepositoryIFace
		amount        int
		expected      []models.Artist
		expectedError bool
	}{
		{
			name: "get artists",
			dbMock: &mock.MockMusicRepositoryIFace{
				CreateArtistsDefaultRequestFunc: func(amount int) string {
					return `SELECT artists.id, artists.name, artists.avatar FROM artists
									WHERE artists.avatar != 'frank_sinatra.jpg' LIMIT 4`
				},
				GetArtistsFunc: func(request string) ([]models.Artist, error) {
					var artists []models.Artist
					for i := 0; i < 4; i++ {
						artists = append(artists, artist)
					}
					return artists, nil
				},
			},
			amount:        4,
			expected:      []models.Artist{artist, artist, artist, artist},
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			usecase := NewMusicUseCase(test.dbMock)
			result, err := usecase.GetArtistsForCollection(test.amount)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestMusicUseCase_GetPlaylistsForCollection(t *testing.T) {
	playlist := models.Playlist{
		Id:   1,
		Name: "awa",
		User: 1,
	}

	tests := []struct {
		name          string
		dbMock        *mock.MockMusicRepositoryIFace
		amount        int
		expected      []models.Playlist
		expectedError bool
	}{
		{
			name: "get playlists",
			dbMock: &mock.MockMusicRepositoryIFace{
				CreatePlaylistsDefaultRequestFunc: func(amount int) string {
					return `SELECT playlists.id, playlists.title, playlists.user FROM playlists LIMIT 4`
				},
				GetPlaylistsFunc: func(request string) ([]models.Playlist, error) {
					var playlists []models.Playlist
					for i := 0; i < 4; i++ {
						playlists = append(playlists, playlist)
					}
					return playlists, nil
				},
			},
			amount:        4,
			expected:      []models.Playlist{playlist, playlist, playlist, playlist},
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			usecase := NewMusicUseCase(test.dbMock)
			result, err := usecase.GetPlaylistsForCollection(test.amount)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestMusicUseCase_GetMusicCollection(t *testing.T) {
	playlist := models.Playlist{
		Id:   1,
		Name: "awa",
		User: 1,
	}

	artist := models.Artist{
		Id:     1,
		Name:   "awa",
		Avatar: "awa",
	}

	album := models.Album{
		Id:             1,
		Title:          "awa",
		Year:           1,
		Artist:         "awa",
		ArtWork:        "awa",
		TracksCount:    1,
		TracksDuration: 1,
	}

	track := models.Track{
		Id:          1,
		Title:       "awa",
		Artist:      "awa",
		Album:       "awa",
		Explicit:    true,
		Genre:       "awa",
		Number:      1,
		ListenCount: 1,
		Duration:    1,
		Cover:       "awa",
	}

	collection := models.MusicCollection{
		Tracks:    []models.Track{track, track, track, track},
		Albums:    []models.Album{album, album, album, album},
		Artists:   []models.Artist{artist, artist, artist, artist},
		Playlists: []models.Playlist{playlist, playlist, playlist, playlist},
	}

	tests := []struct {
		name          string
		dbMock        *mock.MockMusicRepositoryIFace
		amount        int
		expected      models.MusicCollection
		expectedError bool
	}{
		{
			name: "get playlists",
			dbMock: &mock.MockMusicRepositoryIFace{
				CreatePlaylistsDefaultRequestFunc: func(amount int) string {
					return `SELECT playlists.id, playlists.title, playlists.user FROM playlists LIMIT 4`
				},
				GetPlaylistsFunc: func(request string) ([]models.Playlist, error) {
					var playlists []models.Playlist
					for i := 0; i < 4; i++ {
						playlists = append(playlists, playlist)
					}
					return playlists, nil
				},
				CreateArtistsDefaultRequestFunc: func(amount int) string {
					return `SELECT artists.id, artists.name, artists.avatar FROM artists
									WHERE artists.avatar != 'frank_sinatra.jpg' LIMIT 4`
				},
				GetArtistsFunc: func(request string) ([]models.Artist, error) {
					var artists []models.Artist
					for i := 0; i < 4; i++ {
						artists = append(artists, artist)
					}
					return artists, nil
				},
				CreateAlbumsDefaultRequestFunc: func(amount int) string {
					return `SELECT a.id, a.title, a.year, art.name,
									a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a
									LEFT JOIN artists art ON art.id = a.artist
									JOIN tracks t on t.album = a.id
									WHERE art.name = 'Земфира'
									GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count
									LIMIT 4`
				},
				GetAlbumsFunc: func(request string) ([]models.Album, error) {
					var albums []models.Album
					for i := 0; i < 4; i++ {
						albums = append(albums, album)
					}
					return albums, nil
				},
				CreateTracksRequestWithParametersFunc: func(gettingWith uint8, parameters []string, distinctOn uint8) (string, error) {
					switch gettingWith {
					case repository.GettingWithGenres:
						request := `SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE g.name IN ('Pop', 'Rock')`
						return request, nil
					case repository.GettingWithID:
						request := `SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE tracks.id IN (1, 2, 3, 4)`
						return request, nil
					default:
						return "", nil
					}
				},
				GetTracksFunc: func(request string) ([]models.Track, error) {
					switch request {
					case `SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE g.name IN ('Pop', 'Rock')`:
						var tracks []models.Track
						for i := 0; i < 50; i++ {
							tracks = append(tracks, track)
							track.Id += 1
						}
						return tracks, nil
					case `SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE tracks.id IN (1, 2, 3, 4)`:
						var tracks []models.Track
						track.Id = 1
						for i := 0; i < 4; i++ {
							tracks = append(tracks, track)
						}
						return tracks, nil
					default:
						return nil, nil
					}

				},
			},
			expected:      collection,
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			usecase := NewMusicUseCase(test.dbMock)
			result, err := usecase.GetMusicCollection()
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, *result)
			}
		})
	}
}
