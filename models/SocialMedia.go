package models

import "time"

// Social Media represent the model for a Social Media
type SocialMedia struct {
	ID             uint   `gorm:"not null;primaryKey"`
	Name           string `gorm:"not null"`
	SocialMediaUrl string `gorm:"not null"`
	UserID         uint
	User           User `gorm:"foreignKey:UserID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
