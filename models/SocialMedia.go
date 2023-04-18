package models

import "time"

type SocialMedia struct {
	ID             uint   `gorm:"not null;primaryKey"`
	Name           string `gorm:"not null"`
	SocialMediaUrl string `gorm:"not null"`
	UserId         uint
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
