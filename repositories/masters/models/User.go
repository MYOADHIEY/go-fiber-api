package models

type User struct {
	ID        string `json:"id"`
	Email     string `json:"_email"`
	Name      string `json:"_name"`
	Address   string `json:"address"`
	Phone     string `json:"_phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	CreatedBy int    `json:"created_by"`
	UpdatedBy int    `json:"updated_by"`
	DeletedBy int    `json:"deleted_by"`
}

type UserParameters struct {
	ID           string `json:"id"`
	By           string `json:"by"`
	Search       string `json:"serch"`
	Page         int    `json:"page"`
	Offset       int    `json:"offset"`
	Limit        int    `json:"limit"`
	Sort         string `json:"sort"`
	ResponseCode int    `json:"response_code"`
	Type         string `json:"_type"`
}

var (
	UserOrderBy = []string{"def.id", "def.created_at", "def.updated_at"}

	UserByString = []string{
		"def.id",
	}

	UserSelectStatement = ` select def.id,def._name,def._email, def._phone, def.address
	, def.created_at, def.updated_at, def.deleted_at, def.created_by, def.updated_by, def.deleted_by
	from _user def `

	UserWhereSelectStatement = ` where def.deleted_at is null `
)
