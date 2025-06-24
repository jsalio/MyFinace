package intefaces

import (
	"Financial/Domains/ports"
	"Financial/intefaces/controllers"

	_ "Financial/docs" // This is important - points to your generated docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router         *gin.Engine
	userUseCase    ports.UserUseCase
	apiControllers []controllers.Controller
}

func NewServer(userUseCase ports.UserUseCase) *Server {
	server := &Server{
		userUseCase: userUseCase,
	}
	server.setupControllers()
	server.setupRouter()
	return server
}

func (s *Server) setupControllers() {
	// Register all controllers here
	s.apiControllers = []controllers.Controller{
		controllers.NewAccountController(s.userUseCase),
		// Add more controllers here as needed
	}
}

func (s *Server) setupRouter() {
	s.router = gin.Default()
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// API v1 routes
	api := s.router.Group("/api")
	{
		// Register all controllers
		for _, controller := range s.apiControllers {
			controller.RegisterRoutes(api)
		}
	}
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
