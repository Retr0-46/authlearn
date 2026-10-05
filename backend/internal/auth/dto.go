package auth

import "time"

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}
