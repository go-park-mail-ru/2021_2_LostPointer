package usecase

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/music/repository"
	"math/rand"
)

const TracksCollectionLimit = 10
const AlbumCollectionLimit = 4
const PlaylistsCollectionLimit = 4
const ArtistsCollectionLimit = 4

var defaultGenres = []string{"Swing", "Jazz", "Rock", "Easy Listening", "Pop"}

type MusicUseCase struct {
	MusicRepository repository.MusicRepository
}

func NewMusicUseCase(musicRepository repository.MusicRepository) MusicUseCase {
	return MusicUseCase{MusicRepository: musicRepository}
}

func (musicUseCase MusicUseCase) GetMusicCollection() (*models.MusicCollection, error) {
	var collection = new(models.MusicCollection)
	var err error

	if collection.Tracks, err = musicUseCase.GetTracksForCollection(TracksCollectionLimit, defaultGenres); err != nil {
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

func (musicUseCase MusicUseCase) GetTracksForCollection(amount int, genres []string) ([]models.Track, error) {
	requestWithGenres := musicUseCase.MusicRepository.CreateTracksRequestWithParameters(repository.GettingWithGenres, genres, repository.DistinctOnAlbums)
	tracksByGenre, err := musicUseCase.MusicRepository.GetTracks(requestWithGenres)
	if err != nil {
		return nil, err
	}

	randomTracksIDMap := make(map[int64]bool, 0)
	for i := 0; i < amount; {
		newRandomTrackId := tracksByGenre[rand.Intn(len(tracksByGenre)-1)].Id
		if randomTracksIDMap[newRandomTrackId] == false {
			randomTracksIDMap[newRandomTrackId] = true
			i++
		}

	}
	randomTracksIDArray := make([]int64, 0)
	for key, _ := range randomTracksIDMap {
		randomTracksIDArray = append(randomTracksIDArray, key)
	}

	request := musicUseCase.MusicRepository.CreateTracksRequestWithParameters(repository.GettingWithID, randomTracksIDArray, repository.DistinctOnAlbums)
	randomTracks, err := musicUseCase.MusicRepository.GetTracks(request)
	if err != nil {
		return nil, err
	}

	return randomTracks, nil
}

func (musicUseCase MusicUseCase) GetAlbumsForCollection(amount int) ([]models.Album, error) {
	request := musicUseCase.MusicRepository.CreateAlbumsDefaultRequest(amount)
	albums, err := musicUseCase.MusicRepository.GetAlbums(request)
	if err != nil {
		return nil, err
	}

	return albums, nil
}

func (musicUseCase MusicUseCase) GetArtistsForCollection(amount int) ([]models.Artist, error) {
	request := musicUseCase.MusicRepository.CreateArtistsDefaultRequest(amount)
	artists, err := musicUseCase.MusicRepository.GetArtists(request)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (musicUseCase MusicUseCase) GetPlaylistsForCollection(amount int) ([]models.Playlist, error) {
	request := musicUseCase.MusicRepository.CreatePlaylistsDefaultRequest(amount)
	playlists, err := musicUseCase.MusicRepository.GetPlaylists(request)
	if err != nil {
		return nil, err
	}

	return playlists, err
}
