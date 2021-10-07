package utils

import (
	"2021_2_LostPointer/models"
	"database/sql"
	"fmt"
	"math/rand"
)

const TracksSelectionLimit = 10
const AlbumsSelectionLimit = 4
const PlaylistsSelectionLimit = 4
const ArtistsSelectionLimit = 4

func isGenresExist(genres []string, db *sql.DB) (bool, error) {
	availableGenres := make(map[string]bool, 0)
	rows, err := db.Query(`SELECT name FROM genres`)
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

func createTracksRequestWithGenres(genres []string) string {
	baseRequest := `SELECT DISTINCT ON(tracks.album) tracks.id FROM tracks
LEFT JOIN genres g on tracks.genre = g.id
WHERE`

	for _, genre := range genres {
		baseRequest += fmt.Sprintf(" g.name='%s' OR", genre)
	}

	request := baseRequest[:len(baseRequest)-3]
	return request
}

func createTracksRequestWithIds(ids []int) string {
	baseRequest := `SELECT t.id,
							t.title,
							art.name,
							alb.title,
							t.explicit,
							g.name,
							t.number,
							t.file,
							t.listen_count,
							t.duration,
							alb.artwork
							FROM tracks t
							LEFT JOIN albums alb ON alb.id = t.album
							LEFT JOIN artists art ON art.id = t.artist
							LEFT JOIN genres g ON g.id = t.genre
							WHERE t.id IN (`

	for _, id := range ids {
		baseRequest += fmt.Sprintf("%d, ", id)
	}

	request := baseRequest[:len(baseRequest)-2] + ")"

	return request
}

func getTracksIdByGenre(db *sql.DB, genres []string) ([]int, error) {
	isGenresCorrect, err := isGenresExist(genres, db)
	if err != nil || isGenresCorrect == false {
		return nil, err
	}

	request := createTracksRequestWithGenres(genres)
	rows, err := db.Query(request)
	if err != nil {
		return nil, err
	}

	tracksId := make([]int, 0)
	var trackId int
	for rows.Next() {
		if err := rows.Scan(&trackId); err != nil {
			return nil, err
		}
		tracksId = append(tracksId, trackId)
	}

	return tracksId, nil
}

func GetTracks(amount int, db *sql.DB, genres ...string) ([]models.Track, error) {
	tracksRandom := make([]models.Track, 0)
	tracksId, err := getTracksIdByGenre(db, genres)
	if err != nil {
		return nil, err
	}

	fmt.Println(tracksId)
	tracksIdRandom := make([]int, 0)
	for i := 0; i < amount; {
		newRandomTrackId := tracksId[rand.Intn(len(tracksId)-1)]
		var isExist bool
		for _, id := range tracksIdRandom {
			if id == newRandomTrackId {
				isExist = true
			}
		}
		if !isExist {
			tracksIdRandom = append(tracksIdRandom, newRandomTrackId)
			i++
		}
	}
	request := createTracksRequestWithIds(tracksIdRandom)

	rows, err := db.Query(request)
	if err != nil {
		return nil, err
	}
	var track models.Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Title, &track.Artist, &track.Album, &track.Explicit, &track.Genre,
			&track.Number, &track.File, &track.ListenCount, &track.Duration, &track.Cover); err != nil {
			return nil, err
		}
		tracksRandom = append(tracksRandom, track)
	}

	return tracksRandom, nil
}

func GetAlbums(amount int, db *sql.DB) ([]models.Album, error) {
	albums := make([]models.Album, 0)
	rows, err := db.Query(fmt.Sprintf(`SELECT a.id,
										       a.title,
										       a.year,
										       art.name,
										       a.artwork,
										       a.track_count,
										       SUM(t.duration) as tracksDuration
										FROM albums a
										         LEFT JOIN artists art ON art.id = a.artist
										         JOIN tracks t on t.album = a.id
										WHERE art.name = 'Земфира'
										GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count
										LIMIT %d`, amount))

	if err != nil {
		return nil, err
	}

	var album models.Album
	for rows.Next() {
		if err := rows.Scan(&album.Id, &album.Title, &album.Year, &album.Artist, &album.ArtWork, &album.TracksCount, &album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

func GetArtists(amount int, db *sql.DB) ([]models.Artist, error) {
	artists := make([]models.Artist, 0)

	rows, err := db.Query(fmt.Sprintf(`SELECT artists.id, artists.name, artists.bio FROM artists LIMIT %d`, amount))
	if err != nil {
		return nil, err
	}

	var artist models.Artist
	for rows.Next() {
		if err := rows.Scan(&artist.Id, &artist.Name, &artist.Bio); err != nil {
			if artist.Bio == "" && artist.Name == "" {
				return nil, err
			}
		}
		artists = append(artists, artist)
	}

	return artists, nil
}

func GetPlaylists(amount int, db *sql.DB) ([]models.Playlist, error) {
	playlists := make([]models.Playlist, 0)
	rows, err := db.Query(fmt.Sprintf(`SELECT playlists.id, playlists.title, playlists.user FROM playlists LIMIT %d`, amount))
	if err != nil {
		return nil, err
	}

	var playlist models.Playlist
	for rows.Next() {
		if err := rows.Scan(&playlist.Id, &playlist.Name, &playlist.User); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func GetSelectionForHomePage(db *sql.DB) (*models.SelectionFroHomePage, error) {
	var selectionForHomePage = new(models.SelectionFroHomePage)
	var err error

	if selectionForHomePage.Tracks, err = GetTracks(TracksSelectionLimit, db, "Swing", "Jazz", "Rock", "Easy Listening", "Pop"); err != nil {
		return nil, err
	}
	if selectionForHomePage.Albums, err = GetAlbums(AlbumsSelectionLimit, db); err != nil {
		return nil, err
	}
	if selectionForHomePage.Playlists, err = GetPlaylists(PlaylistsSelectionLimit, db); err != nil {
		return nil, err
	}
	if selectionForHomePage.Artists, err = GetArtists(ArtistsSelectionLimit, db); err != nil {
		return nil, err
	}

	return selectionForHomePage, nil
}
