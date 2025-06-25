package middleware

import (
	models "Financial/Models"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	secretKey []byte
	Config    *AuthConfig
}

func NewAuthMiddleware() *AuthMiddleware {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		// Clave secreta por defecto (deberías establecer esto en tus variables de entorno)
		secretKey = "tu_clave_secreta_muy_segura"
	}
	return &AuthMiddleware{
		secretKey: []byte(secretKey),
		Config:    NewAuthConfig(),
	}
}

// SkipAuth verifica si la ruta actual está en la lista de rutas que no requieren autenticación
func (m *AuthMiddleware) SkipAuth(c *gin.Context, skipRoutes []string) bool {
	path := c.FullPath()
	method := c.Request.Method

	// Si es una solicitud OPTIONS (preflight de CORS), la dejamos pasar
	if method == "OPTIONS" {
		return true
	}

	for _, route := range skipRoutes {
		// Verificar si la ruta coincide y el método es POST (para el caso de creación de cuenta)
		if strings.HasPrefix(path, route) && (method == "POST" || route == "/swagger/") {
			return true
		}
	}
	return false
}

// AuthMiddleware verifica el token JWT en las cabeceras de la solicitud
func (m *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	// Rutas que no requieren autenticación
	// skipRoutes := []string{
	// 	"/swagger/",
	// 	"/api/auth/login",
	// 	"/api/auth/register",
	// 	"/api/account",
	// 	// Agrega aquí más rutas que no requieran autenticación
	// }

	return func(c *gin.Context) {
		// Verificar si la ruta actual está en la lista de rutas que no requieren autenticación

		var authError = models.AuthError{
			Message: "Se requiere token de autenticación",
		}

		if m.Config.IsPublicRoute(c.Request.Method, c.FullPath()) {
			c.Next()
			return
		}

		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, authError)
			c.Abort()
			return
		}

		// Eliminar el prefijo "Bearer " si está presente
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Validar el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar el algoritmo de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return m.secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			c.Abort()
			return
		}

		// Extraer claims del token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Agregar el ID de usuario al contexto para que esté disponible en los controladores
			c.Set("userID", claims["sub"])
		}

		c.Next()
	}
}

// GenerateToken genera un nuevo token JWT para un usuario
func (m *AuthMiddleware) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		// Puedes agregar más claims según sea necesario
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expira en 24 horas
	})

	// Firmar el token con la clave secreta
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
