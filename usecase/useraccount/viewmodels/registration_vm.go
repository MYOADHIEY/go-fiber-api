package viewmodels

type RegistrationtVM struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"_password"`
	Address  string `json:"address"`
}
