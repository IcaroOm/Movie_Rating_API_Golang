package handlers

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"success message"`
}

type LoginRequest struct {
	Username string `json:"username" example:"user123"`
	Password string `json:"password" example:"password123"`
}

type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}