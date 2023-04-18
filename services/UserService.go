package services

import (
	"final-project-hacktiv8/models"

	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (us *UserService) CreateUser(newUser models.User) (models.User, error) {
	if err := us.DB.Create(&newUser).Error; err != nil {
		return models.User{}, err
	}
	return newUser, nil
}

func (us *UserService) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := us.DB.Where("email = ?", email).First(&user).Error

	if err != nil || user.Email == "" || user.Password == "" {
		return models.User{}, err
	}

	return user, nil
}
