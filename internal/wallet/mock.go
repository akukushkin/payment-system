// Code generated by MockGen. DO NOT EDIT.
// Source: wallet.go

// Package wallet is a generated GoMock package.
package wallet

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	storage "payment-system/internal/storage"
	reflect "reflect"
)

// MockwalletStorage is a mock of walletStorage interface
type MockwalletStorage struct {
	ctrl     *gomock.Controller
	recorder *MockwalletStorageMockRecorder
}

// MockwalletStorageMockRecorder is the mock recorder for MockwalletStorage
type MockwalletStorageMockRecorder struct {
	mock *MockwalletStorage
}

// NewMockwalletStorage creates a new mock instance
func NewMockwalletStorage(ctrl *gomock.Controller) *MockwalletStorage {
	mock := &MockwalletStorage{ctrl: ctrl}
	mock.recorder = &MockwalletStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockwalletStorage) EXPECT() *MockwalletStorageMockRecorder {
	return m.recorder
}

// AddWallet mocks base method
func (m *MockwalletStorage) AddWallet(ctx context.Context, wallet storage.Wallet) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddWallet", ctx, wallet)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddWallet indicates an expected call of AddWallet
func (mr *MockwalletStorageMockRecorder) AddWallet(ctx, wallet interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddWallet", reflect.TypeOf((*MockwalletStorage)(nil).AddWallet), ctx, wallet)
}

// DepositMoney mocks base method
func (m *MockwalletStorage) DepositMoney(ctx context.Context, deposit storage.Deposit) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DepositMoney", ctx, deposit)
	ret0, _ := ret[0].(error)
	return ret0
}

// DepositMoney indicates an expected call of DepositMoney
func (mr *MockwalletStorageMockRecorder) DepositMoney(ctx, deposit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DepositMoney", reflect.TypeOf((*MockwalletStorage)(nil).DepositMoney), ctx, deposit)
}

// TransferMoney mocks base method
func (m *MockwalletStorage) TransferMoney(ctx context.Context, info storage.Transfer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransferMoney", ctx, info)
	ret0, _ := ret[0].(error)
	return ret0
}

// TransferMoney indicates an expected call of TransferMoney
func (mr *MockwalletStorageMockRecorder) TransferMoney(ctx, info interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransferMoney", reflect.TypeOf((*MockwalletStorage)(nil).TransferMoney), ctx, info)
}

// GetOperations mocks base method
func (m *MockwalletStorage) GetOperations(ctx context.Context, filter storage.Filter) ([]storage.Operation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOperations", ctx, filter)
	ret0, _ := ret[0].([]storage.Operation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOperations indicates an expected call of GetOperations
func (mr *MockwalletStorageMockRecorder) GetOperations(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOperations", reflect.TypeOf((*MockwalletStorage)(nil).GetOperations), ctx, filter)
}
