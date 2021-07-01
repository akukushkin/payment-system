package transfer_money

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type TransferDTO struct {
	FromWalletID   int64   `json:"from_wallet_id"`
	ToWalletID     int64   `json:"to_wallet_id"`
	Value          float64 `json:"value"`
	IdempotencyKey string  `json:"idempotency_key"`
}

func validate(r *http.Request) (TransferDTO, error) {
	decoder := json.NewDecoder(r.Body)
	var transfer TransferDTO
	if err := decoder.Decode(&transfer); err != nil {
		return TransferDTO{}, err
	}

	if transfer.FromWalletID == 0 {
		return TransferDTO{}, fmt.Errorf("from_wallet_id is empty")
	}

	if transfer.ToWalletID == 0 {
		return TransferDTO{}, fmt.Errorf("to_wallet_id is empty")
	}

	if transfer.Value == 0 {
		return TransferDTO{}, fmt.Errorf("value is empty")
	}

	if transfer.IdempotencyKey == "" {
		return TransferDTO{}, fmt.Errorf("idempotency_key is empty")
	}

	return transfer, nil
}
