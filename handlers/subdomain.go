package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/models"
	"github.com/svex99/bind-api/utils/path"
)

func ListSubdomains(c *gin.Context) {
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	domain := models.Domain{Id: pathData.DomainId}

	if err := models.DB.First(&domain).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := models.DB.Preload("Subdomains").First(&domain).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var subdomainDatas []models.SubdomainInfo

	for _, sd := range domain.Subdomains {
		subdomainDatas = append(subdomainDatas, models.SubdomainInfo{
			Id:   sd.Id,
			Name: sd.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"subdomains": subdomainDatas,
	})
}

func GetSubdomain(c *gin.Context) {
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	subdomain := models.Subdomain{Id: pathData.SubdomainId}

	if err := models.DB.First(&subdomain, "domain_id = ?", pathData.DomainId).Error; err != nil {
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

	// TODO: Get the domain id directly from the uri
	subdomain.DomainId = pathData.DomainId

	if err := models.DB.Save(&subdomain).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: Add the subdomain to bind

	c.JSON(http.StatusCreated, gin.H{
		"subdomain": subdomain,
	})
}

func UpdateSubdomain(c *gin.Context) {
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

	subdomain := models.Subdomain{Id: pathData.SubdomainId}

	if err := models.DB.First(&subdomain, "domain_id = ?", pathData.DomainId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if subdomainForm.Name != "" {
		subdomain.Name = subdomainForm.Name
	}
	if subdomainForm.Ip != "" {
		subdomain.Ip = subdomainForm.Ip
	}

	if err := models.DB.Save(&subdomain).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO: Update the subdomain in bind

	c.JSON(http.StatusOK, subdomain)
}

func DeleteSubdomain(c *gin.Context) {
	pathData, err := path.ParsePath(c)

	if err != nil {
		return
	}

	subdomain := models.Subdomain{Id: pathData.SubdomainId}

	if err := models.DB.Delete(&subdomain, "domain_id = ?", pathData.DomainId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
