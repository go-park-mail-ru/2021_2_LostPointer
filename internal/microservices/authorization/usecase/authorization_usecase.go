package usecase

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/authorization/proto"
	"2021_2_LostPointer/internal/microservices/authorization/repository"
	"context"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	storage repository.AuthStorage
}

func NewAuthService(storage repository.AuthStorage) AuthService {
	return AuthService{storage: storage}
}

func (service AuthService) CreateSession(ctx context.Context, data *proto.SessionData) (*proto.Empty, error) {
	err := service.storage.CreateSession(data.ID, data.Cookies)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (service AuthService) Login(ctx context.Context, authData *proto.AuthData) (*proto.Cookie, error) {
	userID, err := service.storage.GetUserByPassword(authData.Login, authData.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if userID == 0 {
		return nil, status.Error(codes.Aborted, constants.WrongCredentials)
	}

	cookieValue := uuid.NewV4()
	sessionData := &proto.Cookie{
		Cookies: cookieValue.String(),
	}
	_, err = service.CreateSession(context.Background(), &proto.SessionData{ID: userID, Cookies: cookieValue.String()})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return sessionData, nil
}
