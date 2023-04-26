package controllers

import (
	"final-project-hacktiv8/helpers"
	"final-project-hacktiv8/models"
	"final-project-hacktiv8/services"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type PhotoController struct {
	PhotoService *services.PhotoService
	FileService  *services.FileService
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
		return
	}

	fileSizeInfo := photoDto.Photo.Size

	defer fileData.Close()

	fileBytes, err := io.ReadAll(io.LimitReader(fileData, fileSizeInfo))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error pas read file data",
		})
		return
	}

	newFile := models.File{
		Name: photoDto.Title,
		File: fileBytes,
	}

	_, err = pc.FileService.SaveUploadedFile(newFile)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error saat upload",
		})
		return
	}

	url := fmt.Sprintf("https://mytagram-production.up.railway.app/files/%s", helpers.GenerateTokenForImage(newFile.Name))

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

	data := c.MustGet("data").(map[string]any)

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

func (pc *PhotoController) UpdatePhotoById(c *gin.Context) {
	var (
		err error

		photoDto PhotoDto
	)

	if err = c.ShouldBind(&photoDto); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	data := c.MustGet("data").(map[string]any)

	photo := data["photo"].(models.Photo)
	photoId := data["photoId"].(uint)

	if photoDto.Photo == nil {
		newPhoto := models.Photo{
			Title:   photoDto.Title,
			Caption: photoDto.Caption,
		}

		err = pc.PhotoService.UpdatePhoto(photoId, newPhoto)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {
		photoToken := strings.Split(photo.PhotoUrl, "/")[4]

		currentFile, err := helpers.VerifyImage(photoToken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error pas open file",
			})
		}

		photoName := currentFile.(jwt.MapClaims)

		fileData, err := photoDto.Photo.Open()

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error pas open file",
			})
			return
		}

		fileSizeInfo := photoDto.Photo.Size

		defer fileData.Close()

		fileBytes, err := io.ReadAll(io.LimitReader(fileData, fileSizeInfo))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error pas read file data",
			})
			return
		}

		newFileData := models.File{
			Name: photoDto.Photo.Filename,
			File: fileBytes,
		}

		newFile, err := pc.FileService.UpdateFile(photoName["pathToPhoto"].(string), newFileData)

		url := fmt.Sprintf("https://mytagram-production.up.railway.app/files/%s", helpers.GenerateTokenForImage(newFile.Name))

		newPhoto := models.Photo{
			Title:    photoDto.Title,
			Caption:  photoDto.Caption,
			PhotoUrl: url,
		}

		err = pc.PhotoService.UpdatePhoto(photoId, newPhoto)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "data sucessfully updated",
	})
}
