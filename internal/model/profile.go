package model

import "gorm.io/gorm"

type Profile struct {
	gorm.Model
	Name        string `json:"name"`
	FullAddress string `json:"full_address"`
	PhoneNumber string `json:"phone_number"`
	SchoolName  string `json:"school_name"`
	UserID      uint   `json:"-"`
}
