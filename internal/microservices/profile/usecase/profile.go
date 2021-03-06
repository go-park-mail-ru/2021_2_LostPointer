package usecase

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/constants"
	customErrors "2021_2_LostPointer/internal/errors"
	"2021_2_LostPointer/internal/microservices/profile"
	"2021_2_LostPointer/internal/microservices/profile/proto"
	"2021_2_LostPointer/pkg/validation"
)

type ProfileService struct {
	storage profile.UserSettingsStorage
}

func NewProfileService(storage profile.UserSettingsStorage) *ProfileService {
	return &ProfileService{storage: storage}
}

func (service *ProfileService) GetSettings(ctx context.Context, user *proto.GetSettingsOptions) (*proto.UserSettings, error) {
	settings, err := service.storage.GetSettings(user.ID)
	if err != nil {
		if errors.Is(err, customErrors.ErrUserNotFound) {
			return &proto.UserSettings{}, status.Error(codes.NotFound, err.Error())
		}
		return &proto.UserSettings{}, status.Error(codes.Internal, err.Error())
	}

	return settings, nil
}

//nolint:cyclop
func (service *ProfileService) UpdateSettings(ctx context.Context, settings *proto.UpdateSettingsOptions) (*proto.EmptyProfile, error) {
	if strings.ToLower(settings.Email) != settings.OldSettings.Email && len(settings.Email) != 0 {
		isEmailValid := govalidator.IsEmail(settings.Email)
		if !isEmailValid {
			return &proto.EmptyProfile{}, status.Error(codes.InvalidArgument, constants.EmailInvalidSyntaxMessage)
		}

		isEmailUnique, err := service.storage.IsEmailUnique(settings.Email)
		if err != nil {
			return &proto.EmptyProfile{}, status.Error(codes.Internal, err.Error())
		}
		if !isEmailUnique {
			return &proto.EmptyProfile{}, status.Error(codes.InvalidArgument, constants.EmailNotUniqueMessage)
		}

		err = service.storage.UpdateEmail(settings.UserID, settings.Email)
		if err != nil {
			return &proto.EmptyProfile{}, status.Error(codes.Internal, err.Error())
		}
	}

	if settings.Nickname != settings.OldSettings.Nickname && len(settings.Nickname) != 0 {
		isNicknameValid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{`+constants.MinNicknameLength+`,`+constants.MaxNicknameLength+`}$`, settings.Nickname)
		if err != nil {
			return &proto.EmptyProfile{}, status.Error(codes.Internal, err.Error())
		}
		if !isNicknameValid {
			return &proto.EmptyProfile{}, status.Error(codes.InvalidArgument, constants.NicknameInvalidSyntaxMessage)
		}

		isNicknameUnique, err := service.storage.IsNicknameUnique(settings.Nickname)
		if err != nil {
			return &proto.EmptyProfile{}, status.Error(codes.Internal, err.Error())
		}
		if !isNicknameUnique {
			return &proto.EmptyProfile{}, status.Error(codes.InvalidArgument, constants.NicknameNotUniqueMessage)
		}

		err = service.storage.UpdateNickname(settings.UserID, settings.Nickname)
		if err != nil {
			return &proto.EmptyProfile{}, status.Error(codes.Internal, err.Error())
		}
	}

	switch isEmpty := len(settings.OldPassword) == 0; isEmpty {
	case true:
		if len(settings.NewPassword) != 0 {
			return &proto.EmptyProfile{}, status.Error(codes.InvalidArgument, constants.OldPasswordFieldIsEmptyMessage)
		}
	default:
		if len(settings.NewPassword) == 0 {
			return &proto.EmptyProfile{}, status.Error(codes.InvalidArgument, constants.NewPasswordFieldIsEmptyMessage)
		}
		isOldPasswordCorrect, err := service.storage.CheckPasswordByUserID(settings.UserID, settings.OldPassword)
		if err != nil {
			if errors.Is(err, customErrors.ErrWrongCredentials) {
				return &proto.EmptyProfile{}, status.Error(codes.InvalidArgument, constants.WrongPasswordMessage)
			}
			return &proto.EmptyProfile{}, status.Error(codes.Internal, err.Error())
		}
		if !isOldPasswordCorrect {
			return &proto.EmptyProfile{}, status.Error(codes.InvalidArgument, constants.WrongPasswordMessage)
		}

		isNewPasswordValid, msg, err := validation.ValidatePassword(settings.NewPassword)
		if err != nil {
			return &proto.EmptyProfile{}, status.Error(codes.Internal, err.Error())
		}
		if !isNewPasswordValid {
			return &proto.EmptyProfile{}, status.Error(codes.InvalidArgument, msg)
		}

		err = service.storage.UpdatePassword(settings.UserID, settings.NewPassword)
		if err != nil {
			return &proto.EmptyProfile{}, status.Error(codes.Internal, err.Error())
		}
	}

	if len(settings.AvatarFilename) != 0 {
		err := service.storage.UpdateAvatar(settings.UserID, settings.AvatarFilename)
		if err != nil {
			return &proto.EmptyProfile{}, status.Error(codes.Internal, err.Error())
		}
	}

	return &proto.EmptyProfile{}, nil
}
