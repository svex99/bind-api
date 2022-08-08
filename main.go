package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/handlers"
	"github.com/svex99/bind-api/models"
	"github.com/svex99/bind-api/utils/setting"
)

func main() {
	models.ConnectDatabase()

	router := gin.Default()

	public := router.Group("/api")
	public.POST("/register", handlers.Register)
	public.POST("/login", handlers.Login)
	public.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"secret":   setting.App.JwtSecret,
			"lifespan": setting.App.TokenHourLifespan,
		})
	})

	protected := router.Group("/api")
	// protected.Use(middlewares.JWTAuth())
	protected.GET("/domains", handlers.ListDomains)
	protected.GET("/domain/:name", handlers.GetDomain)
	protected.POST("/domains", handlers.NewDomain)
	protected.PATCH("/domains/:name", handlers.UpdateDomain)
	protected.DELETE("/domains/:name", handlers.DeleteDomain)
	protected.POST("/subdomains", handlers.NewSubdomain)

	router.Run(":2020")
}
