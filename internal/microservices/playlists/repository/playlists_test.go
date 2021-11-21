package repository

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/playlists/proto"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"log"
	"regexp"
	"testing"
)

func TestPlaylistsStorage_CreatePlaylist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	playlistResponse := &proto.CreatePlaylistResponse{PlaylistID: 1}

	const (
		title        = "testTitle"
		userID       = 1
		artWork      = "testArtWork"
		artWorkColor = "testArtWorkColor"
		isPublic     = true
	)

	tests := []struct {
		name          string
		mock          func()
		expected      *proto.CreatePlaylistResponse
		expectedError bool
	}{
		{
			name: "create playlist success",
			mock: func() {
				row := mock.NewRows([]string{"id"})
				row.AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO playlists(title, user_id, artwork, artwork_color, is_public) VALUES ($1, $2, $3, $4, $5) RETURNING id`)).
					WithArgs(driver.Value(title), driver.Value(userID), driver.Value(artWork), driver.Value(artWorkColor), driver.Value(isPublic)).
					WillReturnRows(row)
			},
			expected: playlistResponse,
		},
		{
			name: "create playlist fail",
			mock: func() {
				row := mock.NewRows([]string{"id"})
				row.AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO playlists(title, user_id, artwork, artwork_color, is_public) VALUES ($1, $2, $3, $4, $5) RETURNING id`)).
					WithArgs(driver.Value(title), driver.Value(userID), driver.Value(artWork), driver.Value(artWorkColor), driver.Value(isPublic)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.CreatePlaylist(userID, title, artWork, artWorkColor, isPublic)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestPlaylistsStorage_GetOldPlaylistSettings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID = 1
		artwork    = "testArtWork"
	)

	tests := []struct {
		name          string
		mock          func()
		expected      string
		expectedError bool
	}{
		{
			name: "get old settings success",
			mock: func() {
				row := mock.NewRows([]string{"artwork"})
				row.AddRow(artwork)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT artwork FROM playlists WHERE id=$1`)).
					WithArgs(driver.Value(playlistID)).
					WillReturnRows(row)
			},
			expected: artwork,
		},
		{
			name: "get old settings fail",
			mock: func() {
				row := mock.NewRows([]string{"artwork"})
				row.AddRow(artwork)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT artwork FROM playlists WHERE id=$1`)).
					WithArgs(driver.Value(playlistID)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.GetOldPlaylistSettings(playlistID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestPlaylistsStorage_DeletePlaylist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID = 1
	)

	tests := []struct {
		name          string
		mock          func()
		expectedError bool
	}{
		{
			name: "delete playlist success",
			mock: func() {
				row := mock.NewRows([]string{"success"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM playlists WHERE id=$1`)).
					WithArgs(driver.Value(playlistID)).
					WillReturnRows(row)
			},
		},
		{
			name: "delete playlist fail",
			mock: func() {
				row := mock.NewRows([]string{"fail"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM playlists WHERE id=$1`)).
					WithArgs(driver.Value(playlistID)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.DeletePlaylist(playlistID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistsStorage_AddTrack(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID = 1
		trackID    = 2
	)

	tests := []struct {
		name          string
		mock          func()
		expected      string
		expectedError bool
	}{
		{
			name: "add track to playlist success",
			mock: func() {
				row := mock.NewRows([]string{"success"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO playlist_tracks(playlist, track) VALUES ($1, $2)`)).
					WithArgs(driver.Value(playlistID), driver.Value(trackID)).
					WillReturnRows(row)
			},
		},
		{
			name: "add track to playlist fail",
			mock: func() {
				row := mock.NewRows([]string{"fail"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO playlist_tracks(playlist, track) VALUES ($1, $2)`)).
					WithArgs(driver.Value(playlistID), driver.Value(trackID)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.AddTrack(playlistID, trackID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistsStorage_IsAdded(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID = 1
		trackID    = 2
	)

	tests := []struct {
		name          string
		mock          func()
		expected      bool
		expectedError bool
	}{
		{
			name: "added",
			mock: func() {
				row := mock.NewRows([]string{"added"})
				row.AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlist_tracks WHERE playlist=$1 AND track=$2`)).
					WithArgs(driver.Value(playlistID), driver.Value(trackID)).
					WillReturnRows(row)
			},
			expected: true,
		},
		{
			name: "not added",
			mock: func() {
				row := mock.NewRows([]string{"not added"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlist_tracks WHERE playlist=$1 AND track=$2`)).
					WithArgs(driver.Value(playlistID), driver.Value(trackID)).
					WillReturnRows(row)
			},
		},
		{
			name: "query returns error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlist_tracks WHERE playlist=$1 AND track=$2`)).
					WithArgs(driver.Value(playlistID), driver.Value(trackID)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.IsAdded(playlistID, trackID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestPlaylistsStorage_DeleteTrack(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID = 1
		trackID    = 2
	)

	tests := []struct {
		name          string
		mock          func()
		expectedError bool
	}{
		{
			name: "delete track success",
			mock: func() {
				row := mock.NewRows([]string{"success"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM playlist_tracks WHERE playlist=$1 AND track=$2`)).
					WithArgs(driver.Value(playlistID), driver.Value(trackID)).
					WillReturnRows(row)
			},
		},
		{
			name: "delete track fail",
			mock: func() {
				row := mock.NewRows([]string{"fail"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`DELETE FROM playlist_tracks WHERE playlist=$1 AND track=$2`)).
					WithArgs(driver.Value(playlistID), driver.Value(trackID)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.DeleteTrack(playlistID, trackID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistsStorage_IsOwner(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID = 1
		trackID    = 2
	)

	tests := []struct {
		name          string
		mock          func()
		expected      bool
		expectedError bool
	}{
		{
			name: "is owner",
			mock: func() {
				row := mock.NewRows([]string{"added"})
				row.AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1 AND (user_id=$2 OR is_public=true)`)).
					WithArgs(driver.Value(playlistID), driver.Value(trackID)).
					WillReturnRows(row)
			},
			expected: true,
		},
		{
			name: "is not owner",
			mock: func() {
				row := mock.NewRows([]string{"not added"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1 AND (user_id=$2 OR is_public=true)`)).
					WithArgs(driver.Value(playlistID), driver.Value(trackID)).
					WillReturnRows(row)
			},
		},
		{
			name: "query returns error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1 AND (user_id=$2 OR is_public=true)`)).
					WithArgs(driver.Value(playlistID), driver.Value(trackID)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.IsOwner(playlistID, trackID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestPlaylistsStorage_DoesPlaylistExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID = 1
	)

	tests := []struct {
		name          string
		mock          func()
		expected      bool
		expectedError bool
	}{
		{
			name: "playlist exists",
			mock: func() {
				row := mock.NewRows([]string{"added"})
				row.AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1`)).
					WithArgs(driver.Value(playlistID)).
					WillReturnRows(row)
			},
			expected: true,
		},
		{
			name: "playlist does not exist",
			mock: func() {
				row := mock.NewRows([]string{"not added"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1`)).
					WithArgs(driver.Value(playlistID)).
					WillReturnRows(row)
			},
		},
		{
			name: "query returns error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM playlists WHERE id=$1`)).
					WithArgs(driver.Value(playlistID)).
					WillReturnError(errors.New("error"))
			},
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

func TestPlaylistsStorage_UpdatePlaylistTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID = 1
		title      = "testTitle"
	)

	tests := []struct {
		name          string
		mock          func()
		expectedError bool
	}{
		{
			name: "update playlist title success",
			mock: func() {
				row := mock.NewRows([]string{"success"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE playlists SET title=$1 WHERE id=$2`)).
					WithArgs(driver.Value(title), driver.Value(playlistID)).
					WillReturnRows(row)
			},
		},
		{
			name: "update playlist title fail",
			mock: func() {
				row := mock.NewRows([]string{"fail"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE playlists SET title=$1 WHERE id=$2`)).
					WithArgs(driver.Value(title), driver.Value(playlistID)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.UpdatePlaylistTitle(playlistID, title)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistsStorage_UpdatePlaylistArtwork(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID = 1
		artWork    = "testArtWork"
		artWorkColor = "testArtWorkColor"
	)

	tests := []struct {
		name          string
		mock          func()
		expectedError bool
	}{
		{
			name: "update playlist artwork success",
			mock: func() {
				row := mock.NewRows([]string{"success"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE playlists SET artwork=$1, artwork_color=$2 WHERE id=$3`)).
					WithArgs(driver.Value(artWork), driver.Value(artWorkColor), driver.Value(playlistID)).
					WillReturnRows(row)
			},
		},
		{
			name: "update playlist artwork fail",
			mock: func() {
				row := mock.NewRows([]string{"fail"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE playlists SET artwork=$1, artwork_color=$2 WHERE id=$3`)).
					WithArgs(driver.Value(artWork), driver.Value(artWorkColor), driver.Value(playlistID)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.UpdatePlaylistArtwork(playlistID, artWork, artWorkColor)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistsStorage_DeletePlaylistArtwork(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID   = 1
		artWork      = constants.PlaylistArtworkDefaultFilename
		artWorkColor = constants.PlaylistArtworkDefaultColor
	)

	tests := []struct {
		name          string
		mock          func()
		expectedError bool
	}{
		{
			name: "delete playlist artwork success",
			mock: func() {
				row := mock.NewRows([]string{"success"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE playlists SET artwork=$1, artwork_color=$2 WHERE id=$3`)).
					WithArgs(driver.Value(artWork), driver.Value(artWorkColor), driver.Value(playlistID)).
					WillReturnRows(row)
			},
		},
		{
			name: "delete playlist artwork fail",
			mock: func() {
				row := mock.NewRows([]string{"fail"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE playlists SET artwork=$1, artwork_color=$2 WHERE id=$3`)).
					WithArgs(driver.Value(artWork), driver.Value(artWorkColor), driver.Value(playlistID)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.DeletePlaylistArtwork(playlistID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPlaylistsStorage_UpdatePlaylistAccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewPlaylistsStorage(db)

	const (
		playlistID = 1
		isPublic   = true
	)

	tests := []struct {
		name          string
		mock          func()
		expectedError bool
	}{
		{
			name: "delete playlist artwork success",
			mock: func() {
				row := mock.NewRows([]string{"success"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE playlists SET is_public=$1 WHERE id=$2`)).
					WithArgs(driver.Value(isPublic), driver.Value(playlistID)).
					WillReturnRows(row)
			},
		},
		{
			name: "delete playlist artwork fail",
			mock: func() {
				row := mock.NewRows([]string{"fail"})
				row.AddRow("")
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE playlists SET is_public=$1 WHERE id=$2`)).
					WithArgs(driver.Value(isPublic), driver.Value(playlistID)).
					WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.UpdatePlaylistAccess(playlistID, isPublic)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
