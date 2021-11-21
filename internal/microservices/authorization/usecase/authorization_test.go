package usecase

import (
	"2021_2_LostPointer/internal/microservices/authorization/mock"
	"2021_2_LostPointer/internal/microservices/authorization/proto"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestAuthService_Login(t *testing.T) {
	var ctrl *gomock.Controller
	mockedUsecase := mock.NewMockAuthorizationClient(ctrl)

	tests := []struct {
		name        string
		storageMock *mock.MockAuthStorage
		usecaseMock *gomock.Call
		input       *proto.AuthData
		expectedErr bool
	}{
		{
			name: "Success",
			storageMock: &mock.MockAuthStorage{
				GetUserByPasswordFunc: func(*proto.AuthData) (int64, error) {
					return 1, nil
				},
			},
			usecaseMock: mockedUsecase.EXPECT().CreateSession(
				gomock.Any(),
				gomock.Any(),
			).Return(&proto.Empty{}, nil),
			input: &proto.AuthData{Email: "lahaine@gmail.com", Password: "Avt8430066!"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewAuthService(test.storageMock)
		})
	}
}
