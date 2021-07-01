package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"payment-system/internal/db"
	"payment-system/internal/handlers/add_wallet"
	"payment-system/internal/handlers/deposit_money"
	"payment-system/internal/handlers/get_operations"
	"payment-system/internal/handlers/transfer_money"
	"payment-system/internal/storage"
	"payment-system/internal/wallet"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	database := db.New()
	walletService := wallet.New(storage.New(database))

	srv := http.Server{Addr: fmt.Sprintf(":%s", port)}
	http.Handle("/addWallet", add_wallet.NewHandler(walletService))
	http.Handle("/depositMoney", deposit_money.NewHandler(walletService))
	http.Handle("/transferMoney", transfer_money.NewHandler(walletService))
	http.Handle("/getOperations", get_operations.NewHandler(walletService))

	go func() {
		log.Printf("listening on port %s\n", port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to listen and serve: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %s", err)
	}
}
