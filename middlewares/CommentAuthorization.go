package middlewares

import (
	"final-project-hacktiv8/models"
	"final-project-hacktiv8/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CommentAuth(commentService *services.CommentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		mapData := c.MustGet("data").(map[string]any)

		commentId, err := strconv.Atoi(c.Param("commentId"))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error when parse to int",
			})
			return
		}

		photo := mapData["photo"].(models.Photo)

		comment, err := commentService.GetOneComment(photo.ID, uint(commentId))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error when get comment from DB",
			})
			return
		}

		if comment.ID != photo.ID {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"messasge": "unauthorized",
			})
			return
		}

		c.Set("dataComment", comment)

		c.Next()
	}
}
