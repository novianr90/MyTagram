package controllers

import (
	"final-project-hacktiv8/helpers"
	"final-project-hacktiv8/models"
	"final-project-hacktiv8/services"
	"final-project-hacktiv8/vendor/github.com/golang-jwt/jwt/v5"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SocialMediaController struct {
	Service *services.SocialMediaService
}

type SocialMediaDto struct {
	Name           string `json:"name" binding:"required" form:"name"`
	SocialMediaUrl string `jsonj:"social_media_url" binding:"required" form:"social_media_url"`
}

func (smc *SocialMediaController) CreateSocialMedia(c *gin.Context) {
	var (
		data = c.MustGet("userData").(jwt.MapClaims)

		err error

		dto SocialMediaDto

		contentType = helpers.GetContentType(c)
	)

	if contentType == helpers.AppJson {
		if err = c.ShouldBindJSON(&dto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {
		if err = c.ShouldBind(&dto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	newAccount := models.SocialMedia{
		Name:           dto.Name,
		SocialMediaUrl: dto.SocialMediaUrl,
		UserID:         uint(data["id"].(float64)),
	}

	result, err := smc.Service.CreateSocialMedia(newAccount)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": result,
	})
}
