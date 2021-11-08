package usecase

import (
	"2021_2_LostPointer/internal/microservices/music/proto"
	"2021_2_LostPointer/internal/microservices/music/repository"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MusicService struct {
	storage repository.MusicStorage
}

func NewMusicService(storage repository.MusicStorage) *MusicService {
	return &MusicService{storage: storage}
}

func (service *MusicService) RandomTracks(ctx context.Context, metadata *proto.Metadata) (*proto.Tracks, error) {
	tracks, err := service.storage.RandomTracks(metadata.Amount, metadata.IsAuthorized)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}


	return tracks, nil
}
