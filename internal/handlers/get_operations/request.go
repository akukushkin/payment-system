package get_operations

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FilterDTO struct {
	WalletID  int64  `json:"wallet_id"`
	Date      string `json:"date"`
	Direction int8   `json:"direction"`
}

func validate(r *http.Request) (FilterDTO, error) {
	decoder := json.NewDecoder(r.Body)
	var filter FilterDTO
	if err := decoder.Decode(&filter); err != nil {
		return FilterDTO{}, err
	}

	if filter.WalletID == 0 {
		return FilterDTO{}, fmt.Errorf("wallet_id is empty")
	}

	if filter.Date == "" {
		return FilterDTO{}, fmt.Errorf("date is empty")
	}

	if _, err := time.Parse("2006-01-02", filter.Date); err != nil {
		return FilterDTO{}, fmt.Errorf("date is invalid: %w", err)
	}

	return filter, nil
}
