package entity

type User struct {
	Id       string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Type     int    `json:"type"`
	Password string `json:"password"`
}
