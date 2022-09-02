package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/svex99/bind-api/models"
	"github.com/svex99/bind-api/pkg/path"
)

func ListEmails(c *gin.Context) {
	pathData, err := path.ParsePath(c)
	if err != nil {
		return
	}

	var emails []models.Email

	if err := models.DB.Model(
		&models.ARecord{},
	).Select(
		"mx_records.id, mx_records.priority, mx_records.email_server, a_records.ip",
	).Joins(
		"left join mx_records on mx_records.email_server = a_records.name",
	).Find(
		&emails, "mx_records.domain_id = ?", pathData.DomainId,
	).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"emails": emails})
}

func GetEmail(c *gin.Context) {
	pathData, err := path.ParsePath(c)
	if err != nil {
		return
	}

	var email models.Email

	if err := models.DB.Model(
		&models.ARecord{},
	).Select(
		"mx_records.id, mx_records.priority, mx_records.email_server as name, a_records.ip",
	).Joins(
		"left join mx_records on mx_records.email_server = a_records.name",
	).First(
		&email, "mx_records.domain_id = ?", pathData.DomainId, "mx_records.id = ?", pathData.ResourceId,
	).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, email)
}

func NewEmail(c *gin.Context) {
	pathData, err := path.ParsePath(c)
	if err != nil {
		return
	}

	var email models.Email

	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := email.Create(pathData.DomainId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, email)
}

func UpdateEmail(c *gin.Context) {
	pathData, err := path.ParsePath(c)
	if err != nil {
		return
	}

	var emailForm models.UpdateEmailForm

	if err := c.ShouldBindJSON(&emailForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	email := &models.Email{Id: pathData.ResourceId}

	if err := email.Update(pathData.DomainId, &emailForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, email)
}
