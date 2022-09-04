package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/svex99/bind-api/models"
	"github.com/svex99/bind-api/pkg/path"
)

func ListTXTRecords(c *gin.Context) {
	pathData, err := path.ParsePath(c)
	if err != nil {
		return
	}

	var txtRecords []models.TXTRecord
	whereStruct := &models.TXTRecord{
		Record: models.Record{DomainId: pathData.DomainId},
	}

	if err := models.DB.Where(whereStruct).Find(&txtRecords).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"TXTrecords": txtRecords})
}

func GetTXTRecord(c *gin.Context) {
	pathData, err := path.ParsePath(c)
	if err != nil {
		return
	}

	txtRecord := &models.TXTRecord{
		Record: models.Record{
			Id:       pathData.ResourceId,
			DomainId: pathData.DomainId,
		},
	}

	if err := models.DB.Where(txtRecord).First(txtRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, txtRecord)
}

func NewTXTRecord(c *gin.Context) {
	pathData, err := path.ParsePath(c)
	if err != nil {
		return
	}

	txtRecord := &models.TXTRecord{
		Record: models.Record{
			DomainId: pathData.DomainId,
		},
	}

	if err := c.ShouldBindJSON(txtRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := txtRecord.Create(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, txtRecord)
}

func UpdateTXTRecord(c *gin.Context) {
	pathData, err := path.ParsePath(c)
	if err != nil {
		return
	}

	txtForm := &models.UpdateTXTRecordForm{}

	if err := c.ShouldBindJSON(txtForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txtRecord := &models.TXTRecord{
		Record: models.Record{
			Id:       pathData.ResourceId,
			DomainId: pathData.DomainId,
		},
	}

	if err := txtRecord.Update(txtForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, txtRecord)
}

func DeleteTXTRecord(c *gin.Context) {
	pathData, err := path.ParsePath(c)
	if err != nil {
		return
	}

	txtRecord := &models.TXTRecord{
		Record: models.Record{
			Id:       pathData.ResourceId,
			DomainId: pathData.DomainId,
		},
	}

	if err := txtRecord.Delete(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
