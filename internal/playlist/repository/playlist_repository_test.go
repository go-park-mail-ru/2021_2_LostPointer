package repository

import (
	"2021_2_LostPointer/internal/models"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestPlaylistRepository_GetRandom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewPlaylistRepository(db)

	playlist := models.Playlist{
		Id:   1,
		Name: "awa",
		User: 1,
	}

	tests := []struct {
		name          string
		amount        int
		mock          func()
		expected      []models.Playlist
		expectedError bool
	}{
		{
			name:   "get 4 random playlists",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user"})
				for i := 0; i < 4; i++ {
					rows.AddRow(playlist.Id, playlist.Name, playlist.User)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT playlists.id, playlists.title, playlists.user "+
					"FROM playlists WHERE playlists.user = $1 LIMIT $2")).WithArgs(driver.Value(0), driver.Value(4)).WillReturnRows(rows)
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
			name:   "get 10 random playlists",
			amount: 10,
			mock: func() {
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user"})
				for i := 0; i < 10; i++ {
					rows.AddRow(playlist.Id, playlist.Name, playlist.User)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT playlists.id, playlists.title, playlists.user "+
					"FROM playlists WHERE playlists.user = $1 LIMIT $2")).WithArgs(driver.Value(0), driver.Value(10)).WillReturnRows(rows)
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
			name:   "get 100 random playlists",
			amount: 100,
			mock: func() {
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user"})
				for i := 0; i < 100; i++ {
					rows.AddRow(playlist.Id, playlist.Name, playlist.User)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT playlists.id, playlists.title, playlists.user "+
					"FROM playlists WHERE playlists.user = $1 LIMIT $2")).WithArgs(driver.Value(0), driver.Value(100)).WillReturnRows(rows)
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 100; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: false,
		},
		{
			name:   "query error",
			amount: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user"})
				for i := 0; i < 1; i++ {
					rows.AddRow(playlist.Id, playlist.Name, playlist.User)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT playlists.id, playlists.title, playlists.user "+
					"FROM playlists WHERE playlists.user = $1 LIMIT $2")).WithArgs(driver.Value(0), driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 1; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: true,
		},
		{
			name:   "scan error",
			amount: 1,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(playlist.Id, playlist.Name, playlist.User, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT playlists.id, playlists.title, playlists.user "+
					"FROM playlists WHERE playlists.user = $1 LIMIT $2")).WithArgs(driver.Value(0), driver.Value(1)).WillReturnRows(rows)
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 1; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.Get(test.amount, 0)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}
