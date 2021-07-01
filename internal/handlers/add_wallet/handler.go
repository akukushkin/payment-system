package add_wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"payment-system/internal/wallet"
)

type walletService interface {
	AddWallet(ctx context.Context, wallet wallet.Wallet) (int64, error)
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
	info := wallet.Wallet{IdempotencyKey: dto.IdempotencyKey}
	walletID, err := h.walletService.AddWallet(ctx, info)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Printf("failed to write error message: %s\n", err)
		}
		return
	}

	response := WalletOutDTO{WalletID: walletID}
	if err = json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
