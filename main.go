package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/svex99/bind-api/handlers"
	"github.com/svex99/bind-api/middlewares"
	"github.com/svex99/bind-api/models"
	"github.com/svex99/bind-api/utils/setting"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading the .env file")
		return
	}

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
	protected.Use(middlewares.JWTAuth())
	protected.GET("/domains", handlers.ListDomains)
	protected.POST("/domains", handlers.NewDomain)
	protected.PUT("/domains/:id", handlers.UpdateDomain)
	protected.DELETE("/domains/:id", handlers.DeleteDomain)
	protected.POST("/subdomains", handlers.NewSubdomain)

	router.Run(":2020")
}
