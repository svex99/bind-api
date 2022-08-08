package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/svex99/bind-api/models"
)

type authInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input authInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := models.User{}
	user.Email = input.Email
	user.Password = input.Password

	if _, err := user.SaveUser(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
	})
}

func Login(c *gin.Context) {
	var input authInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := models.User{}
	user.Email = input.Email
	user.Password = input.Password

	token, err := models.LoginUser(user.Email, user.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email/password combination",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
