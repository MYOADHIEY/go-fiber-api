package models

type Registration struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"_password"`
	Address  string `json:"address"`
}

type RegistrationParameters struct {
	ID           string `json:"id"`
	By           string `json:"by"`
	Search       string `json:"serch"`
	Page         int    `json:"page"`
	Offset       int    `json:"offset"`
	Limit        int    `json:"limit"`
	Sort         string `json:"sort"`
	ResponseCode int    `json:"response_code"`
}
