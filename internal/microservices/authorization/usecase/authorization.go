package usecase

import (
	"2021_2_LostPointer/internal/constants"
	customErrors "2021_2_LostPointer/internal/errors"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/pkg/validation"
	"context"
	"errors"
	"os"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/microservices/authorization/proto"
	"2021_2_LostPointer/internal/microservices/authorization/repository"
)

type AuthService struct {
	storage repository.AuthStorage
}

func NewAuthService(storage repository.AuthStorage) *AuthService {
	return &AuthService{storage: storage}
}

func (service *AuthService) CreateSession(ctx context.Context, data *proto.SessionData) (*proto.Empty, error) {
	err := service.storage.CreateSession(data.ID, data.Cookies)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (service *AuthService) GetUserByCookie(ctx context.Context, cookie *proto.Cookie) (*proto.UserID, error) {
	id, err := service.storage.GetUserByCookie(cookie.Cookies)
	if err != nil {
		return nil, err
	}

	return &proto.UserID{ID: id}, nil
}

func (service *AuthService) DeleteSession(ctx context.Context, cookie *proto.Cookie) (*proto.Empty, error) {
	err := service.storage.DeleteSession(cookie.Cookies)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}

func (service *AuthService) Login(ctx context.Context, authData *proto.AuthData) (*proto.Cookie, error) {
	userID, err := service.storage.GetUserByPassword(
		&models.AuthData{Email: authData.Login, Password: authData.Password})
	if err != nil {
		if errors.Is(err, customErrors.ErrWrongCredentials) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
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

func (service *AuthService) Register(ctx context.Context, registerData *proto.RegisterData) (*proto.Cookie, error) {
	isEmailUnique, err := service.storage.IsEmailUnique(registerData.Login)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isEmailUnique {
		return nil, status.Error(codes.InvalidArgument, constants.NotUniqueEmailMessage)
	}

	isNicknameUnique, err := service.storage.IsNicknameUnique(registerData.Nickname)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isNicknameUnique {
		return nil, status.Error(codes.InvalidArgument, constants.NotUniqueNicknameMessage)
	}

	registerCredentials := &models.RegisterData{
		Email:    registerData.Login,
		Password: registerData.Password,
		Nickname: registerData.Nickname,
	}
	isValidCredentials, message, err := validation.ValidateRegisterCredentials(registerCredentials)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isValidCredentials {
		return nil, status.Error(codes.InvalidArgument, message)
	}

	userID, err := service.storage.CreateUser(registerCredentials)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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

func (service *AuthService) GetAvatar(ctx context.Context, user *proto.UserID) (*proto.Avatar, error) {
	avatar, err := service.storage.GetAvatar(user.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.Avatar{Filename: os.Getenv("USERS_ROOT_PREFIX") + avatar + constants.LittleAvatarPostfix}, nil
}

func (service *AuthService) Logout(ctx context.Context, cookies *proto.Cookie) (*proto.Empty, error) {
	_, err := service.DeleteSession(ctx, cookies)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}
