package infrastructure

import (
	"Financial/Domains/ports"
	"Financial/Models/db"
	"fmt"
	"strconv"

	"github.com/supabase-community/supabase-go"
)

const walletTable = "users"

type SupaBaseWalletRepository struct {
	client *supabase.Client
}

func NewSupaBaseWalletRepository(client *supabase.Client) ports.ExtendedRepository[db.Wallet, int] {
	return &SupaBaseWalletRepository{client: client}
}

func (repo *SupaBaseWalletRepository) Create(model *db.Wallet) (*db.Wallet, error) {
	var result db.Wallet
	_, err := repo.client.From(walletTable).
		Insert(model, false, "", "representation", "").
		Single().
		ExecuteTo(&result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *SupaBaseWalletRepository) Delete(id int) error {
	_, _, err := repo.client.From(table_string).Delete("", "").
		Single().Eq("id", strconv.Itoa(id)).Execute()
	return err
}

func (repo *SupaBaseWalletRepository) FindByField(field string, value any) (*db.Wallet, error) {
	var results []db.Wallet

	var filterValue string
	switch v := value.(type) {
	case string:
		filterValue = v
	case int, int32, int64, uint, uint32, uint64:
		filterValue = fmt.Sprintf("%d", v)
	case float32, float64:
		filterValue = fmt.Sprintf("%f", v)
	case bool:
		filterValue = strconv.FormatBool(v)
	default:
		return nil, fmt.Errorf("unsupported type for field filtering: %T", value)
	}

	// Execute the query and get the results into a slice
	_, err := repo.client.From(walletTable).
		Select("*", "exact", false).
		Filter(field, "eq", filterValue).
		ExecuteTo(&results)

	if err != nil {
		return nil, err
	}

	// If no results found, return not found error
	if len(results) == 0 {
		return nil, ErrNotFound
	}

	// Return the first result
	return &results[0], nil
}

func (repo *SupaBaseWalletRepository) GetAll() ([]db.Wallet, error) {
	var todos []db.Wallet
	_, err := repo.client.From(table_string).Select("*", "exact", false).
		ExecuteTo(&todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (repo *SupaBaseWalletRepository) GetByID(id int) (*db.Wallet, error) {
	var todo db.Wallet
	_, err := repo.client.From(table_string).Select("*", "exact", false).Eq("id", strconv.Itoa(id)).
		Single().ExecuteTo(&todo)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *SupaBaseWalletRepository) Update(todo *db.Wallet) (*db.Wallet, error) {
	var result []db.Wallet
	_, err := r.client.From(table_string).Update(todo, "", "").Eq("id", strconv.Itoa(todo.ID)).
		ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	return &result[0], nil
}

func (r *SupaBaseWalletRepository) GetUserWallet(id int, email string) (*ports.UserWallet, error) {
	var result ports.UserWallet

	_, err := r.client.From(walletTable).
		Select("id, name, type, balance, users.id, users.nickname, users.email", "1", false).
		Eq("users.email", email).
		ExecuteTo(&result)

	if err != nil {
		return nil, fmt.Errorf("error fetching wallet with user: %w", err)
	}
	return &result, nil
}
