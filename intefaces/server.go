package intefaces

import (
	"Financial/Domains/ports"
	models "Financial/Models"
	"Financial/types"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	UserUsercase ports.UserUseCase
}

func NewServer(UserCases ports.UserUseCase) *Server {
	server := &Server{
		UserUsercase: UserCases,
	}
	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	s.router = gin.Default()

	api := s.router.Group("/api")
	{
		account := api.Group("/account")
		{
			account.POST("", s.CreateUserAccount)
			account.PUT("", s.UpdateUserAccount)
			account.DELETE("", s.DeleteUserAccount)
		}
	}
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func (s *Server) CreateUserAccount(c *gin.Context) {
	var request struct {
		Nick     string `json:"nick" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida"})
		return
	}

	account, err := s.UserUsercase.CreateAccount(request.Nick, request.Email, request.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, account)
}

func (s *Server) UpdateUserAccount(c *gin.Context) {
	var request struct {
		ID        int                 `json:"id" binding:"required"`
		FirstName string              `json:"first_name"`                      // No se necesita binding:"optional"
		LastName  string              `json:"last_name"`                       // Corregido a LastName
		Email     string              `json:"email" binding:"omitempty,email"` // Validar formato de email si se proporciona
		Status    types.AccountStatus `json:"status"`                          // Asegúrate de que types.AccountStatus sea compatible con JSON
		Password  string              `json:"password"`                        // Opcional por defecto
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida"})
		return
	}
	account, err := s.UserUsercase.UpdateAccount(models.UpdateAccountRequest{
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

func (s *Server) DeleteUserAccount(c *gin.Context) {

	var request struct {
		ID    int    `json:"id" binding:"required"`
		EMAIL string `json:"email" binding:"omitempty,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Solicitud inválida"})
		return
	}
	err := s.UserUsercase.DestroyAccount(request.EMAIL)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var response struct {
		Message string `json:"message"`
	}
	response.Message = "Account deleted"
	c.JSON(200, response)
}
