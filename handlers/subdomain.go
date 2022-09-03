package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/models"
	"github.com/svex99/bind-api/pkg/path"
)

func ListSubdomains(c *gin.Context) {
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	var subdomains []models.Subdomain

	if err := models.DB.Model(
		&models.ARecord{},
	).Where(
		&models.ARecord{Record: models.Record{DomainId: pathData.DomainId}},
	).Find(&subdomains).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"subdomains": subdomains,
	})
}

func GetSubdomain(c *gin.Context) {
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	var subdomain models.Subdomain

	if err := models.DB.Model(
		&models.ARecord{},
	).Where(
		&models.ARecord{Record: models.Record{Id: pathData.ResourceId, DomainId: pathData.DomainId}},
	).First(&subdomain).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subdomain)
}

func NewSubdomain(c *gin.Context) {
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	var subdomain models.Subdomain

	if err := c.ShouldBindJSON(&subdomain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := subdomain.Create(pathData.DomainId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, subdomain)
}

func UpdateSubdomain(c *gin.Context) {
	// TODO: Handle update to subdomain of the name server
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	var subdomainForm models.UpdateSubdomainForm

	if err := c.ShouldBindJSON(&subdomainForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	subdomain := models.Subdomain{Id: pathData.ResourceId}

	if err := subdomain.Update(pathData.DomainId, &subdomainForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subdomain)
}

func DeleteSubdomain(c *gin.Context) {
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	subdomain := models.Subdomain{Id: pathData.ResourceId}

	if err := subdomain.Delete(pathData.DomainId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
