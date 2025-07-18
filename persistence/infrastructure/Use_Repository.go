package infrastructure

import (
	"Financial/Core/Models/db"
	contracts "Financial/Core/ports"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

const table_string = "users"

type SupaBaseUserRepository struct {
	client *supabase.Client
}

func NewSupaBaseUserRepository(client *supabase.Client) contracts.Repository[db.User, int] {
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

func (repo *SupaBaseUserRepository) Create(model *db.User) (*db.User, error) {
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

	var result db.User
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

func (repo *SupaBaseUserRepository) FindByField(field string, value any) (*db.User, error) {
	var results []db.User

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

func (repo *SupaBaseUserRepository) GetAll() ([]db.User, error) {
	var todos []db.User
	_, err := repo.client.From(table_string).Select("*", "exact", false).
		ExecuteTo(&todos)
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (repo *SupaBaseUserRepository) GetByID(id int) (*db.User, error) {
	var todo db.User
	_, err := repo.client.From(table_string).Select("*", "exact", false).Eq("id", strconv.Itoa(id)).
		Single().ExecuteTo(&todo)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *SupaBaseUserRepository) Update(todo *db.User) (*db.User, error) {
	var result []db.User
	_, err := r.client.From(table_string).Update(todo, "", "").Eq("id", strconv.Itoa(todo.ID)).
		ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	return &result[0], nil
}

// Query executes a custom query and returns the result as interface{}.
// This method provides a flexible way to execute custom queries that don't fit the standard CRUD operations.
func (r *SupaBaseUserRepository) Query(fields string, args contracts.QueryOptions) (interface{}, error) {
	var user []db.User

	query := r.client.From(table_string)
	queryUnfilter := query.Select(fields, "", false)

	for _, filter := range args.Filters {
		if filter.Operator == "eq" {
			queryUnfilter.Eq(filter.Field, filter.Value.(string))
		}
		if filter.Operator == "neq" {
			queryUnfilter.Neq(filter.Field, filter.Value.(string))
		}
		if filter.Operator == "gt" {
			queryUnfilter.Gt(filter.Field, filter.Value.(string))
		}
		if filter.Operator == "gte" {
			queryUnfilter.Gte(filter.Field, filter.Value.(string))
		}
		if filter.Operator == "lt" {
			queryUnfilter.Lt(filter.Field, filter.Value.(string))
		}
		if filter.Operator == "lte" {
			queryUnfilter.Lte(filter.Field, filter.Value.(string))
		}
		if filter.Operator == "like" {
			queryUnfilter.Like(filter.Field, filter.Value.(string))
		}
		if filter.Operator == "ilike" {
			queryUnfilter.Ilike(filter.Field, filter.Value.(string))
		}
		if filter.Operator == "is" {
			queryUnfilter.Is(filter.Field, filter.Value.(string))
		}
		if filter.Operator == "in" {
			queryUnfilter.In(filter.Field, []string{filter.Value.(string)})
		}
	}

	for _, order := range args.OrderBy {
		queryUnfilter.Order(order.Field, &postgrest.OrderOpts{
			Ascending:  order.Ascending,
			NullsFirst: *order.NullsFirst,
		})
	}

	_, err := queryUnfilter.ExecuteTo(&user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
