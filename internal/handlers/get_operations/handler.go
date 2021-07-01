package get_operations

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

	operationDTOs := make([]OperationDTO, 0, len(operations))
	for _, operation := range operations {
		operationDTO := OperationDTO{
			WalletID:  operation.WalletID,
			Value:     operation.Value,
			Direction: operation.Direction,
			Date:      operation.Date,
		}
		operationDTOs = append(operationDTOs, operationDTO)
	}

	if err = json.NewEncoder(w).Encode(operationDTOs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Printf("failed to write error message: %s\n", err)
		}
	}
}
