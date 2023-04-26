package models

import "time"

// Photo represent the model for a photo
type Photo struct {
	ID        uint      `gorm:"not null;primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Caption   string    `json:"caption"`
	PhotoUrl  string    `gorm:"not null" json:"photo_url"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
