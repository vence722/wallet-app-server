package app

import (
	"wallet-app-server/app/controller"

	"github.com/gin-gonic/gin"
)

func configRoutes(g *gin.Engine) {
	apiGroup := g.Group("/api/v1")

	// User endpoints
	userGroup := apiGroup.Group("/auth")
	userGroup.POST("/login", controller.Login)

	// Wallet endpoints (need authentication)
	walletGroup := apiGroup.Group("/wallet")
	walletGroup.POST("/deposit", controller.Deposit)
	walletGroup.POST("/withdraw", controller.Withdraw)
	walletGroup.GET("/checkBalance", controller.CheckBalance)

	// Transaction endpoints (need authentication)
	transactionGroup := apiGroup.Group("/transaction")
	transactionGroup.POST("/transfer", controller.Transfer)
	transactionGroup.POST("/history", controller.History)
}
