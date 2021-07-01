package get_operations

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"payment-system/internal/wallet"
)

type walletService interface {
	GetOperations(ctx context.Context, filter wallet.Filter) ([]wallet.Operation, error)
}

type Handler struct {
	walletService walletService
}

func NewHandler(walletService walletService) *Handler {
	return &Handler{walletService: walletService}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dto, err := validate(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = fmt.Errorf("failed to validate request: %w", err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Printf("failed to write bad request error message: %s\n", err)
		}
		return
	}

	ctx := r.Context()
	filter := wallet.Filter{
		WalletID:  dto.WalletID,
		Date:      dto.Date,
		Direction: dto.Direction,
	}
	operations, err := h.walletService.GetOperations(ctx, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Printf("failed to write error message: %s\n", err)
		}
		return
	}

	records := make([][]string, 0, len(operations)+1)
	records = append(records, []string{"wallet_id", "value", "direction", "date"})
	for _, operation := range operations {
		walletID := strconv.FormatInt(operation.WalletID, 10)
		value := strconv.FormatFloat(operation.Value, 'f', 2, 64)
		direction := strconv.Itoa(int(operation.Direction))
		record := []string{walletID, value, direction, operation.Date}
		records = append(records, record)
	}

	if err = csv.NewWriter(w).WriteAll(records); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Printf("failed to write error message: %s\n", err)
		}
	}
}
