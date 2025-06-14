package controller

import (
	"net/http"
	"wallet-app-server/app/service"

	"github.com/gin-gonic/gin"
)

// User login
// POST /user/login
func Login(c *gin.Context) {
	// Parse request body
	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := c.BindJSON(&req); err != nil {
		respondeWithError(c, http.StatusBadRequest, err)
		return
	}

	// User login
	accessToken, statusCode, err := service.UserService.Login(req.Username, req.Password)
	if err != nil {
		respondeWithError(c, statusCode, err)
		return
	}

	// Return resposne
	resposneWithData(c, gin.H{"access_token": accessToken})
}
