package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       string `gorm:"unique" json:"email"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
	Profile     Profile
}
