package middlewares

import (
	"final-project-hacktiv8/helpers"
	"final-project-hacktiv8/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyToken, err := helpers.VerifyToken(c)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_status":  "Unauthenticated",
				"error_message": err.Error(),
			})
			return
		}

		userData := verifyToken.(models.User)

		c.Set("userData", userData)

		c.Next()
	}
}
