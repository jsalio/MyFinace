package UseCases_test

import (
	"errors"
	"testing"

	// "Financial/Domains/ports"

	"Financial/Domains/ports"
	"Financial/Models/db"
	"Financial/Models/dtos"
	"Financial/infrastructure"
	"Financial/types"

	// "Financial/UseCases"
	usecases "Financial/UseCases"
	mocks "Financial/test"

	"github.com/stretchr/testify/assert"
)

func TestWalletUseCase_CreateWallet(t *testing.T) {
	tests := []mocks.TestSetup[dtos.CreateWalletRequest, db.Wallet, int]{
		{
			Name: "successful wallet creation",
			Req: dtos.CreateWalletRequest{
				Name:       "Savings",
				WalletType: "savings",
				Balance:    1000.0,
				UserID:     1,
			},
			SetupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("GetAll", []db.Wallet{}, nil)
				mock.SetResponse("Create", &db.Wallet{
					ID:      1,
					Name:    "Savings",
					Type:    types.Debit,
					Balance: 1000.0,
					UserID:  1,
				}, nil)
			},
			Verify: func(t *testing.T, wallet *db.Wallet, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Savings", wallet.Name)
				assert.Equal(t, types.Debit, wallet.Type)
				assert.Equal(t, 1000.0, wallet.Balance)
				assert.Equal(t, 1, wallet.UserID)
			},
		},
		{
			Name: "empty wallet name",
			Req: dtos.CreateWalletRequest{
				Name:       "",
				WalletType: "savings",
				Balance:    1000.0,
				UserID:     1,
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("wallet name cannot be empty"),
		},
		{
			Name: "whitespace wallet name",
			Req: dtos.CreateWalletRequest{
				Name:       "  ",
				WalletType: "savings",
				Balance:    1000.0,
				UserID:     1,
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("wallet name cannot be empty"),
		},
		{
			Name: "negative balance",
			Req: dtos.CreateWalletRequest{
				Name:       "Savings",
				WalletType: types.Debit,
				Balance:    -100.0,
				UserID:     1,
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("initial balance cannot be negative"),
		},
		{
			Name: "invalid user ID",
			Req: dtos.CreateWalletRequest{
				Name:       "Savings",
				WalletType: "savings",
				Balance:    1000.0,
				UserID:     0,
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("invalid user ID"),
		},
		{
			Name: "duplicate wallet name for user",
			Req: dtos.CreateWalletRequest{
				Name:       "Savings",
				WalletType: "savings",
				Balance:    1000.0,
				UserID:     1,
			},
			SetupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("GetAll", []db.Wallet{{
					ID:      1,
					Name:    "Savings",
					Type:    "savings",
					Balance: 1000.0,
					UserID:  1,
				}}, nil)
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("a wallet with this name already exists for this user"),
		},
		{
			Name: "error checking wallet existence",
			Req: dtos.CreateWalletRequest{
				Name:       "Savings",
				WalletType: types.Debit,
				Balance:    1000.0,
				UserID:     1,
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("error checking wallet existence"),
			SetupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("GetAll", nil, errors.New("database error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			repo := mocks.NewMockRepository[db.Wallet, int]()
			if tt.SetupMock != nil {
				tt.SetupMock(repo)
			}

			useCase := usecases.NewWalletUseCase(repo)
			wallet, err := useCase.CreateWallet(tt.Req)

			if tt.ExpectErr {
				assert.Error(t, err)
				if tt.ExpectedErr != nil {
					assert.Contains(t, err.Error(), tt.ExpectedErr.Error())
				}
				return
			}

			if tt.Verify != nil {
				tt.Verify(t, wallet, err)
			}
		})
	}
}

// TestWalletUseCase_UpdateWallet tests the UpdateWallet method of the WalletUseCase.
//
//	It checks for various error scenarios and edge cases, such as:
//	- successful wallet update
//	- wallet not found
//	- duplicate wallet name
//	- negative balance
//	- error fetching wallet
//	- error checking wallet names during update
//	- no changes made during update
func TestWalletUseCase_UpdateWallet(t *testing.T) {
	tests := []mocks.TestSetup[dtos.UpdateWalletRequest, db.Wallet, int]{
		{
			Name: "successful wallet update",
			Req: dtos.UpdateWalletRequest{
				WalletID:   1,
				Name:       "Updated Savings",
				WalletType: types.WalletTypePtr(types.Debit),
				Balance:    float64Ptr(2000.0),
			},
			SetupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("FindByField", &db.Wallet{
					ID:      1,
					Name:    "Savings",
					Type:    types.Debit,
					Balance: 1000.0,
					UserID:  1,
				}, nil)
				mock.SetResponse("GetAll", []db.Wallet{}, nil)
				mock.SetResponse("Update", &db.Wallet{
					ID:      1,
					Name:    "Updated Savings",
					Type:    types.Credit,
					Balance: 2000.0,
					UserID:  1,
				}, nil)
			},
			Verify: func(t *testing.T, wallet *db.Wallet, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Updated Savings", wallet.Name)
				assert.Equal(t, types.Credit, wallet.Type)
				assert.Equal(t, 2000.0, wallet.Balance)
			},
		},
		{
			Name: "wallet not found",
			Req: dtos.UpdateWalletRequest{
				WalletID: 999,
				Name:     "Updated",
			},
			SetupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("FindByField", nil, infrastructure.ErrNotFound)
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("wallet not found"),
		},
		{
			Name: "duplicate wallet name",
			Req: dtos.UpdateWalletRequest{
				WalletID: 1,
				Name:     "Existing Wallet",
			},
			SetupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("FindByField", &db.Wallet{
					ID:      1,
					Name:    "Savings",
					Type:    "savings",
					Balance: 1000.0,
					UserID:  1,
				}, nil)
				mock.SetResponse("GetAll", []db.Wallet{{
					ID:      2,
					Name:    "Existing Wallet",
					Type:    types.Debit,
					Balance: 500.0,
					UserID:  1,
				}}, nil)
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("a wallet with this name already exists for this user"),
		},
		{
			Name: "negative balance",
			Req: dtos.UpdateWalletRequest{
				WalletID: 1,
				Balance:  float64Ptr(-100.0),
			},
			SetupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("FindByField", &db.Wallet{
					ID:      1,
					Name:    "Savings",
					Type:    "savings",
					Balance: 1000.0,
					UserID:  1,
				}, nil)
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("balance cannot be negative"),
		},
		{
			Name: "error fetching wallet",
			Req: dtos.UpdateWalletRequest{
				WalletID: 1,
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("error fetching wallet"),
			SetupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("FindByField", nil, errors.New("database error"))
			},
		},
		{
			Name: "Error when wallet id is zero",
			Req: dtos.UpdateWalletRequest{
				WalletID: 0,
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("invalid wallet ID"),
			SetupMock:   func(mr *mocks.MockRepository[db.Wallet, int]) {},
		},
		{
			Name: "error checking wallet names during update",
			Req: dtos.UpdateWalletRequest{
				WalletID: 1,
				Name:     "New Name", // Un nombre diferente al existente para que entre en la validación
			},
			ExpectErr:   true,
			ExpectedErr: errors.New("error checking wallet names"),
			SetupMock: func(mr *mocks.MockRepository[db.Wallet, int]) {
				// Mock para FindByField que devuelve una billetera existente
				mr.SetResponse("FindByField", &db.Wallet{
					ID:      1,
					Name:    "Old Name",
					Type:    "Debit",
					Balance: 1000.0,
					UserID:  1,
				}, nil)
				// Mock para GetAll que devuelve un error
				mr.SetResponse("GetAll", nil, errors.New("database error"))
			},
		},
		{
			Name: "no changes made during update",
			Req: dtos.UpdateWalletRequest{
				WalletID: 1,
				Name:     "Old Name", // Same name as existing
			},
			ExpectErr:   false,
			ExpectedErr: nil,
			SetupMock: func(mr *mocks.MockRepository[db.Wallet, int]) {
				// Mock para FindByField que devuelve una billetera existente
				mr.SetResponse("FindByField", &db.Wallet{
					ID:      1,
					Name:    "Old Name",
					Type:    "Debit",
					Balance: 1000.0,
					UserID:  1,
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			repo := mocks.NewMockRepository[db.Wallet, int]()
			if tt.SetupMock != nil {
				tt.SetupMock(repo)
			}

			useCase := usecases.NewWalletUseCase(repo)
			wallet, err := useCase.UpdateWallet(tt.Req)

			if tt.ExpectErr {
				assert.Error(t, err)
				if tt.ExpectedErr != nil {
					assert.Contains(t, err.Error(), tt.ExpectedErr.Error())
				}
				return
			}

			if tt.Verify != nil {
				tt.Verify(t, wallet, err)
			}
		})
	}
}

func TestWalletUseCase_DeleteWallet(t *testing.T) {
	tests := []struct {
		name        string
		walletID    int
		expectErr   bool
		expectedErr error
		setupMock   func(*mocks.MockRepository[db.Wallet, int])
	}{
		{
			name:     "successful wallet deletion",
			walletID: 1,
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("GetByID", &db.Wallet{
					ID:      1,
					Name:    "Savings",
					Type:    "savings",
					Balance: 1000.0,
					UserID:  1,
				}, nil)
				mock.SetResponse("Delete", nil, nil)
			},
		},
		{
			name:        "invalid wallet ID",
			walletID:    0,
			expectErr:   true,
			expectedErr: errors.New("invalid wallet ID"),
		},
		{
			name:     "wallet not found",
			walletID: 999,
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("GetByID", nil, infrastructure.ErrNotFound)
			},
			expectErr:   true,
			expectedErr: errors.New("wallet not found"),
		},
		{
			name:     "wallet not found",
			walletID: 999,
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("GetByID", nil, infrastructure.ErrNotFound)
			},
			expectErr:   true,
			expectedErr: errors.New("wallet not found"),
		},
		{
			name:     "error fetching wallet",
			walletID: 1,
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("GetByID", nil, errors.New("database error"))
			},
			expectErr:   true,
			expectedErr: errors.New("error fetching wallet: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository[db.Wallet, int]()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			useCase := usecases.NewWalletUseCase(repo)
			err := useCase.DeleteWallet(tt.walletID)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestWalletUseCase_GetUserWallet(t *testing.T) {
	// Setup test cases
	tests := []struct {
		name         string
		mockSetup    func(*mocks.MockRepository[db.Wallet, int])
		userID       int
		email        string
		expectError  bool
		expectWallet *ports.UserWallet
	}{
		{
			name:   "success - user with wallets",
			userID: 1,
			email:  "test@example.com",
			mockSetup: func(mockRepo *mocks.MockRepository[db.Wallet, int]) {
				// Mock the Query method to return wallets for the user
				mockWallets := []db.Wallet{
					{
						ID:      1,
						Name:    "Savings",
						Type:    types.Credit,
						Balance: 1000.50,
						User: &db.User{
							Email: "test@example.com",
						},
					},
					{
						ID:      2,
						Name:    "Checking",
						Type:    types.Debit,
						Balance: 500.75,
						User: &db.User{
							Email: "test@example.com",
						},
					},
				}
				mockRepo.SetResponse("Query", mockWallets, nil)
			},
			expectError: false,
			expectWallet: &ports.UserWallet{
				Email: "test@example.com",
				Wallets: []struct {
					Name    string           `json:"name"`
					Type    types.WalletType `json:"type"`
					Balance float64          `json:"balance"`
				}{
					{Name: "Savings", Type: types.Credit, Balance: 1000.50},
					{Name: "Checking", Type: types.Debit, Balance: 500.75},
				},
			},
		},
		{
			name:   "success - user with no wallets",
			userID: 1,
			email:  "test@example.com",
			mockSetup: func(mockRepo *mocks.MockRepository[db.Wallet, int]) {
				// Mock the Query method to return no wallets
				mockRepo.SetResponse("Query", []db.Wallet{}, nil)
			},
			expectError: false,
			expectWallet: &ports.UserWallet{
				Email: "test@example.com",
				Wallets: []struct {
					Name    string           `json:"name"`
					Type    types.WalletType `json:"type"`
					Balance float64          `json:"balance"`
				}{},
			},
		},
		{
			name:   "error - repository error",
			userID: 1,
			email:  "test@example.com",
			mockSetup: func(mockRepo *mocks.MockRepository[db.Wallet, int]) {
				// Mock the Query method to return an error
				mockRepo.SetResponse("Query", nil, errors.New("database error"))
			},
			expectError: true,
		},
		{
			name:   "error - unexpected type from repository",
			userID: 1,
			email:  "test@example.com",
			mockSetup: func(mockRepo *mocks.MockRepository[db.Wallet, int]) {
				// Mock the Query method to return an unexpected type by using a custom response handler
				mockRepo.SetResponse("Query", []struct{}{}, nil) // This will cause a type assertion failure
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock repository
			mockRepo := mocks.NewMockRepository[db.Wallet, int]()

			// Set up the mock expectations
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			// Create the use case with the mock repository
			uc := usecases.NewWalletUseCase(mockRepo)

			// Call the method being tested
			result, err := uc.GetUserWallet(tt.userID, tt.email)

			// Assert the results
			if tt.expectError {
				assert.Error(t, err, "Expected an error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, tt.expectWallet.Email, result.Email, "Email should match")
				assert.Len(t, result.Wallets, len(tt.expectWallet.Wallets), "Number of wallets should match")

				// Compare each wallet
				for i, expectedWallet := range tt.expectWallet.Wallets {
					assert.Equal(t, expectedWallet.Name, result.Wallets[i].Name, "Wallet name should match")
					assert.Equal(t, expectedWallet.Type, result.Wallets[i].Type, "Wallet type should match")
					assert.Equal(t, expectedWallet.Balance, result.Wallets[i].Balance, "Wallet balance should match")
				}

				// Verify the query was called with the expected parameters
				calls := mockRepo.Calls("Query")
				assert.Len(t, calls, 1, "Query should be called once")

			}
		})
	}
}
