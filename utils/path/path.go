package path

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PathData struct {
	DomainId    uint64 `uri:"domain_id"`
	SubdomainId uint64 `uri:"subdomain_id"`
}

func ParsePath(c *gin.Context) (PathData, error) {
	var pathData PathData

	if err := c.ShouldBindUri(&pathData); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "The requested resource was not found",
		})
		return pathData, err
	}

	return pathData, nil
}
