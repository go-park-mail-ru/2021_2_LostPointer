package usecase

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/music"
	"net/http"
)

const TracksCollectionLimit = 10
const AlbumCollectionLimit = 4
const PlaylistsCollectionLimit = 4
const ArtistsCollectionLimit = 4

const DatabaseNotResponding = "Database not responding"

type MusicUseCase struct {
	MusicRepository music.MusicRepository
}

func NewMusicUseCase(musicRepository music.MusicRepository) MusicUseCase {
	return MusicUseCase{MusicRepository: musicRepository}
}

func (musicUseCase MusicUseCase) GetMusicCollection(isAuthorized bool) (*models.MusicCollection, *models.CustomError) {
	var (
		collection = new(models.MusicCollection)
		err        error
	)

	if collection.Tracks, err = musicUseCase.GetTracksForCollection(TracksCollectionLimit, isAuthorized); err != nil {
		return nil, &models.CustomError{
			ErrorType: http.StatusInternalServerError,
			OriginalError: err,
			Message: DatabaseNotResponding,
		}
	}
	if collection.Albums, err = musicUseCase.GetAlbumsForCollection(AlbumCollectionLimit); err != nil {
		return nil, &models.CustomError{
			ErrorType: http.StatusInternalServerError,
			OriginalError: err,
			Message: DatabaseNotResponding,
		}
	}
	if collection.Artists, err = musicUseCase.GetArtistsForCollection(ArtistsCollectionLimit); err != nil {
		return nil, &models.CustomError{
			ErrorType: http.StatusInternalServerError,
			OriginalError: err,
			Message: DatabaseNotResponding,
		}
	}
	if collection.Playlists, err = musicUseCase.GetPlaylistsForCollection(PlaylistsCollectionLimit); err != nil {
		return nil, &models.CustomError{
			ErrorType: http.StatusInternalServerError,
			OriginalError: err,
			Message: DatabaseNotResponding,
		}
	}

	return collection, nil
}

func (musicUseCase MusicUseCase) GetTracksForCollection(amount int, isAuthorized bool) ([]models.Track, error) {
	tracks, err := musicUseCase.MusicRepository.GetRandomTracks(amount, isAuthorized)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (musicUseCase MusicUseCase) GetAlbumsForCollection(amount int) ([]models.Album, error) {
	albums, err := musicUseCase.MusicRepository.GetRandomAlbums(amount)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (musicUseCase MusicUseCase) GetArtistsForCollection(amount int) ([]models.Artist, error) {
	artists, err := musicUseCase.MusicRepository.GetRandomArtists(amount)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (musicUseCase MusicUseCase) GetPlaylistsForCollection(amount int) ([]models.Playlist, error) {
	playlists, err := musicUseCase.MusicRepository.GetRandomPlaylists(amount)
	if err != nil {
		return nil, err
	}

	return playlists, err
}