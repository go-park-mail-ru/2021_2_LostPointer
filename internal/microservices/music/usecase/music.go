package usecase

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/music/proto"
	"2021_2_LostPointer/internal/microservices/music/repository"
)

type MusicService struct {
	storage repository.MusicStorage
}

func NewMusicService(storage repository.MusicStorage) *MusicService {
	return &MusicService{storage: storage}
}

func (service *MusicService) RandomTracks(ctx context.Context, metadata *proto.RandomTracksOptions) (*proto.Tracks, error) {
	tracks, err := service.storage.RandomTracks(metadata.Amount, metadata.IsAuthorized)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return tracks, nil
}

func (service *MusicService) RandomAlbums(ctx context.Context, metadata *proto.RandomAlbumsOptions) (*proto.Albums, error) {
	albums, err := service.storage.RandomAlbums(metadata.Amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return albums, nil
}

func (service *MusicService) RandomArtists(ctx context.Context, metadata *proto.RandomArtistsOptions) (*proto.Artists, error) {
	artists, err := service.storage.RandomArtists(metadata.Amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return artists, nil
}

func (service *MusicService) ArtistProfile(ctx context.Context, metadata *proto.ArtistProfileOptions) (*proto.Artist, error) {
	artistData, err := service.storage.GetArtistInfo(metadata.ArtistID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	artistData.Tracks, err = service.storage.GetArtistTracks(metadata.ArtistID, metadata.IsAuthorized, constants.ArtistTracksSelectionAmount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	artistData.Albums, err = service.storage.GetArtistAlbums(metadata.ArtistID, constants.ArtistAlbumsSelectionAmount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return artistData, nil
}

func (service *MusicService) IncrementListenCount(ctx context.Context, metadata *proto.IncrementListenCountOptions) (*proto.IncrementListenCountEmpty, error) {
	err := service.storage.IncrementListenCount(metadata.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.IncrementListenCountEmpty{}, nil
}

func (service *MusicService) AlbumPage(ctx context.Context, metadata *proto.AlbumPageOptions) (*proto.AlbumPageResponse, error) {
	album, err := service.storage.AlbumData(metadata.AlbumID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	album.Tracks, err = service.storage.AlbumTracks(metadata.AlbumID, metadata.IsAuthorized)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return album, nil
}

func (service *MusicService) Find(ctx context.Context, data *proto.FindOptions) (*proto.FindResponse, error) {
	tracks, err := service.storage.FindTracksByFullWord(data.Text, data.IsAuthorized)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var FindTracksByPartial []*proto.Track
	if len(tracks) < constants.SearchTracksAmount {
		FindTracksByPartial, err = service.storage.FindTracksByPartial(data.Text, data.IsAuthorized)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	for _, track := range FindTracksByPartial {
		if !contains(tracks, track.ID) {
			tracks = append(tracks, track)
		}
	}

	artists, err := service.storage.FindArtists(data.Text)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	albums, err := service.storage.FindAlbums(data.Text)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := &proto.FindResponse{
		Tracks:  tracks,
		Albums:  albums,
		Artists: artists,
	}

	return result, nil
}

func (service *MusicService) UserPlaylists(ctx context.Context, data *proto.UserPlaylistsOptions) (*proto.PlaylistsData, error) {
	playlists, err := service.storage.GetUserPlaylists(data.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.PlaylistsData{Playlists: playlists}, nil
}

func (service *MusicService) PlaylistPage(ctx context.Context, data *proto.PlaylistPageOptions) (*proto.PlaylistPageResponse, error) {
	doesExist, err := service.storage.DoesPlaylistExist(data.PlaylistID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !doesExist {
		return nil, status.Error(codes.NotFound, constants.PlaylistNotFoundMessage)
	}

	isOwner, err := service.storage.IsPlaylistOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isOwner {
		return nil, status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage)
	}

	playlistInfo, err := service.storage.GetPlaylistInfo(data.PlaylistID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	playlistTracks, err := service.storage.GetPlaylistTracks(data.PlaylistID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	playlistData := &proto.PlaylistPageResponse{
		PlaylistID:   playlistInfo.PlaylistID,
		Title:        playlistInfo.Title,
		Artwork:      playlistInfo.Artwork,
		ArtworkColor: playlistInfo.ArtworkColor,
		Tracks:       playlistTracks,
	}

	return playlistData, nil
}

func contains(tracks []*proto.Track, trackID int64) bool {
	for _, currentTrack := range tracks {
		if currentTrack.ID == trackID {
			return true
		}
	}
	return false
}
