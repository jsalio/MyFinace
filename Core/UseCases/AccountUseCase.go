package usecases

import (
	"Financial/Core/Models/db"
	dtos "Financial/Core/Models/dtos/Request"
	response "Financial/Core/Models/dtos/Response"
	"Financial/Core/ports"
	"Financial/Core/types"
	"Financial/Core/validators"

	"errors"
	"fmt"
	"time"
)

type AccountUseCase struct {
	repository ports.Repository[db.User, int]
}

func NewAccountUseCase(repo ports.Repository[db.User, int]) ports.UserUseCase {
	return &AccountUseCase{
		repository: repo,
	}
}

func (uc *AccountUseCase) validateEmailUniqueness(email string, v *validators.Validator) {
	if email == "" {
		return
	}

	_, err := uc.repository.FindByField("email", email)
	if err == nil {
		v.AddError(validators.ErrEmailExists)
	} else if err != types.ErrNotFound {
		v.AddError(fmt.Sprintf("error checking email existence: %v", err))
	}
}

func (uc *AccountUseCase) validateNickUniqueness(nick string, v *validators.Validator) {
	if nick == "" {
		return
	}

	_, err := uc.repository.FindByField("nick_name", nick)
	if err == nil {
		v.AddError(validators.ErrNickExists)
	} else if err != types.ErrNotFound {
		v.AddError(fmt.Sprintf("error checking nick existence: %v", err))
	}
}

func (uc *AccountUseCase) validateAndGetUser(email string, v *[]response.ErrorResponse) *db.User {
	user, err := uc.repository.FindByField("email", email)
	if err != nil {
		*v = append(*v, response.ErrorResponse{
			Error: "User not found",
		})
		return nil
	}
	return user
}

func (uc *AccountUseCase) CreateAccount(nick string, email string, password string) (*response.SuccessResponse[*response.CreateAccountResponse], *[]response.ErrorResponse) {

	validationsError := []response.ErrorResponse{}
	validator := validators.CreateAccountValidator(dtos.CreateAccountRequest{
		Nick:     nick,
		Email:    email,
		Password: password,
	}, uc.repository)

	if len(validator.Errors) > 0 {
		for _, err := range validator.Errors {
			validationsError = append(validationsError, response.ErrorResponse{
				Error: fmt.Sprintf("%s.", err.Message),
			})
		}
		return nil, &validationsError
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
	result, error := uc.repository.Create(account)

	if error != nil {
		validationsError = append(validationsError, response.ErrorResponse{
			Error: fmt.Sprintf("error checking nick existence: %w", error),
		})
		return nil, &validationsError
	}

	data := &response.CreateAccountResponse{
		ID:    result.ID,
		Nick:  result.Nickname,
		Email: result.Email,
	}

	return &response.SuccessResponse[*response.CreateAccountResponse]{
		Message: "",
		Data:    data,
	}, nil
}

func (uc *AccountUseCase) DestroyAccount(email string) *[]response.ErrorResponse {
	validationsError := []response.ErrorResponse{}
	validator := validators.DestroidAccountValidator(email, uc.repository)

	if len(validator.Errors) > 0 {
		for _, err := range validator.Errors {
			validationsError = append(validationsError, response.ErrorResponse{
				Error: fmt.Sprintf("%s.", err.Message),
			})
		}
		return &validationsError
	}

	user := uc.validateAndGetUser(email, &validationsError)
	if len(validationsError) > 0 {
		return &validationsError
	}

	err := uc.repository.Delete(user.ID)

	if err != nil {
		validationsError = append(validationsError, response.ErrorResponse{
			Error: err.Error(),
		})
		return &validationsError
	}

	return nil
}

func (uc *AccountUseCase) UpdateAccount(req db.UpdateAccountRequest) (*response.SuccessResponse[*response.UpdateAccountResponse], *[]response.ErrorResponse) {

	validationsError := []response.ErrorResponse{}
	validator := validators.UpdateAccountValidator(req, uc.repository)

	if len(validator.Errors) > 0 {
		for _, err := range validator.Errors {
			validationsError = append(validationsError, response.ErrorResponse{
				Error: fmt.Sprintf("%s.", err.Message),
			})
		}
		return nil, &validationsError
	}

	user := uc.validateAndGetUser(req.Email, &validationsError)
	if len(validationsError) > 0 {
		return nil, &validationsError
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
		//return nil, errors.New("user id and Mail not match")
	}

	if !updated {
		return &response.SuccessResponse[*response.UpdateAccountResponse]{
			Message: "NoChanges",
			Data: &response.UpdateAccountResponse{
				ID:    user.ID,
				Email: user.Email,
			},
		}, nil // No hay cambios, devolver el usuario sin actualizar
	}

	data, error := uc.repository.Update(user)

	if error != nil {
		validationsError = append(validationsError, response.ErrorResponse{
			Error: error.Error(),
		})
		return nil, &validationsError
	}

	return &response.SuccessResponse[*response.UpdateAccountResponse]{
		Message: "Updated",
		Data: &response.UpdateAccountResponse{
			ID:    data.ID,
			Email: data.Email,
		},
	}, nil

}

func (uc *AccountUseCase) Login(auth dtos.AuthRequest) (*string, error) {
	v := validators.NewValidator()
	if auth.Email == "" && auth.Nickname == "" {
		v.AddError("nick or email can't be empty")
	}

	if auth.Passwd == "" {
		v.AddError(validators.ErrPasswordRequired)
	}

	if !v.IsValid() {
		return nil, errors.New(v.Error())
	}

	data, err := uc.repository.Query("email, password", ports.QueryOptions{
		Filters: []ports.Filter{
			{
				Field:    "email",
				Operator: "eq",
				Value:    auth.Email,
			},
			{
				Field:    "password",
				Operator: "eq",
				Value:    auth.Passwd,
			},
		},
	})

	if err != nil {
		return nil, errors.New("account not found")
	}

	user, ok := data.([]db.User)

	if !ok {
		return nil, errors.New("account not found")
	}

	return &user[0].Email, nil
}
