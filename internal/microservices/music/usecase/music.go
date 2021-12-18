package usecase

import (
	"2021_2_LostPointer/internal/models"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/music"
	"2021_2_LostPointer/internal/microservices/music/proto"
)

type MusicService struct {
	storage music.Storage
}

func NewMusicService(storage music.Storage) *MusicService {
	return &MusicService{storage: storage}
}

func getRandomID(slice []string, amount int) []string {
	log.Println("Slice:", slice)

	l := len(slice)
	random := make([]string, 0)
	for i := 0; i < amount; i++ {
		random = append(random, slice[rand.Intn(l)])
	}

	return random
}

func (service *MusicService) RandomTracks(ctx context.Context, metadata *proto.RandomTracksOptions) (*proto.Tracks, error) {
	selection, err := service.storage.GetSelections(metadata.UserID)
	if err == redis.Nil {
		favoriteTracks, err := service.storage.GetFavoritesID(metadata.UserID)
		if err != nil {
			return &proto.Tracks{}, status.Error(codes.Internal, err.Error())
		}

		tracksSelection, err := service.storage.GenerateSelections(metadata.UserID, favoriteTracks)
		if err != nil {
			return &proto.Tracks{}, status.Error(codes.Internal, err.Error())
		}

		selectionData := &models.Selection{Tracks: tracksSelection}
		selection.Tracks = tracksSelection
		err = service.storage.StoreSelection(metadata.UserID, selectionData)
		if err != nil {
			return &proto.Tracks{}, status.Error(codes.Internal, err.Error())
		}
	}

	tracksID := getRandomID(selection.Tracks, constants.HomePageTracksSelectionFavoritesAmount)
	tracksSelection, err := service.storage.GetTracksByTrackID(tracksID, metadata.UserID, metadata.IsAuthorized)
	if err != nil {
		log.Println("Error 1")
		return &proto.Tracks{}, status.Error(codes.Internal, err.Error())
	}

	tracksRandom, err := service.storage.RandomTracks(metadata.Amount - int64(len(tracksSelection)), metadata.UserID, metadata.IsAuthorized)
	if err != nil {
		log.Println("Error 2")
		return &proto.Tracks{}, status.Error(codes.Internal, err.Error())
	}

	tracks := &proto.Tracks{Tracks: append(tracksSelection, tracksRandom...)}
	return tracks, nil
}

func (service *MusicService) RandomAlbums(ctx context.Context, metadata *proto.RandomAlbumsOptions) (*proto.Albums, error) {
	albums, err := service.storage.RandomAlbums(metadata.Amount)
	if err != nil {
		return &proto.Albums{}, status.Error(codes.Internal, err.Error())
	}

	return albums, nil
}

func (service *MusicService) RandomArtists(ctx context.Context, metadata *proto.RandomArtistsOptions) (*proto.Artists, error) {
	artists, err := service.storage.RandomArtists(metadata.Amount)
	if err != nil {
		return &proto.Artists{}, status.Error(codes.Internal, err.Error())
	}

	return artists, nil
}

func (service *MusicService) ArtistProfile(ctx context.Context, metadata *proto.ArtistProfileOptions) (*proto.Artist, error) {
	artistData, err := service.storage.ArtistInfo(metadata.ArtistID)
	if err != nil {
		return &proto.Artist{}, status.Error(codes.Internal, err.Error())
	}

	artistData.Tracks, err = service.storage.ArtistTracks(metadata.ArtistID, metadata.UserID, metadata.IsAuthorized, constants.ArtistTracksSelectionAmount)
	if err != nil {
		return &proto.Artist{}, status.Error(codes.Internal, err.Error())
	}

	artistData.Albums, err = service.storage.ArtistAlbums(metadata.ArtistID, constants.ArtistAlbumsSelectionAmount)
	if err != nil {
		return &proto.Artist{}, status.Error(codes.Internal, err.Error())
	}

	return artistData, nil
}

func (service *MusicService) IncrementListenCount(ctx context.Context, metadata *proto.IncrementListenCountOptions) (*proto.IncrementListenCountEmpty, error) {
	err := service.storage.IncrementListenCount(metadata.ID)
	if err != nil {
		return &proto.IncrementListenCountEmpty{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.IncrementListenCountEmpty{}, nil
}

func (service *MusicService) AlbumPage(ctx context.Context, metadata *proto.AlbumPageOptions) (*proto.AlbumPageResponse, error) {
	album, err := service.storage.AlbumData(metadata.AlbumID)
	if err != nil {
		return &proto.AlbumPageResponse{}, status.Error(codes.Internal, err.Error())
	}

	album.Tracks, err = service.storage.AlbumTracks(metadata.AlbumID, metadata.UserID, metadata.IsAuthorized)
	if err != nil {
		return &proto.AlbumPageResponse{}, status.Error(codes.Internal, err.Error())
	}

	return album, nil
}

func (service *MusicService) Find(ctx context.Context, data *proto.FindOptions) (*proto.FindResponse, error) {
	data.Text = strings.TrimSpace(data.Text)
	if len(data.Text) == 0 {
		return &proto.FindResponse{}, nil
	}
	tracks, err := service.storage.FindTracksByFullWord(data.Text, data.UserID, data.IsAuthorized)
	if err != nil {
		return &proto.FindResponse{}, status.Error(codes.Internal, err.Error())
	}

	var FindTracksByPartial []*proto.Track
	if len(tracks) < constants.SearchTracksAmount {
		FindTracksByPartial, err = service.storage.FindTracksByPartial(data.Text, data.UserID, data.IsAuthorized)
		if err != nil {
			return &proto.FindResponse{}, status.Error(codes.Internal, err.Error())
		}
	}

	for _, track := range FindTracksByPartial {
		if !contains(tracks, track.ID) && len(tracks) < constants.SearchTracksAmount {
			tracks = append(tracks, track)
		}
	}

	artists, err := service.storage.FindArtists(data.Text)
	if err != nil {
		return &proto.FindResponse{}, status.Error(codes.Internal, err.Error())
	}

	albums, err := service.storage.FindAlbums(data.Text)
	if err != nil {
		return &proto.FindResponse{}, status.Error(codes.Internal, err.Error())
	}

	result := &proto.FindResponse{
		Tracks:  tracks,
		Albums:  albums,
		Artists: artists,
	}

	return result, nil
}

func (service *MusicService) UserPlaylists(ctx context.Context, data *proto.UserPlaylistsOptions) (*proto.PlaylistsData, error) {
	playlists, err := service.storage.UserPlaylists(data.UserID)
	if err != nil {
		return &proto.PlaylistsData{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.PlaylistsData{Playlists: playlists}, nil
}

func (service *MusicService) PlaylistPage(ctx context.Context, data *proto.PlaylistPageOptions) (*proto.PlaylistPageResponse, error) {
	doesExist, err := service.storage.DoesPlaylistExist(data.PlaylistID)
	if err != nil {
		return &proto.PlaylistPageResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !doesExist {
		return &proto.PlaylistPageResponse{}, status.Error(codes.NotFound, constants.PlaylistNotFoundMessage)
	}

	isOwner, err := service.storage.IsPlaylistOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return &proto.PlaylistPageResponse{}, status.Error(codes.Internal, err.Error())
	}
	isPublic, err := service.storage.IsPlaylistPublic(data.PlaylistID)
	if err != nil {
		return &proto.PlaylistPageResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !isOwner && !isPublic {
		return &proto.PlaylistPageResponse{}, status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage)
	}

	playlistInfo, err := service.storage.PlaylistInfo(data.PlaylistID)
	if err != nil {
		return &proto.PlaylistPageResponse{}, status.Error(codes.Internal, err.Error())
	}

	playlistTracks, err := service.storage.PlaylistTracks(data.PlaylistID, data.UserID)
	if err != nil {
		return &proto.PlaylistPageResponse{}, status.Error(codes.Internal, err.Error())
	}

	playlistData := &proto.PlaylistPageResponse{
		PlaylistID:   playlistInfo.PlaylistID,
		Title:        playlistInfo.Title,
		Artwork:      playlistInfo.Artwork,
		ArtworkColor: playlistInfo.ArtworkColor,
		Tracks:       playlistTracks,
		IsPublic:     playlistInfo.IsPublic,
		IsOwn:        isOwner,
	}

	return playlistData, nil
}

func (service *MusicService) AddTrackToFavorites(ctx context.Context, data *proto.AddTrackToFavoritesOptions) (*proto.AddTrackToFavoritesResponse, error) {
	isExist, err := service.storage.IsTrackInFavorites(data.UserID, data.TrackID)
	if err != nil {
		return &proto.AddTrackToFavoritesResponse{}, status.Error(codes.Internal, err.Error())
	}
	if isExist {
		return &proto.AddTrackToFavoritesResponse{}, status.Error(codes.PermissionDenied, constants.TrackAlreadyInFavorites)
	}

	err = service.storage.AddTrackToFavorite(data.UserID, data.TrackID)
	if err != nil {
		return &proto.AddTrackToFavoritesResponse{}, status.Error(codes.NotFound, constants.TrackNotFound)
	}

	return &proto.AddTrackToFavoritesResponse{}, nil
}

func (service *MusicService) DeleteTrackFromFavorites(ctx context.Context, data *proto.DeleteTrackFromFavoritesOptions) (*proto.DeleteTrackFromFavoritesResponse, error) {
	isExist, err := service.storage.IsTrackInFavorites(data.UserID, data.TrackID)
	if err != nil {
		return &proto.DeleteTrackFromFavoritesResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !isExist {
		return &proto.DeleteTrackFromFavoritesResponse{}, status.Error(codes.PermissionDenied, constants.TrackNotInFavorites)
	}

	err = service.storage.DeleteTrackFromFavorites(data.UserID, data.TrackID)
	if err != nil {
		return &proto.DeleteTrackFromFavoritesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.DeleteTrackFromFavoritesResponse{}, nil
}

func (service *MusicService) GetFavoriteTracks(ctx context.Context, data *proto.UserFavoritesOptions) (*proto.Tracks, error) {
	tracks := new(proto.Tracks)
	var err error

	tracks.Tracks, err = service.storage.GetFavorites(data.UserID)
	if err != nil {
		return &proto.Tracks{}, status.Error(codes.Internal, err.Error())
	}

	return tracks, nil
}

func contains(tracks []*proto.Track, trackID int64) bool {
	for _, currentTrack := range tracks {
		if currentTrack.ID == trackID {
			return true
		}
	}
	return false
}
