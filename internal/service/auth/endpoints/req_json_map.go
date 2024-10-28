package endpoints

// LoginRequest ???
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse ???
type LoginResponse struct {
	RefreshToken string `json:"refresh_token"`
}

// GetRefreshTokenRequest ???
type GetRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// GetRefreshTokenResponse ???
type GetRefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
}

// GetAccessTokenRequest ???
type GetAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// GetAccessTokenResponse ???
type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
