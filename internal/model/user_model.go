package model

type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginResponse struct {
	User         *UserResponse `json:"user,omitempty"`
	AccessToken  string        `json:"access_token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}
