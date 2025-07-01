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
package main

import (
	"Financial/Core/ports"
	"Financial/intefaces"
	"Financial/persistence"
	"fmt"
	"log"
	"os"

	UserCases "Financial/Core/UseCases"

	"github.com/joho/godotenv"
	// "github.com/supabase-community/supabase-go"
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

	dbBoostrap, err := persistence.Init()

	if err != nil {
		fmt.Printf("Error al configurar la aplicación: %v\n", err)
		os.Exit(1)
	}

	accountUseCase := UserCases.NewAccountUseCase(dbBoostrap.AccountRepository)
	walletUseCase := UserCases.NewWalletUseCase(dbBoostrap.WalletRepository)

	// Crear e iniciar el servidor web
	server := intefaces.NewServer(accountUseCase, walletUseCase)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor iniciado en http://localhost:%s\n", port)
	if err := server.Start(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
