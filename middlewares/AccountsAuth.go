package middlewares

import (
	"final-project-hacktiv8/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AccountAuthorization(service *services.SocialMediaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := c.MustGet("userData").(jwt.MapClaims)

		accountId, err := strconv.Atoi(c.Param("accountId"))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		userId, err := service.GetUserIdByAccountId(uint(accountId))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		if userId == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "no record found",
			})
			return
		}

		if userId != uint(data["id"].(float64)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "login first",
			})
			return
		}

		userAndAccountId := map[string]uint{
			"userId":    userId,
			"accountId": uint(accountId),
		}

		c.Set("userAndAccountId", userAndAccountId)

		c.Next()
	}
}
