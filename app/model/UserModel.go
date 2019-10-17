package model

type User struct {
	Id       uint   `gorm:"primary_key" json:"user_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
