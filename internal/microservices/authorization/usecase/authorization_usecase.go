package usecase

import (
	session "2021_2_LostPointer/internal/microservices/authorization/delivery"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/sessions"
	"2021_2_LostPointer/internal/users"
	"2021_2_LostPointer/internal/utils/constants"
	"2021_2_LostPointer/internal/utils/validation"
	"context"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthorizationUseCase struct {
	userDB     users.UserRepository
	sessionsDB sessions.SessionRepository
}

func NewAuthorizationUseCase(userDB users.UserRepository, sessionsDB sessions.SessionRepository) AuthorizationUseCase {
	return AuthorizationUseCase{userDB: userDB, sessionsDB: sessionsDB}
}

func (authU AuthorizationUseCase) GetUserBySession(ctx context.Context, auth *session.SessionData) (*session.UserID, error) {
	userID, err := authU.sessionsDB.GetUserIdByCookie(auth.Cookies)
	if userID == 0 || err != nil {
		userID = -1
	}

	return &session.UserID{
		UserID: int32(userID),
	}, nil
}

func (authU AuthorizationUseCase) DeleteSession(ctx context.Context, auth *session.SessionData) (*session.Empty, error) {
	err := authU.sessionsDB.DeleteSession(auth.Cookies)

	return &session.Empty{}, err
}

func (authU AuthorizationUseCase) SignIn(ctx context.Context, auth *session.Auth) (*session.SessionData, error) {
	authData := &models.Auth{
		Email:    auth.Login,
		Password: auth.Password,
	}
	userID, err := authU.userDB.DoesUserExist(authData)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if userID == 0 {
		return nil, status.Error(codes.Aborted, constants.WrongCredentials)
	}

	cookieValue := uuid.NewV4()
	currentSession := &session.SessionData{
		Cookies: cookieValue.String(),
	}

	err = authU.sessionsDB.CreateSession(userID, currentSession.Cookies)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return currentSession, nil
}

func (authU AuthorizationUseCase) Signup(ctx context.Context, register *session.SignUpData) (*session.SessionData, error) {
	isEmailUnique, err := authU.userDB.IsEmailUnique(register.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isEmailUnique {
		return nil, status.Error(codes.Aborted, constants.NotUniqueEmailMessage)
	}

	isNicknameUnique, err := authU.userDB.IsNicknameUnique(register.Nickname)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isNicknameUnique {
		return nil, status.Error(codes.Aborted, constants.NotUniqueNicknameMessage)
	}

	userData := &models.User{
		Email:    register.Email,
		Password: register.Password,
		Nickname: register.Nickname,
	}
	isValidCredentials, message, err := validation.ValidateRegisterCredentials(userData)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !isValidCredentials {
		return nil, status.Error(codes.Aborted, message)
	}

	userID, err := authU.userDB.CreateUser(userData)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cookieValue := uuid.NewV4()
	currentSession := &session.SessionData{
		Cookies: cookieValue.String(),
	}

	err = authU.sessionsDB.CreateSession(userID, currentSession.Cookies)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return currentSession, nil
}
