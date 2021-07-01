package wallet

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"payment-system/internal/storage"
)

func TestService_AddWallet_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletStorage := NewMockwalletStorage(ctrl)
	mockWalletStorage.EXPECT().AddWallet(gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("something went wrong"))
	service := New(mockWalletStorage)
	_, err := service.AddWallet(context.Background(), Wallet{})
	require.Error(t, err)
}

func TestService_AddWallet_ReturnsWalletID(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletStorage := NewMockwalletStorage(ctrl)
	mockWalletStorage.EXPECT().AddWallet(gomock.Any(), gomock.Any()).Return(int64(1), nil)
	service := New(mockWalletStorage)
	walletID, err := service.AddWallet(context.Background(), Wallet{})
	require.NoError(t, err)
	require.Equal(t, int64(1), walletID)
}

func TestService_DepositMoney_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletStorage := NewMockwalletStorage(ctrl)
	mockWalletStorage.EXPECT().DepositMoney(gomock.Any(), gomock.Any()).Return(fmt.Errorf("something went wrong"))
	service := New(mockWalletStorage)
	err := service.DepositMoney(context.Background(), Deposit{})
	require.Error(t, err)
}

func TestService_DepositMoney_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletStorage := NewMockwalletStorage(ctrl)
	mockWalletStorage.EXPECT().DepositMoney(gomock.Any(), gomock.Any()).Return(nil)
	service := New(mockWalletStorage)
	err := service.DepositMoney(context.Background(), Deposit{})
	require.NoError(t, err)
}

func TestService_TransferMoney_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletStorage := NewMockwalletStorage(ctrl)
	mockWalletStorage.EXPECT().TransferMoney(gomock.Any(), gomock.Any()).Return(fmt.Errorf("something went wrong"))
	service := New(mockWalletStorage)
	err := service.TransferMoney(context.Background(), Transfer{})
	require.Error(t, err)
}

func TestService_TransferMoney_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletStorage := NewMockwalletStorage(ctrl)
	mockWalletStorage.EXPECT().TransferMoney(gomock.Any(), gomock.Any()).Return(nil)
	service := New(mockWalletStorage)
	err := service.TransferMoney(context.Background(), Transfer{})
	require.NoError(t, err)
}

func TestService_GetOperationsReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletStorage := NewMockwalletStorage(ctrl)
	mockWalletStorage.EXPECT().GetOperations(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("something went wrong"))
	service := New(mockWalletStorage)
	_, err := service.GetOperations(context.Background(), Filter{})
	require.Error(t, err)
}

func TestService_GetOperationsReturnsOperations(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWalletStorage := NewMockwalletStorage(ctrl)
	mockWalletStorage.EXPECT().GetOperations(gomock.Any(), gomock.Any()).Return([]storage.Operation{{
		WalletID:  2,
		Value:     115,
		Direction: 1,
		Date:      "2021-06-30",
	}, {
		WalletID:  2,
		Value:     1102,
		Direction: 0,
		Date:      "2021-05-22",
	}}, nil)
	service := New(mockWalletStorage)
	operations, err := service.GetOperations(context.Background(), Filter{})
	require.NoError(t, err)
	require.Equal(t, []Operation{{
		WalletID:  2,
		Value:     1.15,
		Direction: 1,
		Date:      "2021-06-30",
	},
		{
			WalletID:  2,
			Value:     11.02,
			Direction: 0,
			Date:      "2021-05-22",
		}}, operations)
}

func Test_dollarsToCents(t *testing.T) {
	type args struct {
		dollars float64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "dollars without cents",
			args: args{
				dollars: 100,
			},
			want: 10000,
		},
		{
			name: "dollars with cents",
			args: args{
				dollars: 1.53,
			},
			want: 153,
		},
		{
			name: "dollars with cents to round",
			args: args{
				dollars: 1.553,
			},
			want: 155,
		},
		{
			name: "only cents",
			args: args{
				dollars: 0.55,
			},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dollarsToCents(tt.args.dollars); got != tt.want {
				t.Errorf("dollarsToCents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_centsToDollars(t *testing.T) {
	type args struct {
		cents int64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "dollar without cents",
			args: args{
				cents: 100,
			},
			want: 1,
		},
		{
			name: "only cents",
			args: args{
				cents: 1,
			},
			want: 0.01,
		},
		{
			name: "dollar with cents",
			args: args{
				cents: 115,
			},
			want: 1.15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := centsToDollars(tt.args.cents); got != tt.want {
				t.Errorf("centsToDollars() = %v, want %v", got, tt.want)
			}
		})
	}
}
