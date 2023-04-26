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

func (fs *FileService) DeleteFile(name string) error {
	var file models.File

	result := fs.DB.Where("name = ?", name).Delete(&file)

	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return errors.New("no data to delete")
	}

	return nil
}
