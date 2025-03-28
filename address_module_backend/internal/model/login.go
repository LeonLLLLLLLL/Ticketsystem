package model

type LoginRequest struct {
	Identifier string `json:"identifier"` // Email or Username
	Password   string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
