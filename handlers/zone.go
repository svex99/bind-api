package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/schemas"
	"github.com/svex99/bind-api/services/bind"
	"github.com/svex99/bind-api/services/bind/parser"
)

func ListZones(c *gin.Context) {
	bind.Service.Mutex.Lock()
	defer bind.Service.Mutex.Unlock()

	zones := []*parser.ZoneConf{}

	for _, zone := range bind.Service.Zones {
		zones = append(zones, zone)
	}

	c.JSON(http.StatusOK, zones)
}

func GetZone(c *gin.Context) {
	origin := c.Param("origin")

	zConf, ok := bind.Service.Zones[origin]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("zone %s does not exist", origin)})
		return
	}

	c.JSON(http.StatusOK, zConf)
}

func NewZone(c *gin.Context) {
	var data schemas.ZoneData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	zConf, err := bind.Service.CreateZone(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, zConf)
}

func PatchZone(c *gin.Context) {
	var data schemas.ZoneData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dConf, err := bind.Service.UpdateZone(data.Origin, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dConf)
}

func DeleteZone(c *gin.Context) {
	origin := c.Param("origin")

	if err := bind.Service.DeleteZone(origin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
