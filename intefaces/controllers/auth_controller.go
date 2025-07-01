package controllers

import (
	// "Financial/Domains/ports"
	// "Financial/Models/dtos"
	request "Financial/Core/Models/dtos/Request"
	response "Financial/Core/Models/dtos/Response"
	contract "Financial/Core/ports"
	"Financial/intefaces/middleware"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication related HTTP requests
// @Description Controller for handling user authentication operations
type AuthController struct {
	*BaseController
	userUseCase    contract.UserUseCase
	authMiddleware *middleware.AuthMiddleware
}

// NewAuthController creates a new instance of AuthController
// @title Auth Controller
// @version 1.0
// @description This is the authentication controller for the Financial API
// @contact.name API Support
// @license.name Apache 2.0
// @host localhost:8080
// @BasePath /api
func NewAuthController(userUseCase contract.UserUseCase, authMiddlerware *middleware.AuthMiddleware) *AuthController {
	return &AuthController{
		BaseController: NewBaseController("/auth"),
		userUseCase:    userUseCase,
		authMiddleware: authMiddlerware,
	}
}

// RegisterRoutes sets up the routes for authentication endpoints
func (ac *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	ac.authMiddleware.Config.AddPublicRoute("POST", "/api/auth")

	public := router.Group("/auth")
	{
		public.POST("", ac.Login)
	}
}

// Login authenticates a user and returns a token
// @Summary Authenticate user
// @Description Authenticates a user with email/username and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   auth  body      dtos.AuthRequest  true  "Login credentials"
// @Success 200 {string} string "Authentication successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request format"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth [post]
func (ac *AuthController) Login(c *gin.Context) {
	var request request.AuthRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorRes := response.ErrorResponse{
			Error: "Invalid request format: " + err.Error(),
		}
		c.JSON(400, errorRes)
		return
	}

	// Validate required fields
	if (request.Email == "" && request.Nickname == "") || request.Passwd == "" {
		errorRes := response.ErrorResponse{
			Error: "Email/nickname and password are required",
		}
		c.JSON(400, errorRes)
		return
	}

	// Authenticate user
	email, err := ac.userUseCase.Login(request)
	if err != nil {
		errorRes := response.ErrorResponse{
			Error: "Authentication failed: " + err.Error(),
		}
		c.JSON(401, errorRes)
		return
	}

	token, err := ac.authMiddleware.GenerateToken(*email)
	if err != nil {
		var errorRes = response.ErrorResponse{
			Error: err.Error(),
		}
		c.JSON(500, errorRes)
		return
	}
	c.JSON(200, token)
}
