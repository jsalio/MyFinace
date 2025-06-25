package intefaces

import (
	"Financial/Domains/ports"
	"Financial/intefaces/controllers"
	"Financial/intefaces/middleware"

	_ "Financial/docs" // This is important - points to your generated docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router         *gin.Engine
	userUseCase    ports.UserUseCase
	apiControllers []controllers.Controller
	authMiddleware *middleware.AuthMiddleware
}

func NewServer(userUseCase ports.UserUseCase) *Server {
	server := &Server{
		userUseCase:    userUseCase,
		authMiddleware: middleware.NewAuthMiddleware(),
	}
	server.setupControllers()
	server.setupRouter()
	return server
}

func (s *Server) setupControllers() {
	// Register all controllers here
	s.apiControllers = []controllers.Controller{
		controllers.NewAccountController(s.userUseCase, s.authMiddleware),
		// Add more controllers here as needed
	}
}

// server.go
func (s *Server) setupRouter() {
	s.router = gin.Default()

	// Configuración de Swagger
	url := ginSwagger.URL("/swagger/doc.json") // La URL para el archivo JSON generado
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Configuración CORS
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Grupo de rutas de la API
	api := s.router.Group("/api")
	{
		// Rutas públicas
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", func(ctx *gin.Context) {})
			authGroup.POST("/register", func(ctx *gin.Context) {})
		}

		// Rutas protegidas
		api.Use(s.authMiddleware.AuthMiddleware())
		{
			// Registrar controladores
			for _, controller := range s.apiControllers {
				controller.RegisterRoutes(api)
			}
		}
	}
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
