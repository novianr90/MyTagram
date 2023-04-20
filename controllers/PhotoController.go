package controllers

import (
	"final-project-hacktiv8/models"
	"final-project-hacktiv8/services"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type PhotoController struct {
	PhotoService *services.PhotoService
}

type PhotoDto struct {
	Title   string                `form:"title"`
	Caption string                `form:"caption"`
	Photo   *multipart.FileHeader `form:"photo"`
}

type PhotoResponse struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
}

func (pc *PhotoController) CreatePhoto(c *gin.Context) {
	var (
		userData = c.MustGet("userData").(jwt.MapClaims)

		userId = uint(userData["id"].(float64))

		photoDto PhotoDto

		err error
	)

	if err = c.ShouldBind(&photoDto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "error saat bind",
		})
		return
	}

	fileData, err := photoDto.Photo.Open()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error pas open file",
		})
	}

	fileSizeInfo := photoDto.Photo.Size

	defer fileData.Close()

	fileBytes, err := io.ReadAll(io.LimitReader(fileData, fileSizeInfo))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error pas read file data",
		})
	}

	newFile := models.File{
		Name: photoDto.Title,
		File: fileBytes,
	}

	url := fmt.Sprintf("https://mytagram-production.up.railway.app/files/%d/%s", newFile.ID, newFile.Name)

	photo := models.Photo{
		Title:    photoDto.Title,
		Caption:  photoDto.Caption,
		PhotoUrl: url,
		UserID:   userId,
	}

	result, err := pc.PhotoService.CreatePhoto(photo)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error saat create",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photo": result,
	})
}

func (pc *PhotoController) GetAll(c *gin.Context) {
	var (
		userData = c.MustGet("userData").(jwt.MapClaims)

		userId = uint(userData["id"].(float64))

		photosResponse []PhotoResponse

		err error
	)

	photos, err := pc.PhotoService.GetAll(userId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, value := range photos {
		photosResponse = append(photosResponse, PhotoResponse{
			Title:    value.Title,
			Caption:  value.Caption,
			PhotoUrl: value.PhotoUrl,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"photos": photosResponse,
	})
}

func (pc *PhotoController) GetPhotoById(c *gin.Context) {
	var (
		err error

		photo models.Photo
	)

	data := c.MustGet("mapData").(map[string]any)

	photo = data["photo"].(models.Photo)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photo": photo,
	})
}
