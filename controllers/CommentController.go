package controllers

import (
	"final-project-hacktiv8/helpers"
	"final-project-hacktiv8/models"
	"final-project-hacktiv8/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	CommentService *services.CommentService
}

type CommentDto struct {
	Message string `json:"message" form:"message"`
}

func (cc *CommentController) CreateComment(c *gin.Context) {
	var (
		data = c.MustGet("data").(map[string]any)

		commentDto CommentDto

		err error

		contentType = helpers.GetContentType(c)
	)

	if contentType == appJson {
		if err = c.ShouldBindJSON(&commentDto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {
		if err = c.ShouldBind(&commentDto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	userData := data["user"].(models.User)
	photoData := data["photo"].(models.Photo)

	result, err := cc.CommentService.CreateComment(userData.ID, photoData.ID, commentDto.Message)

	if err != nil {
		if err = c.ShouldBind(&commentDto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	fmt.Println(result)

	c.JSON(http.StatusOK, gin.H{
		"comment": result,
	})
}
