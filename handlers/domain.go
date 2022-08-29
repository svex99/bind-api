package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/models"
	"github.com/svex99/bind-api/pkg/path"
)

func ListDomains(c *gin.Context) {
	// TODO: implement pagination
	var domains []models.Domain

	if err := models.DB.Model(&models.Domain{}).Find(&domains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"domains": domains,
	})
}

func GetDomain(c *gin.Context) {
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	domain := models.Domain{Id: pathData.DomainId}

	if err := models.DB.First(&domain).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain)
}

func NewDomain(c *gin.Context) {
	var domain models.Domain

	if err := c.ShouldBindJSON(&domain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := domain.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, domain)
}

func UpdateDomain(c *gin.Context) {
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	domain := models.Domain{Id: pathData.DomainId}

	if err := models.DB.First(&domain).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var domainForm models.UpdateDomainForm

	if err := c.ShouldBindJSON(&domainForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := domain.Update(&domainForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: Update the domain in bind

	c.JSON(http.StatusOK, domain)
}

func DeleteDomain(c *gin.Context) {
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	domain := models.Domain{Id: pathData.DomainId}

	if err := domain.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
