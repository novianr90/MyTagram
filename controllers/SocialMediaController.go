package controllers

import (
	"final-project-hacktiv8/helpers"
	"final-project-hacktiv8/models"
	"final-project-hacktiv8/services"
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

type SocialMediaController struct {
	Service *services.SocialMediaService
}

type SocialMediaDto struct {
	ID             uint   `json:"id"`
	Name           string `json:"name" form:"name"`
	SocialMediaUrl string `jsonj:"social_media_url" form:"social_media_url"`
	UserID         uint   `json:"user_id"`
}

// CreateSocialMedia godoc
// @Security BearerAuth
// @Summary Create new Social Media
// @Description Create new social media with specific user
// @Tags Social Media
// @Accept mpfd
// @Produce json
// @Param name formData string true "Social Media Name"
// @Param social_media_url formData string true "Social Media URL"
// @Success 200 {object} models.SocialMedia
// @Router /accounts [post]
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
		"message": models.SocialMedia{
			ID:             result.ID,
			Name:           result.Name,
			SocialMediaUrl: result.SocialMediaUrl,
			UserID:         result.UserID,
		},
	})
}

// GetAllSocialMedia godoc
// @Security BearerAuth
// @Summary Get All Social Media
// @Description Get all social media with specific user
// @Tags Social Media
// @Accept mpfd
// @Produce json
// @Success 200 {array} models.SocialMedia
// @Router /accounts [get]
func (smc *SocialMediaController) GetAllAccounts(c *gin.Context) {
	var (
		data = c.MustGet("userData").(jwt.MapClaims)

		accounts []models.SocialMedia

		response []SocialMediaDto

		err error
	)

	accounts, err = smc.Service.GetAllSocialMedia(uint(data["id"].(float64)))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	for _, value := range accounts {
		response = append(response, SocialMediaDto{
			ID:             value.ID,
			Name:           value.Name,
			SocialMediaUrl: value.SocialMediaUrl,
			UserID:         value.UserID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": response,
	})
}

// GetSocialMedia godoc
// @Security BearerAuth
// @Summary Get Social Media
// @Description Get social media for specific user
// @Tags Social Media
// @Accept mpfd
// @Produce json
// @Param id path int true "Social Media ID"
// @Success 200 {object} models.SocialMedia
// @Router /accounts/{id} [get]
func (smc *SocialMediaController) GetAccountById(c *gin.Context) {
	var (
		data = c.MustGet("userAndAccountId").(map[string]uint)

		account models.SocialMedia

		response SocialMediaDto

		err error
	)

	account, err = smc.Service.GetAccountById(data["accountId"])

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	response = SocialMediaDto{
		ID:             account.ID,
		Name:           account.Name,
		SocialMediaUrl: account.SocialMediaUrl,
		UserID:         data["userId"],
	}

	c.JSON(http.StatusOK, gin.H{
		"account": response,
	})
}

// GetSocialMedia godoc
// @Security BearerAuth
// @Summary Update Social Media
// @Description Update social media for specific user
// @Tags Social Media
// @Accept mpfd
// @Produce json
// @Param id path int true "Social Media ID"
// @Param name formData string false "Name social media to update"
// @Param social_media_url formData string false "url to update"
// @Success 200 {string} string "data sucesfully updated"
// @Router /accounts/{id} [put]
func (smc *SocialMediaController) UpdateAccount(c *gin.Context) {
	var (
		data = c.MustGet("userAndAccountId").(map[string]uint)

		dto SocialMediaDto

		err error

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
		if err := c.ShouldBind(&dto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	if err = smc.Service.UpdateAccounts(dto.Name, dto.SocialMediaUrl, data["accountId"]); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "data sucesfully updated",
	})
}

// DeleteSocialMedia godoc
// @Security BearerAuth
// @Summary Delete Social Media
// @Description Delete social media for specific user
// @Tags Social Media
// @Accept mpfd
// @Produce json
// @Param id path int true "Social Media ID"
// @Success 200 {string} string "data sucesfully deleted"
// @Router /accounts/{id} [delete]
func (smc *SocialMediaController) DeleteAccount(c *gin.Context) {
	var (
		data = c.MustGet("userAndAccountId").(map[string]uint)

		err error
	)

	if err = smc.Service.DeleteAccount(data["accountId"]); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "data sucesffully deleted",
	})
}
