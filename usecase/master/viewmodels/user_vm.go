package viewmodels

type UserVM struct {
	ID        string  `json:"id"`
	Email     string  `json:"_email"`
	Name      string  `json:"_name"`
	Address   string  `json:"address"`
	Phone     string  `json:"_phone"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
	DeletedAt *string `json:"deleted_at"`
	CreatedBy int     `json:"created_by"`
	UpdatedBy *int    `json:"updated_by"`
	DeletedBy *int    `json:"deleted_by"`
}
