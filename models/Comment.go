package models

import "time"

type Comment struct {
	ID        uint      `gorm:"not null;primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	PhotoID   uint      `json:"photo_id"`
	Message   string    `gorm:"not null" json:"message"`
	User      User      `gorm:"foreignKey:UserID"`
	Photo     Photo     `gorm:"foreignKey:PhotoID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
