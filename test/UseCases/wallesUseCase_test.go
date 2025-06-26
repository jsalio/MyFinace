package UseCases_test

import (
	"errors"
	"testing"

	// "Financial/Domains/ports"
	"Financial/Models/db"
	"Financial/Models/dtos"
	"Financial/infrastructure"
	"Financial/types"

	// "Financial/UseCases"
	usecases "Financial/UseCases"
	mocks "Financial/test"

	"github.com/stretchr/testify/assert"
)

type TestSetup[T any] struct {
	name        string
	req         T
	expectErr   bool
	expectedErr error
	setupMock   func(*mocks.MockRepository[db.Wallet, int])
	verify      func(t *testing.T, wallet *db.Wallet, err error)
}

func TestWalletUseCase_CreateWallet(t *testing.T) {
	tests := []TestSetup[dtos.CreateWalletRequest]{
		{
			name: "successful wallet creation",
			req: dtos.CreateWalletRequest{
				Name:       "Savings",
				WalletType: "savings",
				Balance:    1000.0,
				UserID:     1,
			},
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("GetAll", []db.Wallet{}, nil)
				mock.SetResponse("Create", &db.Wallet{
					ID:      1,
					Name:    "Savings",
					Type:    types.Debit,
					Balance: 1000.0,
					UserID:  1,
				}, nil)
			},
			verify: func(t *testing.T, wallet *db.Wallet, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Savings", wallet.Name)
				assert.Equal(t, types.Debit, wallet.Type)
				assert.Equal(t, 1000.0, wallet.Balance)
				assert.Equal(t, 1, wallet.UserID)
			},
		},
		{
			name: "empty wallet name",
			req: dtos.CreateWalletRequest{
				Name:       "",
				WalletType: "savings",
				Balance:    1000.0,
				UserID:     1,
			},
			expectErr:   true,
			expectedErr: errors.New("wallet name cannot be empty"),
		},
		{
			name: "whitespace wallet name",
			req: dtos.CreateWalletRequest{
				Name:       "  ",
				WalletType: "savings",
				Balance:    1000.0,
				UserID:     1,
			},
			expectErr:   true,
			expectedErr: errors.New("wallet name cannot be empty"),
		},
		{
			name: "negative balance",
			req: dtos.CreateWalletRequest{
				Name:       "Savings",
				WalletType: types.Debit,
				Balance:    -100.0,
				UserID:     1,
			},
			expectErr:   true,
			expectedErr: errors.New("initial balance cannot be negative"),
		},
		{
			name: "invalid user ID",
			req: dtos.CreateWalletRequest{
				Name:       "Savings",
				WalletType: "savings",
				Balance:    1000.0,
				UserID:     0,
			},
			expectErr:   true,
			expectedErr: errors.New("invalid user ID"),
		},
		{
			name: "duplicate wallet name for user",
			req: dtos.CreateWalletRequest{
				Name:       "Savings",
				WalletType: "savings",
				Balance:    1000.0,
				UserID:     1,
			},
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("GetAll", []db.Wallet{{
					ID:      1,
					Name:    "Savings",
					Type:    "savings",
					Balance: 1000.0,
					UserID:  1,
				}}, nil)
			},
			expectErr:   true,
			expectedErr: errors.New("a wallet with this name already exists for this user"),
		},
		{
			name: "error checking wallet existence",
			req: dtos.CreateWalletRequest{
				Name:       "Savings",
				WalletType: types.Debit,
				Balance:    1000.0,
				UserID:     1,
			},
			expectErr:   true,
			expectedErr: errors.New("error checking wallet existence"),
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("GetAll", nil, errors.New("database error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository[db.Wallet, int]()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			useCase := usecases.NewWalletUseCase(repo)
			wallet, err := useCase.CreateWallet(tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				}
				return
			}

			if tt.verify != nil {
				tt.verify(t, wallet, err)
			}
		})
	}
}

func TestWalletUseCase_UpdateWallet(t *testing.T) {
	tests := []TestSetup[dtos.UpdateWalletRequest]{
		{
			name: "successful wallet update",
			req: dtos.UpdateWalletRequest{
				WalletID:   1,
				Name:       "Updated Savings",
				WalletType: types.WalletTypePtr(types.Debit),
				Balance:    float64Ptr(2000.0),
			},
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
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
			verify: func(t *testing.T, wallet *db.Wallet, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "Updated Savings", wallet.Name)
				assert.Equal(t, types.Credit, wallet.Type)
				assert.Equal(t, 2000.0, wallet.Balance)
			},
		},
		{
			name: "wallet not found",
			req: dtos.UpdateWalletRequest{
				WalletID: 999,
				Name:     "Updated",
			},
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("FindByField", nil, infrastructure.ErrNotFound)
			},
			expectErr:   true,
			expectedErr: errors.New("wallet not found"),
		},
		{
			name: "duplicate wallet name",
			req: dtos.UpdateWalletRequest{
				WalletID: 1,
				Name:     "Existing Wallet",
			},
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
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
			expectErr:   true,
			expectedErr: errors.New("a wallet with this name already exists for this user"),
		},
		{
			name: "negative balance",
			req: dtos.UpdateWalletRequest{
				WalletID: 1,
				Balance:  float64Ptr(-100.0),
			},
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("FindByField", &db.Wallet{
					ID:      1,
					Name:    "Savings",
					Type:    "savings",
					Balance: 1000.0,
					UserID:  1,
				}, nil)
			},
			expectErr:   true,
			expectedErr: errors.New("balance cannot be negative"),
		},
		{
			name: "error fetching wallet",
			req: dtos.UpdateWalletRequest{
				WalletID: 1,
			},
			expectErr:   true,
			expectedErr: errors.New("error fetching wallet"),
			setupMock: func(mock *mocks.MockRepository[db.Wallet, int]) {
				mock.SetResponse("FindByField", nil, errors.New("database error"))
			},
		},
		{
			name: "Error when wallet id is zero",
			req: dtos.UpdateWalletRequest{
				WalletID: 0,
			},
			expectErr:   true,
			expectedErr: errors.New("invalid wallet ID"),
			setupMock:   func(mr *mocks.MockRepository[db.Wallet, int]) {},
		},
		{
			name: "error checking wallet names during update",
			req: dtos.UpdateWalletRequest{
				WalletID: 1,
				Name:     "New Name", // Un nombre diferente al existente para que entre en la validación
			},
			expectErr:   true,
			expectedErr: errors.New("error checking wallet names"),
			setupMock: func(mr *mocks.MockRepository[db.Wallet, int]) {
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
			name: "no changes made during update",
			req: dtos.UpdateWalletRequest{
				WalletID: 1,
				Name:     "Old Name", // Same name as existing
			},
			expectErr:   false,
			expectedErr: nil,
			setupMock: func(mr *mocks.MockRepository[db.Wallet, int]) {
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
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockRepository[db.Wallet, int]()
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			useCase := usecases.NewWalletUseCase(repo)
			wallet, err := useCase.UpdateWallet(tt.req)

			if tt.expectErr {
				assert.Error(t, err)
				if tt.expectedErr != nil {
					assert.Contains(t, err.Error(), tt.expectedErr.Error())
				}
				return
			}

			if tt.verify != nil {
				tt.verify(t, wallet, err)
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
