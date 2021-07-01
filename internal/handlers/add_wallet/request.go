package add_wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WalletInDTO struct {
	IdempotencyKey string `json:"idempotency_key"`
}

func validate(r *http.Request) (WalletInDTO, error) {
	decoder := json.NewDecoder(r.Body)
	var wallet WalletInDTO
	if err := decoder.Decode(&wallet); err != nil {
		return WalletInDTO{}, err
	}

	if wallet.IdempotencyKey == "" {
		return WalletInDTO{}, fmt.Errorf("idempotency_key is empty")
	}

	return wallet, nil
}
