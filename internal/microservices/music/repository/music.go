package repository

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/music/proto"
	"2021_2_LostPointer/pkg/wrapper"
	"database/sql"
	"log"
	"os"
)

type MusicStorage struct {
	db *sql.DB
}

func NewMusicStorage(db *sql.DB) *MusicStorage {
	return &MusicStorage{db: db}
}

func (storage *MusicStorage) RandomTracks(amount int64, isAuthorized bool) (*proto.Tracks, error) {
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
		wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb") + ", " +
		wrapper.Wrapper([]string{"id", "name"}, "art") + ", " +
		wrapper.Wrapper([]string{"name"}, "g") +
		`
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		ORDER BY RANDOM() DESC LIMIT $1`

	rows, err := storage.db.Query(query, amount)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	tracks := make([]*proto.Track, 0, amount)
	//nolint:dupl
	for rows.Next() {
		track := &proto.Track{}
		track.Album = &proto.Album{}
		track.Artist = &proto.Artist{}
		if err = rows.Scan(&track.ID, &track.Title, &track.Explicit, &track.Number, &track.File, &track.ListenCount,
			&track.Duration, &track.Lossless, &track.Album.ID, &track.Album.Title, &track.Album.Artwork, &track.Album.ArtworkColor, &track.Artist.ID,
			&track.Artist.Name, &track.Genre); err != nil {
			return nil, err
		}
		if !isAuthorized {
			track.File = ""
		}
		tracks = append(tracks, track)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	tracksList := &proto.Tracks{Tracks: tracks}
	return tracksList, nil
}

func (storage *MusicStorage) RandomAlbums(amount int64) (*proto.Albums, error) {
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb") + ", " +
		wrapper.Wrapper([]string{"name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
		`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY RANDOM()
		LIMIT $1
		`

	rows, err := storage.db.Query(query, amount)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	albums := make([]*proto.Album, 0, amount)
	for rows.Next() {
		album := &proto.Album{}
		if err = rows.Scan(&album.ID, &album.Title, &album.Year, &album.Artwork, &album.TracksAmount, &album.ArtworkColor, &album.Artist,
			&album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	albumsList := &proto.Albums{Albums: albums}
	return albumsList, nil
}

func (storage *MusicStorage) RandomArtists(amount int64) (*proto.Artists, error) {
	query := `
		SELECT
		id, name, avatar
		FROM artists
		ORDER BY RANDOM()
		LIMIT $1
	`

	rows, err := storage.db.Query(query, amount)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	artists := make([]*proto.Artist, 0, amount)
	for rows.Next() {
		artist := &proto.Artist{}
		artist.Tracks = []*proto.Track{}
		artist.Albums = []*proto.Album{}
		if err = rows.Scan(&artist.ID, &artist.Name, &artist.Avatar); err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	artistsList := &proto.Artists{Artists: artists}
	return artistsList, nil
}

func (storage *MusicStorage) ArtistInfo(artistID int64) (*proto.Artist, error) {
	query := `
		SELECT id, name, avatar, video
		FROM artists
		WHERE id = $1
	`
	var video string
	artist := &proto.Artist{}
	err := storage.db.QueryRow(query, artistID).Scan(&artist.ID, &artist.Name, &artist.Avatar, &video)
	if err != nil {
		return nil, err
	}

	artist.Avatar = os.Getenv("ARTISTS_ROOT_PREFIX") + artist.Avatar + constants.ImageExtension

	if len(video) > 1 {
		artist.Video = os.Getenv("MOV_ROOT_PREFIX") + video + constants.VideoExtension
	}

	return artist, nil
}

func (storage *MusicStorage) ArtistTracks(artistID int64, isAuthorized bool, amount int64) ([]*proto.Track, error) {
	var err error
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
		wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb") + ", " +
		wrapper.Wrapper([]string{"name"}, "g") +
		`
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		WHERE t.artist = $1
		ORDER BY t.listen_count DESC LIMIT $2
		`
	rows, err := storage.db.Query(query, artistID, amount)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	tracks := make([]*proto.Track, 0, amount)
	for rows.Next() {
		track := &proto.Track{}
		track.Album = &proto.Album{}
		track.Artist = &proto.Artist{}
		if err = rows.Scan(&track.ID, &track.Title, &track.Explicit, &track.Number, &track.File, &track.ListenCount,
			&track.Duration, &track.Lossless, &track.Album.ID, &track.Album.Title, &track.Album.Artwork, &track.Album.ArtworkColor, &track.Genre); err != nil {
			return nil, err
		}
		if !isAuthorized {
			track.File = ""
		}
		tracks = append(tracks, track)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (storage *MusicStorage) ArtistAlbums(artistID int64, amount int64) ([]*proto.Album, error) {
	var err error
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "year", "artwork", "artwork_color"}, "alb") + ", SUM(t.duration) AS tracksDuration" +
		`
		FROM albums alb
		JOIN tracks t ON alb.id = t.album
		WHERE alb.artist = $1
		GROUP BY alb.id, alb.title, alb.year, alb.artwork, alb.track_count
		ORDER BY alb.year DESC LIMIT $2
		`

	rows, err := storage.db.Query(query, artistID, amount)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	albums := make([]*proto.Album, 0, amount)
	for rows.Next() {
		album := &proto.Album{}
		if err = rows.Scan(&album.ID, &album.Title, &album.Year, &album.Artwork, &album.ArtworkColor, &album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (storage *MusicStorage) IncrementListenCount(trackID int64) error {
	query := `UPDATE tracks SET listen_count = listen_count + 1 WHERE id=$1`

	err := storage.db.QueryRow(query, trackID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (storage *MusicStorage) AlbumData(albumID int64) (*proto.AlbumPageResponse, error) {
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "year", "artwork", "artwork_color", "track_count"}, "alb") + ", " +
		wrapper.Wrapper([]string{"id", "name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
		`
		FROM albums alb
		JOIN artists art ON alb.artist = art.id
		JOIN tracks t ON alb.id = t.album
		WHERE alb.id = $1
		GROUP BY alb.id, art.id
		`

	album := &proto.AlbumPageResponse{}
	album.Artist = &proto.Artist{}
	err := storage.db.QueryRow(query, albumID).Scan(&album.AlbumID, &album.Title, &album.Year, &album.Artwork, &album.ArtworkColor,
		&album.TracksCount, &album.Artist.ID, &album.Artist.Name, &album.TracksDuration)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (storage *MusicStorage) AlbumTracks(albumID int64, isAuthorized bool) ([]*proto.AlbumTrack, error) {
	var err error
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
		wrapper.Wrapper([]string{"name"}, "g") +
		`
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		WHERE alb.id = $1
		ORDER BY t.number
		`
	rows, err := storage.db.Query(query, albumID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	tracks := make([]*proto.AlbumTrack, 0)
	for rows.Next() {
		track := &proto.AlbumTrack{}
		if err = rows.Scan(&track.ID, &track.Title, &track.Explicit, &track.Number, &track.File, &track.ListenCount,
			&track.Duration, &track.Lossless, &track.Genre); err != nil {
			return nil, err
		}
		if !isAuthorized {
			track.File = ""
		}
		tracks = append(tracks, track)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

//nolint:dupl
func (storage *MusicStorage) FindTracksByFullWord(text string, isAuthorized bool) ([]*proto.Track, error) {
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
		wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb") + ", " +
		wrapper.Wrapper([]string{"id", "name"}, "art") + ", " +
		wrapper.Wrapper([]string{"name"}, "g") +
		`
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		WHERE t.id IN (
			SELECT track
		    FROM search
		    WHERE concatenation @@ plainto_tsquery($1)
		)
		ORDER BY t.listen_count DESC LIMIT $2`
	rows, err := storage.db.Query(query, text, constants.SearchTracksAmount)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	tracks := make([]*proto.Track, 0, constants.SearchTracksAmount)
	for rows.Next() {
		track := &proto.Track{}
		track.Album = &proto.Album{}
		track.Artist = &proto.Artist{}
		if err = rows.Scan(&track.ID, &track.Title, &track.Explicit, &track.Number, &track.File, &track.ListenCount,
			&track.Duration, &track.Lossless, &track.Album.ID, &track.Album.Title, &track.Album.Artwork, &track.Album.ArtworkColor,
			&track.Artist.ID, &track.Artist.Name, &track.Genre); err != nil {
			return nil, err
		}
		if !isAuthorized {
			track.File = ""
		}
		tracks = append(tracks, track)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

//nolint:dupl
func (storage *MusicStorage) FindTracksByPartial(text string, isAuthorized bool) ([]*proto.Track, error) {
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
		wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb") + ", " +
		wrapper.Wrapper([]string{"id", "name"}, "art") + ", " +
		wrapper.Wrapper([]string{"name"}, "g") +
		`
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		WHERE t.id IN (
		    SELECT track
		    FROM test
		    WHERE concat ILIKE $1
		)
		ORDER BY t.listen_count DESC LIMIT $2`
	rows, err := storage.db.Query(query, "%"+text+"%", constants.SearchTracksAmount)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	tracks := make([]*proto.Track, 0, constants.SearchTracksAmount)
	for rows.Next() {
		track := &proto.Track{}
		track.Album = &proto.Album{}
		track.Artist = &proto.Artist{}
		if err = rows.Scan(&track.ID, &track.Title, &track.Explicit, &track.Number, &track.File, &track.ListenCount,
			&track.Duration, &track.Lossless, &track.Album.ID, &track.Album.Title, &track.Album.Artwork,
			&track.Album.ArtworkColor, &track.Artist.ID, &track.Artist.Name, &track.Genre); err != nil {
			return nil, err
		}
		if !isAuthorized {
			track.File = ""
		}
		tracks = append(tracks, track)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (storage *MusicStorage) FindArtists(text string) ([]*proto.Artist, error) {
	query := `
		SELECT id, name, avatar
		FROM artists
		WHERE name ILIKE $1
		LIMIT $2
	`
	rows, err := storage.db.Query(query, "%"+text+"%", constants.SearchArtistsAmount)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	artists := make([]*proto.Artist, 0, constants.SearchArtistsAmount)
	for rows.Next() {
		artist := &proto.Artist{}
		artist.Tracks = []*proto.Track{}
		artist.Albums = []*proto.Album{}
		if err = rows.Scan(&artist.ID, &artist.Name, &artist.Avatar); err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (storage *MusicStorage) FindAlbums(text string) ([]*proto.Album, error) {
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "year", "artwork", "track_count", "artwork_color"}, "alb") + ", " +
		wrapper.Wrapper([]string{"name"}, "art") + ", SUM(t.duration) AS tracksDuration" +
		`
		FROM albums alb
		JOIN artists art ON art.id = alb.artist
		JOIN tracks t ON alb.id = t.album
		WHERE alb.title ILIKE $1
		GROUP BY alb.id, alb.title, alb.year, art.name, alb.artwork, alb.track_count
		ORDER BY year DESC
		LIMIT $2
		`
	rows, err := storage.db.Query(query, "%"+text+"%", constants.SearchAlbumsAmount)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	albums := make([]*proto.Album, 0, constants.SearchAlbumsAmount)
	for rows.Next() {
		album := &proto.Album{}
		if err = rows.Scan(&album.ID, &album.Title, &album.Year, &album.Artwork, &album.TracksAmount, &album.ArtworkColor, &album.Artist,
			&album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (storage *MusicStorage) IsPlaylistOwner(playlistID int64, userID int64) (bool, error) {
	query := `SELECT * FROM playlists WHERE id=$1 AND user_id=$2`

	rows, err := storage.db.Query(query, playlistID, userID)
	if err != nil {
		return true, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	if rows.Next() {
		return true, nil
	}
	if err = rows.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func (storage *MusicStorage) IsPlaylistPublic(playlistID int64) (bool, error) {
	query := `SELECT * FROM playlists WHERE id=$1 AND is_public=true`

	rows, err := storage.db.Query(query, playlistID)
	if err != nil {
		return true, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	if rows.Next() {
		return true, nil
	}
	if err = rows.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func (storage *MusicStorage) PlaylistTracks(playlistID int64) ([]*proto.Track, error) {
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
		wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb") + ", " +
		wrapper.Wrapper([]string{"id", "name"}, "art") + ", " +
		wrapper.Wrapper([]string{"name"}, "g") +
		`
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		WHERE t.id IN (
			SELECT track
			FROM playlist_tracks
			WHERE playlist=$1
		)`

	rows, err := storage.db.Query(query, playlistID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	tracks := make([]*proto.Track, 0, constants.SearchTracksAmount)
	for rows.Next() {
		track := &proto.Track{}
		track.Album = &proto.Album{}
		track.Artist = &proto.Artist{}
		if err = rows.Scan(&track.ID, &track.Title, &track.Explicit, &track.Number, &track.File, &track.ListenCount,
			&track.Duration, &track.Lossless, &track.Album.ID, &track.Album.Title, &track.Album.Artwork,
			&track.Album.ArtworkColor, &track.Artist.ID, &track.Artist.Name, &track.Genre); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (storage *MusicStorage) PlaylistInfo(playlistID int64) (*proto.PlaylistData, error) {
	query := `SELECT id, title, artwork, artwork_color, is_public FROM playlists WHERE id=$1`

	playlistInfo := &proto.PlaylistData{}
	err := storage.db.QueryRow(query, playlistID).Scan(
		&playlistInfo.PlaylistID, &playlistInfo.Title, &playlistInfo.Artwork, &playlistInfo.ArtworkColor, &playlistInfo.IsPublic)
	if err != nil {
		return nil, err
	}

	playlistInfo.Artwork = os.Getenv("PLAYLIST_ROOT_PREFIX") + playlistInfo.Artwork + constants.PlaylistArtworkExtension384px

	return playlistInfo, nil
}

func (storage *MusicStorage) UserPlaylists(userID int64) ([]*proto.PlaylistData, error) {
	query := `SELECT id, title, artwork, is_public, user_id=$1 AS is_own FROM playlists WHERE user_id=$1 OR is_public=true`

	rows, err := storage.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	playlists := make([]*proto.PlaylistData, 0)
	for rows.Next() {
		playlist := &proto.PlaylistData{}
		if err = rows.Scan(&playlist.PlaylistID, &playlist.Title, &playlist.Artwork, &playlist.IsPublic, &playlist.IsOwn); err != nil {
			return nil, err
		}
		playlist.Artwork = os.Getenv("PLAYLIST_ROOT_PREFIX") + playlist.Artwork + constants.PlaylistArtworkExtension100px
		playlists = append(playlists, playlist)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return playlists, nil
}

func (storage *MusicStorage) DoesPlaylistExist(playlistID int64) (bool, error) {
	query := `SELECT * FROM playlists WHERE id=$1`

	rows, err := storage.db.Query(query, playlistID)
	if err != nil {
		return true, err
	}
	err = rows.Err()
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()
	if rows.Next() {
		return true, nil
	}

	return false, nil
}

func (storage *MusicStorage) AddTrackToFavorite(userID int64, trackID int64) error {
	query := `INSERT INTO likes(user_id, track_id) VALUES ($1, $2)`

	err := storage.db.QueryRow(query, userID, trackID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (storage *MusicStorage) DeleteTrackFromFavorites(userID int64, trackID int64) error {
	query := `DELETE FROM likes WHERE user_id = $1 and track_id = $2`

	_, err := storage.db.Exec(query, userID, trackID)
	if err != nil {
		return err
	}

	return nil
}

func (storage *MusicStorage) GetFavorites(userID int64) ([]*proto.Track, error) {
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
		wrapper.Wrapper([]string{"id", "title", "artwork", "artwork_color"}, "alb") + ", " +
		wrapper.Wrapper([]string{"id", "name"}, "art") + ", " +
		wrapper.Wrapper([]string{"name"}, "g") +
		`
		FROM tracks t
		JOIN genres g ON t.genre = g.id
		JOIN albums alb ON t.album = alb.id
		JOIN artists art ON t.artist = art.id
		WHERE t.id IN (
			SELECT track_id
			FROM likes
			WHERE user_id = $1
		)`

	rows, err := storage.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	tracks := make([]*proto.Track, 0)
	for rows.Next() {
		track := &proto.Track{}
		track.Album = &proto.Album{}
		track.Artist = &proto.Artist{}
		if err = rows.Scan(&track.ID, &track.Title, &track.Explicit, &track.Number, &track.File, &track.ListenCount,
			&track.Duration, &track.Lossless, &track.Album.ID, &track.Album.Title, &track.Album.Artwork,
			&track.Album.ArtworkColor, &track.Artist.ID, &track.Artist.Name, &track.Genre); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (storage *MusicStorage) IsTrackInFavorites(userID int64, trackID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND track_id = $2)`
	var isExist bool

	err := storage.db.QueryRow(query, userID, trackID).Scan(isExist)
	if err != nil {
		return false, err
	}

	return isExist, nil
}
