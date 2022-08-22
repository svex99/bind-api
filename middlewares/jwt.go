package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/svex99/bind-api/pkg/token"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := token.ValidateToken(c); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
