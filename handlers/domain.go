package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/models"
)

func ListDomains(c *gin.Context) {
	domains := [...]string{
		"example.com",
		"domain.com",
	}

	c.JSON(http.StatusOK, domains)
}

func NewDomain(c *gin.Context) {
	var domain models.Domain

	if err := c.ShouldBindJSON(&domain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: Add the domain

	c.JSON(http.StatusOK, domain)
}

func UpdateDomain(c *gin.Context) {
	// TODO: Get the domain identifier from URL
	// TODO: Use domain name as identifier?
	id := c.Param("id")

	var domain models.Domain

	if err := c.ShouldBindJSON(&domain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: Get the domain by its identifier

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "The domain has been updated",
	})
}

func DeleteDomain(c *gin.Context) {
	// TODO: Get the domain identifier from URL
	id := c.Param("id")

	// TODO: Try to delete the domain if exists

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "The domain has been deleted",
	})
}
