package usecase

import (
	session "2021_2_LostPointer/internal/microservices/authorization/delivery"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/sessions"
	"2021_2_LostPointer/internal/users"
	"2021_2_LostPointer/internal/utils/validation"
	"context"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const WrongCredentials = "Wrong credentials"
const NotUniqueEmail = "Email is not unique"
const NotUniqueNickname = "Nickname is not unique"

type AuthorizationUseCase struct {
	userDB users.UserRepository
	sessionsDB sessions.SessionRepository
}

func NewAuthorizationUseCase(userDB users.UserRepository, sessionsDB sessions.SessionRepository) AuthorizationUseCase {
	return AuthorizationUseCase{userDB: userDB, sessionsDB: sessionsDB}
}

func (authU AuthorizationUseCase) CheckSession(ctx context.Context, auth *session.SessionData) (*session.UserID, error) {
	userID, err := authU.sessionsDB.GetUserIdByCookie(auth.Cookies)
	if userID == 0 || err != nil {
		userID = -1
	}

	return &session.UserID{
		UserID: int32(userID),
	}, nil
}

func (authU AuthorizationUseCase) DeleteSession(ctx context.Context, auth *session.SessionData) (*session.UserID, error) {
	err := authU.sessionsDB.DeleteSession(auth.Cookies)

	return &session.UserID{
		UserID: -1,
	}, err
}

func (authU AuthorizationUseCase) SignIn(ctx context.Context, auth *session.Auth) (*session.SessionData, error) {
	authData := &models.Auth{
		Email: auth.Login,
		Password: auth.Password,
	}
	userID, err := authU.userDB.DoesUserExist(authData)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	if userID == 0 {
		return nil, status.Error(codes.Aborted, WrongCredentials)
	}

	cookieValue := uuid.NewV4()
	currentSession := &session.SessionData{
		Cookies: cookieValue.String(),
	}

	err = authU.sessionsDB.CreateSession(userID, currentSession.Cookies)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	return currentSession, nil
}

func (authU AuthorizationUseCase) Signup(ctx context.Context, register *session.SignUpData) (*session.SessionData, error) {
	isEmailUnique, err := authU.userDB.IsEmailUnique(register.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	if !isEmailUnique {
		return nil, status.Error(codes.Aborted, NotUniqueEmail)
	}

	isNicknameUnique, err := authU.userDB.IsNicknameUnique(register.Nickname)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	if !isNicknameUnique {
		return nil, status.Error(codes.Aborted, NotUniqueNickname)
	}

	userData := &models.User{
		Email: register.Email,
		Password: register.Password,
		Nickname: register.Nickname,
	}
	isValidCredentials, message, err := validation.ValidateRegisterCredentials(userData)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}
	if !isValidCredentials {
		return nil, status.Error(codes.Aborted, message)
	}

	userID, err := authU.userDB.CreateUser(userData)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	cookieValue := uuid.NewV4()
	currentSession := &session.SessionData{
		Cookies: cookieValue.String(),
	}

	err = authU.sessionsDB.CreateSession(userID, currentSession.Cookies)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	return currentSession, nil
}
