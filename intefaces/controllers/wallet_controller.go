// @title           Financial App API
// @version         1.0
// @description     This is a financial application server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.yourdomain.com/support
// @contact.email  support@yourdomain.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8085
// @BasePath  /api
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

package controllers

import (
	request "Financial/Core/Models/dtos/Request"
	response "Financial/Core/Models/dtos/Response"
	contracts "Financial/Core/ports"
	"Financial/intefaces/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WalletController handles wallet related operations
// @Summary Wallet management
// @Description Provides endpoints for managing user wallets
type WalletController struct {
	*BaseController
	wallet          contracts.WalletUseCase
	authMiddlerware *middleware.AuthMiddleware
}

func NewWalletController(walletUseCase contracts.WalletUseCase, auth *middleware.AuthMiddleware) *WalletController {
	return &WalletController{
		BaseController:  NewBaseController("/wallet"),
		wallet:          walletUseCase,
		authMiddlerware: auth,
	}
}

func (wc *WalletController) RegisterRoutes(router *gin.RouterGroup) {
	wc.authMiddlerware.Config.AddPublicRoute("GET", "/api/wallet/:email")
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

// getUserWallets godoc
// @Summary Get user wallets
// @Description Get all wallets for the authenticated user
// @Tags wallets
// @Accept  json
// @Produce  json
// @Param email path string true "User email"
// @Security Bearer
// @Success 200 {object} ports.UserWallet
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Router /wallet/{email} [get]
// @Router /wallet [get]
func (wc *WalletController) getUserWallets(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email is empty"})
		return
	}
	wallet, err := wc.wallet.GetUserWallet(0, email)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, wallet)
}

// createWallet godoc
// @Summary Create a new wallet
// @Description Create a new wallet for the authenticated user
// @Tags wallets
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param wallet body dtos.CreateWalletRequest true "Wallet creation data"
// @Success 201 {object} db.Wallet
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Router /wallet [post]
func (wc *WalletController) createWallet(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var request request.CreateWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Set the user ID from the authenticated user
	request.UserID = userID.(int)

	wallet, err := wc.wallet.CreateWallet(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, wallet)
}

// updateWallet godoc
// @Summary Update a wallet
// @Description Update an existing wallet's information
// @Tags wallets
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path int true "Wallet ID"
// @Param wallet body dtos.UpdateWalletRequest true "Wallet update data"
// @Success 200 {object} db.Wallet
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Router /wallet/{id} [put]
func (wc *WalletController) updateWallet(c *gin.Context) {
	var request request.UpdateWalletRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	updatedWallet, err := wc.wallet.UpdateWallet(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, updatedWallet)
}

// deleteWallet godoc
// @Summary Delete a wallet
// @Description Delete a wallet by ID
// @Tags wallets
// @Accept  json
// @Produce  json
// @Security Bearer
// @Param id path int true "Wallet ID"
// @Success 204 "No Content"
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /wallet/{id} [delete]
func (wc *WalletController) deleteWallet(c *gin.Context) {
	var request request.DeleteWalletRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		var errorRes = response.ErrorResponse{
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
