package middleware

import (
	"net/http"
	"strings"
	"wallet-app-server/app/logger"
	"wallet-app-server/app/redis"
	"wallet-app-server/app/service"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
)

// Authentication middleware
// This is used to protect wallet/transaction endpoints
// that require a authenticated user session
func Authentication(c *gin.Context) {
	// Extract authorization header
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Authorization header is not valid",
		})
		return
	}
	// Extract access token
	accessToken := authHeader[7:]
	// Fetch user ID from Redis
	currentUserID, err := redis.Client.Get(accessToken)
	if err != nil {
		// Record not found error
		if err == goredis.Nil {
			logger.Warn("Failed to fetch currentUserID, accessToken: %s, err: %s", accessToken, err.Error())
			serviceErr := service.ServiceError{ErrType: service.ErrTypeAuthenticationFailed, ErrMessage: service.ErrMessageInvalidAccessToken}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   serviceErr.Error(),
			})
			return
		}
		// Other error
		serviceErr := service.ServiceError{ErrType: service.ErrTypeInternalServerError, ErrMessage: service.ErrMessageDBError, Cause: err}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   serviceErr.Error(),
		})
		return
	}
	// Set current_user_id
	c.Set("current_user_id", currentUserID)
	// Process next handler
	c.Next()
}
