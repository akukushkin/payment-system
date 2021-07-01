package transfer_money

import (
	"context"
	"log"
	"net/http"

	"payment-system/internal/wallet"
)

type walletService interface {
	TransferMoney(ctx context.Context, transfer wallet.Transfer) error
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
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Printf("failed to write bad request error message: %s\n", err)
		}
		return
	}

	ctx := r.Context()
	transfer := wallet.Transfer{
		FromWalletID:   dto.FromWalletID,
		ToWalletID:     dto.ToWalletID,
		Value:          dto.Value,
		IdempotencyKey: dto.IdempotencyKey,
	}
	if err := h.walletService.TransferMoney(ctx, transfer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Printf("failed to write error message: %s\n", err)
		}
	}
}
