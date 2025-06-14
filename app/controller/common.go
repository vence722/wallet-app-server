package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func resposneWithData(c *gin.Context, data gin.H) {
	// Define response map
	resp := gin.H{
		"success": true,
	}
	// Merge data into response map
	for k, v := range data {
		resp[k] = v
	}
	// Return success response
	c.JSON(http.StatusOK, resp)
}

func respondeWithError(c *gin.Context, statusCode int, err error) {
	// Return error resposne
	c.AbortWithStatusJSON(statusCode, gin.H{
		"success": false,
		"error":   err.Error(),
	})
}
