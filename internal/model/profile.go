package model

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	Name        string `json:"name"`
	Email       string `gorm:"unique" json:"email"`
	PhoneNumber string `json:"phone_number"`
	SchoolName  string `json:"school_name"`
	Password    string `json:"password"`
	UserID      uint   `json:"-"`
}
