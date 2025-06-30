package model

import "time"

type Profile struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	FullAddress string    `json:"full_address"`
	PhoneNumber string    `json:"phone_number"`
	UserID      uint      `json:"-"`
}
