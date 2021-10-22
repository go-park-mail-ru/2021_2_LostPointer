package usecase

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"errors"
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
		File:        "awa",
		ListenCount: 1,
		Duration:    1,
		Lossless:    true,
		Cover:       "awa",
	}

	trackWithoutFile := track
	trackWithoutFile.File = ""

	tests := []struct {
		name          string
		amount        int
		isAuthorized  bool
		dbMock        *mock.MockMusicRepository
		expected      []models.Track
		expectedError bool
	}{
		{
			name:         "get 4 tracks",
			amount:       4,
			isAuthorized: true,
			dbMock: &mock.MockMusicRepository{
				GetRandomTracksFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					var tracks []models.Track
					for i := 0; i < amount; i++ {
						tracks = append(tracks, track)
					}
					return tracks, nil
				},
			},
			expected: func() []models.Track {
				var tracks []models.Track
				for i := 0; i < 4; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: false,
		},
		{
			name:         "get 10 tracks",
			amount:       10,
			isAuthorized: true,
			dbMock: &mock.MockMusicRepository{
				GetRandomTracksFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					var tracks []models.Track
					for i := 0; i < amount; i++ {
						tracks = append(tracks, track)
					}
					return tracks, nil
				},
			},
			expected: func() []models.Track {
				var tracks []models.Track
				for i := 0; i < 10; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: false,
		},
		{
			name:         "get 10 tracks unauthorized",
			amount:       10,
			isAuthorized: true,
			dbMock: &mock.MockMusicRepository{
				GetRandomTracksFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					var tracks []models.Track
					for i := 0; i < amount; i++ {
						tracks = append(tracks, trackWithoutFile)
					}
					return tracks, nil
				},
			},
			expected: func() []models.Track {
				var tracks []models.Track
				for i := 0; i < 10; i++ {
					tracks = append(tracks, trackWithoutFile)
				}
				return tracks
			}(),
			expectedError: false,
		},
		{
			name:         "get tracks error",
			amount:       1,
			isAuthorized: true,
			dbMock: &mock.MockMusicRepository{
				GetRandomTracksFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					return nil, errors.New("error")
				},
			},
			expected: func() []models.Track {
				var tracks []models.Track
				for i := 0; i < 1; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useCase := NewMusicUseCase(test.dbMock)
			result, err := useCase.GetTracksForCollection(test.amount, test.isAuthorized)
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
		amount        int
		dbMock        *mock.MockMusicRepository
		expected      []models.Album
		expectedError bool
	}{
		{
			name:   "get 4 albums",
			amount: 4,
			dbMock: &mock.MockMusicRepository{
				GetRandomAlbumsFunc: func(amount int) ([]models.Album, error) {
					var albums []models.Album
					for i := 0; i < amount; i++ {
						albums = append(albums, album)
					}
					return albums, nil
				},
			},
			expected: func() []models.Album {
				var albums []models.Album
				for i := 0; i < 4; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: false,
		},
		{
			name:   "get 10 albums",
			amount: 10,
			dbMock: &mock.MockMusicRepository{
				GetRandomAlbumsFunc: func(amount int) ([]models.Album, error) {
					var albums []models.Album
					for i := 0; i < amount; i++ {
						albums = append(albums, album)
					}
					return albums, nil
				},
			},
			expected: func() []models.Album {
				var albums []models.Album
				for i := 0; i < 10; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: false,
		},
		{
			name:   "get albums error",
			amount: 10,
			dbMock: &mock.MockMusicRepository{
				GetRandomAlbumsFunc: func(amount int) ([]models.Album, error) {
					return nil, errors.New("error")
				},
			},
			expected: func() []models.Album {
				var albums []models.Album
				for i := 0; i < 10; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useCase := NewMusicUseCase(test.dbMock)
			result, err := useCase.GetAlbumsForCollection(test.amount)
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
		amount        int
		dbMock        *mock.MockMusicRepository
		expected      []models.Artist
		expectedError bool
	}{
		{
			name:   "get 4 artists",
			amount: 4,
			dbMock: &mock.MockMusicRepository{
				GetRandomArtistsFunc: func(amount int) ([]models.Artist, error) {
					var artists []models.Artist
					for i := 0; i < amount; i++ {
						artists = append(artists, artist)
					}
					return artists, nil
				},
			},
			expected: func() []models.Artist {
				var artists []models.Artist
				for i := 0; i < 4; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: false,
		},
		{
			name:   "get 10 artists",
			amount: 10,
			dbMock: &mock.MockMusicRepository{
				GetRandomArtistsFunc: func(amount int) ([]models.Artist, error) {
					var artists []models.Artist
					for i := 0; i < amount; i++ {
						artists = append(artists, artist)
					}
					return artists, nil
				},
			},
			expected: func() []models.Artist {
				var artists []models.Artist
				for i := 0; i < 10; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: false,
		},
		{
			name:   "get artists error",
			amount: 10,
			dbMock: &mock.MockMusicRepository{
				GetRandomArtistsFunc: func(amount int) ([]models.Artist, error) {
					return nil, errors.New("error")
				},
			},
			expected: func() []models.Artist {
				var artists []models.Artist
				for i := 0; i < 10; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useCase := NewMusicUseCase(test.dbMock)
			result, err := useCase.GetArtistsForCollection(test.amount)
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
		amount        int
		dbMock        *mock.MockMusicRepository
		expected      []models.Playlist
		expectedError bool
	}{
		{
			name:   "get 4 playlists",
			amount: 4,
			dbMock: &mock.MockMusicRepository{
				GetRandomPlaylistsFunc: func(amount int) ([]models.Playlist, error) {
					var playlists []models.Playlist
					for i := 0; i < amount; i++ {
						playlists = append(playlists, playlist)
					}
					return playlists, nil
				},
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 4; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: false,
		},
		{
			name:   "get 10 playlists",
			amount: 10,
			dbMock: &mock.MockMusicRepository{
				GetRandomPlaylistsFunc: func(amount int) ([]models.Playlist, error) {
					var playlists []models.Playlist
					for i := 0; i < amount; i++ {
						playlists = append(playlists, playlist)
					}
					return playlists, nil
				},
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 10; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: false,
		},
		{
			name:   "get 10 error",
			amount: 10,
			dbMock: &mock.MockMusicRepository{
				GetRandomPlaylistsFunc: func(amount int) ([]models.Playlist, error) {
					return nil, errors.New("error")
				},
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 10; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useCase := NewMusicUseCase(test.dbMock)
			result, err := useCase.GetPlaylistsForCollection(test.amount)
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
	track := models.Track{
		Id:          1,
		Title:       "awa",
		Artist:      "awa",
		Album:       "awa",
		Explicit:    true,
		Genre:       "awa",
		Number:      1,
		File:        "awa",
		ListenCount: 1,
		Duration:    1,
		Lossless:    true,
		Cover:       "awa",
	}

	trackWithoutFile := track
	trackWithoutFile.File = ""

	album := models.Album{
		Id:             1,
		Title:          "awa",
		Year:           1,
		Artist:         "awa",
		ArtWork:        "awa",
		TracksCount:    1,
		TracksDuration: 1,
	}

	artist := models.Artist{
		Id:     1,
		Name:   "awa",
		Avatar: "awa",
	}

	playlist := models.Playlist{
		Id:   1,
		Name: "awa",
		User: 1,
	}

	amountForTracks := TracksCollectionLimit
	amountForAlbums := AlbumCollectionLimit
	amountForArtists := ArtistsCollectionLimit
	amountForPlaylists := PlaylistsCollectionLimit

	tests := []struct {
		name          string
		isAuthorized  bool
		dbMock        *mock.MockMusicRepository
		expected      models.MusicCollection
		expectedError bool
	}{
		{
			name: "get full collection",
			isAuthorized: true,
			dbMock: &mock.MockMusicRepository{
				GetRandomTracksFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					var tracks []models.Track
					for i := 0; i < amount; i++ {
						tracks = append(tracks, track)
					}
					return tracks, nil
				},
				GetRandomAlbumsFunc: func(amount int) ([]models.Album, error) {
					var albums []models.Album
					for i := 0; i < amount; i++ {
						albums = append(albums, album)
					}
					return albums, nil
				},
				GetRandomArtistsFunc: func(amount int) ([]models.Artist, error) {
					var artists []models.Artist
					for i := 0; i < amount; i++ {
						artists = append(artists, artist)
					}
					return artists, nil
				},
				GetRandomPlaylistsFunc: func(amount int) ([]models.Playlist, error) {
					var playlists []models.Playlist
					for i := 0; i < amount; i++ {
						playlists = append(playlists, playlist)
					}
					return playlists, nil
				},
			},
			expected: func() models.MusicCollection {
				var tracks []models.Track
				for i := 0; i < amountForTracks; i++ {
					tracks = append(tracks, track)
				}
				var albums []models.Album
				for i := 0; i < amountForAlbums; i++ {
					albums = append(albums, album)
				}
				var artists []models.Artist
				for i := 0; i < amountForArtists; i++ {
					artists = append(artists, artist)
				}
				var playlists []models.Playlist
				for i := 0; i < amountForPlaylists; i++ {
					playlists = append(playlists, playlist)
				}
				var collection = models.MusicCollection{
					Tracks:    tracks,
					Albums:    albums,
					Artists:   artists,
					Playlists: playlists,
				}
				return collection
			}(),
			expectedError: false,
		},
		{
			name: "get tracks error",
			isAuthorized: true,
			dbMock: &mock.MockMusicRepository{
				GetRandomTracksFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					return nil, errors.New("error")
				},
				GetRandomAlbumsFunc: func(amount int) ([]models.Album, error) {
					var albums []models.Album
					for i := 0; i < amount; i++ {
						albums = append(albums, album)
					}
					return albums, nil
				},
				GetRandomArtistsFunc: func(amount int) ([]models.Artist, error) {
					var artists []models.Artist
					for i := 0; i < amount; i++ {
						artists = append(artists, artist)
					}
					return artists, nil
				},
				GetRandomPlaylistsFunc: func(amount int) ([]models.Playlist, error) {
					var playlists []models.Playlist
					for i := 0; i < amount; i++ {
						playlists = append(playlists, playlist)
					}
					return playlists, nil
				},
			},
			expected: func() models.MusicCollection {
				var tracks []models.Track
				for i := 0; i < amountForTracks; i++ {
					tracks = append(tracks, track)
				}
				var albums []models.Album
				for i := 0; i < amountForAlbums; i++ {
					albums = append(albums, album)
				}
				var artists []models.Artist
				for i := 0; i < amountForArtists; i++ {
					artists = append(artists, artist)
				}
				var playlists []models.Playlist
				for i := 0; i < amountForPlaylists; i++ {
					playlists = append(playlists, playlist)
				}
				var collection = models.MusicCollection{
					Tracks:    tracks,
					Albums:    albums,
					Artists:   artists,
					Playlists: playlists,
				}
				return collection
			}(),
			expectedError: true,
		},
		{
			name: "get albums error",
			isAuthorized: true,
			dbMock: &mock.MockMusicRepository{
				GetRandomTracksFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					var tracks []models.Track
					for i := 0; i < amount; i++ {
						tracks = append(tracks, track)
					}
					return tracks, nil
				},
				GetRandomAlbumsFunc: func(amount int) ([]models.Album, error) {
					return nil, errors.New("error")
				},
				GetRandomArtistsFunc: func(amount int) ([]models.Artist, error) {
					var artists []models.Artist
					for i := 0; i < amount; i++ {
						artists = append(artists, artist)
					}
					return artists, nil
				},
				GetRandomPlaylistsFunc: func(amount int) ([]models.Playlist, error) {
					var playlists []models.Playlist
					for i := 0; i < amount; i++ {
						playlists = append(playlists, playlist)
					}
					return playlists, nil
				},
			},
			expected: func() models.MusicCollection {
				var tracks []models.Track
				for i := 0; i < amountForTracks; i++ {
					tracks = append(tracks, track)
				}
				var albums []models.Album
				for i := 0; i < amountForAlbums; i++ {
					albums = append(albums, album)
				}
				var artists []models.Artist
				for i := 0; i < amountForArtists; i++ {
					artists = append(artists, artist)
				}
				var playlists []models.Playlist
				for i := 0; i < amountForPlaylists; i++ {
					playlists = append(playlists, playlist)
				}
				var collection = models.MusicCollection{
					Tracks:    tracks,
					Albums:    albums,
					Artists:   artists,
					Playlists: playlists,
				}
				return collection
			}(),
			expectedError: true,
		},
		{
			name: "get artists error",
			isAuthorized: true,
			dbMock: &mock.MockMusicRepository{
				GetRandomTracksFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					var tracks []models.Track
					for i := 0; i < amount; i++ {
						tracks = append(tracks, track)
					}
					return tracks, nil
				},
				GetRandomAlbumsFunc: func(amount int) ([]models.Album, error) {
					var albums []models.Album
					for i := 0; i < amount; i++ {
						albums = append(albums, album)
					}
					return albums, nil
				},
				GetRandomArtistsFunc: func(amount int) ([]models.Artist, error) {
					return nil, errors.New("error")
				},
				GetRandomPlaylistsFunc: func(amount int) ([]models.Playlist, error) {
					var playlists []models.Playlist
					for i := 0; i < amount; i++ {
						playlists = append(playlists, playlist)
					}
					return playlists, nil
				},
			},
			expected: func() models.MusicCollection {
				var tracks []models.Track
				for i := 0; i < amountForTracks; i++ {
					tracks = append(tracks, track)
				}
				var albums []models.Album
				for i := 0; i < amountForAlbums; i++ {
					albums = append(albums, album)
				}
				var artists []models.Artist
				for i := 0; i < amountForArtists; i++ {
					artists = append(artists, artist)
				}
				var playlists []models.Playlist
				for i := 0; i < amountForPlaylists; i++ {
					playlists = append(playlists, playlist)
				}
				var collection = models.MusicCollection{
					Tracks:    tracks,
					Albums:    albums,
					Artists:   artists,
					Playlists: playlists,
				}
				return collection
			}(),
			expectedError: true,
		},
		{
			name: "get playlists error",
			isAuthorized: true,
			dbMock: &mock.MockMusicRepository{
				GetRandomTracksFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					var tracks []models.Track
					for i := 0; i < amount; i++ {
						tracks = append(tracks, track)
					}
					return tracks, nil
				},
				GetRandomAlbumsFunc: func(amount int) ([]models.Album, error) {
					var albums []models.Album
					for i := 0; i < amount; i++ {
						albums = append(albums, album)
					}
					return albums, nil
				},
				GetRandomArtistsFunc: func(amount int) ([]models.Artist, error) {
					var artists []models.Artist
					for i := 0; i < amount; i++ {
						artists = append(artists, artist)
					}
					return artists, nil
				},
				GetRandomPlaylistsFunc: func(amount int) ([]models.Playlist, error) {
					return nil, errors.New("error")
				},
			},
			expected: func() models.MusicCollection {
				var tracks []models.Track
				for i := 0; i < amountForTracks; i++ {
					tracks = append(tracks, track)
				}
				var albums []models.Album
				for i := 0; i < amountForAlbums; i++ {
					albums = append(albums, album)
				}
				var artists []models.Artist
				for i := 0; i < amountForArtists; i++ {
					artists = append(artists, artist)
				}
				var playlists []models.Playlist
				for i := 0; i < amountForPlaylists; i++ {
					playlists = append(playlists, playlist)
				}
				var collection = models.MusicCollection{
					Tracks:    tracks,
					Albums:    albums,
					Artists:   artists,
					Playlists: playlists,
				}
				return collection
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useCase := NewMusicUseCase(test.dbMock)
			result, err := useCase.GetMusicCollection(test.isAuthorized)
			if test.expectedError {
				assert.Error(t, err.OriginalError)
			} else {
				assert.Equal(t, test.expected, *result)
			}
		})
	}
}

