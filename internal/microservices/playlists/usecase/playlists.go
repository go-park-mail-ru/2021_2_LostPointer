package usecase

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/playlists"
	"2021_2_LostPointer/internal/microservices/playlists/proto"
	"2021_2_LostPointer/pkg/validation"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlaylistsService struct {
	storage playlists.PlaylistsStorage
}

func NewPlaylistsService(storage playlists.PlaylistsStorage) *PlaylistsService {
	return &PlaylistsService{storage: storage}
}

func (service *PlaylistsService) CreatePlaylist(ctx context.Context, data *proto.CreatePlaylistOptions) (*proto.CreatePlaylistResponse, error) {
	isTitleValid, msg, err := validation.ValidatePlaylistTitle(data.Title)
	if err != nil {
		return &proto.CreatePlaylistResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !isTitleValid {
		return &proto.CreatePlaylistResponse{}, status.Error(codes.InvalidArgument, msg)
	}

	if len(data.Artwork) == 0 {
		data.Artwork = constants.PlaylistArtworkDefaultFilename
		data.ArtworkColor = constants.PlaylistArtworkDefaultColor
	}
	response, err := service.storage.CreatePlaylist(data.UserID, data.Title, data.Artwork, data.ArtworkColor, data.IsPublic)
	if err != nil {
		return &proto.CreatePlaylistResponse{}, status.Error(codes.Internal, err.Error())
	}

	return response, nil
}

//nolint:cyclop
func (service *PlaylistsService) UpdatePlaylist(ctx context.Context, data *proto.UpdatePlaylistOptions) (*proto.UpdatePlaylistResponse, error) {
	doesExist, err := service.storage.DoesPlaylistExist(data.PlaylistID)
	if err != nil {
		return &proto.UpdatePlaylistResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !doesExist {
		return &proto.UpdatePlaylistResponse{}, status.Error(codes.NotFound, constants.PlaylistNotFoundMessage)
	}

	isOwner, err := service.storage.IsOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return &proto.UpdatePlaylistResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !isOwner {
		return &proto.UpdatePlaylistResponse{}, status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage)
	}

	oldArtwork, err := service.storage.GetOldPlaylistSettings(data.PlaylistID)
	if len(data.Title) != 0 {
		var (
			isTitleValid bool
			msg          string
		)
		isTitleValid, msg, err = validation.ValidatePlaylistTitle(data.Title)
		if err != nil {
			return &proto.UpdatePlaylistResponse{}, status.Error(codes.Internal, err.Error())
		}
		if !isTitleValid {
			return &proto.UpdatePlaylistResponse{}, status.Error(codes.InvalidArgument, msg)
		}
		if err = service.storage.UpdatePlaylistTitle(data.PlaylistID, data.Title); err != nil {
			return &proto.UpdatePlaylistResponse{}, status.Error(codes.Internal, err.Error())
		}
	}

	if err = service.storage.UpdatePlaylistAccess(data.PlaylistID, data.IsPublic); err != nil {
		return &proto.UpdatePlaylistResponse{}, status.Error(codes.Internal, err.Error())
	}
	response := &proto.UpdatePlaylistResponse{
		ArtworkColor: "",
	}
	if len(data.Artwork) != 0 {
		response.OldArtworkFilename = oldArtwork
		response.ArtworkColor = data.ArtworkColor
		if err = service.storage.UpdatePlaylistArtwork(data.PlaylistID, data.Artwork, data.ArtworkColor); err != nil {
			return &proto.UpdatePlaylistResponse{}, status.Error(codes.Internal, err.Error())
		}
	}

	return response, nil
}

//nolint:dupl
func (service *PlaylistsService) DeletePlaylist(ctx context.Context, data *proto.DeletePlaylistOptions) (*proto.DeletePlaylistResponse, error) {
	doesExist, err := service.storage.DoesPlaylistExist(data.PlaylistID)
	if err != nil {
		return &proto.DeletePlaylistResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !doesExist {
		return &proto.DeletePlaylistResponse{}, status.Error(codes.NotFound, constants.PlaylistNotFoundMessage)
	}

	isOwner, err := service.storage.IsOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return &proto.DeletePlaylistResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !isOwner {
		return &proto.DeletePlaylistResponse{}, status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage)
	}

	oldArtwork, err := service.storage.GetOldPlaylistSettings(data.PlaylistID)
	if err != nil {
		return &proto.DeletePlaylistResponse{}, status.Error(codes.Internal, err.Error())
	}
	err = service.storage.DeletePlaylist(data.PlaylistID)
	if err != nil {
		return &proto.DeletePlaylistResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.DeletePlaylistResponse{OldArtworkFilename: oldArtwork}, nil
}

func (service *PlaylistsService) AddTrack(ctx context.Context, data *proto.AddTrackOptions) (*proto.AddTrackResponse, error) {
	doesExist, err := service.storage.DoesPlaylistExist(data.PlaylistID)
	if err != nil {
		return &proto.AddTrackResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !doesExist {
		return &proto.AddTrackResponse{}, status.Error(codes.NotFound, constants.PlaylistNotFoundMessage)
	}

	isOwner, err := service.storage.IsOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return &proto.AddTrackResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !isOwner {
		return &proto.AddTrackResponse{}, status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage)
	}

	isAdded, err := service.storage.IsAdded(data.PlaylistID, data.TrackID)
	if err != nil {
		return &proto.AddTrackResponse{}, status.Error(codes.Internal, err.Error())
	}
	if isAdded {
		return &proto.AddTrackResponse{}, status.Error(codes.InvalidArgument, constants.TrackAlreadyInPlaylistMessage)
	}

	err = service.storage.AddTrack(data.PlaylistID, data.TrackID)
	if err != nil {
		return &proto.AddTrackResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.AddTrackResponse{}, nil
}

func (service *PlaylistsService) DeleteTrack(ctx context.Context, data *proto.DeleteTrackOptions) (*proto.DeleteTrackResponse, error) {
	doesExist, err := service.storage.DoesPlaylistExist(data.PlaylistID)
	if err != nil {
		return &proto.DeleteTrackResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !doesExist {
		return &proto.DeleteTrackResponse{}, status.Error(codes.NotFound, constants.PlaylistNotFoundMessage)
	}

	isOwner, err := service.storage.IsOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return &proto.DeleteTrackResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !isOwner {
		return &proto.DeleteTrackResponse{}, status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage)
	}

	err = service.storage.DeleteTrack(data.PlaylistID, data.TrackID)
	if err != nil {
		return &proto.DeleteTrackResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.DeleteTrackResponse{}, nil
}

//nolint:dupl
func (service *PlaylistsService) DeletePlaylistArtwork(ctx context.Context, data *proto.DeletePlaylistArtworkOptions) (*proto.DeletePlaylistArtworkResponse, error) {
	doesExist, err := service.storage.DoesPlaylistExist(data.PlaylistID)
	if err != nil {
		return &proto.DeletePlaylistArtworkResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !doesExist {
		return &proto.DeletePlaylistArtworkResponse{}, status.Error(codes.NotFound, constants.PlaylistNotFoundMessage)
	}

	isOwner, err := service.storage.IsOwner(data.PlaylistID, data.UserID)
	if err != nil {
		return &proto.DeletePlaylistArtworkResponse{}, status.Error(codes.Internal, err.Error())
	}
	if !isOwner {
		return &proto.DeletePlaylistArtworkResponse{}, status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage)
	}

	oldArtwork, err := service.storage.GetOldPlaylistSettings(data.PlaylistID)
	if err != nil {
		return &proto.DeletePlaylistArtworkResponse{}, status.Error(codes.Internal, err.Error())
	}
	err = service.storage.DeletePlaylistArtwork(data.PlaylistID)
	if err != nil {
		return &proto.DeletePlaylistArtworkResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.DeletePlaylistArtworkResponse{OldArtworkFilename: oldArtwork}, nil
}
