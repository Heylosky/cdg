package middlewares

import (
	"github.com/ComfortDelgro/redis"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func TokenVerify(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"data":   "",
			"status": "No token found",
		})
		c.Abort()
		return
	}
	_, exist := redis.GetToken(token)
	if !exist {
		zap.L().Info("Token authorization failed.")
		c.JSON(http.StatusUnauthorized, gin.H{
			"data":   "",
			"status": "Token verification failed",
		})
		c.Abort()
	}
}
