package intefaces

import (
	"Financial/Domains/ports"

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
			account.PUT("", func(ctx *gin.Context) {})
			account.DELETE("", func(ctx *gin.Context) {})
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

	c.JSON(200, nil)
}

func (s *Server) DeleteUserAccount(c *gin.Context) {

	c.JSON(200, nil)
}
