package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       string  `gorm:"unique" json:"email"`
	Password    string  `json:"password"`
	Activation  string  `json:"activation"`
	AccessToken string  `json:"access_token"`
	Profile     Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
