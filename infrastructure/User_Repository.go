// package infrastructure

// import (
// 	"context"
// 	"log"

// 	"Financial/Domains/ports"
// 	"Financial/Models"
// 	"github.com/nedpals/supabase-go"
// )

// type SupaBaseUserRepository struct {
// 	client *supabase.Client
// }

// func NewSupaBaseUserRepository(supabaseURL, supabaseKey string) *SupaBaseUserRepository {
// 	client := supabase.CreateClient(supabaseURL, supabaseKey, false)
// 	return &SupaBaseUserRepository{
// 		client: client,
// 	}
// }

// // Implement the Repository[Models.User, int] interface methods below
// // For example:
// func (r *SupaBaseUserRepository) Create(ctx context.Context, user *Models.User) (*Models.User, error) {
// 	// Implementation for creating a user in Supabase
// 	return nil, nil
// }

// // Add other required methods from the Repository interface
// // Get, Update, Delete, etc.
