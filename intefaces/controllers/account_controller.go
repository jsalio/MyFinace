package controllers

import (
	"Financial/Domains/ports"
	models "Financial/Models"
	"Financial/intefaces/controllers/dtos"
	"Financial/types"

	"github.com/gin-gonic/gin"
)

// @title Financial App API
// 	@version 1.0
// 	@description This is a financial application server.
// 	@termsOfService http://swagger.io/terms/
// 	@contact.name API Support
// 	@contact.email support@financialapp.com
// 	@license.name Apache 2.0
// 	@license.url http://www.apache.org/licenses/LICENSE-2.0.html
// 	@host localhost:8080
// 	@BasePath /api

// AccountController handles HTTP requests related to account operations
//
//	@Summary Account operations
//	@Description Endpoints for managing user accounts
//	@Tags Account
//	@Produce json
//	@Consume json
type AccountController struct {
	*BaseController
	userUseCase ports.UserUseCase
}

func NewAccountController(userUseCase ports.UserUseCase) *AccountController {
	return &AccountController{
		BaseController: NewBaseController("/account"),
		userUseCase:    userUseCase,
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

// CreateUserAccount crea una nueva cuenta de usuario
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario con la información proporcionada
// @Tags usuarios
// @Accept json
// @Produce json
// @param request body dtos.CreateAccountRequest true ""
// @Success 200 "Usuario creado exitosamente"
// @Failure 400 {object} map[string]string "Error en la solicitud"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /account [post]
func (ac *AccountController) CreateUserAccount(c *gin.Context) {
	var request dtos.CreateAccountRequest

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

// UpdateUserAccount actualiza la información de un usuario existente
// @Summary Actualizar usuario
// @Description Actualiza la información de un usuario existente
// @Tags usuarios
// @Accept json
// @Produce json
// @Param request body dtos.UpdateAccountRequest true "Datos actualizados del usuario"
// @Success 200 "Usuario actualizado exitosamente"
// @Failure 400 {object} map[string]string "Error en la solicitud"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /account [put]
func (ac *AccountController) UpdateUserAccount(c *gin.Context) {
	var request dtos.UpdateAccountRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida"})
		return
	}

	account, err := ac.userUseCase.UpdateAccount(models.UpdateAccountRequest{
		ID:        request.ID,
		FirstName: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Status:    types.AccountStatus(request.Status),
		Password:  request.Password,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, account)
}

// DeleteUserAccount elimina una cuenta de usuario
// @Summary Eliminar usuario
// @Description Elimina un usuario existente por su email
// @Tags usuarios
// @Accept json
// @Produce json
// @Param request body dtos.DeleteAccountRequest true "Email del usuario a eliminar"
// @Success 200 {object} map[string]string "Mensaje de éxito"
// @Failure 400 {object} map[string]string "Error en la solicitud"
// @Failure 500 {object} map[string]string "Error interno del servidor"
// @Router /account [delete]
func (ac *AccountController) DeleteUserAccount(c *gin.Context) {
	var request dtos.DeleteAccountRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida"})
		return
	}

	err := ac.userUseCase.DestroyAccount(request.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Account deleted"})
}
