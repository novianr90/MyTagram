package middlewares

import (
	"final-project-hacktiv8/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CommentAuth(commentService *services.CommentService) gin.HandlerFunc {
	return func(c *gin.Context) {

		userData := c.MustGet("userData").(jwt.MapClaims)

		userId := uint(userData["id"].(float64))

		commentId, err := strconv.Atoi(c.Param("commentId"))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error when parse to int",
			})
			return
		}

		comment, err := commentService.GetOneComment(userId, uint(commentId))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error when get comment from DB",
			})
			return
		}

		if comment.UserID != userId {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"messasge": "unauthorized",
			})
			return
		}

		c.Set("dataComment", comment)

		c.Next()
	}
}
