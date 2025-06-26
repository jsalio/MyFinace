package usecases

import (
	"Financial/Domains/ports"
	"Financial/Models/db"
	"Financial/infrastructure"
	"Financial/types"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

var ErrNotFound = infrastructure.ErrNotFound

type AccountUseCase struct {
	repository ports.Repository[db.User, int]
}

func NewAccountUseCase(repo ports.Repository[db.User, int]) ports.UserUseCase {
	return &AccountUseCase{
		repository: repo,
	}
}

func (uc *AccountUseCase) CreateAccount(nick string, email string, password string) (*db.User, error) {

	// Validaciones de entrada
	if strings.TrimSpace(nick) == "" {
		return nil, errors.New("nickname cannot be empty")
	}
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("email cannot be empty")
	}
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("password cannot be empty")
	}

	if !isValidEmail(email) {
		return nil, errors.New("invalid email format")
	}

	// Check if email already exists
	_, err := uc.repository.FindByField("email", email)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if err != ErrNotFound {
		return nil, fmt.Errorf("error checking email existence: %w", err)
	}

	_, err_nick := uc.repository.FindByField("nick_name", nick)
	if err_nick == nil {
		return nil, errors.New("nickname already exists")
	} else if err != ErrNotFound {
		return nil, fmt.Errorf("error checking nick existence: %w", err)
	}

	account := &db.User{
		Nickname:  nick,
		FirstName: "",
		Lastname:  "",
		Email:     email,
		Status:    types.Inactive,
		CreatedAt: time.Now(),
		Password:  password,
	}
	return uc.repository.Create(account)
}

func isValidEmail(email string) bool {
	// Basic email regex: allows alphanumeric, dots, and common domains
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func (uc *AccountUseCase) DestroyAccount(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("email cannot be empty")
	}

	user, err := uc.repository.FindByField("email", email)
	if err != nil {
		return errors.New("account not found")
	}

	return uc.repository.Delete(user.ID)
}

func (uc *AccountUseCase) UpdateAccount(req db.UpdateAccountRequest) (*db.User, error) {
	fmt.Printf("%v", req)
	if strings.TrimSpace(req.Email) == "" {
		return nil, errors.New("email cannot be empty")
	}

	user, err := uc.repository.FindByField("email", req.Email)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, errors.New("account not found")
	}

	// Actualizar solo los campos proporcionados
	updated := false
	if req.FirstName != "" {
		user.FirstName = req.FirstName
		updated = true
	}
	if req.Lastname != "" {
		user.Lastname = req.Lastname
		updated = true
	}
	if req.Password != "" {
		user.Password = req.Password // En un caso real, hashear la contraseña
		updated = true
	}
	if req.Status != "" {
		user.Status = req.Status
		updated = true
	}

	if user.ID != req.ID {
		return nil, errors.New("user id and Mail not match")
	}

	// Validar si el nuevo email ya existe (si se proporciona)
	if req.Email != "" && req.Email != user.Email {
		if _, err := uc.repository.FindByField("Email", req.Email); err == nil {
			return nil, errors.New("new email already exists")
		}
		user.Email = req.Email
		updated = true
	}

	if !updated {
		return user, nil // No hay cambios, devolver el usuario sin actualizar
	}

	return uc.repository.Update(user)
}
