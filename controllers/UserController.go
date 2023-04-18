package controllers

import (
	"final-project-hacktiv8/helpers"
	"final-project-hacktiv8/models"
	"final-project-hacktiv8/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

type UserDto struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Age      uint   `josn:"age" form:"age"`
}

var appJson = "application/json"

func (uc *UserController) Register(c *gin.Context) {
	var (
		userDto     UserDto
		contentType = helpers.GetContentType(c)
	)

	if contentType == appJson {
		if err := c.ShouldBindJSON(&userDto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&userDto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	userPassword := helpers.HashPass(userDto.Password)

	newUser := models.User{
		Email:    userDto.Email,
		Username: userDto.Username,
		Password: userPassword,
		Age:      userDto.Age,
	}

	result, err := uc.UserService.CreateUser(newUser)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"email":    result.Email,
		"username": result.Username,
		"age":      result.Age,
	})
}

func (uc *UserController) Login(c *gin.Context) {

	var (
		contentType = helpers.GetContentType(c)
		userDto     UserDto
	)

	if contentType == appJson {
		if err := c.ShouldBindJSON(&userDto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	} else {
		if err := c.ShouldBind(&userDto); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	fmt.Println("password", userDto.Password)

	user, err := uc.UserService.GetUserByEmail(userDto.Email)

	fmt.Println("user", user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	hashedPassword := helpers.HashPass(userDto.Password)

	ok := helpers.CompareHashAndPassword([]byte(user.Password), []byte(userDto.Password))

	fmt.Println("conds", ok)

	fmt.Println(hashedPassword, user.Password)

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
	}

	token := helpers.GenerateToken(user.ID, user.Email)

	fmt.Println("token", token)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
