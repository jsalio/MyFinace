package UseCases_test

import (
	"errors"
	"testing"

	models "Financial/Models"
	"Financial/infrastructure"
	"Financial/types"

	// "Financial/ports"
	usecases "Financial/UseCases"
	mocks "Financial/test"

	"github.com/stretchr/testify/assert"
)

func TestAccountUseCase_CreateAccount(t *testing.T) {
	tests := []struct {
		name        string
		nickname    string
		email       string
		password    string
		expectErr   bool
		expectedErr error
		setupMock   func(*mocks.MockRepository[models.User, int])
		verify      func(t *testing.T, user *models.User, err error)
	}{
		{
			name:     "successful account creation",
			nickname: "alice",
			email:    "alice@example.com",
			password: "securepassword123!",
			setupMock: func(mock *mocks.MockRepository[models.User, int]) {
				// Mock FindByField to return ErrNotFound for both email and nickname checks
				mock.SetResponse("FindByField", nil, infrastructure.ErrNotFound)
			},
			verify: func(t *testing.T, user *models.User, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "alice", user.Nickname)
				assert.Equal(t, "alice@example.com", user.Email)
				assert.NotEmpty(t, user.Password)
				assert.Equal(t, types.Inactive, user.Status)
			},
		},
		{
			name:        "empty nickname",
			nickname:    "",
			email:       "alice@example.com",
			password:    "securepassword123!",
			expectErr:   true,
			expectedErr: errors.New("nickname cannot be empty"),
		},
		{
			name:        "whitespace nickname",
			nickname:    "   ",
			email:       "alice@example.com",
			password:    "securepassword123!",
			expectErr:   true,
			expectedErr: errors.New("nickname cannot be empty"),
		},
		{
			name:        "invalid email format",
			nickname:    "alice",
			email:       "invalid-email",
			password:    "securepassword123!",
			expectErr:   true,
			expectedErr: errors.New("invalid email format"),
		},
		{
			name:        "empty email",
			nickname:    "alice",
			email:       "",
			password:    "securepassword123!",
			expectErr:   true,
			expectedErr: errors.New("email cannot be empty"),
		},
		{
			name:        "whitespace email",
			nickname:    "alice",
			email:       "   ",
			password:    "securepassword123!",
			expectErr:   true,
			expectedErr: errors.New("email cannot be empty"),
		},
		{
			name:        "empty password",
			nickname:    "alice",
			email:       "alice@example.com",
			password:    "",
			expectErr:   true,
			expectedErr: errors.New("password cannot be empty"),
		},
		{
			name:        "whitespace password",
			nickname:    "alice",
			email:       "alice@example.com",
			password:    "   ",
			expectErr:   true,
			expectedErr: errors.New("password cannot be empty"),
		},
		{
			name:     "duplicate email",
			nickname: "alice",
			email:    "duplicate@example.com",
			password: "securepassword123!",
			setupMock: func(mock *mocks.MockRepository[models.User, int]) {
				mock.SetResponse("FindByField", &models.User{Email: "duplicate@example.com"}, nil)
			},
			expectErr:   true,
			expectedErr: errors.New("email already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository[models.User, int]()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			useCase := usecases.NewAccountUseCase(repo)
			user, err := useCase.CreateAccount(tt.nickname, tt.email, tt.password)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				}
				return
			}

			if tt.verify != nil {
				tt.verify(t, user, err)
			}
		})
	}
}
