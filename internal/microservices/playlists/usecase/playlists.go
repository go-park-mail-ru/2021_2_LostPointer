package usecase

import (
	"2021_2_LostPointer/internal/constants"
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/microservices/playlists/proto"
	"2021_2_LostPointer/internal/microservices/playlists/repository"
	"2021_2_LostPointer/pkg/validation"
)

type PlaylistsService struct {
	storage repository.PlaylistsStorage
}

func NewPlaylistsService(storage repository.PlaylistsStorage) *PlaylistsService {
	return &PlaylistsService{storage: storage}
}

func (service *PlaylistsService) CreatePlaylist(ctx context.Context, data *proto.CreatePlaylistOptions) (*proto.CreatePlaylistResponse, error) {
	isTitleValid, msg, err := validation.ValidatePlaylistTitle(data.Title)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isTitleValid {
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	if len(data.Artwork) == 0 {
		data.Artwork = "default_playlist_artwork"
	}
	response, err := service.storage.CreatePlaylist(data.UserID, data.Title, data.Artwork, data.ArtworkColor)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	log.Println("usecase:", response)

	return response, nil
}

func (service *PlaylistsService) UpdatePlaylist(ctx context.Context, data *proto.UpdatePlaylistOptions) (*proto.UpdatePlaylistResponse, error) {
	isOwner, err := service.storage.IsOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isOwner {
		return nil, status.Error(codes.InvalidArgument, constants.NotPlaylistOwnerMessage)
	}

	isTitleValid, msg, err := validation.ValidatePlaylistTitle(data.Title)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isTitleValid {
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	response :=  &proto.UpdatePlaylistResponse{}
	oldArtwork, oldArtworkColor, err := service.storage.GetOldArtwork(data.PlaylistID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(data.Artwork) == 0 {
		data.Artwork = oldArtwork
		data.ArtworkColor = oldArtworkColor
		response.OldArtworkFilename = ""
	} else {
		response.OldArtworkFilename = oldArtwork
	}
	err = service.storage.UpdatePlaylist(data.PlaylistID, data.Title, data.Artwork, data.ArtworkColor)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

func (service *PlaylistsService) DeletePlaylist(ctx context.Context, data *proto.DeletePlaylistOptions) (*proto.DeletePlaylistResponse, error) {
	isOwner, err := service.storage.IsOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isOwner {
		return nil, status.Error(codes.InvalidArgument, constants.NotPlaylistOwnerMessage)
	}

	oldArtwork, _, err := service.storage.GetOldArtwork(data.PlaylistID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = service.storage.DeletePlaylist(data.PlaylistID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.DeletePlaylistResponse{OldArtworkFilename: oldArtwork}, nil
}

func (service *PlaylistsService) AddTrack(ctx context.Context, data *proto.AddTrackOptions) (*proto.AddTrackResponse, error) {
	isOwner, err := service.storage.IsOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isOwner {
		return nil, status.Error(codes.InvalidArgument, constants.NotPlaylistOwnerMessage)
	}

	isAdded, err := service.storage.IsAdded(data.PlaylistID, data.TrackID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if isAdded {
		return nil, status.Error(codes.InvalidArgument, constants.TrackAlreadyInPlaylistMessage)
	}

	err = service.storage.AddTrack(data.PlaylistID, data.TrackID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.AddTrackResponse{}, nil
}

func (service *PlaylistsService) DeleteTrack(ctx context.Context, data *proto.DeleteTrackOptions) (*proto.DeleteTrackResponse, error) {
	isOwner, err := service.storage.IsOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isOwner {
		return nil, status.Error(codes.InvalidArgument, constants.NotPlaylistOwnerMessage)
	}

	err = service.storage.DeleteTrack(data.PlaylistID, data.TrackID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.DeleteTrackResponse{}, nil
}