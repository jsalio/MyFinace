package controllers

import (
	"Financial/Domains/ports"
	models "Financial/Models"
	"Financial/types"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	*BaseController
	userUseCase ports.UserUseCase
}

func NewAccountController(userUseCase ports.UserUseCase) *AccountController {
	return &AccountController{
		BaseController: NewBaseController("/account"),
		userUseCase:   userUseCase,
	}
}

func (ac *AccountController) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group(ac.Path)
	{
		group.POST("", ac.CreateUserAccount)
		group.PUT("", ac.UpdateUserAccount)
		group.DELETE("", ac.DeleteUserAccount)
	}
}

func (ac *AccountController) CreateUserAccount(c *gin.Context) {
	var request struct {
		Nick     string `json:"nick" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida"})
		return
	}

	account, err := ac.userUseCase.CreateAccount(request.Nick, request.Email, request.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, account)
}

func (ac *AccountController) UpdateUserAccount(c *gin.Context) {
	var request struct {
		ID        int                 `json:"id" binding:"required"`
		FirstName string              `json:"first_name"`
		LastName  string              `json:"last_name"`
		Email     string              `json:"email" binding:"omitempty,email"`
		Status    types.AccountStatus `json:"status"`
		Password  string              `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida"})
		return
	}

	account, err := ac.userUseCase.UpdateAccount(models.UpdateAccountRequest{
		ID:        request.ID,
		FirstName: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Status:    request.Status,
		Password:  request.Password,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, account)
}

func (ac *AccountController) DeleteUserAccount(c *gin.Context) {
	var request struct {
		ID    int    `json:"id" binding:"required"`
		EMAIL string `json:"email" binding:"omitempty,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida"})
		return
	}

	err := ac.userUseCase.DestroyAccount(request.EMAIL)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Account deleted"})
}
