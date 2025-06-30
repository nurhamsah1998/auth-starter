package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `gorm:"unique" json:"email"`
	Password     string    `json:"password"`
	Activation   string    `json:"activation"`
	RefreshToken string    `json:"refresh_token"`
	Profile      Profile   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
