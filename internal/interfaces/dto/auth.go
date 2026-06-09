package dto

type LoginRequest struct {
	Email    string `json:"email"    example:"agent@example.com" description:"Agent email address"`
	Password string `json:"password" example:"secret123"         description:"Agent password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" description:"Refresh token issued on login"`
}

type AuthResponse struct {
	TokenType             string `json:"token_type"               example:"Bearer"  description:"Token type"`
	AccessToken           string `json:"access_token"                               description:"JWT access token"`
	RefreshToken          string `json:"refresh_token"                              description:"JWT refresh token"`
	AccessTokenExpiresIn  int64  `json:"access_token_expires_in"  example:"900"    description:"Access token TTL in seconds"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in" example:"604800" description:"Refresh token TTL in seconds"`
	SessionID             string `json:"session_id"               format:"uuid"    description:"Active session ID"`
	AgentID               string `json:"agent_id"                 format:"uuid"    description:"Authenticated agent ID"`
}

type StatusResponse struct {
	Status string `json:"status" example:"logged_out" description:"Operation status"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials" description:"Error message"`
}
