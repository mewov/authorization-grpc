package models

type (
	DefaultResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}
	StatusData struct {
		StatusAuth    string `json:"authorization"`
		StatusSession string `json:"session"`
	}
	TokenData struct {
		AcessToken   string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	UserData struct {
		UserID int64  `json:"user_id"`
		Login  string `json:"login"`
		Email  string `json:"email"`
		Role   string `json:"role"`
	}
)
