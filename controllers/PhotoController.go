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
	Title   string                `form:"title" binding:"required"`
	Caption string                `form:"caption" binding:"required"`
	Photo   *multipart.FileHeader `form:"photo"`
}

type PhotoResponse struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
}

// CreatePhoto godoc
// @Security ApiKeyAuth
// @Summary Create new Photo
// @Description Create new photo with specific user
// @Tags Photos
// @Accept mpfd
// @Produce json
// @Param title formData string true "Title"
// @Param caption formData string true "Caption"
// @Param photo formData file true "Photo to upload"
// @Success 200 {object} models.Photo
// @Router /photos [post]
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

// GetAll godoc
// @Security ApiKeyAuth
// @Summary Get All Photos
// @Description Get all photos for specific user
// @Tags Photos
// @Accept mpfd
// @Produce json
// @Success 200 {array} models.Photo
// @Router /photos [get]
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

// GetPhoto godoc
// @Security ApiKeyAuth
// @Summary Get Photo by ID
// @Description Get photo by id for specific user
// @Tags Photos
// @Accept mpfd
// @Produce json
// @Param id path int true "Photo ID"
// @Success 200 {object} models.Photo
// @Router /photos/{id} [get]
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

// Update godoc
// @Security ApiKeyAuth
// @Summary Update Photo by ID
// @Description Update photo by id for specific user
// @Tags Photos
// @Accept mpfd
// @Produce json
// @Param id path int true "Photo ID"
// @Param title formData string false "Title to update"
// @Param caption formData string false "Caption to update"
// @Param photo formData file false "Photo to update"
// @Success 200 {string} string "data sucesfully updated"
// @Router /photos/{id} [put]
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

	photoId := data["photoId"].(int)

	if photoDto.Photo == nil {
		newPhoto := models.Photo{
			Title:   photoDto.Title,
			Caption: photoDto.Caption,
		}

		err = pc.PhotoService.UpdatePhoto(uint(photoId), newPhoto)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {

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

		_, err = pc.FileService.SaveUploadedFile(newFileData)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error saat upload",
			})
			return
		}

		url := fmt.Sprintf("https://mytagram-production.up.railway.app/files/%s", helpers.GenerateTokenForImage(newFileData.Name))

		newPhoto := models.Photo{
			Title:    photoDto.Title,
			Caption:  photoDto.Caption,
			PhotoUrl: url,
		}

		err = pc.PhotoService.UpdatePhoto(uint(photoId), newPhoto)

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

// Delete godoc
// @Security ApiKeyAuth
// @Summary Delete Photo by ID
// @Description Delete photo by id for specific user
// @Tags Photos
// @Accept mpfd
// @Produce json
// @Success 200 {string} string "Data sucessfully deleted"
// @Router /photos/{id} [put]
func (pc *PhotoController) DeletePhotoById(c *gin.Context) {

	data := c.MustGet("data").(map[string]any)

	photoId := uint(data["photoId"].(int))

	if err := pc.PhotoService.DeletePhotoById(photoId); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	photo := data["photo"].(models.Photo)

	photoName := strings.Split(photo.PhotoUrl, "/")[4]

	fileData, err := helpers.VerifyImage(photoName)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	file := fileData.(jwt.MapClaims)

	fileName := file["pathToPhoto"].(string)

	err = pc.FileService.DeleteFile(fileName)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messagee": "data sucessfully deleted",
	})
}
