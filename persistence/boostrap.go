package persistence

import (
	"fmt"
	"os"

	"Financial/Core/Models/db"
	port "Financial/Core/ports"
	"Financial/persistence/infrastructure"

	"github.com/supabase-community/supabase-go"
)

type DbBoostrap struct {
	AccountRepository port.Repository[db.User, int]
	WalletRepository  port.Repository[db.Wallet, int]
}

func Init() (*DbBoostrap, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		return nil, fmt.Errorf("error inicializando cliente de Supabase: %w", err)
	}

	return &DbBoostrap{
		AccountRepository: infrastructure.NewSupaBaseUserRepository(client),
		WalletRepository:  infrastructure.NewSupaBaseWalletRepository(client),
	}, nil
}
