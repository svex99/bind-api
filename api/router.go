package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/svex99/bind-api/handlers"
	"github.com/svex99/bind-api/middlewares"
	"github.com/svex99/bind-api/pkg/setting"
)

func SetupRouter() *gin.Engine {
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
	// protected.Use(middlewares.JWTAuth())
	// domain handlers
	protected.GET("/domains", handlers.ListDomains)
	protected.POST("/domains", handlers.NewDomain)
	protected.PATCH("/domains", handlers.UpdateDomain)
	protected.GET("/domains/:origin", handlers.GetDomain)
	protected.DELETE("/domains/:origin", handlers.DeleteDomain)
	// record handlers
	protected.POST("/domains/:origin/record/:type", handlers.PostRecord)
	protected.PATCH("/domains/:origin/record/:type/:hash", handlers.PatchRecord)
	protected.DELETE("/domains/:origin/record/:type", handlers.DeleteRecord)

	return router
}
