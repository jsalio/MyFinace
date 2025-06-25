package middleware

import (
	"strings"
	"sync"
)

type AuthConfig struct {
	PublicRoutes []string
	mu           sync.Mutex
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		PublicRoutes: []string{
			"GET_/swagger/",
			"POST_/api/auth/login",
			"POST_/api/auth/register",
		},
	}
}

// AddPublicRoute agrega una ruta pública
func (ac *AuthConfig) AddPublicRoute(method, path string) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	route := strings.ToUpper(method) + "_" + path
	for _, r := range ac.PublicRoutes {
		if r == route {
			return // Ya existe
		}
	}
	ac.PublicRoutes = append(ac.PublicRoutes, route)
}

// IsPublicRoute verifica si una ruta es pública
func (ac *AuthConfig) IsPublicRoute(method, path string) bool {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	method = strings.ToUpper(method)
	for _, route := range ac.PublicRoutes {
		parts := strings.SplitN(route, "_", 2)
		if len(parts) != 2 {
			continue
		}

		routeMethod := parts[0]
		routePath := parts[1]

		// Si el método es OPTIONS (preflight CORS), siempre permitir
		if method == "OPTIONS" {
			return true
		}

		// Si el método coincide y la ruta coincide con el prefijo
		if (routeMethod == method || routeMethod == "ANY") &&
			(path == routePath || strings.HasPrefix(path, routePath+"/")) {
			return true
		}
	}
	return false
}
