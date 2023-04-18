package models

import (
	"errors"
	"time"

	"net/mail"

	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"not null;primaryKey"`
	Email     string `gorm:"not null;uniqueIndex;type:varchar(100)"`
	Username  string `gorm:"not null;uniqueIndex;type:varchar(100)"`
	Password  string `gorm:"not null;type:varchar(100)"`
	Age       uint   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var user User
	err := tx.Where("email", u.Email).First(&user).Error

	if err == nil {
		return errors.New("already registered")
	}

	_, err = mail.ParseAddress(u.Email)

	if err != nil {
		return errors.New("email invalid format")
	}

	if len(u.Password) < 6 {
		return errors.New("please input password longer")
	}

	if u.Age < 8 {
		return errors.New("you are not old enough to register")
	}

	return nil
}
