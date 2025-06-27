package controllers

import (
	"Financial/Domains/ports"
	"Financial/Models/dtos"
	"Financial/intefaces/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	*BaseController
	wallet          ports.WalletUseCase
	authMiddlerware *middleware.AuthMiddleware
}

func NewWalletController(walletUseCase ports.WalletUseCase, auth *middleware.AuthMiddleware) *WalletController {
	return &WalletController{
		BaseController:  NewBaseController("/wallet"),
		wallet:          walletUseCase,
		authMiddlerware: auth,
	}
}

func (wc *WalletController) RegisterRoutes(router *gin.RouterGroup) {
	wc.authMiddlerware.Config.AddPublicRoute("GET", "api/wallet/:email")
	public := router.Group("/wallet")
	{
		public.GET(":email", wc.getUserWallets)
	}
	protected := router.Group("/wallet")
	protected.Use(wc.authMiddlerware.AuthMiddleware())
	{
		protected.GET("", wc.getUserWallets)
		protected.POST("", wc.createWallet)
		protected.PUT(":id", wc.updateWallet)
		protected.DELETE(":id", wc.deleteWallet)
	}
}

// getUserWallets handles GET /wallets - Get all wallets for the authenticated user
func (wc *WalletController) getUserWallets(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	wallet, err := wc.wallet.GetUserWallet(0, email.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, wallet)
}

// createWallet handles POST /wallet - Create a new wallet
func (wc *WalletController) createWallet(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request dtos.CreateWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set the user ID from the authenticated user
	request.UserID = userID.(int)

	wallet, err := wc.wallet.CreateWallet(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

// updateWallet handles PUT /wallet/:id - Update a wallet
func (wc *WalletController) updateWallet(c *gin.Context) {
	var request dtos.UpdateWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	updatedWallet, err := wc.wallet.UpdateWallet(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedWallet)
}

// deleteWallet handles DELETE /wallet/:id - Delete a wallet
func (wc *WalletController) deleteWallet(c *gin.Context) {
	var request dtos.DeleteWalletRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		var errorRes = dtos.ErrorResponse{
			Error: "Solicitud inválida",
		}
		c.JSON(400, errorRes)
		return
	}

	if err := wc.wallet.DeleteWallet(request.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete wallet"})
		return
	}

	c.Status(http.StatusNoContent)
}
