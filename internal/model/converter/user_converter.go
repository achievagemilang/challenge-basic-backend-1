package converter

import (
	"challenge-backend-1/internal/entity"
	"challenge-backend-1/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}

func UserToTokenResponse(user *entity.User, accessToken, refreshToken string) *model.LoginResponse {
	return &model.LoginResponse{
		User: &model.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func UserToEvent(user *entity.User) *model.UserEvent {
	return &model.UserEvent{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.UnixMilli(),
		UpdatedAt: user.UpdatedAt.UnixMilli(),
	}
}
