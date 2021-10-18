package usecase

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/music"
	"github.com/labstack/echo"
)

const TracksCollectionLimit = 10
const AlbumCollectionLimit = 4
const PlaylistsCollectionLimit = 4
const ArtistsCollectionLimit = 4

type MusicUseCase struct {
	MusicRepository music.MusicRepositoryIFace
}

func NewMusicUseCase(musicRepository music.MusicRepositoryIFace) MusicUseCase {
	return MusicUseCase{MusicRepository: musicRepository}
}

func (musicUseCase MusicUseCase) GetMusicCollection(ctx echo.Context) (*models.MusicCollection, error) {
	var (
		collection = new(models.MusicCollection)
		err        error
	)

	isAuthorized := ctx.Get("is_authorized").(bool)

	if collection.Tracks, err = musicUseCase.GetTracksForCollection(TracksCollectionLimit, isAuthorized); err != nil {
		return nil, err
	}
	if collection.Albums, err = musicUseCase.GetAlbumsForCollection(AlbumCollectionLimit); err != nil {
		return nil, err
	}
	if collection.Artists, err = musicUseCase.GetArtistsForCollection(ArtistsCollectionLimit); err != nil {
		return nil, err
	}
	if collection.Playlists, err = musicUseCase.GetPlaylistsForCollection(PlaylistsCollectionLimit); err != nil {
		return nil, err
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
