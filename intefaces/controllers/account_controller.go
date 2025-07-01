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
	"Financial/Core/Models/db"
	request "Financial/Core/Models/dtos/Request"
	response "Financial/Core/Models/dtos/Response"
	contracts "Financial/Core/ports"
	types "Financial/Core/types"
	"Financial/intefaces/middleware"

	"github.com/gin-gonic/gin"
)

// AccountController handles HTTP requests related to account operations
//
//	@Summary Account operations
//	@Description Endpoints for managing user accounts
//	@Tags Account
//	@Produce json
//	@Consume json
type AccountController struct {
	*BaseController
	userUseCase    contracts.UserUseCase
	authMiddleware *middleware.AuthMiddleware
}

func NewAccountController(userUseCase contracts.UserUseCase, authMiddlerware *middleware.AuthMiddleware) *AccountController {
	return &AccountController{
		BaseController: NewBaseController("/account"),
		userUseCase:    userUseCase,
		authMiddleware: authMiddlerware,
	}
}

func (ac *AccountController) RegisterRoutes(router *gin.RouterGroup) {
	ac.authMiddleware.Config.AddPublicRoute("POST", "/api/account")

	public := router.Group("/account")
	{
		public.POST("", ac.CreateUserAccount)
	}

	protected := router.Group("/account")
	protected.Use(ac.authMiddleware.AuthMiddleware())
	{
		protected.PUT("", ac.UpdateUserAccount)
		protected.DELETE("", ac.DeleteUserAccount)
	}
}

// CreateUserAccount crea una nueva cuenta de usuario
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario con la información proporcionada
// @Tags Account
// @Accept json
// @Produce json
// @param request body dtos.CreateAccountRequest true "Datos del usuario nuevo"
// @Success 200 {object} dtos.CreateAccountResponse "Usuario creado exitosamente"
// @Failure 400 {object} dtos.ErrorResponse "Error en la solicitud"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /account [post]
func (ac *AccountController) CreateUserAccount(c *gin.Context) {
	var request request.CreateAccountRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		var errorRes = response.ErrorResponse{
			Error: "Solicitud inválida",
		}
		c.JSON(400, errorRes)
		return
	}
	account, err := ac.userUseCase.CreateAccount(request.Nick, request.Email, request.Password)
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, account)
}

// UpdateUserAccount actualiza la información de un usuario existente
// @Summary Actualizar usuario
// @Description Actualiza la información de un usuario existente
// @Tags Account
// @Accept json
// @Produce json
// @Param request body dtos.UpdateAccountRequest true "Datos actualizados del usuario"
// @Success 200 {object} dtos.UpdateAccountResponse "Usuario actualizado exitosamente"
// @Failure 400 {object} dtos.ErrorResponse "Error en la solicitud"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /account [put]
func (ac *AccountController) UpdateUserAccount(c *gin.Context) {
	var request request.UpdateAccountRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		var errorRes = response.ErrorResponse{
			Error: "Solicitud inválida",
		}
		c.JSON(400, errorRes)
		return
	}

	account, err := ac.userUseCase.UpdateAccount(db.UpdateAccountRequest{
		ID:        request.ID,
		FirstName: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Status:    types.AccountStatus(request.Status),
		Password:  request.Password,
	})

	if err != nil {

		c.JSON(500, err)
		return
	}
	c.JSON(200, account)
}

// DeleteUserAccount elimina una cuenta de usuario
// @Summary Eliminar usuario
// @Description Elimina un usuario existente por su email
// @Tags Account
// @Accept json
// @Produce json
// @Param request body dtos.DeleteAccountRequest true "Email del usuario a eliminar"
// @Success 200 {object} map[string]string "Mensaje de éxito"
// @Failure 400 {object} dtos.ErrorResponse "Error en la solicitud"
// @Failure 500 {object} dtos.ErrorResponse "Error interno del servidor"
// @Router /account [delete]
func (ac *AccountController) DeleteUserAccount(c *gin.Context) {
	var request request.DeleteAccountRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		var errorRes = response.ErrorResponse{
			Error: "Solicitud inválida",
		}
		c.JSON(400, errorRes)
		return
	}

	err := ac.userUseCase.DestroyAccount(request.Email)
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, gin.H{"message": "Account deleted"})
}
