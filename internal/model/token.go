package model

type AuthTokens struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type ValidateTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id,omitempty"`
}
