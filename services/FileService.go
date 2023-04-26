package services

import (
	"errors"
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

func (fs *FileService) UpdateFile(fileId uint, file models.File) (models.File, error) {
	result := fs.DB.Model(models.File{}).Where("id = ?", fileId).Updates(&file)

	if result.Error != nil {
		return models.File{}, result.Error
	}

	if result.RowsAffected == 0 {
		return models.File{}, errors.New("no record to update")
	}

	return file, nil
}
