package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/schemas"
	"github.com/svex99/bind-api/services/bind"
	"github.com/svex99/bind-api/services/bind/parser"
)

func ListDomains(c *gin.Context) {
	bind.Service.Mutex.Lock()
	defer bind.Service.Mutex.Unlock()

	domains := []*parser.DomainConf{}

	for _, domain := range bind.Service.Domains {
		domains = append(domains, domain)
	}

	c.JSON(http.StatusOK, gin.H{
		"domains": domains,
	})
}

func GetDomain(c *gin.Context) {
	bind.Service.Mutex.Lock()
	defer bind.Service.Mutex.Unlock()

	origin := c.Param("origin")

	dConf, ok := bind.Service.Domains[origin]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("domain %s does not exist", origin)})
		return
	}

	c.JSON(http.StatusOK, dConf)
}

func NewDomain(c *gin.Context) {
	var data schemas.DomainData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dConf, err := bind.Service.CreateDomain(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dConf)
}

func UpdateDomain(c *gin.Context) {
	var data schemas.DomainData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dConf, err := bind.Service.UpdateDomain(data.Origin, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dConf)
}

func DeleteDomain(c *gin.Context) {
	origin := c.Param("origin")

	if err := bind.Service.DeleteDomain(origin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
