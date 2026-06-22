package viewmodels

type JWTVM struct {
	Token              string `json:"token"`
	ExpiredDate        string `json:"expired_date"`
	RefreshToken       string `json:"refresh_token"`
	RefreshExpiredDate string `json:"refresh_expired_date"`
	LatestAction       string `json:"latest_action"`
	UserID             string `json:"user_id"`
}
