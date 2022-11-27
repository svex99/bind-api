package api

import (
	"github.com/gin-gonic/gin"
	"github.com/svex99/bind-api/handlers"
	"github.com/svex99/bind-api/middlewares"
)

func SetupRouter(logRequests bool) *gin.Engine {
	router := gin.New()

	router.UseRawPath = true

	router.Use(gin.Recovery())
	router.Use(middlewares.CORS())
	if logRequests {
		router.Use(gin.Logger())
	}

	api := router.Group("/api")
	// domain handlers
	api.GET("/zones", handlers.ListZones)
	api.GET("/zones/:origin", handlers.GetZone)
	api.POST("/zones", handlers.NewZone)
	api.PATCH("/zones", handlers.PatchZone)
	api.DELETE("/zones/:origin", handlers.DeleteZone)
	// record handlers
	api.POST("/zones/:origin/records", handlers.PostRecord)
	api.PATCH("/zones/:origin/records/:target", handlers.PatchRecord)
	api.DELETE("/zones/:origin/records", handlers.DeleteRecord)

	return router
}
