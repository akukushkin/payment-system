package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const defaultOperationsCapacity = 1000

const (
	insertWalletQuery     = "INSERT INTO wallet(idempotency_key) VALUES (:idempotency_key) RETURNING id"
	insertOperationQuery  = "INSERT INTO operation(wallet_id, value, Direction, idempotency_key) VALUES ($1, $2, $3, $4)"
	updateWalletQuery     = "UPDATE wallet SET value = value + $2 WHERE id = $1"
	selectOperationsQuery = "SELECT wallet_id, value, direction, to_char(date, 'YYYY-MM-DD') as date FROM operation " +
		"WHERE wallet_id = $1 AND date = $2 AND direction = $3"
)

type Direction int8

const (
	deposit    Direction = 0
	withdrawal Direction = 1
)

type Wallet struct {
	IdempotencyKey string `db:"idempotency_key"`
}

type Deposit struct {
	WalletID       int64
	Value          int64
	IdempotencyKey string
}

type Transfer struct {
	FromWalletID   int64
	ToWalletID     int64
	Value          int64
	IdempotencyKey string
}

type Filter struct {
	WalletID  int64
	Date      string
	Direction Direction
}

type Operation struct {
	WalletID  int64     `db:"wallet_id"`
	Value     int64     `db:"value"`
	Direction Direction `db:"direction"`
	Date      string    `db:"date"`
}

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) AddWallet(ctx context.Context, wallet Wallet) (int64, error) {
	rows, err := s.db.NamedQueryContext(ctx, insertWalletQuery, wallet)
	if err != nil {
		return 0, fmt.Errorf("inserting wallet: %w", err)
	}

	if rows != nil {
		defer func(rows *sqlx.Rows) {
			if err = rows.Close(); err != nil {
				log.Printf("failed to close rows: %s\n", err)
			}
		}(rows)
	}

	if !rows.Next() {
		return 0, fmt.Errorf("preparing to scan inserted wallet id: %w", rows.Err())
	}

	var walletID int64
	if err := rows.Scan(&walletID); err != nil {
		return 0, fmt.Errorf("scanning inserted wallet id: %w", err)
	}

	return walletID, nil
}

func (s *Storage) DepositMoney(ctx context.Context, info Deposit) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("beginning deposit money tx: %w", err)
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("failed to rollback deposit money tx: %s\n", err)
			}
			return
		}

		if err = tx.Commit(); err != nil {
			err = fmt.Errorf("commiting deposit money tx: %w", err)
		}
	}()

	_, err = tx.ExecContext(ctx, insertOperationQuery, info.WalletID, info.Value, deposit, info.IdempotencyKey)
	if err != nil {
		err = fmt.Errorf("executing inserting deposit money operation: %w", err)
		return
	}

	_, err = tx.ExecContext(ctx, updateWalletQuery, info.WalletID, info.Value)
	if err != nil {
		err = fmt.Errorf("executing updating wallet: %w", err)
	}

	return
}

func (s *Storage) TransferMoney(ctx context.Context, info Transfer) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = fmt.Errorf("beginning transfer money tx: %w", err)
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("failed to rollback transfer money tx: %s\n", err)
			}
			return
		}

		if err = tx.Commit(); err != nil {
			err = fmt.Errorf("commiting transfer money tx: %w", err)
		}
	}()

	_, err = tx.ExecContext(ctx, insertOperationQuery, info.FromWalletID, info.Value, withdrawal, info.IdempotencyKey)
	if err != nil {
		err = fmt.Errorf("executing inserting withdrawal money operation: %w", err)
		return
	}

	_, err = tx.ExecContext(ctx, insertOperationQuery, info.ToWalletID, info.Value, deposit, info.IdempotencyKey)
	if err != nil {
		err = fmt.Errorf("executing inserting deposit money operation: %w", err)
		return
	}

	_, err = tx.ExecContext(ctx, updateWalletQuery, info.FromWalletID, -info.Value)
	if err != nil {
		err = fmt.Errorf("executing updating deposit wallet: %w", err)
		return
	}

	_, err = tx.ExecContext(ctx, updateWalletQuery, info.ToWalletID, info.Value)
	if err != nil {
		err = fmt.Errorf("executing updating withdrawal wallet: %w", err)
		return
	}

	return
}

func (s *Storage) GetOperations(ctx context.Context, filter Filter) ([]Operation, error) {
	operations := make([]Operation, 0, defaultOperationsCapacity)
	err := s.db.SelectContext(ctx, &operations, selectOperationsQuery, filter.WalletID, filter.Date, filter.Direction)
	if err != nil {
		return nil, fmt.Errorf("getting operations from storage: %w", err)
	}

	return operations, nil
}
