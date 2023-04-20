package controllers

import (
	"final-project-hacktiv8/helpers"
	"final-project-hacktiv8/models"
	"final-project-hacktiv8/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileController struct {
	FileService *services.FileService
}

func (fc *FileController) GetImages(c *gin.Context) {
	fileId := c.Param("image")

	var (
		file models.File
		err  error
	)

	fileName, err := helpers.VerifyImage(fileId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	photoName := fileName.(string)

	file, err = fc.FileService.GetUploadedFile(photoName)

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename="+file.Name)
	c.Header("Content-Type", http.DetectContentType(file.File))
	c.Header("Content-Length", strconv.Itoa(len(file.File)))
	c.Writer.Write(file.File)
}
