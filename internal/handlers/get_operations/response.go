package get_operations

type OperationDTO struct {
	WalletID  int64   `json:"wallet_id"`
	Value     float64 `json:"value"`
	Direction int8    `json:"direction"`
	Date      string  `json:"date"`
}
