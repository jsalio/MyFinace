package main

import (
	"Financial/Domains/ports"
	"fmt"
	"log"
	"os"

	UserCases "Financial/UseCases"
	SupaBaseUserRepository "Financial/infrastructure"
	"Financial/intefaces"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

var (
	supabaseURL  string
	supabaseKey  string
	UserUseCases ports.UserUseCase
)

func PassPrerequirements() bool {
	errEnv := godotenv.Load()

	if errEnv != nil {
		fmt.Printf("Error cargando archivo .env: %v\n", errEnv)
		return false
	}

	supabaseURL = os.Getenv("SUPABASE_URL")
	supabaseKey = os.Getenv("SUPABASE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		fmt.Println("Error: SUPABASE_URL o SUPABASE_KEY no están definidas en el archivo .env")
		return false
	}
	return true
}

func main() {
	passRequirements := PassPrerequirements()
	if !passRequirements {
		os.Exit(1)
	}

	// Inicializar cliente de Supabase
	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		fmt.Printf("Error al inicializar cliente de Supabase: %v\n", err)
		os.Exit(1)
	}

	repo := SupaBaseUserRepository.NewSupaBaseUserRepository(client)
	todoUseCase := UserCases.NewAccountUseCase(repo)

	// Crear e iniciar el servidor web
	server := intefaces.NewServer(todoUseCase)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor iniciado en http://localhost:%s\n", port)
	if err := server.Start(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
