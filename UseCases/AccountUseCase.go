package usecases

import (
	"Financial/Domains/ports"
	models "Financial/Models"
	"Financial/types"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type AccountUseCase struct {
	repository ports.Repository[models.User, int]
}

func NewAccountUseCase(repo ports.Repository[models.User, int]) ports.UserUseCase {
	return &AccountUseCase{
		repository: repo,
	}
}

func (uc *AccountUseCase) CreateAccount(nick string, email string, password string) (*models.User, error) {

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

	// Verificar si el email ya existe
	if _, err := uc.repository.FindByField("Email", email); err == nil {
		return nil, errors.New("email already exists")
	}

	account := &models.User{
		Nickname:  nick,
		FirstName: "",
		Lastname:  "",
		Email:     email,
		Status:    types.Inactive,
		CreatedAt: time.Now(),
		Password:  password,
	}
	fmt.Printf("%v", account)
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

	user, err := uc.repository.FindByField("Email", email)
	if err != nil {
		return errors.New("account not found")
	}

	return uc.repository.Delete(user.ID)
}

func (uc *AccountUseCase) UpdateAccount(req models.UpdateAccountRequest) (*models.User, error) {
	if strings.TrimSpace(req.Email) == "" {
		return nil, errors.New("email cannot be empty")
	}

	user, err := uc.repository.FindByField("Email", req.Email)
	if err != nil {
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
