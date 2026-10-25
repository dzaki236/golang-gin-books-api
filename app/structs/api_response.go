package structs

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
type AuthResponse struct {
	User        User   `json:"user"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}
