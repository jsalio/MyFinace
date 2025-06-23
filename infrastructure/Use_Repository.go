package infrastructure

import (
	"Financial/Domains/ports"
	models "Financial/Models"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/supabase-community/supabase-go"
)

const table_string = "users"

type SupaBaseUserRepository struct {
	client *supabase.Client
}

func NewSupaBaseUserRepository(client *supabase.Client) ports.Repository[models.User, int] {
	return &SupaBaseUserRepository{client: client}
}

// CreateUser is a helper struct that matches the database schema
type CreateUser struct {
	Nickname  string    `json:"nick_name"`
	FirstName string    `json:"first_name"`
	Lastname  string    `json:"last_name"`
	Email     string    `json:"email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Password  string    `json:"password"`
}

func (repo *SupaBaseUserRepository) Create(model *models.User) (*models.User, error) {
	// Create a new user without the ID field
	newUser := CreateUser{
		Nickname:  model.Nickname,
		FirstName: model.FirstName,
		Lastname:  model.Lastname,
		Email:     model.Email,
		Status:    string(model.Status),
		CreatedAt: model.CreatedAt,
		Password:  model.Password,
	}

	var result models.User
	_, err := repo.client.From(table_string).
		Insert(newUser, false, "", "representation", "").
		Single().
		ExecuteTo(&result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *SupaBaseUserRepository) Delete(id int) error {
	_, _, err := repo.client.From(table_string).Delete("", "").
		Single().Eq("id", strconv.Itoa(id)).Execute()
	return err
}

// ErrNotFound is returned when a record is not found
var ErrNotFound = errors.New("record not found")

func (repo *SupaBaseUserRepository) FindByField(field string, value any) (*models.User, error) {
	var results []models.User

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
	_, err := repo.client.From(table_string).
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

func (repo *SupaBaseUserRepository) GetAll() ([]models.User, error) {
	var todos []models.User
	_, err := repo.client.From(table_string).Select("*", "exact", false).
		ExecuteTo(&todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (repo *SupaBaseUserRepository) GetByID(id int) (*models.User, error) {
	var todo models.User
	_, err := repo.client.From(table_string).Select("*", "exact", false).Eq("id", strconv.Itoa(id)).
		Single().ExecuteTo(&todo)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *SupaBaseUserRepository) Update(todo *models.User) (*models.User, error) {
	var result models.User
	_, err := r.client.From(table_string).Update(todo, "", "").Eq("id", strconv.Itoa(todo.ID)).
		ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
