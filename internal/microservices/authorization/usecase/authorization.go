//nolint:contextcheck
package usecase

import (
	"2021_2_LostPointer/internal/constants"
	customErrors "2021_2_LostPointer/internal/errors"
	"2021_2_LostPointer/internal/microservices/authorization"
	"2021_2_LostPointer/pkg/validation"
	"context"
	"errors"
	"os"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/microservices/authorization/proto"
)

type AuthService struct {
	storage authorization.AuthStorage
}

func NewAuthService(storage authorization.AuthStorage) *AuthService {
	return &AuthService{storage: storage}
}

func (service *AuthService) CreateSession(ctx context.Context, data *proto.SessionData) (*proto.Empty, error) {
	err := service.storage.CreateSession(data.ID, data.Cookies)
	if err != nil {
		return &proto.Empty{}, err
	}
	return &proto.Empty{}, nil
}

func (service *AuthService) GetUserByCookie(ctx context.Context, cookie *proto.Cookie) (*proto.UserID, error) {
	id, err := service.storage.GetUserByCookie(cookie.Cookies)
	if err != nil {
		return &proto.UserID{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.UserID{ID: id}, nil
}

func (service *AuthService) DeleteSession(ctx context.Context, cookie *proto.Cookie) (*proto.Empty, error) {
	err := service.storage.DeleteSession(cookie.Cookies)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}

func (service *AuthService) Login(ctx context.Context, authData *proto.AuthData) (*proto.Cookie, error) {
	userID, err := service.storage.GetUserByPassword(authData)
	if err != nil {
		if errors.Is(err, customErrors.ErrWrongCredentials) {
			return &proto.Cookie{}, status.Error(codes.InvalidArgument, err.Error())
		}
		return &proto.Cookie{}, status.Error(codes.Internal, err.Error())
	}

	cookieValue := uuid.NewV4()
	sessionData := &proto.Cookie{
		Cookies: cookieValue.String(),
	}
	_, err = service.CreateSession(context.Background(), &proto.SessionData{ID: userID, Cookies: cookieValue.String()})
	if err != nil {
		return &proto.Cookie{}, status.Error(codes.Internal, err.Error())
	}

	return sessionData, nil
}

func (service *AuthService) Register(ctx context.Context, registerData *proto.RegisterData) (*proto.Cookie, error) {
	isEmailUnique, err := service.storage.IsEmailUnique(registerData.Email)
	if err != nil {
		return &proto.Cookie{}, status.Error(codes.Internal, err.Error())
	}
	if !isEmailUnique {
		return &proto.Cookie{}, status.Error(codes.InvalidArgument, constants.EmailNotUniqueMessage)
	}

	isNicknameUnique, err := service.storage.IsNicknameUnique(registerData.Nickname)
	if err != nil {
		return &proto.Cookie{}, status.Error(codes.Internal, err.Error())
	}
	if !isNicknameUnique {
		return &proto.Cookie{}, status.Error(codes.InvalidArgument, constants.NicknameNotUniqueMessage)
	}

	isValidCredentials, message, err := validation.ValidateRegisterCredentials(registerData.Email, registerData.Password, registerData.Nickname)
	if err != nil {
		return &proto.Cookie{}, status.Error(codes.Internal, err.Error())
	}
	if !isValidCredentials {
		return &proto.Cookie{}, status.Error(codes.InvalidArgument, message)
	}

	userID, err := service.storage.CreateUser(registerData)
	if err != nil {
		return &proto.Cookie{}, status.Error(codes.Internal, err.Error())
	}

	cookieValue := uuid.NewV4()
	sessionData := &proto.Cookie{
		Cookies: cookieValue.String(),
	}
	_, err = service.CreateSession(context.Background(), &proto.SessionData{ID: userID, Cookies: cookieValue.String()})
	if err != nil {
		return &proto.Cookie{}, status.Error(codes.Internal, err.Error())
	}

	return sessionData, nil
}

func (service *AuthService) GetAvatar(ctx context.Context, user *proto.UserID) (*proto.Avatar, error) {
	avatar, err := service.storage.GetAvatar(user.ID)
	if err != nil {
		return &proto.Avatar{}, status.Error(codes.Internal, err.Error())
	}

	return &proto.Avatar{Filename: os.Getenv("USERS_ROOT_PREFIX") + avatar + constants.UserAvatarExtension150px}, nil
}

func (service *AuthService) Logout(ctx context.Context, cookies *proto.Cookie) (*proto.Empty, error) {
	_, err := service.DeleteSession(ctx, cookies)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &proto.Empty{}, nil
}
