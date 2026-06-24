package useraccount

type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"_password"`
}
