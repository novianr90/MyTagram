package middlewares

import (
	"final-project-hacktiv8/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func PhotoAuthorization(photoService *services.PhotoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		photoId, err := strconv.Atoi(c.Param("photoId"))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error_status":  "Bad request",
				"error_message": "Invalid parameter",
			})
			return
		}

		photo, err := photoService.GetPhotoById(uint(photoId))

		data := c.MustGet("userData").(jwt.MapClaims)

		userId := uint(data["id"].(float64))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		if photo.UserID != userId {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error_status":  "Unauthorized",
				"error_message": fmt.Sprintf("you do not have photo with id : %d", photoId),
			})
			return
		}

		mapData := map[string]any{
			"user":  data,
			"photo": photo,
		}

		c.Set("data", mapData)

		c.Next()
	}
}
