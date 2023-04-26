package services

import (
	"final-project-hacktiv8/models"

	"gorm.io/gorm"
)

type CommentService struct {
	DB *gorm.DB
}

func (cs *CommentService) CreateComment(userId, photoId uint, message string) (models.Comment, error) {
	newComment := models.Comment{
		UserID:  userId,
		PhotoID: photoId,
		Message: message,
	}

	if err := cs.DB.Create(&newComment).Error; err != nil {
		return models.Comment{}, err
	}

	return newComment, nil
}

func (cs *CommentService) GetAllComments(userId uint) ([]models.Comment, error) {
	var comments []models.Comment

	if err := cs.DB.Where("user_id = ?", userId).Find(&comments).Error; err != nil {
		return []models.Comment{}, err
	}

	return comments, nil
}

func (cs *CommentService) GetOneComment(photoId, messageId uint) (models.Comment, error) {
	var comment models.Comment

	if err := cs.DB.Where("photo_id = ?", photoId).Where("id = ?", messageId).First(&comment).Error; err != nil {
		return models.Comment{}, err
	}

	return comment, nil
}

func (cs *CommentService) UpdateComment(id uint, comment models.Comment) (models.Comment, error) {

	if err := cs.DB.Model(models.Comment{}).Where("id = ?", id).Updates(&comment).Error; err != nil {
		return models.Comment{}, err
	}

	return comment, nil
}

func (cs *CommentService) DeleteComment(id uint) error {
	var deleteComment models.Comment
	if err := cs.DB.Model(models.Comment{}).Where("id = ?", id).Delete(&deleteComment).Error; err != nil {
		return err
	}

	return nil
}
