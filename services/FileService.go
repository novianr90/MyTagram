package services

import (
	"final-project-hacktiv8/models"

	"gorm.io/gorm"
)

type FileService struct {
	DB *gorm.DB
}

func (fs *FileService) SaveUploadedFile(file models.File) (models.File, error) {
	err := fs.DB.Create(&file).Error

	if err != nil {
		return models.File{}, err
	}

	return file, nil
}

func (fs *FileService) GetUploadedFile(name string) (models.File, error) {
	var file models.File
	err := fs.DB.Where("name = ?", name).First(&file).Error

	if err != nil {
		return models.File{}, err
	}

	return file, nil
}
