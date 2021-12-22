package repository

/*import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"

	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/music/proto"
	"2021_2_LostPointer/pkg/wrapper"
	_ "database/sql"
	"database/sql/driver"
	"errors"
	"log"
	"os"
	"regexp"
	"testing"
)

//nolint:cyclop
func TestMusicStorage_RandomTracks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const userID = 1

	track := &proto.Track{
		ID:          1,
		Title:       "testTrackTitle",
		Explicit:    true,
		Genre:       "testGenre",
		Number:      2,
		File:        "testFile",
		ListenCount: 3,
		Duration:    4,
		Lossless:    true,
		Album: &proto.Album{
			ID:           5,
			Title:        "testAlbumTitle",
			Artwork:      "testAlbumArtwork",
			ArtworkColor: "testArtWOrkColor",
		},
		Artist: &proto.Artist{
			ID:   6,
			Name: "testArtistName",
		},
		IsInFavorites: true,
	}
	trackWithoutFile := new(proto.Track)
	_ = copier.Copy(trackWithoutFile, track)
	trackWithoutFile.File = ""

	tests := []struct {
		name          string
		amount        int64
		isAuthorized  bool
		mock          func()
		expected      *proto.Tracks
		expectedError bool
	}{
		{
			name:         "get 4 random tracks",
			amount:       4,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		ORDER BY RANDOM() DESC LIMIT $2`)).WithArgs(driver.Value(userID), driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() *proto.Tracks {
				var amount = 4
				var tracks = new(proto.Tracks)
				tracks.Tracks = make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks.Tracks = append(tracks.Tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "get 10 random tracks",
			amount:       10,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < 10; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		ORDER BY RANDOM() DESC LIMIT $2`)).WithArgs(driver.Value(userID), driver.Value(10)).WillReturnRows(rows)
			},
			expected: func() *proto.Tracks {
				var amount = 10
				var tracks = new(proto.Tracks)
				tracks.Tracks = make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks.Tracks = append(tracks.Tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "get 100 random tracks",
			amount:       100,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < 100; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		ORDER BY RANDOM() DESC LIMIT $2`)).WithArgs(driver.Value(userID), driver.Value(100)).WillReturnRows(rows)
			},
			expected: func() *proto.Tracks {
				var amount = 100
				var tracks = new(proto.Tracks)
				tracks.Tracks = make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks.Tracks = append(tracks.Tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "query returns error",
			amount:       1,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < 1; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		ORDER BY RANDOM() DESC LIMIT $2`)).WithArgs(driver.Value(userID), driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected: func() *proto.Tracks {
				var tracks = new(proto.Tracks)
				tracks.Tracks = make([]*proto.Track, 0, 1)
				tracks.Tracks = append(tracks.Tracks, track)
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:         "scan returns error",
			amount:       1,
			isAuthorized: true,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		ORDER BY RANDOM() DESC LIMIT $2`)).WithArgs(driver.Value(userID), driver.Value(1)).WillReturnRows(rows)
			},
			expected: func() *proto.Tracks {
				var tracks = new(proto.Tracks)
				tracks.Tracks = make([]*proto.Track, 0, 1)
				tracks.Tracks = append(tracks.Tracks, track)
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:   "get 4 random tracks unauthorized",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, "", track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		ORDER BY RANDOM() DESC LIMIT $2`)).WithArgs(driver.Value(userID), driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() *proto.Tracks {
				var amount = 4
				var tracks = new(proto.Tracks)
				tracks.Tracks = make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks.Tracks = append(tracks.Tracks, trackWithoutFile)
				}
				return tracks
			}(),
		},
		{
			name:         "rows.Err() returns error",
			amount:       4,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"}).RowError(0, errors.New("error"))
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		ORDER BY RANDOM() DESC LIMIT $2`)).WithArgs(driver.Value(userID), driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() *proto.Tracks {
				var amount = 4
				var tracks = new(proto.Tracks)
				tracks.Tracks = make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks.Tracks = append(tracks.Tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.RandomTracks(currentTest.amount, userID, currentTest.isAuthorized)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_RandomAlbums(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	album := &proto.Album{
		ID:             1,
		Title:          "testTitle",
		Year:           2,
		Artist:         "testArtist",
		Artwork:        "testArtwork",
		TracksAmount:   3,
		TracksDuration: 4,
		ArtworkColor:   "testArtWorkColor",
	}

	tests := []struct {
		name          string
		amount        int64
		mock          func()
		expected      *proto.Albums
		expectedError bool
	}{
		{
			name:   "get 4 random albums",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.track_count", "artwork_color", "art.name", "tracksDuration"})
				for i := 0; i < 4; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.TracksAmount, album.ArtworkColor, album.Artist, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb") + ", " +
					wrapper.Wrapper([]string{"name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
					`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY RANDOM()
		LIMIT $1
		`)).WithArgs(driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() *proto.Albums {
				var amount = 4
				var albums = new(proto.Albums)
				albums.Albums = make([]*proto.Album, 0, amount)
				for i := 0; i < amount; i++ {
					albums.Albums = append(albums.Albums, album)
				}
				return albums
			}(),
		},
		{
			name:   "get 10 random albums",
			amount: 10,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.track_count", "artwork_color", "art.name", "tracksDuration"})
				for i := 0; i < 10; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.TracksAmount, album.ArtworkColor, album.Artist, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb") + ", " +
					wrapper.Wrapper([]string{"name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
					`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY RANDOM()
		LIMIT $1
		`)).WithArgs(driver.Value(10)).WillReturnRows(rows)
			},
			expected: func() *proto.Albums {
				var amount = 10
				var albums = new(proto.Albums)
				albums.Albums = make([]*proto.Album, 0, amount)
				for i := 0; i < amount; i++ {
					albums.Albums = append(albums.Albums, album)
				}
				return albums
			}(),
		},
		{
			name:   "get 100 random albums",
			amount: 100,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.track_count", "artwork_color", "art.name", "tracksDuration"})
				for i := 0; i < 100; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.TracksAmount, album.ArtworkColor, album.Artist, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb") + ", " +
					wrapper.Wrapper([]string{"name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
					`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY RANDOM()
		LIMIT $1
		`)).WithArgs(driver.Value(100)).WillReturnRows(rows)
			},
			expected: func() *proto.Albums {
				var amount = 100
				var albums = new(proto.Albums)
				albums.Albums = make([]*proto.Album, 0, amount)
				for i := 0; i < amount; i++ {
					albums.Albums = append(albums.Albums, album)
				}
				return albums
			}(),
		},
		{
			name:   "query returns error",
			amount: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.track_count", "artwork_color", "art.name", "tracksDuration"})
				for i := 0; i < 1; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.TracksAmount, album.ArtworkColor, album.Artist, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb") + ", " +
					wrapper.Wrapper([]string{"name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
					`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY RANDOM()
		LIMIT $1
		`)).WithArgs(driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected: func() *proto.Albums {
				var albums = new(proto.Albums)
				albums.Albums = make([]*proto.Album, 0, 1)
				albums.Albums = append(albums.Albums, album)
				return albums
			}(),
			expectedError: true,
		},
		{
			name:   "scan returns error",
			amount: 1,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.track_count", "artwork_color", "art.name", "tracksDuration", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.TracksAmount, album.ArtworkColor, album.Artist, album.TracksDuration, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb") + ", " +
					wrapper.Wrapper([]string{"name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
					`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY RANDOM()
		LIMIT $1
		`)).WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			expected: func() *proto.Albums {
				var albums = new(proto.Albums)
				albums.Albums = make([]*proto.Album, 0, 1)
				albums.Albums = append(albums.Albums, album)
				return albums
			}(),
			expectedError: true,
		},
		{
			name:   "rows returns error",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.track_count", "artwork_color", "art.name", "tracksDuration"}).RowError(1, errors.New("error"))
				for i := 0; i < 4; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.TracksAmount, album.ArtworkColor, album.Artist, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb") + ", " +
					wrapper.Wrapper([]string{"name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
					`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY RANDOM()
		LIMIT $1
		`)).WithArgs(driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() *proto.Albums {
				var amount = 4
				var albums = new(proto.Albums)
				albums.Albums = make([]*proto.Album, 0, amount)
				for i := 0; i < amount; i++ {
					albums.Albums = append(albums.Albums, album)
				}
				return albums
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.RandomAlbums(currentTest.amount)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_RandomArtists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	tracks := make([]*proto.Track, 0)
	albums := make([]*proto.Album, 0)

	artist := &proto.Artist{
		ID:     1,
		Name:   "testName",
		Avatar: "testAvatar",
		Video:  "",
		Tracks: tracks,
		Albums: albums,
	}

	tests := []struct {
		name          string
		amount        int64
		mock          func()
		expected      *proto.Artists
		expectedError bool
	}{
		{
			name:   "get 4 random artists",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				for i := 0; i < 4; i++ {
					rows.AddRow(artist.ID, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
		id, name, avatar
		FROM artists
		ORDER BY RANDOM()
		LIMIT $1
	`)).WithArgs(driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() *proto.Artists {
				var amount = 4
				var artists = new(proto.Artists)
				artists.Artists = make([]*proto.Artist, 0, amount)
				for i := 0; i < amount; i++ {
					artists.Artists = append(artists.Artists, artist)
				}
				return artists
			}(),
		},
		{
			name:   "get 10 random artists",
			amount: 10,
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				for i := 0; i < 10; i++ {
					rows.AddRow(artist.ID, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
		id, name, avatar
		FROM artists
		ORDER BY RANDOM()
		LIMIT $1
	`)).WithArgs(driver.Value(10)).WillReturnRows(rows)
			},
			expected: func() *proto.Artists {
				var amount = 10
				var artists = new(proto.Artists)
				artists.Artists = make([]*proto.Artist, 0, amount)
				for i := 0; i < amount; i++ {
					artists.Artists = append(artists.Artists, artist)
				}
				return artists
			}(),
		},
		{
			name:   "get 100 random artists",
			amount: 100,
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				for i := 0; i < 100; i++ {
					rows.AddRow(artist.ID, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
		id, name, avatar
		FROM artists
		ORDER BY RANDOM()
		LIMIT $1
	`)).WithArgs(driver.Value(100)).WillReturnRows(rows)
			},
			expected: func() *proto.Artists {
				var amount = 100
				var artists = new(proto.Artists)
				artists.Artists = make([]*proto.Artist, 0, amount)
				for i := 0; i < amount; i++ {
					artists.Artists = append(artists.Artists, artist)
				}
				return artists
			}(),
		},
		{
			name:   "query returns error",
			amount: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				for i := 0; i < 1; i++ {
					rows.AddRow(artist.ID, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
		id, name, avatar
		FROM artists
		ORDER BY RANDOM()
		LIMIT $1
	`)).WithArgs(driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected: func() *proto.Artists {
				var artists = new(proto.Artists)
				artists.Artists = make([]*proto.Artist, 0, 1)
				artists.Artists = append(artists.Artists, artist)
				return artists
			}(),
			expectedError: true,
		},
		{
			name:   "scan returns error",
			amount: 1,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(artist.ID, artist.Name, artist.Avatar, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
		id, name, avatar
		FROM artists
		ORDER BY RANDOM()
		LIMIT $1
	`)).WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			expected: func() *proto.Artists {
				var artists = new(proto.Artists)
				artists.Artists = make([]*proto.Artist, 0, 1)
				artists.Artists = append(artists.Artists, artist)
				return artists
			}(),
			expectedError: true,
		},
		{
			name:   "rows.Err() returns error",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"}).RowError(1, errors.New("error"))
				for i := 0; i < 4; i++ {
					rows.AddRow(artist.ID, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
		id, name, avatar
		FROM artists
		ORDER BY RANDOM()
		LIMIT $1
	`)).WithArgs(driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() *proto.Artists {
				var amount = 4
				var artists = new(proto.Artists)
				artists.Artists = make([]*proto.Artist, 0, amount)
				for i := 0; i < amount; i++ {
					artists.Artists = append(artists.Artists, artist)
				}
				return artists
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.RandomArtists(currentTest.amount)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestMusicStorage_ArtistInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	artist := &proto.Artist{
		ID:     1,
		Name:   "testName",
		Avatar: "testAvatar",
		Video:  "testVideo",
	}
	expectedArtistWithVideo := &proto.Artist{
		ID:     1,
		Name:   "testName",
		Avatar: os.Getenv("ARTISTS_ROOT_PREFIX") + "testAvatar" + constants.ImageExtension,
		Video:  os.Getenv("MOV_ROOT_PREFIX") + "testVideo" + constants.VideoExtension,
	}
	expectedArtistWithoutVideo := &proto.Artist{
		ID:     1,
		Name:   "testName",
		Avatar: os.Getenv("ARTISTS_ROOT_PREFIX") + "testAvatar" + constants.ImageExtension,
		Video:  "",
	}
	var artistID int64 = 1

	tests := []struct {
		name          string
		mock          func()
		expected      *proto.Artist
		expectedError bool
	}{
		{
			name: "get artist info with video",
			mock: func() {
				row := mock.NewRows([]string{"id", "name", "avatar", "video"})
				row.AddRow(artist.ID, artist.Name, artist.Avatar, artist.Video)
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, avatar, video
		FROM artists
		WHERE id = $1
	`)).WillReturnRows(row)
			},
			expected: expectedArtistWithVideo,
		},
		{
			name: "get artist info without video",
			mock: func() {
				video := ""
				row := mock.NewRows([]string{"id", "name", "avatar", "video"})
				row.AddRow(artist.ID, artist.Name, artist.Avatar, video)
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, avatar, video
		FROM artists
		WHERE id = $1
	`)).WillReturnRows(row)
			},
			expected: expectedArtistWithoutVideo,
		},
		{
			name: "query returns error",
			mock: func() {
				row := mock.NewRows([]string{"id", "name", "avatar", "video"})
				row.AddRow(artist.ID, artist.Name, artist.Avatar, artist.Video)
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, avatar, video
		FROM artists
		WHERE id = $1
	`)).WillReturnError(errors.New("error"))
			},
			expected:      expectedArtistWithoutVideo,
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.ArtistInfo(artistID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_ArtistTracks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const userID = 1

	track := &proto.Track{
		ID:          1,
		Title:       "testTrackTitle",
		Explicit:    true,
		Genre:       "testGenre",
		Number:      2,
		File:        "testFile",
		ListenCount: 3,
		Duration:    4,
		Lossless:    true,
		Album: &proto.Album{
			ID:           5,
			Title:        "testAlbumTitle",
			Artwork:      "testAlbumArtwork",
			ArtworkColor: "testArtworkColor",
		},
		Artist:        new(proto.Artist),
		IsInFavorites: true,
	}
	trackWithoutFile := new(proto.Track)
	_ = copier.Copy(trackWithoutFile, track)
	trackWithoutFile.File = ""

	var artistID int64 = 6

	tests := []struct {
		name          string
		artistID      int64
		amount        int64
		isAuthorized  bool
		mock          func()
		expected      []*proto.Track
		expectedError bool
	}{
		{
			name:         "get 4 random tracks by artist ID",
			amount:       4,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "g.name", "favorites"})
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.artist = $2
		ORDER BY t.listen_count DESC LIMIT $3
		`)).WithArgs(driver.Value(userID), driver.Value(artistID), driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				var amount = 4
				tracks := make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "get 10 random tracks by artist ID",
			amount:       10,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "g.name", "favorites"})
				for i := 0; i < 10; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.artist = $2
		ORDER BY t.listen_count DESC LIMIT $3
		`)).WithArgs(driver.Value(userID), driver.Value(artistID), driver.Value(10)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				var amount = 10
				tracks := make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "get 100 random tracks by artist ID",
			amount:       100,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "g.name", "favorites"})
				for i := 0; i < 100; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.artist = $2
		ORDER BY t.listen_count DESC LIMIT $3
		`)).WithArgs(driver.Value(userID), driver.Value(artistID), driver.Value(100)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				var amount = 100
				tracks := make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "query returns error",
			amount:       1,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "g.name", "favorites"})
				for i := 0; i < 1; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.artist = $2
		ORDER BY t.listen_count DESC LIMIT $3
		`)).WithArgs(driver.Value(userID), driver.Value(artistID), driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, 1)
				tracks = append(tracks, track)
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:         "scan returns error",
			amount:       1,
			isAuthorized: true,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "g.name", "favorites", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Genre, track.IsInFavorites, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.artist = $2
		ORDER BY t.listen_count DESC LIMIT $3
		`)).WithArgs(driver.Value(userID), driver.Value(artistID), driver.Value(1)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, 1)
				tracks = append(tracks, track)
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:   "get 4 random tracks unauthorized",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "g.name", "favorites"})
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, "", track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.artist = $2
		ORDER BY t.listen_count DESC LIMIT $3
		`)).WithArgs(driver.Value(userID), driver.Value(artistID), driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				var amount = 4
				tracks := make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, trackWithoutFile)
				}
				return tracks
			}(),
		},
		{
			name:         "rows.Err() returns error",
			amount:       4,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "g.name", "favorites"}).RowError(1, errors.New("error"))
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.artist = $2
		ORDER BY t.listen_count DESC LIMIT $3
		`)).WithArgs(driver.Value(userID), driver.Value(artistID), driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				var amount = 4
				tracks := make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.ArtistTracks(artistID, userID, currentTest.isAuthorized, currentTest.amount)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_ArtistAlbums(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	album := &proto.Album{
		ID:             1,
		Title:          "testTitle",
		Year:           2,
		Artwork:        "testArtwork",
		ArtworkColor:   "testArtWorkColor",
		TracksDuration: 3,
	}

	var artistID int64 = 1

	tests := []struct {
		name          string
		amount        int64
		mock          func()
		expected      []*proto.Album
		expectedError bool
	}{
		{
			name:   "get 4 random albums",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.artwork_color", "tracksDuration"})
				for i := 0; i < 4; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.ArtworkColor, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "artwork_color"}, "alb")+", SUM(t.duration) AS tracksDuration"+
					`
		FROM albums alb
		JOIN tracks t ON alb.id = t.album
		WHERE alb.artist = $1
		GROUP BY alb.id, alb.title, alb.year, alb.artwork, alb.track_count
		ORDER BY alb.year DESC LIMIT $2
		`)).WithArgs(driver.Value(artistID), driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() []*proto.Album {
				var amount = 4
				albums := make([]*proto.Album, 0, amount)
				for i := 0; i < amount; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
		},
		{
			name:   "get 10 random albums",
			amount: 10,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.artwork_color", "tracksDuration"})
				for i := 0; i < 10; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.ArtworkColor, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "artwork_color"}, "alb")+", SUM(t.duration) AS tracksDuration"+
					`
		FROM albums alb
		JOIN tracks t ON alb.id = t.album
		WHERE alb.artist = $1
		GROUP BY alb.id, alb.title, alb.year, alb.artwork, alb.track_count
		ORDER BY alb.year DESC LIMIT $2
		`)).WithArgs(driver.Value(artistID), driver.Value(10)).WillReturnRows(rows)
			},
			expected: func() []*proto.Album {
				var amount = 10
				albums := make([]*proto.Album, 0, amount)
				for i := 0; i < amount; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
		},
		{
			name:   "get 100 random albums",
			amount: 100,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.artwork_color", "tracksDuration"})
				for i := 0; i < 100; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.ArtworkColor, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "artwork_color"}, "alb")+", SUM(t.duration) AS tracksDuration"+
					`
		FROM albums alb
		JOIN tracks t ON alb.id = t.album
		WHERE alb.artist = $1
		GROUP BY alb.id, alb.title, alb.year, alb.artwork, alb.track_count
		ORDER BY alb.year DESC LIMIT $2
		`)).WithArgs(driver.Value(artistID), driver.Value(100)).WillReturnRows(rows)
			},
			expected: func() []*proto.Album {
				var amount = 100
				albums := make([]*proto.Album, 0, amount)
				for i := 0; i < amount; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
		},
		{
			name:   "query returns error",
			amount: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.artwork_color", "tracksDuration"})
				for i := 0; i < 1; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.ArtworkColor, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "artwork_color"}, "alb")+", SUM(t.duration) AS tracksDuration"+
					`
		FROM albums alb
		JOIN tracks t ON alb.id = t.album
		WHERE alb.artist = $1
		GROUP BY alb.id, alb.title, alb.year, alb.artwork, alb.track_count
		ORDER BY alb.year DESC LIMIT $2
		`)).WithArgs(driver.Value(artistID), driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected: func() []*proto.Album {
				albums := make([]*proto.Album, 0, 1)
				albums = append(albums, album)
				return albums
			}(),
			expectedError: true,
		},
		{
			name:   "scan returns error",
			amount: 1,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.artwork_color", "tracksDuration", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.ArtworkColor, album.TracksDuration, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "artwork_color"}, "alb")+", SUM(t.duration) AS tracksDuration"+
					`
		FROM albums alb
		JOIN tracks t ON alb.id = t.album
		WHERE alb.artist = $1
		GROUP BY alb.id, alb.title, alb.year, alb.artwork, alb.track_count
		ORDER BY alb.year DESC LIMIT $2
		`)).WithArgs(driver.Value(artistID), driver.Value(1)).WillReturnRows(rows)
			},
			expected: func() []*proto.Album {
				albums := make([]*proto.Album, 0, 1)
				albums = append(albums, album)
				return albums
			}(),
			expectedError: true,
		},
		{
			name:   "rows returns error",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.artwork_color", "tracksDuration"}).RowError(1, errors.New("error"))
				for i := 0; i < 4; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.ArtworkColor, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "artwork_color"}, "alb")+", SUM(t.duration) AS tracksDuration"+
					`
		FROM albums alb
		JOIN tracks t ON alb.id = t.album
		WHERE alb.artist = $1
		GROUP BY alb.id, alb.title, alb.year, alb.artwork, alb.track_count
		ORDER BY alb.year DESC LIMIT $2
		`)).WithArgs(driver.Value(artistID), driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() []*proto.Album {
				var amount = 4
				albums := make([]*proto.Album, 0, amount)
				for i := 0; i < amount; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.ArtistAlbums(artistID, currentTest.amount)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestMusicStorage_IncrementListenCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	var trackID int64 = 1

	tests := []struct {
		name          string
		mock          func()
		expectedError bool
	}{
		{
			name: "increment listen count",
			mock: func() {
				row := mock.NewRows([]string{})
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE tracks SET listen_count = listen_count + 1 WHERE id=$1`)).WithArgs(driver.Value(trackID)).WillReturnRows(row)
			},
		},
		{
			name: "query returns error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE tracks SET listen_count = listen_count + 1 WHERE id=$1`)).WithArgs(driver.Value(trackID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.IncrementListenCount(trackID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMusicStorage_AlbumData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	album := &proto.AlbumPageResponse{
		AlbumID:      1,
		Title:        "testTitle",
		Year:         2,
		Artwork:      "testArtwork",
		ArtworkColor: "testArtWOrkColor",
		TracksCount:  3,
		Artist: &proto.Artist{
			ID:   4,
			Name: "testArtistName",
		},
		TracksDuration: 5,
	}

	var albumID int64 = 1

	tests := []struct {
		name          string
		mock          func()
		expected      *proto.AlbumPageResponse
		expectedError bool
	}{
		{
			name: "get album data",
			mock: func() {
				row := mock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "artwork_color", "a.track_count", "art.id", "art.name", "tracksDuration"})
				row.AddRow(album.AlbumID, album.Title, album.Year, album.Artwork, album.ArtworkColor, album.TracksCount, album.Artist.ID, album.Artist.Name, album.TracksDuration)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "artwork_color", "track_count"}, "alb") + ", " +
					wrapper.Wrapper([]string{"id", "name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
					`
		FROM albums alb
		JOIN artists art ON alb.artist = art.id
		JOIN tracks t ON alb.id = t.album
		WHERE alb.id = $1
		GROUP BY alb.id, art.id
		`)).WithArgs(driver.Value(albumID)).WillReturnRows(row)
			},
			expected: album,
		},
		{
			name: "query returns error",
			mock: func() {
				row := mock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "artwork_color", "a.track_count", "art.id", "art.name", "tracksDuration"})
				row.AddRow(album.AlbumID, album.Title, album.Year, album.Artwork, album.ArtworkColor, album.TracksCount, album.Artist.ID, album.Artist.Name, album.TracksDuration)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "artwork_color", "track_count"}, "alb") + ", " +
					wrapper.Wrapper([]string{"id", "name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
					`
		FROM albums alb
		JOIN artists art ON alb.artist = art.id
		JOIN tracks t ON alb.id = t.album
		WHERE alb.id = $1
		GROUP BY alb.id, art.id
		`)).WithArgs(driver.Value(albumID)).WillReturnError(errors.New("error"))
			},
			expected:      album,
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.AlbumData(albumID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, currentTest.expected, result)
				assert.NoError(t, err)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_AlbumTracks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const userID = 1

	track := &proto.AlbumTrack{
		ID:            1,
		Title:         "testTrackTitle",
		Explicit:      true,
		Genre:         "testGenre",
		Number:        2,
		File:          "testFile",
		ListenCount:   3,
		Duration:      4,
		Lossless:      true,
		IsInFavorites: true,
	}
	trackWithoutFile := new(proto.AlbumTrack)
	_ = copier.Copy(trackWithoutFile, track)
	trackWithoutFile.File = ""

	var albumID int64 = 1

	tests := []struct {
		name          string
		amount        int64
		isAuthorized  bool
		mock          func()
		expected      []*proto.AlbumTrack
		expectedError bool
	}{
		{
			name:         "get tracks by album",
			amount:       4,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "g.name", "favorite"})
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE alb.id = $2
		ORDER BY t.number
		`)).WithArgs(driver.Value(userID), driver.Value(albumID)).WillReturnRows(rows)
			},
			expected: func() []*proto.AlbumTrack {
				var amount = 4
				tracks := make([]*proto.AlbumTrack, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "query returns error",
			amount:       1,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "g.name", "favorite"})
				for i := 0; i < 1; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE alb.id = $2
		ORDER BY t.number
		`)).WithArgs(driver.Value(userID), driver.Value(albumID)).WillReturnError(errors.New("error"))
			},
			expected: func() []*proto.AlbumTrack {
				tracks := make([]*proto.AlbumTrack, 0, 1)
				tracks = append(tracks, track)
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:         "scan returns error",
			amount:       1,
			isAuthorized: true,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "g.name", "favorite", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Genre, track.IsInFavorites, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE alb.id = $2
		ORDER BY t.number
		`)).WithArgs(driver.Value(userID), driver.Value(albumID)).WillReturnRows(rows)
			},
			expected: func() []*proto.AlbumTrack {
				tracks := make([]*proto.AlbumTrack, 0, 1)
				tracks = append(tracks, track)
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:   "get tracks by album unauthorized",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "g.name", "favorite"})
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, "", track.ListenCount, track.Duration, track.Lossless, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE alb.id = $2
		ORDER BY t.number
		`)).WithArgs(driver.Value(userID), driver.Value(albumID)).WillReturnRows(rows)
			},
			expected: func() []*proto.AlbumTrack {
				var amount = 4
				tracks := make([]*proto.AlbumTrack, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, trackWithoutFile)
				}
				return tracks
			}(),
		},
		{
			name:         "rows.Err() returns error",
			amount:       4,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "g.name", "favorite"}).RowError(1, errors.New("error"))
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE alb.id = $2
		ORDER BY t.number
		`)).WithArgs(driver.Value(userID), driver.Value(albumID)).WillReturnRows(rows)
			},
			expected: func() []*proto.AlbumTrack {
				var amount = 4
				tracks := make([]*proto.AlbumTrack, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.AlbumTracks(albumID, userID, currentTest.isAuthorized)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_FindTracksByFullWord(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const userID = 1

	track := &proto.Track{
		ID:          1,
		Title:       "testTrackTitle",
		Explicit:    true,
		Genre:       "testGenre",
		Number:      2,
		File:        "testFile",
		ListenCount: 3,
		Duration:    4,
		Lossless:    true,
		Album: &proto.Album{
			ID:           5,
			Title:        "testAlbumTitle",
			Artwork:      "testAlbumArtwork",
			ArtworkColor: "testArtWorkColor",
		},
		Artist: &proto.Artist{
			ID:   6,
			Name: "testArtistName",
		},
		IsInFavorites: true,
	}
	trackWithoutFile := new(proto.Track)
	_ = copier.Copy(trackWithoutFile, track)
	trackWithoutFile.File = ""

	const text = "testText"

	tests := []struct {
		name          string
		isAuthorized  bool
		mock          func()
		expected      []*proto.Track
		expectedError bool
	}{
		{
			name:         "find constants.SearchPageTracksAmount tracks with text",
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
			SELECT track
		    FROM search
		    WHERE concatenation @@ plainto_tsquery($2)
		)
		ORDER BY t.listen_count DESC LIMIT $3`)).WithArgs(driver.Value(userID), driver.Value(text), driver.Value(constants.SearchPageTracksAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "query returns error",
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
			SELECT track
		    FROM search
		    WHERE concatenation @@ plainto_tsquery($2)
		)
		ORDER BY t.listen_count DESC LIMIT $3`)).WithArgs(driver.Value(userID), driver.Value(text), driver.Value(constants.SearchPageTracksAmount)).WillReturnError(errors.New("error"))
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:         "scan returns error",
			isAuthorized: true,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite", "newArg"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
			SELECT track
		    FROM search
		    WHERE concatenation @@ plainto_tsquery($2)
		)
		ORDER BY t.listen_count DESC LIMIT $3`)).WithArgs(driver.Value(userID), driver.Value(text), driver.Value(constants.SearchPageTracksAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
		{
			name: "find tracks unauthorized",
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, "", track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
			SELECT track
		    FROM search
		    WHERE concatenation @@ plainto_tsquery($2)
		)
		ORDER BY t.listen_count DESC LIMIT $3`)).WithArgs(driver.Value(userID), driver.Value(text), driver.Value(constants.SearchPageTracksAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, trackWithoutFile)
				}
				return tracks
			}(),
		},
		{
			name:         "rows.Err() returns error",
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"}).RowError(1, errors.New("error"))
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
			SELECT track
		    FROM search
		    WHERE concatenation @@ plainto_tsquery($2)
		)
		ORDER BY t.listen_count DESC LIMIT $3`)).WithArgs(driver.Value(userID), driver.Value(text), driver.Value(constants.SearchPageTracksAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.FindTracksByFullWord(text, userID, currentTest.isAuthorized)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_FindTracksByPartial(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const userID = 1

	track := &proto.Track{
		ID:          1,
		Title:       "testTrackTitle",
		Explicit:    true,
		Genre:       "testGenre",
		Number:      2,
		File:        "testFile",
		ListenCount: 3,
		Duration:    4,
		Lossless:    true,
		Album: &proto.Album{
			ID:      5,
			Title:   "testAlbumTitle",
			Artwork: "testAlbumArtwork",
		},
		Artist: &proto.Artist{
			ID:   6,
			Name: "testArtistName",
		},
		IsInFavorites: true,
	}
	trackWithoutFile := new(proto.Track)
	_ = copier.Copy(trackWithoutFile, track)
	trackWithoutFile.File = ""

	const text = "testText"
	textArgument := "%" + text + "%"

	tests := []struct {
		name          string
		isAuthorized  bool
		mock          func()
		expected      []*proto.Track
		expectedError bool
	}{
		{
			name:         "find constants.SearchPageTracksAmount tracks with text",
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
		    SELECT track
		    FROM test
		    WHERE concat ILIKE $2
		)
		ORDER BY t.listen_count DESC LIMIT $3`)).WithArgs(driver.Value(userID), driver.Value(textArgument), driver.Value(constants.SearchPageTracksAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "query returns error",
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
		    SELECT track
		    FROM test
		    WHERE concat ILIKE $2
		)
		ORDER BY t.listen_count DESC LIMIT $3`)).WithArgs(driver.Value(userID), driver.Value(textArgument), driver.Value(constants.SearchPageTracksAmount)).WillReturnError(errors.New("error"))
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:         "scan returns error",
			isAuthorized: true,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite", "newArg"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
		    SELECT track
		    FROM test
		    WHERE concat ILIKE $2
		)
		ORDER BY t.listen_count DESC LIMIT $3`)).WithArgs(driver.Value(userID), driver.Value(textArgument), driver.Value(constants.SearchPageTracksAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
		{
			name: "find tracks unauthorized",
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, "", track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
		    SELECT track
		    FROM test
		    WHERE concat ILIKE $2
		)
		ORDER BY t.listen_count DESC LIMIT $3`)).WithArgs(driver.Value(userID), driver.Value(textArgument), driver.Value(constants.SearchPageTracksAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, trackWithoutFile)
				}
				return tracks
			}(),
		},
		{
			name:         "rows.Err() returns error",
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"}).RowError(1, errors.New("error"))
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
		    SELECT track
		    FROM test
		    WHERE concat ILIKE $2
		)
		ORDER BY t.listen_count DESC LIMIT $3`)).WithArgs(driver.Value(userID), driver.Value(textArgument), driver.Value(constants.SearchPageTracksAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.FindTracksByPartial(text, userID, currentTest.isAuthorized)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_FindArtists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	tracks := make([]*proto.Track, 0)
	albums := make([]*proto.Album, 0)

	artist := &proto.Artist{
		ID:     1,
		Name:   "testName",
		Avatar: "testAvatar",
		Video:  "",
		Tracks: tracks,
		Albums: albums,
	}

	text := "testText"
	textArgument := "%" + text + "%"

	tests := []struct {
		name          string
		mock          func()
		expected      []*proto.Artist
		expectedError bool
	}{
		{
			name: "find artists",
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				for i := 0; i < constants.SearchPageArtistsAmount; i++ {
					rows.AddRow(artist.ID, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, avatar
		FROM artists
		WHERE name ILIKE $1
		LIMIT $2
	`)).WithArgs(driver.Value(textArgument), driver.Value(constants.SearchPageArtistsAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Artist {
				artists := make([]*proto.Artist, 0, constants.SearchPageArtistsAmount)
				for i := 0; i < constants.SearchPageArtistsAmount; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
		},
		{
			name: "query returns error",
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				for i := 0; i < constants.SearchPageArtistsAmount; i++ {
					rows.AddRow(artist.ID, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, avatar
		FROM artists
		WHERE name ILIKE $1
		LIMIT $2
	`)).WithArgs(driver.Value(textArgument), driver.Value(constants.SearchPageArtistsAmount)).WillReturnError(errors.New("error"))
			},
			expected: func() []*proto.Artist {
				artists := make([]*proto.Artist, 0, constants.SearchPageArtistsAmount)
				for i := 0; i < constants.SearchPageArtistsAmount; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: true,
		},
		{
			name: "scan returns error",
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar", "newArg"})
				for i := 0; i < constants.SearchPageArtistsAmount; i++ {
					rows.AddRow(artist.ID, artist.Name, artist.Avatar, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, avatar
		FROM artists
		WHERE name ILIKE $1
		LIMIT $2
	`)).WithArgs(driver.Value(textArgument), driver.Value(constants.SearchPageArtistsAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Artist {
				artists := make([]*proto.Artist, 0, constants.SearchPageArtistsAmount)
				for i := 0; i < constants.SearchPageArtistsAmount; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: true,
		},
		{
			name: "rows.Err() returns error",
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"}).RowError(1, errors.New("error"))
				for i := 0; i < constants.SearchPageArtistsAmount; i++ {
					rows.AddRow(artist.ID, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, avatar
		FROM artists
		WHERE name ILIKE $1
		LIMIT $2
	`)).WithArgs(driver.Value(textArgument), driver.Value(constants.SearchPageArtistsAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Artist {
				artists := make([]*proto.Artist, 0, constants.SearchPageArtistsAmount)
				for i := 0; i < constants.SearchPageArtistsAmount; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.FindArtists(text)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_FindAlbums(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	album := &proto.Album{
		ID:             1,
		Title:          "testTitle",
		Year:           2,
		Artist:         "testArtist",
		Artwork:        "testArtwork",
		TracksAmount:   3,
		TracksDuration: 4,
		ArtworkColor:   "testArtWorkColor",
	}

	text := "testText"
	textArgument := "%" + text + "%"

	tests := []struct {
		name          string
		mock          func()
		expected      []*proto.Album
		expectedError bool
	}{
		{
			name: "find albums",
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.track_count", "artwork_color", "art.name", "tracksDuration"})
				for i := 0; i < constants.SearchPageAlbumsAmount; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.TracksAmount, album.ArtworkColor, album.Artist, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "art")+", SUM(t.duration) AS tracksDuration"+
					`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		WHERE alb.title ILIKE $1
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY year DESC
		LIMIT $2
		`)).WithArgs(driver.Value(textArgument), driver.Value(constants.SearchPageAlbumsAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Album {
				albums := make([]*proto.Album, 0, constants.SearchPageAlbumsAmount)
				for i := 0; i < constants.SearchPageAlbumsAmount; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
		},
		{
			name: "query returns error",
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.track_count", "artwork_color", "art.name", "tracksDuration"})
				for i := 0; i < constants.SearchPageAlbumsAmount; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.TracksAmount, album.ArtworkColor, album.Artist, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "art")+", SUM(t.duration) AS tracksDuration"+
					`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		WHERE alb.title ILIKE $1
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY year DESC
		LIMIT $2
		`)).WithArgs(driver.Value(textArgument), driver.Value(constants.SearchPageAlbumsAmount)).WillReturnError(errors.New("error"))
			},
			expected: func() []*proto.Album {
				albums := make([]*proto.Album, 0, constants.SearchPageAlbumsAmount)
				for i := 0; i < constants.SearchPageAlbumsAmount; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: true,
		},
		{
			name: "scan returns error",
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.track_count", "artwork_color", "art.name", "tracksDuration", "newArg"})
				for i := 0; i < constants.SearchPageAlbumsAmount; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.TracksAmount, album.ArtworkColor, album.Artist, album.TracksDuration, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "art")+", SUM(t.duration) AS tracksDuration"+
					`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		WHERE alb.title ILIKE $1
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY year DESC
		LIMIT $2
		`)).WithArgs(driver.Value(textArgument), driver.Value(constants.SearchPageAlbumsAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Album {
				albums := make([]*proto.Album, 0, constants.SearchPageAlbumsAmount)
				for i := 0; i < constants.SearchPageAlbumsAmount; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: true,
		},
		{
			name: "rows returns error",
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "a.artwork", "a.track_count", "artwork_color", "art.name", "tracksDuration"}).RowError(1, errors.New("error"))
				for i := 0; i < 4; i++ {
					rows.AddRow(album.ID, album.Title, album.Year, album.Artwork, album.TracksAmount, album.ArtworkColor, album.Artist, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"name"}, "art")+", SUM(t.duration) AS tracksDuration"+
					`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		WHERE alb.title ILIKE $1
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY year DESC
		LIMIT $2
		`)).WithArgs(driver.Value(textArgument), driver.Value(constants.SearchPageAlbumsAmount)).WillReturnRows(rows)
			},
			expected: func() []*proto.Album {
				albums := make([]*proto.Album, 0, constants.SearchPageAlbumsAmount)
				for i := 0; i < constants.SearchPageAlbumsAmount; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.FindAlbums(text)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestMusicStorage_IsPlaylistOwner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	var (
		playlistID int64
		userID     int64
	)

	tests := []struct {
		name          string
		isOwner       bool
		mock          func()
		expected      bool
		expectedError bool
	}{
		{
			name:    "user is owner",
			isOwner: true,
			mock: func() {
				rows := mock.NewRows([]string{"isOwner"})
				rows.AddRow("true")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1 AND user_id=$2`)).WithArgs(driver.Value(userID), driver.Value(playlistID)).WillReturnRows(rows)
			},
			expected: true,
		},
		{
			name:    "user is not owner",
			isOwner: true,
			mock: func() {
				rows := mock.NewRows([]string{"isOwner"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1 AND user_id=$2`)).WithArgs(driver.Value(userID), driver.Value(playlistID)).WillReturnRows(rows)
			},
		},
		{
			name:    "query returns error",
			isOwner: true,
			mock: func() {
				rows := mock.NewRows([]string{"isOwner"})
				rows.AddRow("true")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1 AND user_id=$2`)).WithArgs(driver.Value(userID), driver.Value(playlistID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.IsPlaylistOwner(playlistID, userID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_PlaylistTracks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const userID = 1

	track := &proto.Track{
		ID:          1,
		Title:       "testTrackTitle",
		Explicit:    true,
		Genre:       "testGenre",
		Number:      2,
		File:        "testFile",
		ListenCount: 3,
		Duration:    4,
		Lossless:    true,
		Album: &proto.Album{
			ID:           5,
			Title:        "testAlbumTitle",
			Artwork:      "testArtWork",
			ArtworkColor: "testArtWorkColor",
		},
		Artist: &proto.Artist{
			ID:   6,
			Name: "testArtistName",
		},
		IsInFavorites: true,
	}
	trackWithoutFile := new(proto.Track)
	_ = copier.Copy(trackWithoutFile, track)
	trackWithoutFile.File = ""

	var playlistID int64 = 1

	tests := []struct {
		name          string
		isAuthorized  bool
		mock          func()
		expected      []*proto.Track
		expectedError bool
	}{
		{
			name:         "get playlist tracks",
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
			SELECT track
			FROM playlist_tracks
			WHERE playlist=$2
		)`)).WithArgs(driver.Value(userID), driver.Value(playlistID)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "query returns error",
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
			SELECT track
			FROM playlist_tracks
			WHERE playlist=$2
		)`)).WithArgs(driver.Value(userID), driver.Value(playlistID)).WillReturnError(errors.New("error"))
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:         "scan returns error",
			isAuthorized: true,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite", "newArg"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
			SELECT track
			FROM playlist_tracks
			WHERE playlist=$2
		)`)).WithArgs(driver.Value(userID), driver.Value(playlistID)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
		{
			name: "find tracks unauthorized",
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, "", track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
			SELECT track
			FROM playlist_tracks
			WHERE playlist=$2
		)`)).WithArgs(driver.Value(userID), driver.Value(playlistID)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, trackWithoutFile)
				}
				return tracks
			}(),
		},
		{
			name:         "rows.Err() returns error",
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"}).RowError(1, errors.New("error"))
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT `+
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t")+", "+
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb")+", "+
					wrapper.Wrapper([]string{"id", "name"}, "art")+", "+
					wrapper.Wrapper([]string{"name"}, "g")+", "+
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		LEFT JOIN likes l on t.id = l.track_id and l.user_id = $1
		WHERE t.id IN (
			SELECT track
			FROM playlist_tracks
			WHERE playlist=$2
		)`)).WithArgs(driver.Value(userID), driver.Value(playlistID)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				tracks := make([]*proto.Track, 0, constants.SearchPageTracksAmount)
				for i := 0; i < constants.SearchPageTracksAmount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.PlaylistTracks(playlistID, userID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestMusicStorage_PlaylistInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	playlist := &proto.PlaylistData{
		PlaylistID:   1,
		Title:        "testTitle",
		Artwork:      "testArtWork",
		ArtworkColor: "testArtWorkColor",
		IsPublic:     true,
	}

	expectedPlaylist := &proto.PlaylistData{
		PlaylistID:   1,
		Title:        "testTitle",
		Artwork:      os.Getenv("PLAYLIST_ROOT_PREFIX") + "testArtWork" + constants.PlaylistArtworkExtension384px,
		ArtworkColor: "testArtWorkColor",
		IsPublic:     true,
	}

	var playlistID int64 = 1

	tests := []struct {
		name          string
		mock          func()
		expected      *proto.PlaylistData
		expectedError bool
	}{
		{
			name: "get playlist info",
			mock: func() {
				row := mock.NewRows([]string{"playlistID", "title", "artwork", "artworkColor", "isPublic"})
				row.AddRow(playlist.PlaylistID, playlist.Title, playlist.Artwork, playlist.ArtworkColor, playlist.IsPublic)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, artwork, artwork_color, is_public FROM playlists WHERE id=$1`)).WithArgs(driver.Value(playlistID)).WillReturnRows(row)
			},
			expected: expectedPlaylist,
		},
		{
			name: "query returns error",
			mock: func() {
				row := mock.NewRows([]string{"playlistID", "title", "artwork", "artworkColor", "isPublic"})
				row.AddRow(playlist.PlaylistID, playlist.Title, playlist.Artwork, playlist.ArtworkColor, playlist.IsPublic)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, artwork, artwork_color, is_public FROM playlists WHERE id=$1`)).WithArgs(driver.Value(playlistID)).WillReturnError(errors.New("error"))
			},
			expected:      expectedPlaylist,
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.PlaylistInfo(playlistID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_UserPlaylists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	playlist := &proto.PlaylistData{
		PlaylistID: 1,
		Title:      "testTitle",
		Artwork:    "testArtWork",
		IsPublic:   true,
		IsOwn:      true,
	}

	expectedPlaylist := &proto.PlaylistData{
		PlaylistID: 1,
		Title:      "testTitle",
		Artwork:    os.Getenv("PLAYLIST_ROOT_PREFIX") + "testArtWork" + constants.PlaylistArtworkExtension100px,
		IsPublic:   true,
		IsOwn:      true,
	}

	var (
		userID          int64 = 1
		playlistsAmount       = 4
	)

	tests := []struct {
		name          string
		mock          func()
		expected      []*proto.PlaylistData
		expectedError bool
	}{
		{
			name: "get playlist info",
			mock: func() {
				row := mock.NewRows([]string{"playlistID", "title", "artwork", "isPublic", "is_own"})
				for i := 0; i < playlistsAmount; i++ {
					row.AddRow(playlist.PlaylistID, playlist.Title, playlist.Artwork, playlist.IsPublic, playlist.IsOwn)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, artwork, is_public, user_id=$1 AS is_own FROM playlists WHERE user_id=$1 OR is_public=true`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expected: func() []*proto.PlaylistData {
				playlists := make([]*proto.PlaylistData, 0, playlistsAmount)
				for i := 0; i < playlistsAmount; i++ {
					playlists = append(playlists, expectedPlaylist)
				}
				return playlists
			}(),
		},
		{
			name: "query returns error",
			mock: func() {
				row := mock.NewRows([]string{"playlistID", "title", "artwork", "isPublic", "is_own"})
				for i := 0; i < playlistsAmount; i++ {
					row.AddRow(playlist.PlaylistID, playlist.Title, playlist.Artwork, playlist.IsPublic, playlist.IsOwn)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, artwork, is_public, user_id=$1 AS is_own FROM playlists WHERE user_id=$1 OR is_public=true`)).WithArgs(driver.Value(userID)).WillReturnError(errors.New("error"))
			},
			expected: func() []*proto.PlaylistData {
				playlists := make([]*proto.PlaylistData, 0, playlistsAmount)
				for i := 0; i < playlistsAmount; i++ {
					playlists = append(playlists, expectedPlaylist)
				}
				return playlists
			}(),
			expectedError: true,
		},
		{
			name: "scan returns error",
			mock: func() {
				newArg := 1
				row := mock.NewRows([]string{"playlistID", "title", "artwork", "isPublic", "is_own", "newArg"})
				for i := 0; i < playlistsAmount; i++ {
					row.AddRow(playlist.PlaylistID, playlist.Title, playlist.Artwork, playlist.IsPublic, playlist.IsOwn, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, artwork, is_public, user_id=$1 AS is_own FROM playlists WHERE user_id=$1 OR is_public=true`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expected: func() []*proto.PlaylistData {
				playlists := make([]*proto.PlaylistData, 0, playlistsAmount)
				for i := 0; i < playlistsAmount; i++ {
					playlists = append(playlists, expectedPlaylist)
				}
				return playlists
			}(),
			expectedError: true,
		},
		{
			name: "rows.Err() returns error",
			mock: func() {
				row := mock.NewRows([]string{"playlistID", "title", "artwork", "isPublic", "is_own"}).RowError(1, errors.New("error"))
				for i := 0; i < playlistsAmount; i++ {
					row.AddRow(playlist.PlaylistID, playlist.Title, playlist.Artwork, playlist.IsPublic, playlist.IsOwn)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, title, artwork, is_public, user_id=$1 AS is_own FROM playlists WHERE user_id=$1 OR is_public=true`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expected: func() []*proto.PlaylistData {
				playlists := make([]*proto.PlaylistData, 0, playlistsAmount)
				for i := 0; i < playlistsAmount; i++ {
					playlists = append(playlists, expectedPlaylist)
				}
				return playlists
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.UserPlaylists(userID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestMusicStorage_DoesPlaylistExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const playlistID int64 = 1

	tests := []struct {
		name          string
		mock          func()
		expected      bool
		expectedError bool
	}{
		{
			name: "playlist exist",
			mock: func() {
				row := mock.NewRows([]string{"playlist_exist"})
				row.AddRow("true")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1`)).WithArgs(driver.Value(playlistID)).WillReturnRows(row)
			},
			expected: true,
		},
		{
			name: "playlist does not exist",
			mock: func() {
				row := mock.NewRows([]string{"playlist_exist"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1`)).WithArgs(driver.Value(playlistID)).WillReturnRows(row)
			},
		},
		{
			name: "query returns error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1`)).WithArgs(driver.Value(playlistID)).WillReturnError(errors.New("error"))
			},
			expected:      false,
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.DoesPlaylistExist(playlistID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestMusicStorage_AddTrackToFavorite(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const (
		userID = iota
		trackID
	)

	tests := []struct {
		name          string
		mock          func()
		expectedError bool
	}{
		{
			name: "successful adding track to favorites",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO likes(user_id, track_id) VALUES ($1, $2)`)).WithArgs(driver.Value(userID), driver.Value(trackID)).WillReturnRows()
			},
		},
		{
			name: "failed adding track to favorites",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO likes(user_id, track_id) VALUES ($1, $2)`)).WithArgs(driver.Value(userID), driver.Value(trackID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.AddTrackToFavorites(userID, trackID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMusicStorage_DeleteTrackFromFavorites(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const (
		userID = iota
		trackID
	)

	tests := []struct {
		name          string
		mock          func()
		expectedError bool
	}{
		{
			name: "successful deleting track from favorites",
			mock: func() {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM likes WHERE user_id = $1 and track_id = $2`)).WithArgs(driver.Value(userID), driver.Value(trackID)).WillReturnResult(driver.ResultNoRows)
			},
		},
		{
			name: "failed deleting track from favorites",
			mock: func() {
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM likes WHERE user_id = $1 and track_id = $2`)).WithArgs(driver.Value(userID), driver.Value(trackID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.DeleteTrackFromFavorites(userID, trackID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMusicStorage_IsTrackInFavorites(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const (
		userID = iota
		trackID
	)

	tests := []struct {
		name              string
		mock              func()
		expectedExistence bool
		expectedError     bool
	}{
		{
			name: "track in favorites",
			mock: func() {
				row := mock.NewRows([]string{"exist"})
				row.AddRow(true)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND track_id = $2)`)).WithArgs(driver.Value(userID), driver.Value(trackID)).WillReturnRows(row)
			},
			expectedExistence: true,
		},
		{
			name: "track not in favorites",
			mock: func() {
				row := mock.NewRows([]string{"exist"})
				row.AddRow(false)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND track_id = $2)`)).WithArgs(driver.Value(userID), driver.Value(trackID)).WillReturnRows(row)
			},
		},
		{
			name: "query error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND track_id = $2)`)).WithArgs(driver.Value(userID), driver.Value(trackID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			isExist, err := repository.IsTrackInFavorites(userID, trackID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expectedExistence, isExist)
			}
		})
	}
}

//nolint:cyclop
func TestMusicStorage_GetFavorites(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const userID = 1

	track := &proto.Track{
		ID:          1,
		Title:       "testTrackTitle",
		Explicit:    true,
		Genre:       "testGenre",
		Number:      2,
		File:        "testFile",
		ListenCount: 3,
		Duration:    4,
		Lossless:    true,
		Album: &proto.Album{
			ID:           5,
			Title:        "testAlbumTitle",
			Artwork:      "testAlbumArtwork",
			ArtworkColor: "testArtWOrkColor",
		},
		Artist: &proto.Artist{
			ID:   6,
			Name: "testArtistName",
		},
		IsInFavorites: true,
	}

	tests := []struct {
		name          string
		amount        int64
		isAuthorized  bool
		mock          func()
		expected      []*proto.Track
		expectedError bool
	}{
		{
			name:         "successful getting favorites",
			amount:       4,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb") + ", " +
					wrapper.Wrapper([]string{"id", "name"}, "art") + ", " +
					wrapper.Wrapper([]string{"name"}, "g") + ", " +
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		JOIN likes l on t.id = l.track_id and l.user_id = $1
        ORDER BY l.id`)).WithArgs(driver.Value(userID)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				var amount = 4
				tracks := make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
		},
		{
			name:         "query returns error",
			amount:       1,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"})
				for i := 0; i < 1; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb") + ", " +
					wrapper.Wrapper([]string{"id", "name"}, "art") + ", " +
					wrapper.Wrapper([]string{"name"}, "g") + ", " +
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		JOIN likes l on t.id = l.track_id and l.user_id = $1
        ORDER BY l.id`)).WithArgs(driver.Value(userID)).WillReturnError(errors.New("error"))
			},
			expected: func() []*proto.Track {
				var amount = 1
				tracks := make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:         "scan returns error",
			amount:       1,
			isAuthorized: true,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb") + ", " +
					wrapper.Wrapper([]string{"id", "name"}, "art") + ", " +
					wrapper.Wrapper([]string{"name"}, "g") + ", " +
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		JOIN likes l on t.id = l.track_id and l.user_id = $1
        ORDER BY l.id`)).WithArgs(driver.Value(userID)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				var amount = 1
				tracks := make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
		{
			name:         "rows.Err() returns error",
			amount:       4,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "alb.artwork_color", "art.id", "art.name", "g.name", "favorite"}).RowError(0, errors.New("error"))
				for i := 0; i < 4; i++ {
					rows.AddRow(track.ID, track.Title, track.Explicit, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.ID, track.Album.Title, track.Album.Artwork, track.Album.ArtworkColor, track.Artist.ID, track.Artist.Name, track.Genre, track.IsInFavorites)
				}
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT ` +
					wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
					wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb") + ", " +
					wrapper.Wrapper([]string{"id", "name"}, "art") + ", " +
					wrapper.Wrapper([]string{"name"}, "g") + ", " +
					`
		l.id IS NOT NULL as favorite
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		JOIN likes l on t.id = l.track_id and l.user_id = $1
        ORDER BY l.id`)).WithArgs(driver.Value(userID)).WillReturnRows(rows)
			},
			expected: func() []*proto.Track {
				var amount = 4
				tracks := make([]*proto.Track, 0, amount)
				for i := 0; i < amount; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.GetFavorites(userID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestMusicStorage_IsPlaylistPublic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewMusicStorage(db)

	const playlistID = 1

	tests := []struct {
		name              string
		mock              func()
		expectedPublicity bool
		expectedError     bool
	}{
		{
			name: "playlist is public",
			mock: func() {
				row := mock.NewRows([]string{"public"})
				row.AddRow(true)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1 AND is_public=true`)).WithArgs(driver.Value(playlistID)).WillReturnRows(row)
			},
			expectedPublicity: true,
		},
		{
			name: "playlist is not public",
			mock: func() {
				row := mock.NewRows([]string{"public"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1 AND is_public=true`)).WithArgs(driver.Value(playlistID)).WillReturnRows(row)
			},
		},
		{
			name: "query error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1 AND is_public=true`)).WithArgs(driver.Value(playlistID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
		{
			name: "row.Err() returns error",
			mock: func() {
				row := mock.NewRows([]string{"public"}).RowError(0, errors.New("error"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1 AND is_public=true`)).WithArgs(driver.Value(playlistID)).WillReturnRows(row)
			},
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			isPublic, err := repository.IsPlaylistPublic(playlistID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expectedPublicity, isPublic)
			}
		})
	}
}
*/