package usecase

import (
	"2021_2_LostPointer/internal/constants"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	artistData.Tracks, err = service.storage.GetArtistTracks(metadata.ArtistID, metadata.IsAuthorized, constants.TracksDefaultAmountForArtist)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	artistData.Albums, err = service.storage.GetArtistAlbums(metadata.ArtistID, constants.AlbumsDefaultAmountForArtist)
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