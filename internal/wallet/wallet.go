//go:generate mockgen -source=wallet.go -destination mock.go -package $GOPACKAGE
package wallet

import (
	"context"
	"fmt"

	"payment-system/internal/storage"
)

type Wallet struct {
	IdempotencyKey string
}

type Deposit struct {
	WalletID       int64
	Value          float64
	IdempotencyKey string
}

type Transfer struct {
	FromWalletID   int64
	ToWalletID     int64
	Value          float64
	IdempotencyKey string
}

type Filter struct {
	WalletID  int64
	Date      string
	Direction int8
}

type Operation struct {
	WalletID  int64
	Value     float64
	Direction int8
	Date      string
}

type walletStorage interface {
	AddWallet(ctx context.Context, wallet storage.Wallet) (int64, error)
	DepositMoney(ctx context.Context, deposit storage.Deposit) error
	TransferMoney(ctx context.Context, info storage.Transfer) error
	GetOperations(ctx context.Context, filter storage.Filter) ([]storage.Operation, error)
}

type Service struct {
	storage walletStorage
}

func New(storage walletStorage) *Service {
	return &Service{storage: storage}
}

func (s *Service) AddWallet(ctx context.Context, wallet Wallet) (int64, error) {
	w := storage.Wallet{
		IdempotencyKey: wallet.IdempotencyKey,
	}
	walletID, err := s.storage.AddWallet(ctx, w)
	if err != nil {
		return 0, fmt.Errorf("adding wallet into storage: %w", err)
	}

	return walletID, nil
}

func (s *Service) DepositMoney(ctx context.Context, deposit Deposit) error {
	d := storage.Deposit{
		WalletID:       deposit.WalletID,
		Value:          dollarsToCents(deposit.Value),
		IdempotencyKey: deposit.IdempotencyKey,
	}
	if err := s.storage.DepositMoney(ctx, d); err != nil {
		return fmt.Errorf("depositing money into storage: %w", err)
	}

	return nil
}

func (s *Service) TransferMoney(ctx context.Context, transfer Transfer) error {
	t := storage.Transfer{
		FromWalletID:   transfer.FromWalletID,
		ToWalletID:     transfer.ToWalletID,
		Value:          dollarsToCents(transfer.Value),
		IdempotencyKey: transfer.IdempotencyKey,
	}
	if err := s.storage.TransferMoney(ctx, t); err != nil {
		return fmt.Errorf("transferring money into storage: %w", err)
	}

	return nil
}

func (s *Service) GetOperations(ctx context.Context, filter Filter) ([]Operation, error) {
	f := storage.Filter{
		WalletID:  filter.WalletID,
		Date:      filter.Date,
		Direction: storage.Direction(filter.Direction),
	}
	storageOperations, err := s.storage.GetOperations(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("getting operations from storage: %w", err)
	}

	operations := make([]Operation, 0, len(storageOperations))
	for _, storageOperation := range storageOperations {
		operation := Operation{
			WalletID:  storageOperation.WalletID,
			Value:     centsToDollars(storageOperation.Value),
			Direction: int8(storageOperation.Direction),
			Date:      storageOperation.Date,
		}
		operations = append(operations, operation)
	}

	return operations, nil
}

func dollarsToCents(dollars float64) int64 {
	return int64(dollars * 100)
}

func centsToDollars(cents int64) float64 {
	return float64(cents) / 100.0
}
