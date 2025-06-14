package middleware

import (
	"net/http"
	"strings"
	"wallet-app-server/app/redis"
	"wallet-app-server/app/service"

	"github.com/gin-gonic/gin"
)

// Authentication middleware
// This is used to protect wallet/transaction endpoints
// that require a authenticated user session
func Authentication(c *gin.Context) {
	// Extract authorization header
	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
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
		serviceErr := service.ServiceError{ErrType: service.ErrTypeAuthenticationFailed, ErrMessage: service.ErrMessageInvalidAccessToken}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
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
