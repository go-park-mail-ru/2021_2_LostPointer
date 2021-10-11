package repository

import (
	"2021_2_LostPointer/pkg/models"
	"database/sql"
	"fmt"
	"strings"
)

const (
	DistinctOnNone = iota
	DistinctOnAlbums
	DistinctOnArtists
	GettingWithGenres
	GettingWithID
)

type MusicRepository struct {
	Database *sql.DB
}

func NewMusicRepository(db *sql.DB) MusicRepository {
	return MusicRepository{Database: db}
}

func (musicRepository MusicRepository) GetTracks(request string) ([]models.Track, error) {
	rows, err := musicRepository.Database.Query(request)
	if err != nil {
		return nil, err
	}

	tracks := make([]models.Track, 0)
	var track models.Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Title, &track.Artist, &track.Album, &track.Explicit, &track.Genre,
			&track.Number, &track.Cover, &track.ListenCount, &track.Duration); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (musicRepository MusicRepository) CreateTracksRequestWithParameters(gettingWith uint8, parameters interface{}, distinctOn uint8) string {
	var baseRequest = `SELECT `
	switch distinctOn {
	case DistinctOnAlbums:
		baseRequest += `DISTINCT ON(alb.title) `
	case DistinctOnArtists:
		baseRequest += `DISTINCT ON(art.name) `
	case DistinctOnNone:
	default:
	}

	baseRequest += `tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id `

	var params []string
	var attribute string
	var parameterTemplate string
	switch gettingWith {
	case GettingWithGenres:
		params = parameters.([]string)
		attribute = `g.name`
		parameterTemplate = `'%s'`
	case GettingWithID:
		paramsInt := parameters.([]int64)
		params = strings.SplitN(strings.Trim(fmt.Sprint(paramsInt), "[]"), " ", len(paramsInt))
		attribute = `tracks.id`
		parameterTemplate = `%s`
	}

	if len(params) != 0 {
		baseRequest += `WHERE ` + attribute + ` IN (`
		for _, param := range params {
			baseRequest += fmt.Sprintf(parameterTemplate, param) + `, `
		}

		baseRequest = baseRequest[:len(baseRequest)-2] + `)`
	}

	return baseRequest
}

func (musicRepository MusicRepository) IsGenreExist(genres []string) (bool, error) {
	availableGenres := make(map[string]bool, 0)
	rows, err := musicRepository.Database.Query(`SELECT name FROM genres`)
	if err != nil {
		return false, err
	}

	var genre string
	for rows.Next() {
		if err := rows.Scan(&genre); err != nil {
			return false, err
		}
		availableGenres[genre] = true
	}

	for _, genre := range genres {
		if availableGenres[genre] == false {
			return false, nil
		}
	}

	return true, nil
}

func (musicRepository MusicRepository) CreateAlbumsDefaultRequest(amount int) string {
	request := fmt.Sprintf(`SELECT a.id, a.title, a.year, art.name,
									a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a
									LEFT JOIN artists art ON art.id = a.artist
									JOIN tracks t on t.album = a.id
									WHERE art.name = 'Земфира'
									GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count
									LIMIT %d`, amount)

	return request
}

func (musicRepository MusicRepository) GetAlbums(request string) ([]models.Album, error) {
	rows, err := musicRepository.Database.Query(request)
	if err != nil {
		return nil, err
	}

	albums := make([]models.Album, 0)
	var album models.Album
	for rows.Next() {
		if err := rows.Scan(&album.Id, &album.Title, &album.Year, &album.Artist, &album.ArtWork, &album.TracksCount, &album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

func (musicRepository MusicRepository) CreateArtistsDefaultRequest(amount int) string {
	request := fmt.Sprintf(`SELECT artists.id, artists.name, artists.avatar FROM artists
									WHERE artists.avatar != 'frank_sinatra.jpg' LIMIT %d`, amount)

	return request
}

func (musicRepository MusicRepository) GetArtists(request string) ([]models.Artist, error) {
	rows, err := musicRepository.Database.Query(request)
	if err != nil {
		return nil, err
	}

	artists := make([]models.Artist, 0)
	var artist models.Artist
	for rows.Next() {
		if err := rows.Scan(&artist.Id, &artist.Name, &artist.Avatar); err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}

	return artists, nil
}

func (musicRepository MusicRepository) CreatePlaylistsDefaultRequest(amount int) string {
	request := fmt.Sprintf(`SELECT playlists.id, playlists.title, playlists.user FROM playlists LIMIT %d`, amount)

	return request
}

func (musicRepository MusicRepository) GetPlaylists(request string) ([]models.Playlist, error) {
	rows, err := musicRepository.Database.Query(request)
	if err != nil {
		return nil, err
	}

	playlists := make([]models.Playlist, 0)
	var playlist models.Playlist
	for rows.Next() {
		if err := rows.Scan(&playlist.Id, &playlist.Name, &playlist.User); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}
