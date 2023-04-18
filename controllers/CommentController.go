package controllers

import (
	"final-project-hacktiv8/services"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	CommentService *services.CommentService
}

func (cc *CommentController) CreateComment(c *gin.Context) {

}
