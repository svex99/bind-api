package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/handlers"
	"github.com/svex99/bind-api/middlewares"
	"github.com/svex99/bind-api/models"
	"github.com/svex99/bind-api/utils/setting"
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
	protected.GET("/domains/:domain_id", handlers.GetDomain)
	protected.PATCH("/domains/:domain_id", handlers.UpdateDomain)
	protected.DELETE("/domains/:domain_id", handlers.DeleteDomain)
	// subdomain handlers
	protected.GET("/domains/:domain_id/subdomains", handlers.ListSubdomains)
	protected.POST("/domains/:domain_id/subdomains", handlers.NewSubdomain)
	protected.GET("/domains/:domain_id/subdomains/:subdomain_id", handlers.GetSubdomain)
	protected.PATCH("/domains/:domain_id/subdomains/:subdomain_id", handlers.UpdateSubdomain)
	protected.DELETE("/domains/:domain_id/subdomains/:subdomain_id", handlers.DeleteSubdomain)

	router.Run(":2020")
}
