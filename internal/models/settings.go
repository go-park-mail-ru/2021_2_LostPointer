package models

import "2021_2_LostPointer/internal/microservices/profile/proto"

type UserSettings struct {
	Email       string `json:"email,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	SmallAvatar string `json:"small_avatar,omitempty"`
	BigAvatar   string `json:"big_avatar,omitempty"`
}

func (u *UserSettings) BindProto(proto *proto.UserSettings) {
	bindedData := UserSettings{
		Email:       proto.Email,
		Nickname:    proto.Nickname,
		SmallAvatar: proto.SmallAvatar,
		BigAvatar:   proto.BigAvatar,
	}

	*u = bindedData
}
