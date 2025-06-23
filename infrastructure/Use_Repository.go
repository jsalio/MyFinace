package infrastructure

import (
	"Financial/Domains/ports"
	models "Financial/Models"
	"strconv"

	"github.com/supabase-community/supabase-go"
)

const table_string = "users"

type SupaBaseUserRepository struct {
	client *supabase.Client
}

func NewSupaBaseUserRepository(client *supabase.Client) ports.Repository[models.User, int] {
	return &SupaBaseUserRepository{client: client}
}

func (repo *SupaBaseUserRepository) Create(model *models.User) (*models.User, error) {
	var result models.User

	insertData := map[string]interface{}{}
	_, err := repo.client.From(table_string).
		Insert(insertData, false, "", "representation", "").
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

func (repo *SupaBaseUserRepository) FindByField(field string, value any) (*models.User, error) {
	var target models.User

	_, err := repo.client.From(table_string).Select("*", "", false).
		Filter(field, "Equal to", string(value)).
		ExecuteTo(target)

	if err != nil {
		return nil, err
	}
	return &target, nil
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
