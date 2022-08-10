package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/models"
	"github.com/svex99/bind-api/utils/path"
)

func ListDomains(c *gin.Context) {
	// TODO: implement pagination
	var domains []models.DomainInfo

	if err := models.DB.Model(&models.Domain{}).Select("id", "name").Find(&domains).Error; err != nil {
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

	if err := models.DB.Create(&domain).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: Add the domain to bind

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

	if domainForm.Name != "" {
		domain.Name = domainForm.Name
	}
	if domainForm.NameServer != "" {
		domain.NameServer = domainForm.NameServer
	}
	if domainForm.NSIp != "" {
		domain.NSIp = domainForm.NSIp
	}
	if domainForm.Ttl != "" {
		domain.Ttl = domainForm.Ttl
	}

	if err := models.DB.Save(&domain).Error; err != nil {
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

	if err := models.DB.Delete(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: Remove the domain from bind

	c.JSON(http.StatusNoContent, gin.H{})
}
