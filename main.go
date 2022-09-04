package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/handlers"
	"github.com/svex99/bind-api/middlewares"
	"github.com/svex99/bind-api/models"
	"github.com/svex99/bind-api/pkg/setting"
)

func main() {
	models.ConnectDatabase()

	router := gin.Default()

	public := router.Group("/api")
	router.Use(middlewares.CORS())
	public.POST("/register", handlers.Register)
	public.POST("/login", handlers.Login)
	public.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"secret":   setting.App.JwtSecret,
			"lifespan": setting.App.TokenHourLifespan,
		})
	})

	protected := router.Group("/api")
	protected.Use(middlewares.JWTAuth())
	// domain handlers
	protected.GET("/domains", handlers.ListDomains)
	protected.POST("/domains", handlers.NewDomain)
	protected.GET("/domains/:domainId", handlers.GetDomain)
	protected.PATCH("/domains/:domainId", handlers.UpdateDomain)
	protected.DELETE("/domains/:domainId", handlers.DeleteDomain)
	// subdomain handlers
	protected.GET("/domains/:domainId/subdomains", handlers.ListSubdomains)
	protected.POST("/domains/:domainId/subdomains", handlers.NewSubdomain)
	protected.GET("/domains/:domainId/subdomains/:resourceId", handlers.GetSubdomain)
	protected.PATCH("/domains/:domainId/subdomains/:resourceId", handlers.UpdateSubdomain)
	protected.DELETE("/domains/:domainId/subdomains/:resourceId", handlers.DeleteSubdomain)
	// email handlers
	protected.GET("/domains/:domainId/emails", handlers.ListEmails)
	protected.POST("domains/:domainId/emails", handlers.NewEmail)
	protected.GET("domains/:domainId/emails/:resourceId", handlers.GetEmail)
	protected.PATCH("domains/:domainId/emails/:resourceId", handlers.UpdateEmail)
	protected.DELETE("domains/:domainId/emails/:resourceId", handlers.DeleteEmail)
	// TXT handlers
	protected.GET("/domains/:domainId/TXTrecords", handlers.ListTXTRecords)
	protected.POST("domains/:domainId/TXTrecords", handlers.NewTXTRecord)
	protected.GET("domains/:domainId/TXTrecords/:resourceId", handlers.GetTXTRecord)
	protected.PATCH("domains/:domainId/TXTrecords/:resourceId", handlers.UpdateTXTRecord)
	protected.DELETE("domains/:domainId/TXTrecords/:resourceId", handlers.DeleteTXTRecord)

	router.Run(":2020")
}
