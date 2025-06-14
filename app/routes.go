package app

import (
	"wallet-app-server/app/controller"
	"wallet-app-server/app/middleware"

	"github.com/gin-gonic/gin"
)

func configRoutes(g *gin.Engine) {
	apiGroup := g.Group("/api/v1")

	// User endpoints
	userGroup := apiGroup.Group("/user")
	userGroup.POST("/login", controller.Login)

	// Wallet endpoints (need authentication)
	walletGroup := apiGroup.Group("/wallet", middleware.Authentication)
	walletGroup.GET("/list", controller.ListWallets)
	walletGroup.POST("/checkBalance", controller.CheckWalletBalance)
	walletGroup.POST("/deposit", controller.Deposit)
	walletGroup.POST("/withdraw", controller.Withdraw)

	// Transaction endpoints (need authentication)
	transactionGroup := apiGroup.Group("/transaction", middleware.Authentication)
	transactionGroup.POST("/transfer", controller.Transfer)
	transactionGroup.POST("/history", controller.History)
}
