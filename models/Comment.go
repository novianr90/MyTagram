package models

import "time"

type Comment struct {
	ID        uint `gorm:"not null;primaryKey"`
	UserID    uint
	PhotoID   uint
	Message   string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
