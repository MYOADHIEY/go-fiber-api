package models

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
}

var (
	UserOrderBy = []string{"def.id", "def.created_at", "def.updated_at"}

	UserByString = []string{
		"def.id",
	}

	UserSelectStatement = ` select def.id,def._name,def._email, def._phone, def.address 
	from _user def `

	UserWhereSelectStatement = ` where def.deleted_at is null `
)
