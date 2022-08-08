package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/models"
)

func NewSubdomain(c *gin.Context) {
	var subdomain models.Subdomain

	if err := c.ShouldBindJSON(&subdomain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: Add the subdomain

	c.JSON(http.StatusOK, gin.H{
		"message": "Added new subdomain " + subdomain.Name,
	})
}
