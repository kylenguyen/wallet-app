package service_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/kylenguyen/wallet-app/internal/model"
	"github.com/kylenguyen/wallet-app/internal/service"
	walletmocks "github.com/kylenguyen/wallet-app/internal/service/mocks" // Assuming you'll create this mock
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Define some UUIDs for testing consistency
var (
	testUser1UUIDString   = "d5b706b8-c331-4670-9199-a773d10878d5"
	testWallet1UUIDString = "f8b3f7a7-1b3b-4b3f-8f3b-3b3b3b3b3b3b"
	testUser1UUID         = uuid.MustParse(testUser1UUIDString)
	testWallet1UUID       = uuid.MustParse(testWallet1UUIDString)
)

func TestWalletServiceImpl_GetWalletInfo(t *testing.T) {
	type fields struct {
		wRepo func() service.WalletRepo
	}
	type args struct {
		ctx      context.Context
		userId   string
		walletId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Wallet
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success - get wallet info",
			fields: fields{
				wRepo: func() service.WalletRepo {
					m := new(walletmocks.WalletRepoMock)
					m.On("GetWalletInfo", mock.Anything, testUser1UUIDString, testWallet1UUIDString).Return(&model.Wallet{ID: testWallet1UUID, UserID: testUser1UUID, Balance: decimal.NewFromInt(100)}, nil)
					return m
				},
			},
			args: args{
				ctx:      context.Background(),
				userId:   testUser1UUIDString,
				walletId: testWallet1UUIDString,
			},
			want:    &model.Wallet{ID: testWallet1UUID, UserID: testUser1UUID, Balance: decimal.NewFromInt(100)},
			wantErr: assert.NoError,
		},
		{
			name: "error - repository returns error",
			fields: fields{
				wRepo: func() service.WalletRepo {
					m := new(walletmocks.WalletRepoMock)
					m.On("GetWalletInfo", mock.Anything, testUser1UUIDString, testWallet1UUIDString).Return(nil, errors.New("repo error"))
					return m
				},
			},
			args: args{
				ctx:      context.Background(),
				userId:   testUser1UUIDString,
				walletId: testWallet1UUIDString,
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wRepo := tt.fields.wRepo().(*walletmocks.WalletRepoMock)
			ws := service.NewWalletImpl(wRepo)
			got, err := ws.GetWalletInfo(tt.args.ctx, tt.args.userId, tt.args.walletId)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
			wRepo.AssertExpectations(t)
		})
	}
}
