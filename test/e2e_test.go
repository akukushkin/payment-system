package test

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"payment-system/internal/handlers/add_wallet"
	"payment-system/internal/handlers/deposit_money"
	"payment-system/internal/handlers/get_operations"
	"payment-system/internal/handlers/transfer_money"
)

func TestHappyPath(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	// add 1st wallet
	idempotencyKey := uuid.New().String()
	walletInDTO := add_wallet.WalletInDTO{IdempotencyKey: idempotencyKey}
	walletOutDTO, err := addWallet(&httpClient, walletInDTO)
	require.NoError(t, err)

	fromWalletID := walletOutDTO.WalletID

	// imitate duplicate operation
	_, err = addWallet(&httpClient, walletInDTO)
	require.Error(t, err)

	// add 2nd wallet
	idempotencyKey = uuid.New().String()
	walletInDTO.IdempotencyKey = idempotencyKey
	walletOutDTO, err = addWallet(&httpClient, walletInDTO)
	require.NoError(t, err)

	toWalletID := walletOutDTO.WalletID

	// deposit to 1st wallet
	idempotencyKey = uuid.New().String()
	depositDTO := deposit_money.DepositDTO{
		IdempotencyKey: idempotencyKey,
		WalletID:       fromWalletID,
		Value:          100.53,
	}
	err = depositMoney(&httpClient, depositDTO)
	require.NoError(t, err)

	// imitate duplicate deposit to 1st wallet
	err = depositMoney(&httpClient, depositDTO)
	require.Error(t, err)

	// transfer part of money from 1st wallet to 2nd wallet
	idempotencyKey = uuid.New().String()
	transferDTO := transfer_money.TransferDTO{
		FromWalletID:   fromWalletID,
		ToWalletID:     toWalletID,
		Value:          50.51,
		IdempotencyKey: idempotencyKey,
	}
	err = transferMoney(&httpClient, transferDTO)
	require.NoError(t, err)

	// imitate duplicate transfer
	err = transferMoney(&httpClient, transferDTO)
	require.Error(t, err)

	// transfer part of money from 1st wallet to 2nd wallet
	idempotencyKey = uuid.New().String()
	transferDTO.IdempotencyKey = idempotencyKey
	transferDTO.Value = 50
	err = transferMoney(&httpClient, transferDTO)
	require.NoError(t, err)

	// transfer part of money from 1st wallet to 2nd wallet, but 1st wallet does not have enough money
	idempotencyKey = uuid.New().String()
	transferDTO.IdempotencyKey = idempotencyKey
	err = transferMoney(&httpClient, transferDTO)
	require.Error(t, err)

	// get income operation by 1st wallet
	filterDTO := get_operations.FilterDTO{
		WalletID:  fromWalletID,
		Date:      time.Now().UTC().Format("2006-01-02"),
		Direction: 0,
	}
	operations, err := getOperations(&httpClient, filterDTO)
	require.NoError(t, err)
	require.Len(t, operations, 2)

	// get income operation by 2nd wallet
	filterDTO.WalletID = toWalletID
	operations, err = getOperations(&httpClient, filterDTO)
	require.NoError(t, err)
	require.Len(t, operations, 3)

	// get outcome operation by 1st wallet
	filterDTO.WalletID = fromWalletID
	filterDTO.Direction = 1
	operations, err = getOperations(&httpClient, filterDTO)
	require.NoError(t, err)
	require.Len(t, operations, 3)

	// get outcome operation by 2nd wallet
	filterDTO.WalletID = toWalletID
	operations, err = getOperations(&httpClient, filterDTO)
	require.NoError(t, err)
	require.Len(t, operations, 1)
}

func getOperations(client *http.Client, in get_operations.FilterDTO) ([][]string, error) {
	marshaled, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/getOperations",
		bytes.NewReader(marshaled),
	)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unsuccess status code")
	}

	return csv.NewReader(resp.Body).ReadAll()
}

func transferMoney(client *http.Client, in transfer_money.TransferDTO) error {
	marshaled, err := json.Marshal(in)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/transferMoney",
		bytes.NewReader(marshaled),
	)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccess status code")
	}

	return nil
}

func depositMoney(client *http.Client, in deposit_money.DepositDTO) error {
	marshaled, err := json.Marshal(in)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/depositMoney",
		bytes.NewReader(marshaled),
	)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccess status code")
	}

	return nil
}

func addWallet(client *http.Client, in add_wallet.WalletInDTO) (add_wallet.WalletOutDTO, error) {
	marshaled, err := json.Marshal(in)
	if err != nil {
		return add_wallet.WalletOutDTO{}, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/addWallet",
		bytes.NewReader(marshaled),
	)
	if err != nil {
		return add_wallet.WalletOutDTO{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return add_wallet.WalletOutDTO{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return add_wallet.WalletOutDTO{}, fmt.Errorf("unsuccess status code")
	}

	var out add_wallet.WalletOutDTO
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return add_wallet.WalletOutDTO{}, err
	}

	return out, nil
}
