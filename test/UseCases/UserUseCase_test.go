package UseCases_test

import (
	"testing"

	"Financial/Core/Models/db"
	response "Financial/Core/Models/dtos/Response"
	usecases "Financial/Core/UseCases"
	"Financial/Core/types"
	mocks "Financial/Test"

	"github.com/stretchr/testify/assert"
)

func TestAccountUseCase_CreateAccount(t *testing.T) {
	tests := []struct {
		name        string
		nickname    string
		email       string
		password    string
		expectErr   bool
		expectedErr interface{}
		setupMock   func(*mocks.MockRepository[db.User, int])
		verify      func(t *testing.T, user *response.CreateAccountResponse, err error)
	}{
		{
			name:     "successful account creation",
			nickname: "alice_serat",
			email:    "alice@example.com",
			password: "securepassword123!",
			setupMock: func(mock *mocks.MockRepository[db.User, int]) {
				mock.SetResponse("FindByField", nil, types.ErrNotFound)
				mock.SetFindByFieldNotExists(true)
				mock.SetResponse("FindByField", nil, nil)
				mock.SetResponse("FindByField", nil, nil)
			},
			verify: func(t *testing.T, user *response.CreateAccountResponse, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "alice_serat", user.Nick)
				assert.Equal(t, "alice@example.com", user.Email)
			},
		},
		{
			name:      "empty nickname",
			nickname:  "",
			email:     "alice@example.com",
			password:  "securepassword123!",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Nickname Is Empty.",
			},
		},
		{
			name:      "whitespace nickname",
			nickname:  "   ",
			email:     "alice@example.com",
			password:  "securepassword123!",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Nickname contains space.",
			},
		},
		{
			name:      "Nickname length required",
			nickname:  "   ",
			email:     "alice@example.com",
			password:  "securepassword123!",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Nickname not have length.",
			},
		},
		{
			name:      "invalid email format",
			nickname:  "alice",
			email:     "invalid-email",
			password:  "securepassword123!",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Email not match.",
			},
		},
		{
			name:      "Valid email should pass",
			nickname:  "alice_morat",
			email:     "my_personal@mail.com",
			password:  "securepassword123!",
			expectErr: true,
			setupMock: func(mock *mocks.MockRepository[db.User, int]) {

				mock.SetFindByFieldNotExists(true)
				mock.SetResponse("FindByField", nil, nil)
				mock.SetResponse("FindByField", nil, nil)
			},
		},
		{
			name:      "empty email",
			nickname:  "alice",
			email:     "",
			password:  "securepassword123!",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Email Is Empty.",
			},
		},
		{
			name:      "whitespace email",
			nickname:  "alice",
			email:     "   ",
			password:  "securepassword123!",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Email not match.",
			},
		},
		{
			name:      "Math with lenght required",
			nickname:  "alice",
			email:     "1@mail.com",
			password:  "securepassword123!",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Email not have length.",
			},
		},
		{
			name:      "Match with lenght required",
			nickname:  "alice",
			email:     "my_personal@mail.com",
			password:  "securepassword123!",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Email already exists.",
			},
			setupMock: func(mock *mocks.MockRepository[db.User, int]) {
				mock.SetResponse("FindByField", &db.User{Email: "my_personal@mail.com"}, nil)
			},
		},
		{
			name:      "empty password",
			nickname:  "alice",
			email:     "alice@example.com",
			password:  "",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Password Is Empty.",
			},
		},
		{
			name:      "whitespace password",
			nickname:  "alice",
			email:     "alice@example.com",
			password:  "1   3",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Password contains space.",
			},
		},
		{
			name:      "Password length required",
			nickname:  "alice",
			email:     "alice@example.com",
			password:  "1   3",
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Password not have length.",
			},
		},
		{
			name:     "duplicate email",
			nickname: "alice",
			email:    "duplicate@example.com",
			password: "securepassword123!",
			setupMock: func(mock *mocks.MockRepository[db.User, int]) {
				mock.SetResponse("FindByField", &db.User{Email: "duplicate@example.com"}, nil)
			},
			expectErr: true,
			expectedErr: response.ErrorResponse{
				Error: "Email already exists.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository[db.User, int]()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			useCase := usecases.NewAccountUseCase(repo)
			newUser, err := useCase.CreateAccount(tt.nickname, tt.email, tt.password)

			if tt.expectErr {
				assert.IsType(t, &[]response.ErrorResponse{}, err)
				if tt.expectedErr != nil {
					assert.Contains(t, *err, tt.expectedErr)
				}
				return
			} else {
				if err != nil {
					if len(*err) > 0 {
						assert.Fail(t, "This test suppouse not have errors")
					}
				}
			}

			if tt.verify != nil {
				if newUser == nil {
					t.Error("Expected newUser to not be nil")
					return
				}
				tt.verify(t, newUser.Data, nil)
			}
		})
	}
}
