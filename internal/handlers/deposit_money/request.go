package deposit_money

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DepositDTO struct {
	IdempotencyKey string  `json:"idempotency_key"`
	WalletID       int64   `json:"wallet_id"`
	Value          float64 `json:"value"`
}

func validate(r *http.Request) (DepositDTO, error) {
	decoder := json.NewDecoder(r.Body)
	var deposit DepositDTO
	if err := decoder.Decode(&deposit); err != nil {
		return DepositDTO{}, err
	}

	if deposit.IdempotencyKey == "" {
		return DepositDTO{}, fmt.Errorf("idempotency_key is empty")
	}

	if deposit.WalletID == 0 {
		return DepositDTO{}, fmt.Errorf("wallet_id is empty")
	}

	if deposit.Value == 0 {
		return DepositDTO{}, fmt.Errorf("deposit_value is empty")
	}

	return deposit, nil
}
