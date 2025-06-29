package model

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	Name        string `json:"name"`
	FullAddress string `json:"full_address"`
	PhoneNumber string `json:"phone_number"`
	UserID      uint   `json:"-"`
}
