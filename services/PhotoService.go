package services

import (
	"errors"
	"final-project-hacktiv8/models"

	"gorm.io/gorm"
)

type PhotoService struct {
	DB *gorm.DB
}

func (ps *PhotoService) GetAll(userId uint) ([]models.Photo, error) {
	var photos []models.Photo
	if err := ps.DB.Where("user_id = ?", userId).Find(&photos).Error; err != nil {
		return []models.Photo{}, err
	}

	return photos, nil
}

func (ps *PhotoService) GetPhotoById(photoId uint) (models.Photo, error) {
	var photo models.Photo

	if err := ps.DB.Where("id = ?", photoId).First(&photo).Error; err != nil {
		return models.Photo{}, err
	}

	return photo, nil
}

func (ps *PhotoService) CreatePhoto(photo models.Photo) (models.Photo, error) {
	if err := ps.DB.Create(&photo).Error; err != nil {
		return models.Photo{}, err
	}

	return photo, nil
}

func (ps *PhotoService) UpdatePhoto(photoId uint, updatePhoto models.Photo) error {

	result := ps.DB.Model(models.Photo{}).Where("id = ?", photoId).Updates(&updatePhoto)

	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return errors.New("no data to updated")
	}

	return nil
}

func (ps *PhotoService) DeletePhotoById(photoId uint) error {

	var photo models.Photo

	result := ps.DB.Where("id = ?", photoId).Delete(&photo)

	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return errors.New("no data to delete")
	}

	return nil
}
