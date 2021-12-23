package models

type RegisterRequest struct {
	Email     string `json:"email"`
	PublicKey string `json:"public_key"`
	Name      string `json:"name"`
}

type RegisterResponse struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}
